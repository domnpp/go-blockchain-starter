# Go P2P Network Learning Project

A Go learning project that implements a simple peer-to-peer network with blockchain-inspired features. This project demonstrates networking fundamentals, concurrent programming, and distributed systems concepts while teaching core Go programming skills.

## 🎯 Learning Goals

This project is designed to help you learn **Go fundamentals** through hands-on network programming:

### **Go Core Concepts Covered**
- **Project Setup**: Go modules, package organization, dependency management
- **Concurrency**: Goroutines, concurrent network operations, background tasks
- **Synchronization**: `sync.RWMutex`, thread-safe data access, concurrent state management
- **Low-Level Networking**: TCP connections, socket programming, network address handling
- **Data Handling**: JSON serialization, struct tags, type safety, error handling
- **Command Line**: Flag parsing, argument handling

### **Network & Distributed Systems**
- Peer-to-peer networking fundamentals
- Node discovery and connection management

## 🚀 Features

- **Peer-to-Peer Networking**: Nodes can discover and connect to each other
- **Network Discovery**: Nodes share peer lists and connect to discovered peers
- **Go Concurrency**: Uses goroutines for concurrent network operations
- **JSON Communication**: Structured message passing between nodes
- **Thread-Safe Operations**: Proper synchronization for multi-threaded access

## 🏗️ Architecture

### Network Layer
- TCP-based peer-to-peer communication
- Network state synchronization
- Concurrent connection handling
- Node identification with UUIDs

### Blockchain Layer
- Currently just an idea, nothing functional is implemented.

## 📦 Installation

### Prerequisites
- Go 1.24.4 or higher

### Setup
```bash
# Clone the repository
git clone https://github.com/yourusername/go-blockchain-lecture.git
cd go-blockchain-lecture

# Install dependencies
go mod tidy

# Run the first node (listens on port 8000)
go run main.go -port 8000

# In another terminal, run a second node (connects to first)
go run main.go -port 8001 -node_url 127.0.0.1:8000
```

## 🎯 Usage

### Starting Nodes

**First Node (Network Seed):**
```bash
go run main.go -port 8000
```

**Additional Nodes:**
```bash
go run main.go -port 8001 -node_url 127.0.0.1:8000
go run main.go -port 8002 -node_url 127.0.0.1:8000
```

### Command Line Options
- `-port`: Port to listen on (default: 8000)
- `-node_url`: URL of existing node to connect to (optional)

## 📚 Go Learning Journey

This project covers essential Go concepts in a practical context:

### **1. Project Organization**
- Module management with `go.mod`
- Package structure and imports
- Dependency resolution

### **2. Concurrency Patterns**
- Goroutines for background tasks
- Concurrent network operations
- Non-blocking I/O operations

### **3. Synchronization**
- `sync.RWMutex` for thread-safe data
- Concurrent access to shared state
- Proper locking strategies

### **4. Network Programming**
- TCP socket programming
- Connection lifecycle management
- Network address parsing
- Error handling for network operations

### **5. Data Serialization**
- JSON marshaling/unmarshaling
- Struct tags for serialization
- Type-safe data exchange

## 🧪 Testing

```bash
# Run tests
go test ./...

# Run with verbose output
go test -v ./...
```

## 🤝 Contributing

This is a learning project, but contributions are welcome! Feel free to:

- Report bugs
- Suggest improvements
- Add new features
- Improve documentation
- Share learning insights

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built for educational purposes to learn Go fundamentals
- Inspired by blockchain technology
- Demonstrates practical application of Go's concurrency primitives
- Shows real-world networking and distributed systems concepts

## 🔮 Future Enhancements

- [ ] Channels for goroutine communication
- [ ] Context for cancellation and timeouts
- [ ] Unit tests and benchmarks
- [ ] Performance profiling
- [ ] Advanced Go patterns (interfaces, generics)
- [ ] Proof of Stake consensus mechanism
- [ ] Token implementation (max supply, current supply, validator rewards)
- [ ] Wallet implementation
- [ ] Web interface
- [ ] Configuration management
- [ ] Logging and monitoring

---

**Note**: This is an educational project designed to teach Go programming fundamentals through network programming. Do not use for production or financial transactions. 