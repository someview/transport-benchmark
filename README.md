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

## 测试结果
