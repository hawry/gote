package buffer

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hawry/gote/helpers"
)

func TestEmptyBuffer(t *testing.T) {
	os.Remove("./.gote_buffer")
	if b := hasBuffer(); b {
		t.Logf("expected %t, returned %t", false, b)
		t.Fail()
	}
}

func TestSaveBuffer(t *testing.T) {
	setFilePath("../../.gote_buffer")

	if hasBuffer() {
		os.Remove("../../.gote_buffer")
	}

	Add(helpers.Issue{Title: "issue title", Body: fmt.Sprintf("issue body (%d)", time.Now().UnixNano())})

	if b := Count(); b != 1 {
		t.Logf("expected %d, returned %d", 1, b)
		t.Fail()
	}
	save()

	load()
	if b := Count(); b != 1 {
		t.Logf("expected %d, returned %d", 1, b)
		t.Fail()
	}
	Add(helpers.Issue{Title: "issue title", Body: fmt.Sprintf("issue body (%d)", time.Now().UnixNano())})
	if b := Count(); b != 2 {
		t.Logf("expected %d, returned %d", 2, b)
		t.Fail()
	}
	save()

	load()
	if b := Count(); b != 2 {
		t.Logf("expected %d, returned %d", 2, b)
		t.Fail()
	}
}
