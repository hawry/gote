package format

import (
	"strings"
	"testing"
)

func TestTrimNewlines(t *testing.T) {
	sExpected := "This is a string to compare against"

	sInput := "This is a string to compare against\n"
	sReturned := strings.TrimRightFunc(sInput, TrimNewlines)
	if n := strings.Compare(sExpected, sReturned); n != 0 {
		t.Logf("expected '%s', returned '%s' (%d)", sExpected, sReturned, n)
		t.Fail()
	}

	sInput = "\nThis is a string to compare against"
	sReturned = strings.TrimLeftFunc(sInput, TrimNewlines)
	if n := strings.Compare(sExpected, sReturned); n != 0 {
		t.Logf("expected '%s', returned '%s' (%d)", sExpected, sReturned, n)
		t.Fail()
	}

	sInput = "\n\n\nThis is a string to compare against\n\n\n"
	sReturned = strings.TrimFunc(sInput, TrimNewlines)
	if n := strings.Compare(sExpected, sReturned); n != 0 {
		t.Logf("expected '%s', returned '%s' (%d)", sExpected, sReturned, n)
		t.Fail()
	}
}
