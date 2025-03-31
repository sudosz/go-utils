# Go Utils - High-Performance Toolkit

[![Go Reference](https://pkg.go.dev/badge/github.com/sudosz/go-utils.svg)](https://pkg.go.dev/github.com/sudosz/go-utils)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/sudosz/go-utils)](https://goreportcard.com/report/github.com/sudosz/go-utils)
[![Coverage Status](https://coveralls.io/repos/github/sudosz/go-utils/badge.svg?branch=main)](https://coveralls.io/github/sudosz/go-utils?branch=main)

**Go Utils** is a curated collection of high-performance, production-ready Go utility functions and packages. Designed for speed, efficiency, and robustness, these utilities have been rigorously tested and optimized for demanding applications.

**Key Features:**

* **Performance-Focused:** Leveraging unsafe optimizations and efficient algorithms for maximum speed.
* **Production-Ready:** Battle-tested across numerous production environments.
* **Comprehensive Suite:** Covering a wide range of utility needs from string manipulation to concurrency management.
* **Well-Documented:** Ensuring reliability and ease of use.

## Packages

| Package    | Description                                                                 |
| ---------- | --------------------------------------------------------------------------- |
| `bytes`    | Optimized byte manipulation utilities.                                       |
| `cache`    | High-performance caching implementations.                                    |
| `channel`  | Robust channel utilities with comprehensive testing.                             |
| `gopool`   | Efficient goroutine pool management for concurrent tasks.                      |
| `ints`     | Optimized integer manipulation functions.                                     |
| `iter`     | Flexible iterator implementations for data processing.                          |
| `net`      | Networking utilities, including HTTP client helpers and user agent parsing. |
| `pool`     | Versatile pool implementations (LRU, limited, recycler) for resource management. |
| `slices`   | Optimized slice manipulation utilities for common operations.              |
| `strings`  | String manipulation with unsafe optimizations for critical performance paths. |
| `terminal` | ANSI terminal utilities for rich command-line interfaces.                 |

## Installation

```bash
go get [github.com/sudosz/go-utils](https://github.com/sudosz/go-utils)
```

## Usage

Import the specific package you need:

```go
import (
        "fmt"
        "[github.com/sudosz/go-utils/strings](https://github.com/sudosz/go-utils/strings)"
)

func main() {
        result := strings.ToLower("EXAMPLE STRING")
        fmt.Println(result) // Output: example string
}
```

Refer to the individual package documentation on [pkg.go.dev](https://pkg.go.dev/github.com/sudosz/go-utils) for detailed usage examples and API references.

## Contributing

Contributions are highly encouraged! Please follow these guidelines:

1.  **Fork** the repository.
2.  Create a **feature branch** (`git checkout -b feature/your-feature`).
3.  Implement your changes and write **thorough tests**.
4.  Ensure code adheres to `gofmt` and `golint`.
5.  Submit a **pull request** with a clear description of your changes.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.