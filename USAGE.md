## How to Use This Implementation

### Running the Standard Version
```sh
go run -tags=standard main.go
```

### Running the Optimized Version
```sh
go run -tags=optimized main.go
```

### Expected Output
#### Standard Version
```
2025/04/01 15:41:08 Node 2 is requesting critical section
2025/04/01 15:41:08 Sending Message[Request] from 2 to 2
2025/04/01 15:41:08 Sending Message[Request] from 2 to 3
2025/04/01 15:41:08 Node 2 received Message[Request] from 2 to 2
2025/04/01 15:41:08 Sending Message[Grant] from 2 to 2
2025/04/01 15:41:08 Node 2 received Message[Grant] from 2 to 2
2025/04/01 15:41:08 Node 2 received Grant from 2
```
#### Optimized Version
```
2025/04/01 15:41:51 Node 2 is requesting critical section
2025/04/01 15:41:51 Sending Message[Request] from 2 to 2
2025/04/01 15:41:51 Sending Message[Request] from 2 to 3
2025/04/01 15:42:24 Node 0 received Heartbeat from 4
2025/04/01 15:42:24 Sending Message[Heartbeat] from 2 to 3
2025/04/01 15:42:24 Sending Message[Heartbeat] from 2 to 4
2025/04/01 15:41:51 Node 3 received Message[Request] from 2 to 3
```

