# hexagonal_go

This project demonstrates the Hexagonal Architecture (Ports and Adapters) pattern in Go. Hexagonal Architecture emphasizes the separation of concerns between the core business logic (the domain) and the outside world, which includes databases, user interfaces, and other systems.

## Architecture Overview

### Hexagonal Architecture: Ports and Adapters

                        +-----------------+
                        |  Primary Port   |
                        | (e.g., REST API)|
                        +--------+--------+
                                 |
                                 |
    +------------+    +----------v----------+    +------------+
    | Secondary  | <--|      Domain         |--->| Secondary  |
    |   Adapter  |    |      Logic          |    |   Adapter  |
    | (e.g., DB) |    +---------------------+    | (e.g., MQ) |
    +------------+                                +------------+

