Containers: 5
Server type: Go
Total requests: 3 million
Parallel Requests: 500 per container
EOF Errors: 1
Force Closed: 0
ECONNABORTED: 0
syscall ECONNRESET: 207
ECONNREFUSED: 1380

The goal of this experiment was to add back parallel containers on a single node.  Reproduced ECONNRESETs.  Next I'll try splitting server and client.