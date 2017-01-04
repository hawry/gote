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
	"strings"

	"golang.org/x/oauth2"

	"github.com/google/go-github/github"
	"github.com/hawry/gote/config"
	"github.com/hawry/gote/helpers"
	"github.com/hawry/gote/helpers/editor"
	"github.com/hawry/gote/helpers/format"
	"github.com/spf13/cobra"
)

var dryRun bool

// noteCmd represents the note command
var noteCmd = &cobra.Command{
	Use:   "note",
	Short: "Creates and pushes a new issue to the repository specified in the configuration file",
	Long:  `The issue will be created using either a secondary-prompt or in any editor specified by $EDITOR (or in your global configuration)`,
	Run: func(cmd *cobra.Command, args []string) {
		globalCfg, useGlobal, err := config.LoadGlobal()
		if useGlobal {
			if err != nil {
				log.Printf("error: something went wrong trying to parse the global configuration (%v)", err)
				return
			}
		}

		localCfg, cfgExists, err := config.LoadLocal()
		if !cfgExists {
			log.Printf("error: The configuration is missing. Please run 'gote init' first!")
			return
		}
		if err != nil {
			log.Printf("error: something went wrong when trying to open the local configuration (%v)", err)
			return
		}

		var accessToken string
		var goteIssue *helpers.Issue

		if useGlobal {
			accessToken = globalCfg.AccessToken()
		} else {
			accessToken = localCfg.AccessToken()
		}

		if b, ed := editor.UseEditor(globalCfg); b {
			ed := editor.New(ed)
			if !ed.Valid {
				log.Printf("warning: empty note, ignoring input")
				return
			}
			goteIssue = ed.Issue
		} else {
			fmt.Println("Please enter your issue text below. A newline character will exit this mode and create the issue. Press Ctrl+C to cancel input.")
			fmt.Print("> ")
			r := bufio.NewReader(os.Stdin)
			raw, rerr := r.ReadString('\n')
			if rerr != nil {
				log.Printf("error: could not read input (%v)", rerr)
				return
			}
			raw = strings.TrimRightFunc(raw, format.TrimNewlines)
			if !(len(raw) > 0) {
				log.Printf("warn: empty note, ignoring input")
				return
			}
			goteIssue = helpers.NewIssue(raw)
		}

		tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
		tokenClient := oauth2.NewClient(oauth2.NoContext, tokenSource)

		cli := github.NewClient(tokenClient)
		defaultLabels := []string{"gote"}
		newIssue := &github.IssueRequest{Title: &goteIssue.Title, Body: &goteIssue.Body, Labels: &defaultLabels}

		if dryRun {
			log.Printf("info: dry run enabled, the following would normally be sent to remote: %+v", goteIssue)
			log.Printf("debug: access token is: %s (using global %t)", accessToken, useGlobal)
			return
		}
		_, response, err := cli.Issues.Create(localCfg.RepoOwner, localCfg.Repository, newIssue)
		if err != nil {
			log.Printf("error: could not create issue for %s (%v)", localCfg.Remote, err)
			return
		}

		switch response.StatusCode {
		case 201:
			log.Printf("info: created new issue '%s' for remote %s", goteIssue.Title, localCfg.Remote)
		default:
			log.Printf("warning: unknown response code from remote: %d", response.StatusCode)
		}
	},
}

func init() {
	RootCmd.AddCommand(noteCmd)

	noteCmd.PersistentFlags().BoolVarP(&dryRun, "dry", "d", false, "Do a dry run, to test configuration settings and credentials without creating any issues")
}
