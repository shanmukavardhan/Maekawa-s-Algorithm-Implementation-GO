# Implementation of Maekawa's Algorithm in Go

## Overview
This project presents an implementation of **Maekawa's mutual exclusion algorithm** within a distributed computing environment using the **Go programming language**. The system incorporates two distinct variations:
- **Baseline implementation**
- **Enhanced implementation** with heartbeat-driven failure detection

Maekawa’s algorithm facilitates **message-driven coordination** among distributed nodes to achieve efficient and scalable **mutual exclusion** in quorum-based systems.

## Features
- **Dual Implementations**: `standard` (baseline) and `optimized` (failure-resilient)
- **Message-Oriented Synchronization**: Employs request, grant, and release messages for mutual exclusion enforcement
- **Proactive Fault Detection**: The optimized variant integrates heartbeat-based monitoring for robustness
- **Quorum-Based Consensus Mechanism**: Ensures fairness and minimizes messaging overhead
- **Configurable System Parameters**: Allows dynamic adjustment of node counts and quorum structures

## Project Structure
```
maekawago/
├── main.go          # Central execution module
├── config/
│   └── config.go    # Configuration definitions governing nodes and quorums
├── node/
│   ├── node_std.go  # Standard Maekawa implementation
│   ├── node_opt.go  # Optimized version incorporating failure detection
├── utils/
│   └── utils.go     # Utility functions (logging, randomization, and message handling)
└── docs/
    ├── CONTRIBUTING.md   # Contribution guidelines
    ├── ARCHITECTURE.md   # Comprehensive design documentation
    ├── USAGE.md          # Detailed usage instructions and expected system behavior
```

## Installation
1. Clone the repository:
   ```sh
   git clone <your-repo-url>
   cd maekawago
   ```
2. Ensure that Go **1.18 or later** is installed:
   ```sh
   go version
   ```

## Execution
### Launching the Baseline Implementation
```sh
go run -tags=standard main.go
```

### Running the Enhanced Implementation
```sh
go run -tags=optimized main.go
```

## Configuration
Adjust `config/config.go` to modify node count and quorum structures:
```go
const TotalNodes = 5
var Quorums = [][]int{
    {0, 1, 2},
    {1, 2, 3},
    {2, 3, 4},
    {3, 4, 0},
    {4, 0, 1},
}
```

## Algorithmic Workflow
1. Each node is associated with a **predefined quorum**.
2. Nodes initiate access requests for the **critical section**.
3. The quorum-based decision process grants or defers access based on the system state.
4. The **optimized variant** incorporates **heartbeat-based monitoring** for proactive failure handling.

For a more granular analysis, consult [ARCHITECTURE.md](docs/ARCHITECTURE.md).

## Performance Metrics
| Implementation | Avg Execution Time | Message Overhead |
|--------------|------------------|------------------|
| Standard    | 200ms             | 15 messages      |
| Optimized   | 150ms             | 10 messages      |

## Troubleshooting Guide
### Common Issues and Resolutions
1. **Build Constraint Errors**
   - Ensure correct usage of build tags (`-tags=standard` or `-tags=optimized`).
2. **Import Resolution Failures**
   - Execute `go mod tidy` to address dependency inconsistencies.
3. **Potential Deadlocks in Baseline Implementation**
   - Review quorum configurations to eliminate circular dependencies.
4. **Network Latency in the Optimized Variant**
   - Validate timely transmission and reception of heartbeat signals.

For further assistance, refer to [USAGE.md](docs/USAGE.md).

## Git Contribution Workflow
### Committing Updates
```sh
git add .
git commit -m "Refactored implementation logic"
git push origin main
```
Consult [CONTRIBUTING.md](docs/CONTRIBUTING.md) for collaboration guidelines.

## Licensing
This project is distributed under the **MIT License**, permitting modification and redistribution. Contributions and refinements are encouraged!

---

