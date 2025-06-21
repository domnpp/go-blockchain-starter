# App code

## Internals

main() → startTCPListener() → handleIncomingConnection(conn) → nodeCommunication(conn) → loop (switch, case)

### Second node

main() → connectToNode() → go nodeCommunication(conn) → loop (switch, case)

## Flow

| Node 1            | Node 2           |
|-------------------|------------------|
| Listen            |                  |
|                   | < Connect        |
| Identify >        |                  |
|                   | < Identification |
|                   | < Network        |
| AddPeer()         |                  |
| networkResponse > |                  |
