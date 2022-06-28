/*
Package errbox provides errors with stack trace, which can also be groupped (boxed) and processed later.
*/
package errbox

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
)

// OmitPrefixFromTrace will SET package level variable filePrefix.
// Later, when errors are printed out (Error() is called), stack trace is inspected.
// Filename where the error occured is searched for the prefix, and everything before the prefix plus the prefix itself
// is dropped from the filename.
//
// Why is this useful: suppose you have package called recombobulator, and you do not want to print out the path
// to the current package in our error. You can achieve this by calling OmitPrefixFromTrace("recombobulator/").
//
// Beware, this variable is not mutex protected, therefore you should only set it ONCE, and then it should NOT be touched!
func OmitPrefixFromTrace(pfx string) {
	pfx = filepath.ToSlash(pfx)
	filePrefix = pfx
}

// filePrefix will always be removed from the stack trace.
// This variable is NOT mutex protected, therefore you should set it once at the beginning of your program, and then
// it should NOT be touched again.
var filePrefix string

// ShowStack will SET package level variable showStack. This variable controls how errors are printed out.
// Beware, this variable is not mutex protected, therefore you should only set it ONCE, and then it should NOT be touched!
func ShowStack(show bool) {
	showStack = show
}

// showStack controls if stack trace is printed out.
var showStack = true

// Box can store multiple errors, and also implements the error interface itself,
// It is a mutex protected storage of other errors. Use it via Append, or directly via PushIf, or PushIfErr.
type Box struct {
	mu     sync.Mutex
	errLis []*StackErr // list of errors encountered so far
}

// Append appends the error to the error of type *Box, and returns it.
//
// If the first parameter is nil, or is of different type, it is converted to the type Box.
//
// If err is nil, nothing happens (returns nil).
func Append(box, err error) error {
	// noop if no error
	if err == nil {
		return box
	}
	newBox := asBox(box)

	// flatten the box
	if errBox, ok := err.(*Box); ok {
		newBox.errLis = append(newBox.errLis, errBox.errLis...)
		return newBox
	}

	// err is not a *Box, convert it to the type *StackErr
	newErr := WithStack(err)
	newBox.errLis = append(newBox.errLis, newErr)
	return newBox
}

// Errors returns copy of slice of errors encountered so far. Nil slice is returned if err is nil.
//
// If err is of type *Box, returns slice with all errors appended to the *Box.
//
// If err is not of type *Box, slice containing the err is returned.
func Errors(err error) []error {
	if err == nil {
		return nil
	}
	b := asBox(err)
	b.mu.Lock()
	defer b.mu.Unlock()
	l := len(b.errLis)
	if l == 0 {
		return nil
	}

	errs := make([]error, l)
	for i := range b.errLis {
		errs[i] = b.errLis[i]
	}
	return errs
}

// NewBox returns a new Box pointer. The Box should never be copied, because it contains a mutex.
func NewBox() *Box {
	box := new(Box)
	return box
}

// asBox converts the error to a *Box. Nil is also converted to a new box.
func asBox(err error) *Box {
	if err == nil {
		return NewBox()
	}
	if b, ok := err.(*Box); ok {
		return b
	}
	// underlying error is not a Box, convert it
	b := NewBox()
	this := WithStack(err)
	b.errLis = append(b.errLis, this)
	return b
}

// PushIf adds the error to the Box, and returns true if the first parameter was not nil. If the error is nil, returns false.
// If you need to return the annotated error, use PushIfErr instead.
func (b *Box) PushIf(err error, message string, args ...interface{}) bool {
	// return on no error
	if err == nil {
		return false
	}
	b.mu.Lock()
	defer b.mu.Unlock()

	// this error and last error, both as boxed errors
	// annotate this error (give it stack trace and additional message
	this := WithStack(err)
	this.annotate(2, message, args...)
	last := b.last()

	// append this error if last error was different
	if this != last {
		b.errLis = append(b.errLis, this)
	}

	// return the error
	return true
}

// PushIfErr adds the error to the Box, and returns the annotated error if the first parameter was not nil. If the error is nil, returns nil.
func (b *Box) PushIfErr(err error, message string, args ...interface{}) error {
	// return on no error
	if err == nil {
		return nil
	}
	b.mu.Lock()
	defer b.mu.Unlock()

	// this error and last error, both as boxed errors
	// annotate this error (give it stack trace and additional message
	this := WithStack(err)
	this.annotate(2, message, args...)
	last := b.last()

	// append this error if last error was different
	if this != last {
		b.errLis = append(b.errLis, this)
	}

	// return the error
	return this
}

// last returns the last error encountered, or nil if no error were encountered yet.
func (b *Box) last() *StackErr {
	if len(b.errLis) == 0 {
		return nil
	}
	return b.errLis[len(b.errLis)-1]
}

// Error implements the error interface
func (b *Box) Error() string {
	b.mu.Lock()
	defer b.mu.Unlock()

	if len(b.errLis) == 0 {
		return ""
	}

	if len(b.errLis) == 1 {
		return b.errLis[0].Error()
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Got %d errors:\n", len(b.errLis)))
	for i, err := range b.errLis {
		sb.WriteString("----------------------------\n")
		sb.WriteString(fmt.Sprintf("# %d\n", i+1))
		sb.WriteString(err.Error())
		sb.WriteString("\n")
	}
	return sb.String()
}

// IsInside checks whether the err is the target (think errors.Is).
// When the err is *Box, it returns true if any of the errors in the err is the target.
func IsInside(err error, target error) bool {
	// if the error is not a *Box, simply use errors package
	// otherwise, cycle through all errors
	if be, ok := err.(*Box); ok {
		for _, se := range be.errLis {
			if errors.Is(se, target) {
				return true
			}
		}
		return false
	} else {
		return errors.Is(err, target)
	}
}

// String implements Stringer interface
func (b Box) String() string {
	return b.Error()
}
