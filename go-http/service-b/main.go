package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"time"
)

// Where Service A lives. Hardcoded for this lab; in a real system this would
// come from config or service discovery, because "the network address of my
// dependency is fixed and known" is one of the eight fallacies.
const serviceAEcho = "http://127.0.0.1:8080/echo"

// The timeout is the whole point of this client. Without it, a hung Service A
// would hang every Service B request too, and one slow dependency would take
// down the entire system. One second is deliberately short so the failure is
// fast and obvious.
var client = &http.Client{Timeout: 1 * time.Second}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func withLogging(service, endpoint string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next(rec, r)
		log.Printf("service=%s endpoint=%s status=%d latency_ms=%d",
			service, endpoint, rec.status, time.Since(start).Milliseconds())
	}
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func callEcho(w http.ResponseWriter, r *http.Request) {
	// url.Values escapes the value for us. Interpolating msg straight into the
	// URL string breaks on spaces and lets an "&" inside msg inject extra
	// query params into the request we send to Service A.
	q := url.Values{}
	q.Set("msg", r.URL.Query().Get("msg"))
	if d := r.URL.Query().Get("delay"); d != "" {
		q.Set("delay", d) // passed through so we can demonstrate the timeout
	}
	target := serviceAEcho + "?" + q.Encode()

	resp, err := client.Get(target)
	if err != nil {
		// Two very different failures land here and they need different fixes,
		// so the log has to tell them apart: "A is gone" (connection refused)
		// vs "A is alive but too slow" (our deadline fired).
		reason := "connection_failed"
		var urlErr *url.Error
		if errors.As(err, &urlErr) && urlErr.Timeout() {
			reason = "timeout"
		}
		log.Printf("service=B endpoint=/call-echo upstream=A reason=%s error=%q", reason, err)

		writeJSON(w, http.StatusServiceUnavailable, map[string]any{
			"service_b": "ok",
			"service_a": "unavailable",
			"reason":    reason,
			"error":     err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	// A answered, but "answered" is not "succeeded". Without this check a 400
	// or 500 from A would be reported to our own caller as a success with an
	// empty body attached.
	if resp.StatusCode != http.StatusOK {
		log.Printf("service=B endpoint=/call-echo upstream=A reason=upstream_status upstream_status=%d", resp.StatusCode)
		writeJSON(w, http.StatusBadGateway, map[string]any{
			"service_b":       "ok",
			"service_a":       "error",
			"reason":          "upstream_status",
			"upstream_status": resp.StatusCode,
		})
		return
	}

	var data map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		// A returned 200 but the body was not the JSON we expect. Trusting the
		// status code alone is not enough.
		log.Printf("service=B endpoint=/call-echo upstream=A reason=bad_body error=%q", err)
		writeJSON(w, http.StatusBadGateway, map[string]any{
			"service_b": "ok",
			"service_a": "error",
			"reason":    "bad_body",
			"error":     err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"service_b": "ok",
		"service_a": data,
	})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", withLogging("B", "/health", health))
	mux.HandleFunc("/call-echo", withLogging("B", "/call-echo", callEcho))

	log.Println("service=B listening on :8081")
	log.Fatal(http.ListenAndServe(":8081", mux))
}
