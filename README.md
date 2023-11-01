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
### 正常情况

```

mode:multi, clientNum:1e2, routineNum: 606-1e4, paylaod: 13字节 总效率(send+recv)/2: 3.5e6/s  ws: 2.0e5
mode:multi, clientNum:1e4, routinueNum:,    payload: 13字节 总效率(send+recv)/2:
mode: 
```
### 异常情况
```
server recv webtransport conn: 252 routines: 12870
server recv webtransport conn: 253 routines: 12872
server recv webtransport conn: 254 routines: 12874
```

## 可行性分析
- 架构
  契合当前im系统结构，只需要少量改动就能满足需求  
- 功能 
  不稳定,go的webtransport实现仍有bug,协程数暴涨，连接不能大量建立
  转发效率远高于websocket，相同cpu消耗下高了一个数量级
- 安全性
  只支持tls,不支持http，内网使用
- 环境
  需要负载均衡器支持quic、http3
- 兼容性问题
  当前草案阶段

