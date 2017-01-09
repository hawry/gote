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
	"github.com/hawry/gote/helpers/buffer"
	"github.com/hawry/gote/helpers/editor"
	"github.com/hawry/gote/helpers/format"
	"github.com/spf13/cobra"
)

var addToBuffer bool

var globalCfg, localCfg config.Configuration
var accessToken string

// noteCmd represents the note command
var noteCmd = &cobra.Command{
	Use:   "note",
	Short: "Creates and pushes a new issue to the repository specified in the configuration file",
	Long:  `The issue will be created using either a secondary-prompt or in any editor specified by $EDITOR (or in your global configuration)`,
	Run: func(cmd *cobra.Command, args []string) {
		doNote(cmd, args)
	},
}

//doNote is a workaround to keeping the application a bit DRY, since this can also be called from the buffer command
func doNote(cmd *cobra.Command, args []string) {

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

	if doSendBuffer {
		log.Printf("do send buffer is true")
		if bc := buffer.Count(); bc > 0 {
			log.Printf("warning: there are %d buffered issues. sending them now!", bc)
			for bc > 0 {
				if !sendIssue(buffer.Remove(), localCfg.RepoOwner, localCfg.Repository, accessToken) {
					//not really any use to keep hacking away at trying to send if one of them fails...
					log.Printf("warning: could not send buffered issue. Aborting. All unsent issues will be sent the next time a successful transmit is made")
					break
				}
				bc = buffer.Count()
			}
		}
		return
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
	if sendIssue(*goteIssue, localCfg.RepoOwner, localCfg.Repository, accessToken); !addToBuffer {
		if bc := buffer.Count(); bc > 0 {
			log.Printf("warning: there are %d buffered issues. sending them now!", bc)
			for bc > 0 {
				if !sendIssue(buffer.Remove(), localCfg.RepoOwner, localCfg.Repository, accessToken) {
					//not really any use to keep hacking away at trying to send if one of them fails...
					log.Printf("warning: could not send buffered issue. Aborting. All unsent issues will be sent the next time a successful transmit is made")
					break
				}
				bc = buffer.Count()
			}
		}
		log.Printf("info: ")
	}
}

func sendIssue(goteIssue helpers.Issue, repoOwner, repoName, accessToken string) bool {
	log.Printf("%+v", goteIssue)
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
	tokenClient := oauth2.NewClient(oauth2.NoContext, tokenSource)

	cli := github.NewClient(tokenClient)
	defaultLabels := []string{"gote"}
	newIssue := &github.IssueRequest{Title: &goteIssue.Title, Body: &goteIssue.Body, Labels: &defaultLabels}

	if addToBuffer {
		log.Printf("info: adding '%s' to buffer", goteIssue.Title)
		buffer.Add(goteIssue)
		buffer.Save()
		return true
	}

	if dryRun {
		log.Printf("sending issue: %+v", goteIssue)
		// log.Printf("info: dry run enabled, the following would normally be sent to remote: %+v", goteIssue)
		// log.Printf("debug: access token is: %s (using global %t)", accessToken, useGlobal)
		return true
	}
	_, response, err := cli.Issues.Create(repoOwner, repoName, newIssue)
	if err != nil {
		log.Printf("error: could not create issue for %s (%v)", fmt.Sprintf("%s/%s", repoOwner, repoName), err)
		//add to buffer
		log.Printf("warning: this issue was added to the send buffer and will be sent the next time a successful transmission is made")
		buffer.Add(goteIssue)
		buffer.Save()
		return false
	}

	switch response.StatusCode {
	case 201:
		log.Printf("info: created new issue '%s' for %s/%s", goteIssue.Title, repoOwner, repoName)
	default:
		log.Printf("warning: unknown response code from remote: %d", response.StatusCode)
	}
	return true
}

func init() {
	RootCmd.AddCommand(noteCmd)

	noteCmd.Flags().BoolVarP(&addToBuffer, "addbuffer", "b", false, "Add issue to buffer instead of sending directly")
	noteCmd.PersistentFlags().BoolVarP(&dryRun, "dry", "d", false, "Do a dry run, to test configuration settings and credentials without creating any issues")

}
