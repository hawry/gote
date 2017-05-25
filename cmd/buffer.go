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
	"fmt"
	"log"

	"github.com/hawry/gote/config"
	"github.com/hawry/gote/helpers/buffer"
	"github.com/hawry/gote/helpers/editor"
	"github.com/spf13/cobra"
)

// flagCmd represents the flag command
var bufferCmd = &cobra.Command{
	Use:   "buffer",
	Short: "Show and edit the issue buffer. Running the command without any options will show the current buffer size",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		globalCfg, useGlobal, err := config.LoadGlobal()
		if useGlobal {
			if err != nil {
				log.Printf("error: something went wrong trying to parse the global configuration (%v)", err)
				return
			}
		}

		if clearBuffer {
			buffer.Empty()
			log.Printf("info: removed all entries from buffer")
			return
		}

		if doSendBuffer {
			doNote(nil, nil)
			return
		}

		if editIssue >= 0 {
			log.Printf("debug: trying to edit %d", editIssue)
			if !buffer.Contains(editIssue) {
				log.Printf("warning: the issue ID you provided doesn't exist in the buffer")
				return
			}
			if b, ed := editor.UseEditor(globalCfg); b {
				issue, err := buffer.Find(editIssue)
				if err != nil {
					log.Printf("error: could not open issue for modification (%v)", err)
					return
				}
				log.Printf("debug: editing issue: %v", issue)
				ed := editor.Edit(ed, issue)
				if !ed.Valid {
					log.Printf("warning: empty issue. no changes were saved")
					return
				}
				buffer.Overwrite(editIssue, *ed.Issue)
				log.Printf("success: edited issue '%s'", ed.Title)
			} else {
				log.Printf("error: you must specify an editor to use to modify buffered issues. no changes were saved")
				return
			}
		}

		if rmSingleIssue >= 0 {
			log.Printf("debug: trying to remove %d", rmSingleIssue)
			if !buffer.Contains(rmSingleIssue) {
				log.Printf("warning: the issue ID you provided doesn't exist in the buffer")
				return
			}
			//delete flag is set, remove from map (if it exists) and re-save
			bufferCopy := buffer.All()
			iTitle := bufferCopy[rmSingleIssue].Title
			delete(bufferCopy, rmSingleIssue)
			buffer.SaveMap(bufferCopy)
			log.Printf("success: removed buffer entry %d - %s", rmSingleIssue, iTitle)
		}

		bufferSize := buffer.Count()
		log.Printf("info: there are %d issues buffered", bufferSize)
		for k, v := range buffer.All() {
			fmt.Printf("[ID: %d]\t%s\n", k, v.Title)
		}
	},
}

var rmSingleIssue, editIssue int
var clearBuffer, doSendBuffer bool

func init() {
	RootCmd.AddCommand(bufferCmd)
	bufferCmd.Flags().IntVarP(&editIssue, "edit", "e", -1, "Edit the issue with given ID (requires an editor)")
	bufferCmd.Flags().IntVarP(&rmSingleIssue, "delete", "d", -1, "Delete the issue with given ID")
	bufferCmd.Flags().BoolVarP(&clearBuffer, "clear", "c", false, "Removes all buffered issues")
	bufferCmd.PersistentFlags().BoolVarP(&doSendBuffer, "send", "s", false, "Try to send all buffered issues now")
}
