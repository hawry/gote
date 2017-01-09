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

	"github.com/hawry/gote/config"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Force update the current configuration from the git-configuration if it have changed",
	Long:  `Update will assume that the access token setting haven't changed since the last initialization of the configuration, and will thus only modify (if needed) the remote address, repository owner and repository name of the configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		localCfg, cfgExists, err := config.LoadLocal()
		if err != nil {
			log.Printf("error: could not load local configuration (%v)", err)
			if !cfgExists {
				log.Printf("error: the configuration doesn't seem to have been initialized yet")
				return
			}
		}
		accessToken := localCfg.AccessTokenString //Save this value
		_, err = config.Create(&config.Local{AccessTokenString: accessToken})
		if err != nil {
			log.Printf("error: could not update configuration file (%v)", err)
			return
		}

		log.Printf("info: configuration have been updated")
	},
}

func init() {
	RootCmd.AddCommand(updateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
