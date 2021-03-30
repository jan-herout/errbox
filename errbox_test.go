package errbox

import (
	"fmt"
	"testing"
)

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
}

func TestAppend(t *testing.T) {
	err := Annotate(fmt.Errorf("kablam"), "")
	err = Append(err, Annotate(fmt.Errorf("boom"), ""))
	// TODO add tests
}
