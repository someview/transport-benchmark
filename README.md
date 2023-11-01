# webtransport基本性能测试
## 测试环境
- server
```
CPU(s):                4
Model name:            Intel(R) Xeon(R) CPU           X5650  @ 2.67GHz
```

- client
```
the same as server
```
- 测试目标
  1. 协程资源消耗 
  2. cpu消耗
  3. 性能分析
  4. 异常分析

- 测试代码
https://github.com/someview/transport-benchmark
```
go run ./server/main.go
go run ./client/main.go
```
mode: multi, 多协程频繁调度
mode: single, 多协程，某个协程频繁调度
mode: silent，多协程，不调度

## 测试结果
```

mode:multi, clientNum:1e2, routineNum: 606, paylaod: 13字节 总效率(send+recv)/2: 3.5e6/s  ws: 2.0e5
mode:multi, clientNum:1e4, routinueNum:,    payload: 13字节 总效率(send+recv)/2:
mode: 
```
