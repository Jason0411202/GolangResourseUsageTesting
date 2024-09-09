# 基於實驗一，進一步探討在限制記憶體的情況下 Golang 的 Call by Value/Reference 效能
## 實驗設計
* 與實驗一唯一不同的是，在使用 docker 執行 benchmark test 時，會透過 `-m` 參數限制記憶體上限
* 本實驗旨在測試不同記憶體上限下，Golang 的 Call by Value/Reference 效能表現

(window 環境下適用，記憶體上限設為 500 MB)
```bash
docker run -m 500m --rm -v ${PWD}:/app myapp
```

(linux 環境下適用，記憶體上限設為 500 MB)
```bash
docker run -m 500m --rm -v "$(pwd):/app" myapp
```

## 記憶體上限不限制 (沿用實驗一的結果，作為對照組)
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

## 記憶體上限設定為 245 MB
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
* 根據多次實驗，記憶體上限設定為 245 MB 約是 CallByValue_Escape() 函式的極限值；若是再限縮記憶體上限的話，便會導致整個程式被強制終止
    ```
    docker run -m 240m --rm -v ${PWD}:/app myapp

    signal: killed
    FAIL    test    6.276s
    ```
* 從實驗結果可以觀察到，CallByValue_Escape() 函式成為了拖累性能的最大瓶頸，推測是因為記憶體不足的原因，Golang 的 garbage collection 機制需要非常頻繁的回收記憶體，導致性能明顯下降
* 為了能繼續實驗，接下來將會註解掉 CallByValue_Escape() 的部分，繼續壓測其他函式

## 記憶體上限設定為 130 MB
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

* 根據多次實驗，記憶體上限設定為 130 MB 約是整個程式的極限值；若是再限縮記憶體上限的話，便會導致整個程式被強制終止
    ```
    docker run -m 125m --rm -v ${PWD}:/app myapp

    encoding/json: /usr/local/go/pkg/tool/linux_amd64/compile: signal: killed
    testing: /usr/local/go/pkg/tool/linux_amd64/compile: signal: killed
    FAIL    test [build failed]
    ```
* 從實驗結果可以觀察到，這次換 CallByValue() 函式成為了拖累性能的最大瓶頸，由於需要複製一份的關係，這樣的記憶體限制已經是極限了

## 小結
* 在不斷壓測記憶體上限時
  * CallByValue_Escape() 函式首先出現性能瓶頸
  * 緊接著是 CallByValue() 函式
* 這顯示在實際的使用情景下，當記憶體不足時，使用 CallByValue() 會導致性能進一步下降
