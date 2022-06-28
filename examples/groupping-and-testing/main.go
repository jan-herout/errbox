/*
 * Output:
 *
-----------------------------
this, however will activate
-----------------------------
Got 2 errors:
----------------------------
# 1
something tragic happened
 +--> because we were careless
    @ /main.go:37 (doSomethingCareless)

----------------------------
# 2
open a non existing file: The system cannot find the file specified.
 +--> obviously, file does not exist: path=a non existing file
 |  @ /main.go:44 (openNonExistentFile)
 +--> the thing broke!
    @ /main.go:50 (doAnotherCarelessThing)
*/
package main

import (
	"fmt"
	"os"

	"github.com/jan-herout/errbox"
)

var ourSentinelErr = fmt.Errorf("something tragic happened")
var stuffMissingErr = fmt.Errorf("stuff is missing")

func doSomethingWhichIsOK() error {
	return nil
}

func doSomethingCareless() error {
	return errbox.Annotate(ourSentinelErr, "because we were careless")
}

func openNonExistentFile() error {
	var err error
	path := "a non existing file"
	fh, err := os.Open(path)
	err = errbox.Annotate(err, "obviously, file does not exist: path=%s", path)
	defer fh.Close()
	return err
}

func doAnotherCarelessThing() error {
	return errbox.Annotate(openNonExistentFile(), "the thing broke!")
}

func main() {
	// Cleanup trace, do not show full paths.
	errbox.OmitPrefixFromTrace("C:/Git/errbox/examples/groupping-and-testing")

	// do stuff
	var err error

	// do things that may or may not work, and collect errors from them
	err = errbox.Append(err, doSomethingCareless())
	err = errbox.Append(err, doSomethingWhichIsOK())
	err = errbox.Append(err, doAnotherCarelessThing())

	// test if a particular error happened
	if errbox.IsInside(err, stuffMissingErr) {
		fmt.Println("you will not see this message")
	}
	// test if the sentinel error is stored inside the error box
	if errbox.IsInside(err, os.ErrNotExist) {
		fmt.Println("-----------------------------")
		fmt.Println("this, however will activate")
		fmt.Println("-----------------------------")
		fmt.Println(err)
	}
}
