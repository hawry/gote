package config

import (
	"fmt"
	"log"
	"strings"
	"testing"
)

func TestParseRemoteSSH(t *testing.T) {
	eUsr := "hawry"
	eRep := "gote"
	remote := fmt.Sprintf("git@github.com:%s/%s", eUsr, eRep)

	usr, rep := parseRemoteInformation(remote)
	log.Printf("%s %s", usr, rep)
	if !strings.EqualFold(usr, "hawry") {
		t.Logf("expected %s, returned %s", eUsr, usr)
		t.Fail()
	}
	if !strings.EqualFold(rep, eRep) {
		t.Logf("expected %s, returned %s", eRep, rep)
		t.Fail()
	}
}

func TestParseRemoteHTTPS(t *testing.T) {
	eUsr := "hawry"
	eRep := "gote"
	remote := fmt.Sprintf("https://github.com/%s/%s", eUsr, eRep)

	usr, rep := parseRemoteInformation(remote)
	log.Printf("%s %s", usr, rep)
	if !strings.EqualFold(usr, "hawry") {
		t.Logf("expected %s, returned %s", eUsr, usr)
		t.Fail()
	}
	if !strings.EqualFold(rep, eRep) {
		t.Logf("expected %s, returned %s", eRep, rep)
		t.Fail()
	}
}
