## Architectural Overview

This document outlines the design and implementation of Maekawa's Algorithm in Go, highlighting the interaction between components.

### Design Considerations
- **Quorum-Based Mutual Exclusion:** Nodes form quorums to manage access.
- **Message-Driven Approach:** Communication is handled through request, grant, and release messages.
- **Optimized Version:** Adds heartbeat messages for failure detection.

### Code Structure
```
maekawago/
├── main.go          # Entry point
├── config/
│   └── config.go    # Configurations
├── node/
│   ├── node_std.go  # Standard implementation
│   ├── node_opt.go  # Optimized version with heartbeat
└── utils/
    └── utils.go     # Logging and helper functions
```

### Sequence Flow
1. Nodes initialize and form their quorums.
2. A node requests access to the critical section.
3. Other nodes grant or defer based on quorum logic.
4. The optimized version introduces a heartbeat to track node availability.

