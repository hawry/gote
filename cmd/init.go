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
	"log"
	"os"

	"github.com/hawry/gote/config"
	"github.com/hawry/gote/helpers"
	"github.com/spf13/cobra"
)

var goteIgnore = `# Ignore .gote-files since they can contain personal access tokens
.gote`
var interactiveMode bool
var reinitMode bool
var modifyGitignore bool

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes your repository to be used with gote",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		if reinitMode {
			modifyGitignore = false
			log.Printf("info: Reinitializing gote-configuration")
			os.Remove("./.gote")
		}

		if interactiveMode {
			newConfig := config.Local{}
			newConfig.AccessTokenString = helpers.Ask("Please provide the personal access token for this repository: (enter to leave this blank to fill in manually)").String()
			newConfig.Remote = helpers.Ask("Remote endpoint (with protocol, eg. https://github.com/hawry/gote). Leave blank to fetch this from git config").String()
			newConfig.RepoOwner = helpers.Ask("Username (leave blank to fetch from git config): ").String()
			newConfig.Repository = helpers.Ask("Repository (leave blank to fetch from git config): ").String()
			log.Printf("debug: modify gitignore: %t", modifyGitignore)
			if !modifyGitignore {
				modifyGitignore = helpers.Ask("Add .gote to .gitignore?").Bool()
				if modifyGitignore {
					appendToGitignore()
				}
			}
			_, err := config.Create(&newConfig)
			if err != nil {
				log.Printf("error: could not save new configuration (%v)", err)
				return
			}
			log.Printf("info: configuration file created")
			return
		}

		createDefault()

		if modifyGitignore {
			appendToGitignore()
		}
	},
}

func createDefault() {
	_, err := config.Default()
	if err != nil {
		log.Printf("warning: could not create default configuration (%v)", err)
		return
	}
}

func appendToGitignore() {
	if _, ferr := os.Stat("./.gitignore"); os.IsNotExist(ferr) {
		b := helpers.Ask("No .gitignore file found. Do you wish to create it?").Bool()
		if b {
			_, err := os.Create("./.gitignore")
			if err != nil {
				log.Printf("warning: could not create .gitignore, make sure to add .gote-files manually to minimize the risk of leaking access tokens to remote endpoints")
			}
		} else {
			return
		}
	}
	f, err := os.OpenFile("./.gitignore", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Printf("error: could not open .gitignore for writing (%v)", err)
		return
	}
	if _, err = f.WriteString(goteIgnore); err != nil {
		log.Printf("error: something went wrong when trying to write to .gitignore, please make sure to add .gote-files to your gitignore to avoid pushing access tokens to any remote endpoints")
	}
}

func init() {
	RootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolVarP(&reinitMode, "reinit", "r", false, "Re-initialize the configuration, overwriting the previous one with the new one")
	initCmd.Flags().BoolVarP(&interactiveMode, "interactive", "i", false, "Use interactive mode to set configuration values, default is false")
	initCmd.Flags().BoolVarP(&modifyGitignore, "no-gitignore", "", false, "Don't modify .gitignore file (default to true when --reinit is used, otherwise defaults to false)")
}
