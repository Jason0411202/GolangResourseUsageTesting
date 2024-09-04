package main

type myStruct struct { // simulate large data
	arr [10000000]int
}

// passing parameters using call by value method
func CallByValue(myStruct myStruct) {

}

// passing parameters using call by reference method
func CallByReference(myStruct *myStruct) {

}

// passing parameters using call by value method
func CallByValue_Escape(myStruct myStruct) *myStruct {

	return &myStruct
}

// passing parameters using call by reference method
func CallByReference_Escape(myStruct *myStruct) *myStruct {

	return myStruct
}

func AllTest() {
	data := myStruct{}

	for i := 0; i < 100; i++ {
		CallByValue(data)
	}
	for i := 0; i < 100; i++ {
		CallByReference(&data)
	}

	for i := 0; i < 100; i++ {
		CallByValue_Escape(data)
	}
	for i := 0; i < 100; i++ {
		CallByReference_Escape(&data)
	}
}

func main() {
}
