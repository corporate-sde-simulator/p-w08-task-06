# Beginner Explanatory Guide: PLATFORM-2984: Refactor load shedding and admission control

> **Task Type**: Product Task  
> **Domain/Focus**: Backend Development in Go (Golang)

---

## 1. The Goal (In-Depth Beginner Explanation)

### The Core Problem
The task at hand involves refactoring the load shedding and admission control system within our application. Load shedding is a technique used to prevent system overload by rejecting requests when the system is under heavy load. While this system is functional, it suffers from poor code quality, which can lead to maintenance challenges and potential bugs in the future. The current implementation contains "magic numbers," which are hard-coded values that lack context, making the code difficult to read and understand. Additionally, the priority logic used to determine which requests to shed is overly complex and not well documented, leading to confusion among developers.

Fixing these issues is crucial for several reasons. First, improving code quality enhances maintainability, making it easier for current and future developers to work with the code. Second, simplifying the priority logic ensures that the system behaves predictably under load, which is vital for user experience. Lastly, addressing redundant mutex double-locking will improve performance and reduce the risk of deadlocks, which can cause the application to hang or crash.

### Jargon Buster (Key Terms Explained)
* **Load Shedding**: This is a strategy used in systems to manage overload by intentionally rejecting some requests. For example, if a server is receiving too many requests, it might drop the least important ones to keep the system running smoothly.
* **Magic Numbers**: These are hard-coded values in the code that appear without explanation. For instance, if you see a number like `42` in the code, it’s unclear what it represents. Instead, using named constants like `MAX_CONNECTIONS` makes the code more readable.
* **Mutex**: Short for "mutual exclusion," a mutex is a synchronization primitive that prevents multiple threads from accessing shared resources simultaneously. For example, if two threads try to update the same variable at the same time, a mutex ensures that only one can do so at a time, preventing data corruption.
* **Refactoring**: This is the process of restructuring existing computer code without changing its external behavior. It improves the nonfunctional attributes of the software, making it easier to understand and maintain.

### Expected Outcome
After implementing the refactor, the system should maintain its current functionality while exhibiting improved code quality. 

**Before**: The code contains magic numbers, complex priority logic, and redundant mutex locks, making it difficult to read and maintain.

**After**: The code will replace magic numbers with named constants, simplify and document the priority logic, and eliminate redundant mutex locks. This will result in cleaner, more maintainable code that performs better under load.

---

## 2. Related Coding Concepts & Syntax (50% Theory, 50% Practice)

### Concept 1: Constants in Go
#### 📘 Theoretical Overview (50%)
* **Why it exists**: Constants are used to define fixed values that do not change throughout the execution of a program. They improve code readability and maintainability. Without constants, developers might use arbitrary numbers (magic numbers), which can lead to confusion and errors.
* **Key Mechanisms**: In Go, constants are defined using the `const` keyword. They can be of various types, including integers, floats, strings, and booleans. Constants are evaluated at compile time, which means they cannot be changed during runtime.

#### 💻 Syntax & Practical Examples (50%)
* **Language Syntax**:
  ```go
  const MaxConnections = 100 // Defines a constant for maximum connections
  ```
* **Real-World Application**:
  ```go
  package main

  import "fmt"

  const MaxRetries = 3 // Maximum number of retries for a request

  func main() {
      for i := 0; i < MaxRetries; i++ {
          fmt.Println("Attempting connection", i+1)
      }
  }
  ```

### Concept 2: Mutex in Go
#### 📘 Theoretical Overview (50%)
* **Why it exists**: Mutexes are essential for managing concurrent access to shared resources in multi-threaded applications. Without mutexes, multiple threads could modify the same data simultaneously, leading to inconsistent or corrupted data.
* **Key Mechanisms**: In Go, the `sync` package provides the `Mutex` type. A mutex can be locked and unlocked, ensuring that only one goroutine can access the critical section of code at a time.

#### 💻 Syntax & Practical Examples (50%)
* **Language Syntax**:
  ```go
  var mu sync.Mutex // Declare a mutex
  mu.Lock()         // Lock the mutex
  // Critical section of code
  mu.Unlock()       // Unlock the mutex
  ```
* **Real-World Application**:
  ```go
  package main

  import (
      "fmt"
      "sync"
  )

  var mu sync.Mutex
  var counter int

  func increment() {
      mu.Lock() // Lock the mutex
      counter++ // Increment the counter
      mu.Unlock() // Unlock the mutex
  }

  func main() {
      var wg sync.WaitGroup
      for i := 0; i < 10; i++ {
          wg.Add(1)
          go func() {
              defer wg.Done()
              increment()
          }()
      }
      wg.Wait()
      fmt.Println("Final counter value:", counter)
  }
  ```

---

## 3. Step-by-Step Logic & Walkthrough

1. **Step 1: Locate and Analyze the Target File**
   * Navigate to the `p-w08-task-06` folder and locate the file responsible for load shedding and admission control. This is likely named something like `load_shedding.go` or `admission_control.go`.
   * Open the file and inspect the lines of code where magic numbers are used, particularly in the logic that determines how requests are handled under load.

2. **Step 2: Input Verification & Validation**
   * Before making changes, check for edge cases in the current implementation. For example, ensure that the system can handle scenarios where the number of incoming requests is zero or extremely high.

3. **Step 3: Core Implementation / Modification**
   * Replace all magic numbers with named constants. For instance, if you see `if requests > 100`, change it to `if requests > MaxConnections`.
   * Simplify the priority logic by breaking it down into smaller, well-documented functions. Ensure that each function has a clear purpose and is easy to understand.
   * Review the mutex usage and remove any redundant double-locking. Ensure that each critical section is protected by a mutex lock only once.

4. **Step 4: Output Verification & Testing**
   * After making the changes, run the existing test suite to ensure that all tests pass. Use the command `go test ./...` in the terminal to execute all tests in the package.
   * Verify that the system behaves as expected under load by simulating various scenarios, such as high request rates.

---

## 4. Detailed Walkthrough of Test Cases

### Test Case 1: Standard / Success Case
* **Description**: This test checks if the system correctly handles a normal load of requests.
* **Inputs**:
  ```json
  {
      "requests": 50
  }
  ```
* **Step-by-Step Execution Trace**:
  1. Input values are received by the load shedding function.
  2. The function checks if the number of requests exceeds `MaxConnections`.
  3. Since 50 is less than `MaxConnections`, the function processes all requests.
  4. Returns a success message indicating that all requests were handled.

* **Expected Output**: 
  ```json
  {
      "status": "success",
      "message": "All requests processed successfully."
  }
  ```

### Test Case 2: Edge Case / Validation Fail
* **Description**: This test checks how the system behaves when the number of requests exceeds the maximum allowed connections.
* **Inputs**:
  ```json
  {
      "requests": 150
  }
  ```
* **Step-by-Step Execution Trace**:
  1. Input values are received by the load shedding function.
  2. The function checks if the number of requests exceeds `MaxConnections`.
  3. Since 150 is greater than `MaxConnections`, the function sheds excess requests.
  4. Returns a message indicating that some requests were rejected.

* **Expected Output**: 
  ```json
  {
      "status": "failure",
      "message": "Some requests were shed due to overload."
  }
  ```