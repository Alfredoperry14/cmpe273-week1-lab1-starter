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


Run these commands on one terminal (also with output):
░▒▓ ~/CMPE273/cmpe273-week1-lab1-starter  main ▓▒░·······································································░▒▓ system Node ▓▒░
❯ cd go-http 
░▒▓ ~/CMPE273/cmpe273-week1-lab1-starter/go-http  main ▓▒░·······························································░▒▓ system Node ▓▒░
❯ cd service-a 
░▒▓ ~/CMPE273/cmpe273-week1-lab1-starter/go-http/service-a  main ▓▒░····························································░▒▓ system Node ▓▒░
❯ go mod init service-a
go: creating new go.mod: module service-a
go: to add module requirements and sums:
        go mod tidy
░▒▓ ~/CMPE273/cmpe273-week1-lab1-starter/go-http/service-a  main ?2 ▓▒░········································░▒▓ 13m 36s  system Node ▓▒░
❯ go run .             
2026/09/03 12:36:30 service=A listening on :8080
2026/09/03 12:36:37 service=A endpoint=/echo status=ok latency_ms=0
^Csignal: interrupt

Second terminal: 
░▒▓ ~/CMPE273/cmpe273-week1-lab1-starter  main ?1 ▓▒░····································································░▒▓ system Node ▓▒░
❯ cd go-http 
░▒▓ ~/CMPE273/cmpe273-week1-lab1-starter/go-http  main ?1 ▓▒░····························································░▒▓ system Node ▓▒░
❯ cd service-b  
░▒▓ ~/CMPE273/cmpe273-week1-lab1-starter/go-http/service-b  main ?1 ▓▒░··················································░▒▓ system Node ▓▒░
❯ go mod init service-b
go: creating new go.mod: module service-b
go: to add module requirements and sums:
        go mod tidy
░▒▓ ~/CMPE273/cmpe273-week1-lab1-starter/go-http/service-b  main ?2 ▓▒░··················································░▒▓ system Node ▓▒░
❯ go run .
2026/09/03 12:23:27 service=B listening on :8081
curl "http://127.0.0.1:8081/call-echo?msg=hello"
^Csignal: interrupt
░▒▓ ~/CMPE273/cmpe273-week1-lab1-starter/go-http/service-b  main ?2 ▓▒░········································░▒▓ 12m 53s  system Node ▓▒░
❯ go run .             
2026/09/03 12:36:23 service=B listening on :8081
2026/09/03 12:36:37 service=B endpoint=/call-echo status=ok latency_ms=0

Third terminal: 
A & B Enabled: 
░▒▓ ~/CMPE273/cmpe273-week1-lab1-starter  main ?2 ▓▒░····································································░▒▓ system Node ▓▒░
❯ curl "http://127.0.0.1:8081/call-echo?msg=hello"
{"service_a":{"echo":"hello"},"service_b":"ok"}
A Disabled: 
░▒▓ ~/CMPE273/cmpe273-week1-lab1-starter  main ?2 ▓▒░····································································░▒▓ system Node ▓▒░
❯ curl "http://127.0.0.1:8081/call-echo?msg=hello"
{"error":"Get \"http://127.0.0.1:8080/echo?msg=hello\": dial tcp 127.0.0.1:8080: connect: connection refused","service_a":"unavailable","service_b":"ok"}
B Disabled:
░▒▓ ~/CMPE273/cmpe273-week1-lab1-starter  main ?2 ▓▒░····································································░▒▓ system Node ▓▒░
❯ curl "http://127.0.0.1:8081/call-echo?msg=hello"
curl: (7) Failed to connect to 127.0.0.1 port 8081 after 0 ms: Couldn't connect to server


What makes this a distributed system is that you can curl system B 
and system A has to be turned on to get the {"echo":"hello"} message. 
When system A is off you still can successfully curl system B. 
Because system B can return a successful curl when system A is offline 
this shows a distributed system.
