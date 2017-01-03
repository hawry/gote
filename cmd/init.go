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

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes your repository to be used with gote",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := config.Default()
		if err != nil {
			log.Printf("warning: could not create default configuration (%v)", err)
			return
		}
		log.Printf("debug: %+v", c)
		if _, ferr := os.Stat("./.gitignore"); os.IsNotExist(ferr) {
			b := helpers.Ask("No .gitignore file found. Do you wish to create it?").Bool()
			if b {
				_, err = os.Create("./.gitignore")
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
	},
}

func init() {
	RootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
