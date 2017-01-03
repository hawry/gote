package helpers

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

//Editor enables packages outside of the helper package to access the default editor and parses content as well as title
type Editor struct {
	File    *os.File
	Title   string
	Content string
	Valid   bool
	Issue   *Issue
}

const (
	titlePattern = "(?P<title>\\A.*)\n*(?ms)(?P<content>.*)\\z"
)

//CanUseEditor returns true if the $EDITOR environment variable is set, otherwise returns false
func CanUseEditor() bool {
	editor := os.Getenv("EDITOR")
	_, err := exec.LookPath(editor)
	if err != nil {
		return false
	}
	return true
}

//NewEditor returns an editor struct, which can't be manipulated
func NewEditor() Editor {
	e := Editor{Valid: false}
	editor := os.Getenv("EDITOR")
	tmpFile, err := ioutil.TempFile(os.TempDir(), "gote_")
	if err != nil {
		log.Printf("error: tempfile could not be created (%v)", err)
		return e
	}
	editorPath, err := exec.LookPath(editor)
	if err != nil {
		log.Printf("error: %s not found in path", editor)
	}
	e.File = tmpFile
	cmd := exec.Command(editorPath, e.File.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		log.Printf("error: could not start command (%v)", err)
		return e
	}
	err = cmd.Wait()
	if err != nil {
		log.Printf("warning: empty response or something went wrong (%v)", err)
	}
	e.parse()
	return e
}

func (e *Editor) parse() {
	s, err := ioutil.ReadAll(e.File)
	if err != nil {
		log.Printf("error: could not parse file (%v)", err)
		return
	}
	sval := string(s)
	sval = strings.TrimSpace(sval)

	r := regexp.MustCompile(titlePattern)
	result := ToMap(r.FindStringSubmatch(sval), r.SubexpNames())

	log.Printf("debug: length of result: %d", len(result))
	e.Title = result["title"]
	e.Content = result["content"]
	if len(e.Title) > 0 {
		if !(len(e.Content) > 0) {
			e.Content = e.Title //Never send an empty body, mainly for other users sake!
		}
		i := Issue{Title: e.Title, Body: e.Content}
		e.Issue = &i
		e.Valid = true
	}
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
