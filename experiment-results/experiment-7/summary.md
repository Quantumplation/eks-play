Containers: 10
Server type: Go
Total requests: 58m
Parallel Requests: 500 per container
EOF Errors: 0
Force Closed: 0
ECONNABORTED: 0
syscall ECONNRESET: 58
ECONNREFUSED: 0

This experiment was run with the following settings tweaked, to measure the impact on ECONNRESET

echo 1 > /proc/sys/net/netfilter/nf_conntrack_tcp_be_liberal
echo 1024 65535 > /proc/sys/net/ipv4/ip_local_port_range
echo 65536 > /proc/sys/net/core/netdev_max_backlog
echo 65536 > /proc/sys/net/core/somaxconn

These are largely based on the recommendations here: https://edenmal.moe/post/2019/My-sysctl-Parameters/
and here: https://serverfault.com/questions/875035/sane-value-for-net-ipv4-tcp-max-syn-backlog-in-sysctl-conf

These settings seemed to help significantly!