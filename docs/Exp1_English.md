# In-depth Exploration of Call by Value/Reference Performance in Different Languages

## Experiment Design

* In this experiment, I will use both Call by Value and Call by Reference to pass a struct containing 10,000,000 elements to a function (simulating the passing of large data). Each method will be executed 100 times.
* The `CallByValueTest()` function is responsible for calling the `CallByValue()` function 100 times, passing large data as an argument using Call by Value.
* The `CallByReferenceTest()` function is responsible for calling the `CallByReference()` function 100 times, passing large data as an argument using Call by Reference.
* During this process, the CPU and memory usage will be recorded.

## Golang

* In the Golang directory of this project, execute the following command to observe the experimental results (it's best to clear the previous experiment's results first).

### Run the benchmark test

```bash
go test -benchmem -bench . -memprofile=memout -cpuprofile=cpuout
```

![alt text](image.png)

* Observing the experimental results:
    * In terms of execution time, Call by Reference is approximately 2.65 times faster than Call by Value. This shows that Call by Reference is more efficient.
    * In terms of memory usage, there is little difference between Call by Reference and Call by Value. This could be because the memory used by Call by Value is released after the function call ends, resulting in no significant difference in total memory usage.
    * Regarding the number of memory allocations, Call by Reference allocates memory significantly fewer times than Call by Value. More memory allocations usually lead to increased execution time, which is further reflected in the difference in execution time.

### Analyze the benchmark test's CPU usage

```bash
go tool pprof cpuout
top -cum
```

![alt text](image-3.png)

* Observing the experimental results:
  * `CallByValueTest()` direct statement time is 840ms, total time is 1530ms.
  * `CallByReferenceTest()` direct statement time is 0ms, total time is 650ms.
  * `CallByValue()` direct statement time is 680ms, total time is 680ms.
  * `CallByReference()` direct statement time is 630ms, total time is 630ms.
  * It can be observed that the total time for `CallByValueTest()` is approximately 2.35 times that of `CallByReferenceTest()`, closely matching the benchmark test results.
  * It can also be observed that the time directly consumed by `CallByValue()` and `CallByReference()` functions themselves is almost identical.
  * Therefore, the time difference between `CallByValueTest()` and `CallByReferenceTest()` mainly lies in the time consumed by the `CallByValueTest()` function itself (840ms).

### Analyze the benchmark test's memory usage

```bash
go tool pprof memout
top -cum
```

![alt text](image-1.png)

* Observing the experimental results:
  * `CallByValueTest()` direct statement memory usage is 76.30MB, total memory usage is 76.30MB.
  * `CallByReferenceTest()` direct statement memory usage is 76.30MB, total memory usage is 76.30MB.
  * It can be seen that there is no significant difference in memory usage, which is consistent with the results observed in the benchmark test.

## Summary

* In Golang, Call by Reference does not require repeated memory allocation, thus it has higher performance than Call by Value.
* In terms of total memory consumption, since memory can be reclaimed and reused, there is little difference between Call by Reference and Call by Value.
