# I/O Model Performance Comparison

This report compares two server implementation approaches under the same benchmark workload:

1. **I/O multiplexing model**, using one thread for all the connections - [server.go](./internal/server/server.go)
2. **Blocking-I/O model**, using one thread per connection - [TcpServer.go](./TcpServer/TcpServer.go)

> The second implementation creates one Go goroutine per connection. A goroutine is scheduled by the Go runtime and is not necessarily a dedicated operating-system thread.

## Benchmark setup

To do the benchmark for both implementation, I used the official benchmark tool from [redis](https://github.com/redis/redis).

I rebuilt the both servers from the current source before testing. Each server was run separately on port 3000. The following command was executed 10 times against each implementation:

```bash
./redis-benchmark -n 10000 -t ping_mbulk -c 200 -h localhost -p 3000 --csv
```

| Setting | Value |
|---|---:|
| Runs per implementation | 10 |
| Requests per run | 10,000 |
| Concurrent clients | 200 |
| Workload | `PING_MBULK` |
| Target | `localhost:3000` |

The tables below report values emitted by `redis-benchmark`. Average percentile values are arithmetic means of the percentile reported in each run; they are not percentiles recomputed from a combined latency distribution.

## I/O multiplexing results

| Run | Requests/sec | Avg (ms) | Min (ms) | p50 (ms) | p95 (ms) | p99 (ms) | Max (ms) |
|---:|---:|---:|---:|---:|---:|---:|---:|
| 1 | 68,493.15 | 2.023 | 0.768 | 1.447 | 5.703 | 7.743 | 8.367 |
| 2 | 83,333.34 | 1.450 | 0.336 | 1.383 | 2.831 | 3.447 | 3.823 |
| 3 | 90,090.09 | 1.370 | 0.312 | 1.367 | 2.271 | 2.743 | 2.831 |
| 4 | 107,526.88 | 1.175 | 0.560 | 1.183 | 1.735 | 2.047 | 2.191 |
| 5 | 91,743.12 | 1.340 | 0.520 | 1.375 | 1.943 | 2.287 | 2.375 |
| 6 | 98,039.22 | 1.319 | 0.392 | 1.327 | 2.471 | 2.687 | 2.767 |
| 7 | 76,923.08 | 1.645 | 0.528 | 1.415 | 3.007 | 3.167 | 3.271 |
| 8 | 83,333.34 | 1.526 | 0.488 | 1.391 | 2.295 | 2.431 | 3.575 |
| 9 | 100,000.00 | 1.278 | 0.408 | 1.343 | 1.959 | 2.399 | 2.487 |
| 10 | 87,719.30 | 1.411 | 0.576 | 1.407 | 2.015 | 2.215 | 2.279 |
| **Average** | **88,720.15** | **1.454** | **0.489** | **1.364** | **2.623** | **3.117** | **3.397** |

## Blocking I/O results

| Run | Requests/sec | Avg (ms) | Min (ms) | p50 (ms) | p95 (ms) | p99 (ms) | Max (ms) |
|---:|---:|---:|---:|---:|---:|---:|---:|
| 1 | 64,935.07 | 2.166 | 0.552 | 1.487 | 6.415 | 8.527 | 8.967 |
| 2 | 71,428.57 | 2.039 | 0.520 | 1.463 | 3.727 | 4.135 | 4.967 |
| 3 | 68,493.15 | 1.864 | 0.536 | 1.527 | 3.239 | 3.415 | 3.863 |
| 4 | 70,422.53 | 1.970 | 0.680 | 1.407 | 4.311 | 4.671 | 4.927 |
| 5 | 73,529.41 | 1.899 | 0.672 | 1.399 | 3.839 | 4.079 | 4.303 |
| 6 | 71,942.45 | 1.906 | 0.120 | 1.455 | 3.415 | 3.943 | 4.527 |
| 7 | 74,626.87 | 1.843 | 0.480 | 1.423 | 3.063 | 3.215 | 3.407 |
| 8 | 69,930.07 | 2.064 | 0.640 | 1.471 | 3.871 | 4.359 | 4.575 |
| 9 | 74,074.07 | 1.843 | 0.416 | 1.415 | 3.023 | 3.303 | 3.615 |
| 10 | 74,074.07 | 1.889 | 0.608 | 1.399 | 3.463 | 3.767 | 4.727 |
| **Average** | **71,345.63** | **1.948** | **0.522** | **1.445** | **3.837** | **4.341** | **4.788** |

## Average performance comparison

| Metric | I/O multiplexing | Blocking I/O | I/O multiplexing difference |
|---|---:|---:|---:|
| Throughput (requests/sec) | 88,720.15 | 71,345.63 | **24.35% higher** |
| Average latency (ms) | 1.454 | 1.948 | **25.39% lower** |
| Minimum latency (ms) | 0.489 | 0.522 | **6.43% lower** |
| p50 latency (ms) | 1.364 | 1.445 | **5.59% lower** |
| p95 latency (ms) | 2.623 | 3.837 | **31.63% lower** |
| p99 latency (ms) | 3.117 | 4.341 | **28.21% lower** |
| Maximum latency (ms) | 3.397 | 4.788 | **29.06% lower** |

Higher throughput is better; lower latency is better.

## Conclusion

Across these 10 runs, the I/O-multiplexing implementation averaged **88,720.15 requests per second**, compared with **71,345.63 requests per second** for the Blocking I/O implementation. I/O multiplexing therefore delivered **24.35% higher throughput**, or approximately **1.24 times the throughput**.

I/O multiplexing also produced lower averages for every reported latency metric. Mean request latency was **25.39% lower**, while averaged p95 and p99 latency were **31.63%** and **28.21% lower**, respectively.

This outcome is consistent with the expected overhead of a thread-per-connection design under high concurrency. That design can incur additional thread scheduling and context-switching costs, whereas I/O multiplexing allows one event loop to manage many connections efficiently.

## Interpretation notes

These results describe these specific implementations and this local `PING_MBULK` workload; they do not establish that one concurrency model is always faster.
