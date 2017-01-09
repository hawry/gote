package buffer

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hawry/gote/helpers"
)

var fpath string

func TestMain(m *testing.M) {
	fpath = newTestFile()
	m.Run()
	os.Remove(fpath)
}

func newTestFile() string {
	return fmt.Sprintf("./.gote_buffer_%d", time.Now().UnixNano())
}

func TestEmptyBuffer(t *testing.T) {
	setFilePath(fpath)
	if b := hasBuffer(); b {
		t.Logf("expected %t, returned %t", false, b)
		t.Fail()
	}
}

func TestSaveBuffer(t *testing.T) {
	setFilePath(fpath)
	if hasBuffer() {
		os.Remove(fpath)
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

	Remove()
	if b := Count(); b != 1 {
		t.Logf("expected %d, returned %d", 1, b)
		t.Fail()
	}
	Remove()
	if b := Count(); b != 0 {
		t.Logf("expected %d, returned %d", 0, b)
		t.Fail()
	}
	Save()
	load()
	if b := Count(); b != 0 {
		t.Logf("expected %d, returned %d", 0, b)
		t.Fail()
	}
}
