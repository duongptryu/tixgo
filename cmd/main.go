package main

import (
	"errors"
	"fmt"
	"tixgo/pkg/syserr"
)

func main() {
	// errWrapped := syserr.Wrap(err, syserr.InternalCode, "test error", syserr.F("test", "test"))
	// fmt.Println(errWrapped.Error())

	// fmt.Println(errWrapped.Code())
	// fmt.Println(errWrapped.Fields())

	// _ = errWrapped.StackFormatted()
	// fmt.Println(x)

	// err = syserr.Wrap(err, syserr.InternalCode, "test error 2")
	// fmt.Println(err.Error())
	fmt.Println("Hello world")
	err := errors.New("test error")
	fmt.Println(err)

	errx := syserr.Wrap(err, syserr.InternalCode, "test error 2")
	fmt.Println(errx.StackFormatted())
}
