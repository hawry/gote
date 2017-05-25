package editor

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/hawry/gote/config"
	"github.com/hawry/gote/helpers"
)

//Editor enables packages outside of the helper package to access the default editor and parses content as well as title
type Editor struct {
	File    *os.File
	Title   string
	Content string
	Valid   bool
	Issue   *helpers.Issue
}

const (
	titlePattern = "(?P<title>\\A.*)\n*(?ms)(?P<content>.*)\\z"
)

//UseEditor returns true if the $EDITOR environment variable is set or if the global configuration specified which editor to use. The method also checks if the chosen editor exists in the PATH, and returns true if the editor will be possible to use - false otherwise. The editor path will also be returned for further use
func UseEditor(cfg config.Global) (bool, string) {
	var editor string
	log.Printf("debug: cfg editor=%s", cfg.Editor)
	if cfg.Editor != "" {
		if strings.HasPrefix(cfg.Editor, "$") {
			editor = os.Getenv(cfg.Editor[1:])
		} else {
			editor = cfg.Editor
		}
	} else {
		editor = os.Getenv("EDITOR")
	}
	p, err := exec.LookPath(editor)
	if err != nil {
		return false, ""
	}
	return true, p
}

//New returns an editor struct, which can't be manipulated
func New(editor string) Editor {
	e := Editor{Valid: false}
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

//Edit returns an editor struct, and puts the given input into the temporary file that's created. The editor struct can't be manipulated
func Edit(editor string, issue helpers.Issue) Editor {
	e := Editor{Valid: false}
	tmpFile, err := ioutil.TempFile(os.TempDir(), "gote_")
	if err != nil {
		log.Printf("error: tempfile could not be created (%v)", err)
		return e
	}
	editorPath, err := exec.LookPath(editor)
	if err != nil {
		log.Printf("error: %s not found in path", editor)
		return e
	}
	if _, err = tmpFile.WriteString(fmt.Sprintf("%s\n", issue.Title)); err != nil {
		log.Printf("warning: could not create modifiable issue title (%v)", err)
	}
	if _, err = tmpFile.WriteString(issue.Body); err != nil {
		log.Printf("warning: could not create modifiable issue body (%v)", err)
	}
	//Reset the I/O offset for the parse command to be able to actually perform the parse. Otherwise it will keep reading from the current cursor location - which will cut the issue body and title in weird ways
	_, err = tmpFile.Seek(0, 0)
	if err != nil {
		log.Printf("error: could not reset IO offset (%v)", err)
	}
	e.File = tmpFile
	cmd := exec.Command(editorPath, e.File.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		log.Printf("error: could not run command (%v)", err)
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

	log.Printf("debug: entire file '%v'", sval)

	r := regexp.MustCompile(titlePattern)
	result := helpers.ToMap(r.FindStringSubmatch(sval), r.SubexpNames())

	log.Printf("debug: length of result: %d", len(result))
	e.Title = result["title"]
	e.Content = result["content"]
	if len(e.Title) > 0 {
		if !(len(e.Content) > 0) {
			e.Content = e.Title //Never send an empty body, mainly for other users sake!
		}
		i := helpers.Issue{Title: e.Title, Body: e.Content}
		e.Issue = &i
		e.Valid = true
	}
}
