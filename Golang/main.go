package main

type myStruct struct {
	arr [10000000]int
}

//go:noinline
func CallByValue(myStruct myStruct) {
	sum := 0
	for i := 0; i < len(myStruct.arr); i++ { // modify the parameters, to prevent the Go compiler from optimizing away the parameter modification
		myStruct.arr[i] += 1
		sum += myStruct.arr[i]
	}
	//fmt.Println(sum)
}

//go:noinline
func CallByReference(myStruct *myStruct) {
	sum := 0
	for i := 0; i < len(myStruct.arr); i++ {
		myStruct.arr[i] += 1
		sum += myStruct.arr[i]
	}
	//fmt.Println(sum)
}

//go:noinline
func CallByValueTest() {
	data := myStruct{}

	for i := 0; i < 100; i++ {
		CallByValue(data)
	}
}

//go:noinline
func CallByReferenceTest() {
	data := myStruct{}

	for i := 0; i < 100; i++ {
		CallByReference(&data)
	}
}

func main() {
	CallByValueTest()
	CallByReferenceTest()
}
