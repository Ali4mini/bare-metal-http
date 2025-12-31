# Bare Metal HTTP Server

A high-performance, multi-threaded web server written from scratch in Go. 

**No frameworks. No `net/http`. No magic.**

## 📖 About The Project

I started this project to demystify backend engineering. As web developers, we often rely on frameworks (Express, Gin, Django) that abstract away the complexity of the web. 

This project is a deep dive into the **Fundamental Mechanics** of the internet. It interacts directly with the TCP layer, parses raw byte streams, manages memory manually, and implements the HTTP/1.1 protocol specification without external dependencies.

**Core Learning Objectives:**
- 🧠 Understanding TCP Handshakes and Sockets (`net` package).
- 🧵 Implementing Concurrency Patterns (Worker Pools, Channels, Goroutines).
- 🔒 Securing file access and preventing Directory Traversal attacks.
- 📄 Parsing raw HTTP text protocols manually.

## ✅ Completed Features

- [x] **TCP Listener:** Establishes raw socket connections on specific ports.
- [x] **HTTP Parser:** Manually extracts Methods, Paths, and Protocols from byte streams.
- [x] **Static File Serving:** Serves HTML, CSS, and Images from the local disk.
- [x] **MIME Type Detection:** Dynamically sets `Content-Type` headers based on file extensions.
- [x] **Robust Error Handling:** Custom 404 pages and safe error recovery.
- [x] **Security:** Implements `filepath.Clean` logic to prevent Path Traversal attacks.
- [x] **High Performance Concurrency:** Uses a **Worker Pool** pattern to limit active goroutines and prevent resource exhaustion under load.
- [x] **CLI Configuration:** Add flags to configure Port and Root Directory via command line (e.g., `--port=9000`).
- [x] **Structured Logging:** Implement a middleware-style logger to track request timing and status codes.
- [x] **Code Refactor:** Separate concerns into distinct packages (`server`, `http`, `config`).

## 🚀 Roadmap

The following features are planned for the next iteration:

- [ ] **Keep-Alive:** Support persistent TCP connections for better performance.
- [ ] Rate limiting: a rate limiter middleware
- [ ] straming large files: open the file but don't read it all at once into the memory
- [ ] incoming headers: parse incoming headers
- [ ] directory listing: list the directory if the index.html file isn't found

## 🛠️ Usage

### Prerequisites
- Go 1.20+

### Running the Server

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/bare-metal-http.git
   cd bare-metal-http
   go run main.go
   ```

2. Open your browser or `curl`:
```
# Visit http://localhost:8080

# Or use curl
curl -v http://localhost:8080/index.html
```


Built as a learning exercise to master Go and System Design.
