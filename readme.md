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

Please, first consider the following **mature** alternatives, which inspired this package:

- https://github.com/hashicorp/go-multierror
- https://github.com/palantir/stacktrace
- https://github.com/emperror/emperror - excellent library, however, seems to be an overkill

## Installation and docs

Install using `go get github.com/jan-herout/errbox`.
Full documentation is available at https://pkg.go.dev/github.com/jan-herout/errbox.

## Usage

See [examples](examples) provided within this repo.
