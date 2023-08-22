package log

import "fmt"

var IsVerbose = false

func Println(a ...any) {
	fmt.Println(a...)
}

func Print(a ...any) {
	fmt.Print(a...)
}

func Verboseln(a ...any) {
	if IsVerbose {
		fmt.Println(a...)
	}
}

func Verbose(a ...any) {
	if IsVerbose {
		fmt.Print(a...)
	}
}

func Errorln(err error) {
	fmt.Println("Error:", err)
}
