package gotegit

import (
	"log"
	"reflect"
	"testing"
)

func TestIdentifyGithub(t *testing.T) {
	expected := reflect.TypeOf(&Github{})
	provider, err := Provider("git@github.com:hawry/gote")
	if err != nil {
		t.Logf("received error (%v)", err)
		t.Fail()
	}
	returned := reflect.TypeOf(provider)
	log.Printf("%+v", returned)
	if returned != expected {
		t.Logf("expected %v, returned %v", expected, returned)
		t.Fail()
	}
}

func TestIdentifyBitbucket(t *testing.T) {
	expected := reflect.TypeOf(&Bitbucket{})
	provider, err := Provider("git@github.com:hawry/gote")
	if err != nil {
		t.Logf("received error (%v)", err)
		t.Fail()
	}
	returned := reflect.TypeOf(provider)
	log.Printf("%+v", returned)
	if returned != expected {
		t.Logf("expected %v, returned %v", expected, returned)
		t.Fail()
	}
}

func TestIdentifyGitlab(t *testing.T) {
	expected := reflect.TypeOf(&Gitlab{})
	provider, err := Provider("git@github.com:hawry/gote")
	if err != nil {
		t.Logf("received error (%v)", err)
		t.Fail()
	}
	returned := reflect.TypeOf(provider)
	log.Printf("%+v", returned)
	if returned != expected {
		t.Logf("expected %v, returned %v", expected, returned)
		t.Fail()
	}
}
