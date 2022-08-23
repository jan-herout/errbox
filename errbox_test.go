package errbox

import (
	"errors"
	"fmt"
	"testing"
)

func TestErrorsIs(t *testing.T) {
	e1 := fmt.Errorf("e1")
	var e2 error
	e2 = Append(e2, e1)
	assertEquals := false
	for _, err := range Errors(e2) {
		if errors.Is(err, e1) {
			assertEquals = true
		}
	}
	if !assertEquals {
		t.Errorf("box does not contain e1")
	}
}

func TestAccumulate(t *testing.T) {
	retErr := fmt.Errorf("error")
	counter := 0
	okFunc := func() error {
		counter = counter + 1
		return nil
	}
	failingFunc := func() error {
		counter = counter + 1
		return retErr
	}
	err := Run(okFunc).Then(okFunc).Then(failingFunc).Then(okFunc).Then(failingFunc).First()
	if counter != 3 {
		t.Errorf("expected counterto be equal to 3")
	}
	if !errors.Is(err, retErr) {
		t.Errorf("call chain did not return expected error")
	}
}

func TestBox(t *testing.T) {
	OmitPrefixFromTrace("errbox/")
	ShowStack(true)

	b := NewBox()
	if b == nil {
		t.FailNow()
	}

	var err error
	err = fmt.Errorf("boom")            // an error from stdlib
	b.PushIf(err, "")                   // push it back
	b.PushIf(b.last(), "two levels up") // annotate it again; since the err is a "stdlib" error, we need to annotate the last error

	// create annotated error and push it back
	err = Annotate(fmt.Errorf("boom"), "with num %d", 10)
	err = Annotate(err, "")
	err = Annotate(err, "another message")
	err = b.PushIfErr(err, "another level up")

	if x := len(Errors(b)); x != 2 {
		t.Errorf("wanted to get 2 errors, got %d", x)
	}

	// at this stage, err is already of type StackErr, because it was created by Annotate
	// now we try to attach additional (not printed) fields to it
	f := WithStack(err).Fields()
	f["what"] = "thingy"
	if WithStack(err).StringField("what") != "thingy" {
		t.Errorf("got: %#v", f)
	}
	if WithStack(err).StringField("n/a") != "" {
		t.Errorf("got: %#v", f)
	}

	var storedErr = fmt.Errorf("this is inside")
	var missingErr = fmt.Errorf("this is missing")

	b.PushIf(storedErr, "")
	if !IsInside(b, storedErr) {
		t.Errorf("expected to see this in the box: %s", storedErr)
	}
	if IsInside(b, missingErr) {
		t.Errorf("expected NOT to see this in th box: %s", storedErr)
	}
}

func TestAppend(t *testing.T) {
	err := Annotate(fmt.Errorf("kablam"), "")
	err = Append(err, Annotate(fmt.Errorf("boom"), ""))
	// TODO add tests
}
