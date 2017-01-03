package helpers

import "fmt"

//Issue is a gote representation of a GitHub issue. Not to be confused by those used by support libs
type Issue struct {
	Title string
	Body  string
}

//NewIssue returns a new issue and also formats the raw issue data provided by the user to better suited lengths
func NewIssue(raw string) *Issue {
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
