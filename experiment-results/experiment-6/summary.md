Containers: 10
Server type: Go
Total requests: 72 million (overnight)
Parallel Requests: 500 per container
EOF Errors: 0
Force Closed: 0
ECONNABORTED: 0
syscall ECONNRESET: 147
ECONNREFUSED: 0

This experiment was run overnight with the suggested /proc/sys/net/netfilter/nf_conntrack_tcp_be_liberal setting flipped, as suggested here: https://kubernetes.io/blog/2019/03/29/kube-proxy-subtleties-debugging-an-intermittent-connection-reset/

We did still receive ECONNRESETs, but they seemed to be drastically reduced.