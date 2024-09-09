# In-Depth Exploration of Golang Call by Value/Reference Performance

## Experiment Design
* In this experiment, I will pass a struct containing an array of 10,000,000 elements to a function (simulating large data transfer) using both Call by Value and Call by Reference. Each method will be executed 100 times.
* The `CallByValue()` function will be called 100 times by the `AllTest()` function, passing large data as an argument via Call by Value.
* Similarly, the `CallByReference()` function will also be called 100 times by the `AllTest()` function, passing large data as an argument via Call by Reference.
* To simulate an escape condition (where the allocated space cannot be freed after the function ends), I designed corresponding Escape versions for both functions to compare.
* During this process, I will record the CPU and memory usage during execution.
* To observe the experimental results, execute the following commands under the Golang directory of this project.

## Run the benchmark test
After switching to the Golang directory, the following command will execute the golang benchmark test and analyze the result using go pprof.

(For Windows environment)
```bash
docker run --rm -v ${PWD}:/app myapp
```

(For Linux environment)
```bash
docker run --rm -v "$(pwd):/app" myapp
```

## Analyze CPU usage in the benchmark test
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
         .          .     42:   for i := 100; i++ {
         .          .     43:           CallByReference_Escape(&data)
         .          .     44:   }
         .          .     45:}
```
* Observing the experimental results:
  * Line 33 shows a direct execution time of 1.04s. This is because Call by Value requires copying data to the function, consuming a lot of time.
  * Line 36 shows a direct execution time of 0.00s since Call by Reference only needs to pass a pointer, thus saving time on copying data.
  * Line 40 shows a direct execution time of 1.09s, also due to copying data with Call by Value. The total time is 2.81s, likely due to the time consumed by garbage collection.
    * Only memory allocated on the heap needs to be garbage collected. Memory allocated on the stack will be freed after the function ends.
  * Line 43 shows both the direct execution time and total time as 0.00s. This is because Call by Reference only requires passing a pointer, so no time is consumed copying data.

## Analyze memory usage in the benchmark test
```
Total: 7.53GB
ROUTINE ======================== test.AllTest in /app/main.go
   76.30MB     7.53GB (flat, cum)   100% of Total
         .          .     29:func AllTest() {
   76.30MB    76.30MB     30:   data := myStruct{}
         .          .     31:
         .          .     32:   for i := 100; i++ {
         .          .     33:           CallByValue(data)
         .          .     34:   }
         .          .     35:   for i := 100; i++ {
         .          .     36:           CallByReference(&data)
         .          .     37:   }
         .          .     38:
         .          .     39:   for i := 100; i++ {
         .     7.45GB     40:           CallByValue_Escape(data)
         .          .     41:   }
         .          .     42:   for i := 100; i++ {
         .          .     43:           CallByReference_Escape(&data)
         .          .     44:   }
         .          .     45:}
```
* Observing the experimental results:
  * Line 33 shows 0 MB of memory usage since memory allocated in the function (copied via Call by Value) does not escape, and Golang's compiler decides to allocate it on the stack, so it won't be detected by pprof.
    * According to the issue https://github.com/golang/go/issues/15848, go pprof currently only analyzes memory allocated on the heap.
  * Line 36 also shows 0 MB since Call by Reference only passes a pointer, requiring no extra memory allocation.
  * Line 40 shows 7.45GB because the memory allocated in the function (copied via Call by Value) escapes, and Golang's compiler decides to allocate it on the heap. As the function is called 100 times, the total memory usage of 7.45GB closely aligns with 100 times the size of data (76.3MB).
  * Line 43 shows 0 MB for the same reason as Line 36 — Call by Reference doesn't require extra memory allocation.

## Conclusion
* In terms of time:
  * Call by Value performs worse because it requires copying data with each function call. If the data escapes, additional time is consumed due to garbage collection.
  * Call by Reference is more efficient because no time is spent allocating memory.
* In terms of memory:
  * When data passed via Call by Value does not escape, Golang's compiler allocates it on the stack, and it is freed after the function ends, making it undetectable by pprof.
  * If the data passed via Call by Value escapes, it is allocated on the heap, and memory usage is detected by pprof.
  * Call by Reference does not require additional memory allocation, making it more efficient in this regard as well.

Note: For more details on Golang's memory allocation, check out this article: https://medium.com/eureka-engineering/understanding-allocations-in-go-stack-heap-memory-9a2631b5035d