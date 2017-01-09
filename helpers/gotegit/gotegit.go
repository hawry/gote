package gotegit

import (
	"fmt"
	"log"
	"regexp"
)

const (
	identifyProviderPattern = "\\A(?P<protocol>git@|git://|https://)(?P<host>[^/:]*)[/|:]{1}(?P<user>[^/]*)/(?P<repo>[^/]*)\\z"
)

type GitProvider interface {
	CreateIssue(accessToken string) (bool, error)
	ProviderName() string
}

func Provider(remoteAddress string) (GitProvider, error) {
	r := regexp.MustCompile(identifyProviderPattern)
	res := parseToMap(r.FindStringSubmatch(remoteAddress), r.SubexpNames())
	log.Printf("debug: remote: '%s'", res["host"])
	switch res["host"] {
	case "github.com":
		g := &Github{}
		return g, nil
	case "bitbucket.org":
		g := &Bitbucket{}
		return g, nil
	case "gitlab.com":
		g := &Gitlab{}
		return g, nil
	default:
		return nil, fmt.Errorf("invalid remote (%s)", remoteAddress)
	}
}

func parseToMap(fss []string, names []string) map[string]string {
	rval := make(map[string]string)
	for i, name := range names {
		if i != 0 {
			rval[name] = fss[i]
		}
	}
	return rval
}
