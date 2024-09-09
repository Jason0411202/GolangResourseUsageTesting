# Exploring the Performance of Golang's Call by Value/Reference under Memory Limitations Based on Experiment 1

## Experiment Design
* The only difference from Experiment 1 is that during the benchmark test run using Docker, the memory limit will be restricted with the `-m` parameter.
* This experiment aims to test the performance of Golang's Call by Value/Reference under different memory limits.

(For Windows environment, set memory limit to 500 MB)
```bash
docker run -m 500m --rm -v ${PWD}:/app myapp
```

(For Linux environment, set memory limit to 500 MB)
```bash
docker run -m 500m --rm -v "$(pwd):/app" myapp
```

## No Memory Limit (Using results from Experiment 1 as a control group)
cpu usage
```
(pprof) Total: 5.84s
ROUTINE ======================== test.AllTest in /app/main.go
     2.14s      3.86s (flat, cum) 66.10% of Total
         .          .     29:func AllTest() {
      10ms       10ms     30:   data := myStruct{}
         .          .     31:
         .          .     32:   for i := 0; i < 100; i++ {
     1.04s      1.04s     33:           CallByValue(data)
         .          .     34:   }
         .          .     35:   for i := 0; i < 100; i++ {
         .          .     36:           CallByReference(&data)
         .          .     37:   }
         .          .     38:
         .          .     39:   for i := 0; i < 100; i++ {
     1.09s      2.81s     40:           CallByValue_Escape(data)
         .          .     41:   }
         .          .     42:   for i := 0; i < 100; i++ {
         .          .     43:           CallByReference_Escape(&data)
         .          .     44:   }
         .          .     45:}
```

memory usage
```
Total: 7.53GB
ROUTINE ======================== test.AllTest in /app/main.go
   76.30MB     7.53GB (flat, cum)   100% of Total
         .          .     29:func AllTest() {
   76.30MB    76.30MB     30:   data := myStruct{}
         .          .     31:
         .          .     32:   for i := 0; i < 100; i++ {
         .          .     33:           CallByValue(data)
         .          .     34:   }
         .          .     35:   for i := 0; i < 100; i++ {
         .          .     36:           CallByReference(&data)
         .          .     37:   }
         .          .     38:
         .          .     39:   for i := 0; i < 100; i++ {
         .     7.45GB     40:           CallByValue_Escape(data)
         .          .     41:   }
         .          .     42:   for i := 0; i < 100; i++ {
         .          .     43:           CallByReference_Escape(&data)
         .          .     44:   }
         .          .     45:}
```

## Memory Limit Set to 245 MB
cpu usage
```
Total: 20.26s
ROUTINE ======================== test.AllTest in /app/main.go
     5.11s     18.13s (flat, cum) 89.49% of Total
         .          .     29:func AllTest() {
         .          .     30:   data := myStruct{}
         .          .     31:
         .          .     32:   for i := 0; i < 100; i++ {
     1.28s      1.28s     33:           CallByValue(data)
         .          .     34:   }
         .          .     35:   for i := 0; i < 100; i++ {
         .          .     36:           CallByReference(&data)
         .          .     37:   }
         .          .     38:
         .          .     39:   for i := 0; i < 100; i++ {
     3.83s     16.85s     40:           CallByValue_Escape(data)
         .          .     41:   }
         .          .     42:   for i := 0; i < 100; i++ {
         .          .     43:           CallByReference_Escape(&data)
         .          .     44:   }
         .          .     45:}
```

memory usage
```
Total: 7.53GB
ROUTINE ======================== test.AllTest in /app/main.go
   76.30MB     7.53GB (flat, cum)   100% of Total
         .          .     29:func AllTest() {
   76.30MB    76.30MB     30:   data := myStruct{}
         .          .     31:
         .          .     32:   for i := 0; i < 100; i++ {
         .          .     33:           CallByValue(data)
         .          .     34:   }
         .          .     35:   for i := 0; i < 100; i++ {
         .          .     36:           CallByReference(&data)
         .          .     37:   }
         .          .     38:
         .          .     39:   for i := 0; i < 100; i++ {
         .     7.45GB     40:           CallByValue_Escape(data)
         .          .     41:   }
         .          .     42:   for i := 0; i < 100; i++ {
         .          .     43:           CallByReference_Escape(&data)
         .          .     44:   }
         .          .     45:}
```

* Based on multiple experiments, a memory limit of 245 MB is approximately the limit for the CallByValue_Escape() function; further reduction in memory limit causes the entire program to be forcefully terminated.
    ```
    docker run -m 240m --rm -v ${PWD}:/app myapp

    signal: killed
    FAIL    test    6.276s
    ```
* From the experimental results, it can be observed that the CallByValue_Escape() function becomes the biggest performance bottleneck, likely due to frequent garbage collection caused by insufficient memory, leading to a noticeable performance drop.
* To continue the experiment, the CallByValue_Escape() portion will be commented out, and other functions will be stress-tested.

## Memory Limit Set to 130 MB
cpu usage
```
Total: 17.72s
ROUTINE ======================== test.AllTest in /app/main.go
    17.08s     17.40s (flat, cum) 98.19% of Total
         .          .     29:func AllTest() {
      10ms       10ms     30:   data := myStruct{}
         .          .     31:
         .          .     32:   for i := 0; i < 100; i++ {
    17.07s     17.39s     33:           CallByValue(data)
         .          .     34:   }
         .          .     35:   for i := 0; i < 100; i++ {
         .          .     36:           CallByReference(&data)
         .          .     37:   }
         .          .     38:
```

memory usage
```
Total: 77.18MB
ROUTINE ======================== test.AllTest in /app/main.go
   76.30MB    76.30MB (flat, cum) 98.86% of Total
         .          .     29:func AllTest() {
   76.30MB    76.30MB     30:   data := myStruct{}
         .          .     31:
         .          .     32:   for i := 0; i < 100; i++ {
         .          .     33:           CallByValue(data)
         .          .     34:   }
         .          .     35:   for i := 0; i < 100; i++ {
```

* Based on multiple experiments, setting the memory limit to 130 MB is approximately the maximum value for the entire program; if the memory limit is further reduced, it will cause the entire program to be forcibly terminated.
    ```
    docker run -m 125m --rm -v ${PWD}:/app myapp

    encoding/json: /usr/local/go/pkg/tool/linux_amd64/compile: signal: killed
    testing: /usr/local/go/pkg/tool/linux_amd64/compile: signal: killed
    FAIL    test [build failed]
    ```
* From the experimental results, it can be observed that the CallByValue() function has become the biggest bottleneck for performance. Due to the need for copying, this memory limit is already at the extreme.

## Summary
* When continuously stress-testing memory limits:
  * The CallByValue_Escape() function first encounters performance bottlenecks.
  * Followed by the CallByValue() function.
* This indicates that in practical usage scenarios, when memory is insufficient, using CallByValue() will further degrade performance.
