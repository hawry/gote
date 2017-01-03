package helpers

import (
	"io/ioutil"
	"log"
	"testing"
)

func TestAskBool(t *testing.T) {
	in, err := ioutil.TempFile("", "")
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()
	_, err = in.WriteString("yes\n")
	if err != nil {
		log.Fatal(err)
	}

	i := Ask("This should be a bool later on")
	i.setInput(in)
	rval := i.Bool()

	if !rval {
		t.Logf("expected %t, returned %t", true, rval)
		t.Fail()
	}
}

//
// func TestAskString(t *testing.T) {
// 	i := Ask("Please supply a string!")
// 	t.Logf("answer: %s", i.String())
// }
