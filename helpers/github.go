package helpers

import (
	"fmt"
	"log"
	"strings"
)

//Issue is a gote representation of a GitHub issue. Not to be confused by those used by support libs
type Issue struct {
	Title string
	Body  string
}

//NewIssue returns a new issue and also formats the raw issue data provided by the user to better suited lengths
//TODO: Regexp instead of this weird hack!
func NewIssue(raw string) *Issue {
	raw = strings.TrimSpace(raw)
	raw = strings.TrimLeftFunc(raw, func(r rune) bool {
		return r == '\n'
	})
	i := Issue{}
	if len(raw) > 50 {
		i.Title = fmt.Sprintf("%s...", raw[:49])
		i.Body = raw
	} else {
		i.Title = raw
		i.Body = raw
	}
	return &i
}

//ToMap is to not make me have to repeat this monstrosity more than neccessary
func ToMap(fss []string, names []string) map[string]string {
	log.Printf("debug: received the following: %+v, %+v", fss, names)
	rval := make(map[string]string)
	for i, name := range names {
		if i != 0 {
			rval[name] = fss[i]
		}
	}
	return rval
}
