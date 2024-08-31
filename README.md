# GolangResourseUsageTesting
## 環境配置
* Golang: 1.21.4

## 實驗


## Golang

### 執行 benchmark test
```bash
go test -benchmem -bench . -memprofile=memout -cpuprofile=cpuout
```

### 分析 tseting 的 cpu usage
```bash
go tool pprof cpuout
top -cum
```

### 分析 tseting 的 memory usage
```bash
go tool pprof memout
top -cum
```