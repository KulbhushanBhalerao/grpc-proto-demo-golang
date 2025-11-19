# Go Basics for Understanding This gRPC Program

A beginner-friendly guide to the Go language concepts used in this gRPC demo project.

## 📚 Table of Contents

1. [Package System](#-package-system)
2. [Imports](#-imports)
3. [Functions](#-functions)
4. [Structs and Methods](#-structs-and-methods)
5. [Pointers](#-pointers)
6. [Error Handling](#-error-handling)
7. [Interfaces](#-interfaces)
8. [Context](#-context)
9. [Goroutines and Concurrency](#-goroutines-and-concurrency)
10. [Defer Statement](#-defer-statement)
11. [For Loops](#-for-loops)
12. [String Formatting](#-string-formatting)

---

## 📦 Package System

### What is a Package?

Every Go file starts with a `package` declaration. It's like a namespace that organizes code.

```go
package main
```

**Rules**:
- `package main` - Special package for executable programs (must have a `main()` function)
- `package greeting` - Regular package for reusable libraries
- One package per directory (all files in same directory must have same package name)

**In this project**:
```go
// server/main.go and client/main.go
package main  // These are executables

// proto/greeting.pb.go
package greeting  // This is a library package
```

---

## 📥 Imports

Import other packages to use their functionality.

```go
import (
    "fmt"           // Standard library: formatting
    "log"           // Standard library: logging
    "context"       // Standard library: context management
    
    // External packages (downloaded via go get)
    "google.golang.org/grpc"
    
    // Alias: Use 'pb' instead of 'greeting' 
    pb "github.com/KulbhushanBhalerao/grpc-proto-demo-golang/proto"
)
```

**Common patterns**:
- `"fmt"` - Standard library (no domain)
- `"google.golang.org/grpc"` - External package (has domain)
- `pb "package/path"` - Import with alias
- `_ "package"` - Import for side effects only (not used in this project)

**In this project**:
```go
// We use 'pb' as shorthand for the protobuf package
client := pb.NewGreetingServiceClient(conn)
```

---

## 🔧 Functions

### Basic Function Syntax

```go
func functionName(parameter type) returnType {
    // function body
    return value
}
```

### Examples from this project:

**1. Simple function**
```go
func main() {
    // Entry point of program
    // No parameters, no return value
}
```

**2. Function with parameters and return value**
```go
func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
    // Takes context and request
    // Returns response and error
    response := &pb.HelloResponse{
        Message: "Hello!",
        Count: 1,
    }
    return response, nil  // nil means no error
}
```

**3. Multiple return values**
```go
conn, err := grpc.NewClient("localhost:50051", options)
// Returns both connection AND error
// Common Go pattern: last return value is error
```

---

## 🏗️ Structs and Methods

### Structs (Data Structures)

Think of structs as classes without inheritance.

```go
// Define a struct type
type Person struct {
    Name string
    Age  int
}

// Create an instance
person := Person{
    Name: "Alice",
    Age:  30,
}

// Access fields
fmt.Println(person.Name)  // "Alice"
```

**In this project**:
```go
// server/main.go
type server struct {
    pb.UnimplementedGreetingServiceServer  // Embedded struct
}
```

### Methods

Methods are functions attached to structs.

```go
// Method syntax: func (receiver) methodName(params) returns
func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
    // 's' is the receiver (like 'this' or 'self')
    // '*server' means pointer receiver
    return &pb.HelloResponse{Message: "Hello"}, nil
}
```

**Key points**:
- `(s *server)` - Receiver (the struct instance)
- `*server` - Pointer receiver (can modify the struct)
- `server` - Value receiver (cannot modify, gets a copy)

---

## 👉 Pointers

Pointers store memory addresses instead of values.

### Syntax

```go
// & = "address of"
// * = "value at address" (dereference)

var x int = 42
var ptr *int = &x    // ptr points to x's address

fmt.Println(x)       // 42 (value)
fmt.Println(&x)      // 0xc0000... (address)
fmt.Println(ptr)     // 0xc0000... (same address)
fmt.Println(*ptr)    // 42 (dereference: value at address)
```

### Why Use Pointers?

1. **Efficiency** - Pass large structs by reference instead of copying
2. **Modification** - Allow functions to modify the original value
3. **nil** - Can represent "no value"

**In this project**:
```go
// Creating pointer to struct
response := &pb.HelloResponse{  // & creates pointer
    Message: "Hello",
}

// Method with pointer receiver
func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
    // req is a pointer to HelloRequest
    // We return a pointer to HelloResponse
}

// Accessing fields (Go auto-dereferences)
name := req.GetName()  // No need for req->GetName() like in C++
```

---

## ⚠️ Error Handling

Go doesn't have exceptions. Instead, functions return errors.

### Pattern

```go
result, err := someFunction()
if err != nil {
    // Handle error
    log.Fatalf("Error: %v", err)  // Print and exit
    return err                     // Or return error to caller
}
// Continue if no error
```

**In this project**:
```go
// Client connection
conn, err := grpc.NewClient("localhost:50051", options)
if err != nil {
    log.Fatalf("Failed to connect: %v", err)
}

// RPC call
response, err := client.SayHello(ctx, request)
if err != nil {
    log.Fatalf("Error calling SayHello: %v", err)
}
```

**Common error values**:
- `nil` - No error (success)
- `io.EOF` - End of file/stream
- Custom errors - Created with `fmt.Errorf()` or `errors.New()`

---

## 🔌 Interfaces

Interfaces define behavior (methods) without implementation.

### Concept

```go
// Define interface
type Speaker interface {
    Speak() string
}

// Any type with Speak() method implements Speaker
type Dog struct{}
func (d Dog) Speak() string { return "Woof!" }

type Cat struct{}
func (c Cat) Speak() string { return "Meow!" }

// Both Dog and Cat implement Speaker (implicitly!)
```

**Key point**: Go interfaces are **implicit**. No need to declare "implements Speaker".

**In this project**:
```go
// Generated in greeting_grpc.pb.go
type GreetingServiceServer interface {
    SayHello(context.Context, *HelloRequest) (*HelloResponse, error)
    SayHelloMultiple(*HelloRequest, GreetingService_SayHelloMultipleServer) error
}

// Our server implements this interface
type server struct {
    pb.UnimplementedGreetingServiceServer
}

// By defining these methods, server automatically implements the interface
func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) { ... }
func (s *server) SayHelloMultiple(req *pb.HelloRequest, stream ...) error { ... }
```

---

## 🎯 Context

Context carries deadlines, cancellation signals, and request-scoped values.

### Why Context?

- **Timeouts** - Automatically cancel long-running operations
- **Cancellation** - Stop work when client disconnects
- **Values** - Pass request metadata

### Common Patterns

```go
// 1. Background context (never canceled)
ctx := context.Background()

// 2. Timeout context (canceled after duration)
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()  // Always call cancel to free resources

// 3. Cancel context (manual cancellation)
ctx, cancel := context.WithCancel(context.Background())
// Call cancel() when you want to stop
```

**In this project**:
```go
// Client sets timeout for RPC call
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

response, err := client.SayHello(ctx, request)
// If server doesn't respond in 5 seconds, ctx cancels the call
```

---

## 🔄 Goroutines and Concurrency

Goroutines are lightweight threads managed by Go runtime.

### Basic Syntax

```go
// Normal function call (blocking)
doSomething()

// Goroutine (non-blocking, runs concurrently)
go doSomething()
```

### Example

```go
func main() {
    // Start 3 goroutines
    go printNumbers()
    go printLetters()
    go printSymbols()
    
    // Wait for them (in real code, use sync.WaitGroup or channels)
    time.Sleep(2 * time.Second)
}
```

**In this project**:
- gRPC server handles each client connection in a separate goroutine automatically
- You don't see explicit `go` keywords, but gRPC uses them internally
- Server can handle multiple clients simultaneously

**Not heavily used in this simple demo, but important for:**
- Handling multiple client connections
- Background tasks
- Streaming operations

---

## ⏱️ Defer Statement

`defer` delays execution of a function until the surrounding function returns.

### Syntax

```go
func example() {
    defer fmt.Println("Third")   // Executes when function returns
    fmt.Println("First")
    fmt.Println("Second")
}
// Output: First, Second, Third
```

### Common Uses

**1. Cleanup resources**
```go
file, err := os.Open("file.txt")
if err != nil {
    return err
}
defer file.Close()  // Ensures file is closed when function exits
```

**2. Unlock mutexes**
```go
mutex.Lock()
defer mutex.Unlock()  // Ensures unlock happens
```

**In this project**:
```go
// Client closes connection when main() exits
conn, err := grpc.NewClient("localhost:50051", options)
defer conn.Close()  // Cleanup connection

// Cancel context when done
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()  // Free context resources
```

**Multiple defers**:
```go
defer fmt.Println("First defer")
defer fmt.Println("Second defer")
defer fmt.Println("Third defer")
// Output: Third defer, Second defer, First defer (LIFO: Last In First Out)
```

---

## 🔁 For Loops

Go has only one loop keyword: `for` (no `while` or `do-while`).

### Variations

**1. Traditional for loop**
```go
for i := 0; i < 5; i++ {
    fmt.Println(i)  // 0, 1, 2, 3, 4
}
```

**2. While-style loop**
```go
i := 0
for i < 5 {
    fmt.Println(i)
    i++
}
```

**3. Infinite loop**
```go
for {
    // Runs forever (until break or return)
    if condition {
        break  // Exit loop
    }
}
```

**4. Range loop (iterate over collections)**
```go
nums := []int{10, 20, 30}
for index, value := range nums {
    fmt.Printf("Index: %d, Value: %d\n", index, value)
}
```

**In this project**:
```go
// Infinite loop to receive streaming responses
for {
    response, err := stream.Recv()
    if err == io.EOF {
        break  // Exit when stream ends
    }
    if err != nil {
        log.Fatalf("Error: %v", err)
    }
    fmt.Printf("Received: %s\n", response.GetMessage())
}
```

---

## 📝 String Formatting

Go uses `fmt` package for formatting strings.

### Common Functions

**1. `fmt.Println()` - Print with newline**
```go
fmt.Println("Hello", "World")  // Hello World\n
```

**2. `fmt.Printf()` - Formatted print**
```go
name := "Alice"
age := 30
fmt.Printf("Name: %s, Age: %d\n", name, age)
// Output: Name: Alice, Age: 30
```

**3. `fmt.Sprintf()` - Return formatted string (not print)**
```go
message := fmt.Sprintf("Hello, %s!", name)
// message = "Hello, Alice!"
```

### Format Verbs

| Verb | Type | Example |
|------|------|---------|
| `%s` | String | `fmt.Printf("%s", "hello")` → hello |
| `%d` | Integer | `fmt.Printf("%d", 42)` → 42 |
| `%f` | Float | `fmt.Printf("%.2f", 3.14159)` → 3.14 |
| `%t` | Boolean | `fmt.Printf("%t", true)` → true |
| `%v` | Any value (default) | `fmt.Printf("%v", anything)` |
| `%+v` | Struct with field names | `fmt.Printf("%+v", person)` |
| `%T` | Type | `fmt.Printf("%T", 42)` → int |
| `%%` | Literal % | `fmt.Printf("100%%")` → 100% |

**In this project**:
```go
// String formatting
log.Printf("Received request from: %s", req.GetName())

// Creating messages
message := fmt.Sprintf("Hello, %s! Welcome to gRPC with Go!", req.GetName())

// Multiple values
fmt.Printf("📨 Received: %s (Count: %d)\n", response.GetMessage(), response.GetCount())
```

---

## 🎓 Quick Reference for This Project

### Client Code Pattern

```go
// 1. Connect to server
conn, err := grpc.NewClient("address", options)
if err != nil { handle error }
defer conn.Close()

// 2. Create client
client := pb.NewServiceClient(conn)

// 3. Make RPC call
response, err := client.Method(context, request)
if err != nil { handle error }

// 4. Use response
fmt.Println(response.GetField())
```

### Server Code Pattern

```go
// 1. Define server struct
type server struct {
    pb.UnimplementedServiceServer
}

// 2. Implement interface methods
func (s *server) Method(ctx context.Context, req *pb.Request) (*pb.Response, error) {
    // Business logic
    return &pb.Response{Field: "value"}, nil
}

// 3. Start server
lis, _ := net.Listen("tcp", ":50051")
s := grpc.NewServer()
pb.RegisterServiceServer(s, &server{})
s.Serve(lis)
```

---

## 📚 Learn More

### Recommended Resources

- **Official Go Tour**: https://go.dev/tour/
- **Go by Example**: https://gobyexample.com/
- **Effective Go**: https://go.dev/doc/effective_go
- **Go Playground**: https://go.dev/play/ (Try code online!)

### Next Steps

1. Try modifying the server methods
2. Add your own struct types
3. Experiment with goroutines
4. Add more error handling
5. Create new RPC methods

---

## 💡 Key Takeaways

1. **Simple syntax** - Go is designed to be easy to read
2. **Explicit errors** - No hidden exceptions
3. **Pointers** - Used for efficiency and modification
4. **Interfaces** - Implicit implementation
5. **defer** - Clean up resources automatically
6. **Context** - Manage timeouts and cancellation
7. **Concurrency** - Built-in with goroutines

You now have the foundation to understand and modify this gRPC project! 🚀
