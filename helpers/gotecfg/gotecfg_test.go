package gotecfg

import (
	"fmt"
	"log"
	"testing"
	"time"
)

var localTestConfig = `access_token: global
remote: git@github.com:hawry/git-note
repository_owner: hawry
repository_name: git-note
provider: github`

var globalTestConfig = `access_token: $GOTE_ACCESS_TOKEN
editor: vim`

var testLocal, testGlobal string

func TestMain(m *testing.M) {
	//create test files
	testLocal = fmt.Sprintf("./.local_%d", time.Now().UnixNano())
	testGlobal = fmt.Sprintf("./.local_%d", time.Now().UnixNano())
	//run tests
	m.Run()
	//remove test files
}

func TestLoadGlobal(t *testing.T) {
	config := New()
	log.Printf("%+v", config.globalCfg)
}
