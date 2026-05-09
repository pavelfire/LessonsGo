package main

import "fmt"

type myError struct {
	code int
}

func (e myError) Error() string {
	return fmt.Sprintf("code: %d", e.code)
}

func run() error {
	var e *myError
	// var e error
	if false {
		e = &myError{code: 123}
	}
	return e
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("failed to run, error: ", err)
	} else {
		fmt.Println("success")
	}
}
