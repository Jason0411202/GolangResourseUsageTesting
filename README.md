# GolangResourseUsageTesting
## Golang
### 執行 benchmark
```bash
go test -benchmem -bench . -memprofile=memout -cpuprofile=cpuout
```

### 分析 cpu usage
```bash
go tool pprof cpuout
top -cum
```

### 分析 memory usage
```bash
go tool pprof memout
top -cum
```