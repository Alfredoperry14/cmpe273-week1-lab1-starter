# CMPE 273 – Week 1 Lab 1: Your First Distributed System (Starter)

This starter provides two implementation tracks:
- `python-http/` (Flask + requests)
- `go-http/` (net/http)

Pick **one** track for Week 1.

## Lab Goal
Build **two services** that communicate over the network:
- **Service A** (port 8080): `/health`, `/echo?msg=...`
- **Service B** (port 8081): `/health`, `/call-echo?msg=...` calls Service A

Minimum requirements:
- Two independent processes
- HTTP (or gRPC if you choose stretch)
- Basic logging per request (service name, endpoint, status, latency)
- Timeout handling in Service B
- Demonstrate independent failure (stop A; B returns 503 and logs error)

## Deliverables
1. Repo link
2. README updates:
   - how to run locally
   - success + failure proof (curl output or screenshot)
   - 1 short paragraph: “What makes this distributed?”
## How to Run Locally

**One-time setup**

```bash
cd go-http/service-a && go mod init service-a
cd ../service-b && go mod init service-b
```

**Terminal 1 — Service A (port 8080)**

```bash
cd go-http/service-a
go run .
```

**Terminal 2 — Service B (port 8081)**

```bash
cd go-http/service-b
go run .
```

**Terminal 3 — test**

```bash
# health checks
curl -s -w "\nHTTP %{http_code}\n" "http://127.0.0.1:8080/health"
curl -s -w "\nHTTP %{http_code}\n" "http://127.0.0.1:8081/health"

# success
curl -s -w "\nHTTP %{http_code}\n" "http://127.0.0.1:8081/call-echo?msg=hello"

# failure: stop Service A with Ctrl+C in Terminal 1, then rerun
curl -s -w "\nHTTP %{http_code}\n" "http://127.0.0.1:8081/call-echo?msg=hello"
```
<h1>A</h1>
<img width="521" height="347" alt="image" src="https://github.com/user-attachments/assets/d5ed9d69-40d7-4045-bd74-8177c0b588ee" />
<h1>B</h1>
<img width="1031" height="442" alt="image" src="https://github.com/user-attachments/assets/ba1db122-6fe2-498a-9555-7219f7e3efab" />
<h1>Output</h1>
<img width="1023" height="163" alt="image" src="https://github.com/user-attachments/assets/ab3cb42d-0160-45a5-968a-c89d0c28ecd2" />

<h2>Paragraph</h2>
<p>What makes this a distributed system is that you can curl system B 
and system A has to be turned on to get the {"echo":"hello"} message. 
When system A is off you still can successfully curl system B. 
Because system B can return a successful curl when system A is offline 
this shows a distributed system.</p>

