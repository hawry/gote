package buffer

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"log"
	"os"

	"github.com/hawry/gote/helpers"
)

var issues map[int]helpers.Issue
var currentItr int

var bufferFile = "./.gote_buffer"

func init() {
	issues = make(map[int]helpers.Issue)
	currentItr = 0
	load()
}

func setFilePath(newpath string) {
	bufferFile = newpath
}

//Empty will remove all issues from the buffer
func Empty() {
	//Delete verything!
	if hasBuffer() {
		os.Remove(bufferFile)
	}
}

func save() {
	Empty()
	//dont save unless we actually have to
	if len(issues) == 0 {
		return
	}

	var b bytes.Buffer
	ge := gob.NewEncoder(&b)
	gob.Register(map[int]helpers.Issue{})
	gob.Register(helpers.Issue{})
	err := ge.Encode(&issues)
	if err != nil {
		log.Printf("error: could not encode buffer (%v)", err)
		return
	}
	f, err := os.OpenFile(bufferFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Printf("error: could not create buffer file (%v)", err)
		return
	}
	defer f.Close()
	_, err = f.Write(b.Bytes())
	if err != nil {
		log.Printf("error: could not write to buffer (%v)", err)
		return
	}
}

func load() {
	if !hasBuffer() {
		return
	}
	if len(issues) > 0 {
		issues = make(map[int]helpers.Issue)
	}
	var b bytes.Buffer
	f, err := os.Open(bufferFile)
	if err != nil {
		log.Printf("error: could not open buffer file (%v)", err)
		return
	}
	defer f.Close()
	bd, err := ioutil.ReadAll(f)
	b.Write(bd)
	gd := gob.NewDecoder(&b)
	gob.Register(map[int]helpers.Issue{})
	// gob.Register(helpers.Issue{})
	err = gd.Decode(&issues)
	if err != nil {
		log.Printf("error: could not decode buffer (%v)", err)
		return
	}
	currentItr = len(issues)
}

//Save is the exported variant of the save
func Save() {
	save()
}

//Add will add the specified issue to the buffer
func Add(v helpers.Issue) {
	issues[currentItr] = v
	currentItr++
}

//Count will return the number of issues in the buffer
func Count() int {
	return len(issues)
}

//Remove will return one of the issues in the buffer
func Remove() helpers.Issue {
	itr := currentItr - 1
	if v, ok := issues[itr]; ok {
		currentItr--
		//actually remove the element from the underlying map
		delete(issues, itr)
		return v
	}
	return helpers.Issue{}
}

//HasEntry returns true if buffer contains an issue with given id, false otherwise
func Contains(id int) bool {
	if _, ok := issues[id]; ok {
		return true
	}
	return false
}

//ALl returns a copy of the issues in the buffer
func All() map[int]helpers.Issue {
	cpy := make(map[int]helpers.Issue)
	for k, v := range issues {
		cpy[k] = v
	}
	return cpy
}

//SaveMap allows changing the underlying issue map in buffer. This will overwrite any existing buffer in both memory and on disk
func SaveMap(v map[int]helpers.Issue) {
	issues = v
	save()
}

func hasBuffer() bool {
	if _, err := os.Stat(bufferFile); os.IsNotExist(err) {
		return false
	}
	return true
}
