package errbox

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

// StackErr is an error with stack trace.
type StackErr struct {
	cause      error                  // the original error
	annotation []stackAnnotation      // annotation of the error
	fields     map[string]interface{} // optional fields attached to the error via Fields.
}

// stackAnnotation is the annotation of the error.
type stackAnnotation struct {
	// what happened?
	message string
	// where did it happen?
	file     string
	function string
	line     int
}

// Annotate returns back an error annotated with stack trace (of type *StackErr), or nil, if the first parameter was nil.
//
// Repeated call of Annotate on the same error only add the annotation to the (already existing) error.
//
// Call of Annotate on error which is a *Box  annotates all errors.
//
// If the message is not empty string, it is added to the stack. The message is formatting string used by fmt.Sprintf,
// and args... is a variadic parameter which is also provided to the fmt.Sprintf.
func Annotate(err error, message string, args ...interface{}) error {
	// return on no error
	if err == nil {
		return nil
	}

	// what if the err is actually *Box?
	// then we annotate all errors in the box
	if b, ok := err.(*Box); ok {
		for i := range b.errLis {
			b.errLis[i].annotate(2, message, args...)
		}
		return b
	}

	// annotate this error (give it stack trace and additional message
	this := WithStack(err)
	this.annotate(2, message, args...)
	return this
}

// WithStack returns the error as StackErr error, or converts the err to a new StackErr if possible.
// Returns nil if err is nil.
func WithStack(err error) *StackErr {
	if err == nil {
		return nil
	}
	if be, ok := err.(*StackErr); ok {
		return be
	}
	be := new(StackErr)
	be.cause = err
	return be
}

// Cause returns cause of the error.
//
// It the error is nil, nil is returned.
//
// If the error is StackErr, then the cause (first annotated error) is returned transparently.
//
// If the error is Box, cause of the first error in the box is returned.
//
// Otherwise, err is returned.
func Cause(err error) error {
	if err == nil {
		return nil
	}
	if e, ok := err.(*StackErr); ok {
		return e.cause
	}
	if e, ok := err.(*Box); ok {
		if len(e.errLis) == 0 {
			return nil
		}
		return Cause(e.errLis[0])
	}
	return err
}

// Fields returns a map, which can be used to store or fetch anything. Typically, you would use it as follows:
//   stacked := WithStack(err)
//   fields := stacked.Fields()
//   fields["whatever"] = "whatever"  // set it to string
//   x := stacked.StringField("whatever")  // get it back
func (b *StackErr) Fields() map[string]interface{} {
	if b.fields == nil {
		b.fields = make(map[string]interface{})
	}
	return b.fields
}

// StringField attempts to access a field, convert it to a string, and return it.
// The function returns empty string if the field was not found, or if the value could not be converted to string.
func (b *StackErr) StringField(name string) string {
	m := b.Fields()
	i, ok := m[name]
	if !ok {
		return ""
	}
	s, ok := i.(string)
	if ok {
		return s
	}
	return ""
}

// Error implements the Error interface
func (b *StackErr) Error() string {
	// if no annotation is found, return the original error
	if len(b.annotation) == 0 {
		return b.cause.Error()
	}

	// otherwise, prepare the string
	var sb strings.Builder
	dNext := " |  "
	dThis := " +--"
	dEmpty := "    "

	ln := len(b.annotation) - 1
	sb.WriteString(fmt.Sprintf("%s\n", b.cause))
	for i, anno := range b.annotation {
		delim := dThis
		if anno.message != "" {
			sb.WriteString(fmt.Sprintf("%s> %s\n", delim, anno.message))
			if i < ln {
				delim = dNext
			} else {
				delim = dEmpty
			}
		}
		if showStack && anno.line > 0 {
			sb.WriteString(fmt.Sprintf("%s@ %s:%d (%s)\n", delim, anno.file, anno.line, anno.function))
		}
	}
	return sb.String()
}

// Unwrap implements errors.Unwrap interface.
func (b *StackErr) Unwrap() error {
	return b.cause
}

// annotate adds the message to the original error
func (b *StackErr) annotate(skip int, message string, args ...interface{}) {
	// User code is two stack frames up, as this is called from Annotate
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return
	}

	// clean the file
	if filePrefix != "" {
		file = filepath.ToSlash(file)
		idx := strings.Index(file, filePrefix)
		if idx > -1 {
			idx = idx + len(filePrefix)
			file = file[idx:]
		}
	}

	// prepare the annotation
	annotation := stackAnnotation{
		message: fmt.Sprintf(message, args...),
		file:    file,
		line:    line,
	}

	// get the function
	f := runtime.FuncForPC(pc)
	if f != nil {
		annotation.function = shortFuncName(f)
	}
	// append it to the error
	b.annotation = append(b.annotation, annotation)
}

// this comes from https://github.com/palantir/stacktrace/blob/master/stacktrace.go
// props to them!
func shortFuncName(f *runtime.Func) string {
	// f.Name() is like one of these:
	// - "github.com/palantir/shield/package.FuncName"
	// - "github.com/palantir/shield/package.Receiver.MethodName"
	// - "github.com/palantir/shield/package.(*PtrReceiver).MethodName"
	longName := f.Name()

	withoutPath := longName[strings.LastIndex(longName, "/")+1:]
	withoutPackage := withoutPath[strings.Index(withoutPath, ".")+1:]

	shortName := withoutPackage
	shortName = strings.Replace(shortName, "(", "", 1)
	shortName = strings.Replace(shortName, "*", "", 1)
	shortName = strings.Replace(shortName, ")", "", 1)

	return shortName
}
