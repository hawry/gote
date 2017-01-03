// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2"

	"github.com/google/go-github/github"
	"github.com/hawry/gote/config"
	"github.com/hawry/gote/helpers"
	"github.com/spf13/cobra"
)

var dryRun bool

// noteCmd represents the note command
var noteCmd = &cobra.Command{
	Use:   "note",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := config.LoadDefault()
		if err != nil {
			log.Printf("error: could not open configuration file (%v)", err)
			return
		}

		log.Printf("debug: %+v", config)
		var goteIssue *helpers.Issue

		if helpers.CanUseEditor() {
			e := helpers.NewEditor()
			if !e.Valid {
				log.Printf("warning: empty response, ignoring")
				return
			}
			goteIssue = e.Issue
		} else {
			//Use secondary prompt
			fmt.Print("> ")
			r := bufio.NewReader(os.Stdin)
			rawBody, rerr := r.ReadString('\n')
			if rerr != nil {
				log.Printf("error: could not read from standard input (%v)", rerr)
				return
			}
			if !(len(rawBody) > 0) {
				log.Printf("warn: no content, ignoring note")
				return
			}
			goteIssue = helpers.NewIssue(rawBody)
		}

		tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: config.AccessToken()})
		tokenClient := oauth2.NewClient(oauth2.NoContext, tokenSource)

		cli := github.NewClient(tokenClient)
		defaultLabels := []string{"gote"}
		newIssue := &github.IssueRequest{Title: &goteIssue.Title, Body: &goteIssue.Body, Labels: &defaultLabels}

		if dryRun {
			log.Printf("info: dry run enabled, the following would normally be sent to remote: %+v", goteIssue)
			log.Printf("debug: access token is: %s", config.AccessToken())
			return
		}
		_, response, err := cli.Issues.Create(config.User, config.Repository, newIssue)
		if err != nil {
			log.Printf("error: could not create issue for %s (%v)", config.Remote, err)
			return
		}

		switch response.StatusCode {
		case 201:
			log.Printf("info: created new issue '%s' for remote %s", goteIssue.Title, config.Remote)
		default:
			log.Printf("warning: unknown response code from remote: %d", response.StatusCode)
		}
	},
}

func init() {
	RootCmd.AddCommand(noteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// noteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// noteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	noteCmd.PersistentFlags().BoolVarP(&dryRun, "dry", "d", false, "Do a dry run, to test configuration settings and credentials without creating any issues")
}
