/*
 * Output:
 *
bang
 +--> with the number: 10
 |  @ /main.go:14 (doSomething)
 +--> with string: recombobulator
    @ /main.go:21 (main)
*/
package main

import (
	"fmt"

	"github.com/jan-herout/errbox"
)

func thisFails() error {
	return fmt.Errorf("bang")
}

func doSomething() error {
	return errbox.Annotate(thisFails(), "with the number: %d", 10)
}

func main() {
	// use this if you want to sanitize names of files in stack trace
	errbox.OmitPrefixFromTrace("C:/Git/errbox/examples/annotate-errors")
	if err := doSomething(); err != nil {
		fmt.Println(errbox.Annotate(err, "with string: %s", "recombobulator"))
	}
}
