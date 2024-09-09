# 深入探討 Golang 的 Call by Value/Reference 效能
## 實驗設計
* 在本實驗中，我將會分別使用 Call by Value/Reference 的方式傳遞一個包含 10000000 個元素的陣列的 struct 給一個函式 (模擬傳遞大資料)，共計執行 100 次
* 其中，CallByValue() 函式會被 AllTest() 函式呼叫 100 次，並透過 Call by Value 的方式傳遞大資料作為參數
* 而 CallByReference() 函式同樣會被 AllTest() 函式呼叫 100 次，並透過 Call by Reference 的方式傳遞大資料作為參數
* 為了模擬 escape 的情況 (函式結束後，用到的空間仍不能被釋放)，本實驗還為上述兩個函式設計了對應的 Escape 版本以供比較
* 在這個過程中，紀錄執行時的 CPU/memory 使用情形
* 於本專案的 Golang 資料夾下執行以下指令即可觀察實驗結果

## 執行 benchmark test
切換進 Golang 資料夾後，以下指令會執行 golang benchmark test，並用 go pprof 分析結果

(window 環境下適用)
```bash
docker run --rm -v ${PWD}:/app myapp
```

(linux 環境下適用)
```bash
docker run --rm -v "$(pwd):/app" myapp

```

## 分析 benchmark test 的 cpu usage
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
* 觀察實驗結果
  * 第 33 行的直接語句與總耗時皆為 1.04s，這是因為 call by value 需要複製一份資料到函式中，故需要消耗大量時間
  * 第 36 行的直接語句與總耗時皆為 0.00s，這是因為 call by reference 只需要傳遞指標，故不需要消耗時間複製資料
  * 第 40 行的直接語句耗時為 1.09s，這同樣是因為 call by value 需要複製一份資料到函式中，故需要消耗大量時間；值得注意的是，總耗時為 2.81s，推測是 garbage collection 所消耗的時間
    * 只有配置在 heap 中的 memory 才需要被 garbage collection 機制所回收，配置在 stack 中的 memory，在函式結束後，就會直接被釋放掉了
  * 第 43 行的直接語句與總耗時皆為 0.00s，這同樣是因為 call by reference 只需要傳遞指標，故不需要消耗時間複製資料

## 分析 benchmark test 的 memory usage
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


* 觀察實驗結果
  * 第 33 行的直接語句與總消耗的記憶體皆為 0 MB，這是因為在函式中配置的記憶體 (call by value 複製的那份) 沒有 escape (函式結束後便不再用到)，故 golang 編譯器決定將其配置在 Stack 中，故不會被 pprof 檢測到
    * 根據 `https://github.com/golang/go/issues/15848` Issue，go pprof 目前只分析配置在 heap 中的記憶體
  * 第 36 行的直接語句與總消耗的記憶體皆為 0 MB，這是因為 call by reference 只需要傳遞指標，故不需要配置記憶體
  * 第 40 行總消耗的記憶體為 7.45GB，這是因為在函式中配置的記憶體 (call by value 複製的那份) 有 escape (函式結束後還可能再用到)，故 golang 編譯器決定將其配置在 Heap 中，故會被 pprof 檢測到；由於函式會被呼叫 100 次，故 7.45GB 在數量集上也十分接近 data (76.3MB) 的 100 倍
  * 第 43 行的直接語句與總消耗的記憶體皆為 0 MB，這是因為 call by reference 只需要傳遞指標，故不需要配置記憶體

## 小結
* 在時間方面
  * call by value 每次函式呼叫，皆需要複製一份資料到函式中，故效能較差；若是因為 escape 的情況需要將資料配置在 Heap 中的話，還會因為後續需要 garbage collection 而消耗更多時間
  * call by reference 不需要消耗時間分配記憶體，故效能顯然比 Call by Value 高
* 在記憶體方面
  * 若是透過 call by value 傳遞的資料沒有 escape 的情形，則 golang 編譯器會將其分配在 Stack 中，函式結束後即釋放掉，故不會被 pprof 檢測到
  * 若是透過 call by value 傳遞的資料有 escape 的情形，則 golang 編譯器會將其分配在 Heap 中，函式結束後不會被立即釋放，故會被 pprof 檢測到
  * 若是透過 call by reference 來傳遞資料，則不需多配置記憶體


註: 關於 golang 記憶體配置的更多細節，可參考一下這篇 `https://medium.com/eureka-engineering/understanding-allocations-in-go-stack-heap-memory-9a2631b5035d`
