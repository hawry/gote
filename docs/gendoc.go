package main

import (
	"log"

	gotecmd "github.com/hawry/gote/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	cmd := gotecmd.RootCmd
	if err := doc.GenMarkdownTree(cmd, "./"); err != nil {
		log.Printf("error: %v", err)
	}
}
