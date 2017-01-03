package helpers

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

//Input describes a "stateful" type that handles user input
type Input struct {
	question string
	answer   string
	input    *os.File
}

//Ask creates a new Input type and returns it, setting the question
func Ask(q string) *Input {
	i := Input{question: q, input: os.Stdin}
	return &i
}

func (i *Input) setInput(f *os.File) {
	i.input = f
}

//Bool will ask the user the same question until they have answered [Y]es or [N]o
func (i *Input) Bool() bool {
	if len(i.answer) > 0 {
		b, err := strconv.ParseBool(i.answer)
		if err != nil {
			log.Printf("error when parsing: %v", err)
			return false
		}
		return b
	}
	r := bufio.NewReader(i.input)
	for {
		showQuestion(fmt.Sprintf("%s [y/n] ", i.question))
		a, err := r.ReadString('\n')
		if err != nil {
			log.Printf("error when reading input: %v", err)
			return false
		}
		a = strings.TrimSpace(a)
		log.Printf("debug: received '%s'", a)
		if strings.HasPrefix(a, "y") {
			i.answer = "true"
			return true
		} else if strings.HasPrefix(a, "n") {
			i.answer = "false"
			return false
		}
	}

}

func (i *Input) String() string {
	if len(i.answer) > 0 {
		return i.answer
	}
	r := bufio.NewReader(os.Stdin)
	showQuestion(i.question)
	a, err := r.ReadString('\n')
	if err != nil {
		return ""
	}
	i.answer = a
	return a
}

// func (i *Input) Int() int {
//
// }
func showQuestion(q string) {
	fmt.Print(q)
}
