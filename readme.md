# errbox

## Synopsis

`errbox` is a package for Go, which provides the ability to:

- represent multiple errors as one error
- annotate original error by adding commented stack trace to it

## Is this for me?

Does the world need `yet-another-error-handling` package? No. 
However, I needed something I could easily tweak for my specific purposes.

Is this fast? Probably not so much.

This is an alpha-version of this package. API can change without warning.

Use this on your own risk - should this package break your program, you get 
to keep both pieces.

Please, first consider the following alternatives, which inspired this package:

- https://github.com/hashicorp/go-multierror
- https://github.com/palantir/stacktrace

## Installation and docs

Install using `go get github.com/jan-herout/errbox`. 
Full documentation is available at https://pkg.go.dev/github.com/jan-herout/errbox.

## Usage

### Annotating an existing error

**This:**

```
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
	errbox.OmitPrefixFromTrace("<---root dir of your repo--->")
	if err := doSomething(); err != nil {
		fmt.Println(errbox.Annotate(err, "with string: %s", "recombobulator"))
	}
}
```

**Prints**

```
bang
 +--> with the number: 10
 |  @ main.go:14 (doSomething)
 +--> with string: recombobulator
    @ main.go:20 (main)
```

### Groupping errors together

**This**

```
package main

import (
	"fmt"

	"github.com/jan-herout/errbox"
)

func thisFails(why string) error {
	return fmt.Errorf(why)
}

func main() {
	// use this if you want to sanitize names of files in stack trace
	errbox.OmitPrefixFromTrace("<---root dir of your repo--->")
	var err error
	err = errbox.Append(err, thisFails("just because"))
	err = errbox.Append(err, thisFails("because why not?"))
	err = errbox.Annotate(err,"this annotates only the last error!")
	fmt.Println(err)
}
```

**Prints**

```
Got 2 errors:
# 1
just because
# 2
because why not?
 +--> this annotates only the last error!
    @ main.go:19 (main)
```