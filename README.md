# Mini Device Fleet Monitor

A lightweight, concurrent HTTP service written in Go to monitor a fleet of edge devices sending periodic heartbeats.

---

## 🌟 Architecture & Design Choices

The application follows a **Clean / Layered Architecture**:

1. **Model Layer (`internal/model`)**: Defines core entities (`Device`, `FleetSummary`) and API DTOs.
2. **Repository Layer (`internal/repository`)**: Provides thread-safe, in-memory data store using Go's `sync.RWMutex`.
3. **Service Layer (`internal/service`)**: Contains business rules, including dynamic 30-second `ONLINE`/`OFFLINE` calculations.
4. **Handler Layer (`internal/handler`)**: REST endpoints parsing input and serving JSON responses.

### Key Highlights
- **Thread Safety**: Complete concurrent map protection via read-write mutexes (`sync.RWMutex`).
- **Dynamic Status Calculation**: Status is computed dynamically on lookup based on time elapsed since the last heartbeat (`<= 30s` = `ONLINE`, `> 30s` or `nil` = `OFFLINE`).
- **Zero External Dependencies**: Built entirely using Go’s Standard Library (`net/http`, `sync`, `testing`).

---

## 🚀 Getting Started

### Prerequisites
- **Go 1.22+** installed on your system, or
- **Docker**

---

## 🛠️ Build & Run Instructions

### 1. Running Server Locally
```bash
go run cmd/main.go