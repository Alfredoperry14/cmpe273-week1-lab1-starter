package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// statusRecorder wraps http.ResponseWriter so the logger can see which status
// code a handler actually wrote. Without it we could only log a hardcoded
// "ok"/"error" guess instead of the real number.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// withLogging emits exactly one line per request: service, endpoint, status,
// latency. Wrapping every route (including /health) is what satisfies the
// lab's "basic logging per request" requirement.
func withLogging(service, endpoint string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next(rec, r)
		log.Printf("service=%s endpoint=%s status=%d latency_ms=%d",
			service, endpoint, rec.status, time.Since(start).Milliseconds())
	}
}

// writeJSON sets Content-Type before WriteHeader. Order matters: headers set
// after WriteHeader are silently dropped.
func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func echo(w http.ResponseWriter, r *http.Request) {
	// ?delay=2s makes Service A deliberately slow. This exists only so the
	// timeout branch in Service B can be demonstrated -- with an instant
	// response, B's 1s client timeout can never fire and "what happens on
	// timeout?" is unanswerable.
	if d := r.URL.Query().Get("delay"); d != "" {
		wait, err := time.ParseDuration(d)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": "invalid delay (want a Go duration like 2s): " + d,
			})
			return
		}
		time.Sleep(wait)
	}

	msg := r.URL.Query().Get("msg")
	if msg == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "missing required query param: msg",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"echo": msg})
}

func main() {
	// An explicit mux instead of the package-level DefaultServeMux: nothing
	// else can register routes on this server behind our back.
	mux := http.NewServeMux()
	mux.HandleFunc("/health", withLogging("A", "/health", health))
	mux.HandleFunc("/echo", withLogging("A", "/echo", echo))

	log.Println("service=A listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
