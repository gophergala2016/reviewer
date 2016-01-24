// Copyright © 2016 See CONTRIBUTORS <ignasi.fosch@gmail.com>
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
	"os"

	"github.com/gophergala2016/reviewer/reviewer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

type config struct {
	Authorization struct {
		Token string
	}
	Repositories map[string]struct {
		Username string
		Status   bool
		Required int
	}
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "reviewer",
	Short: "Code review your pull requests",
	Long: `By running reviewer your repo's pull requests will get merged
according to the configuration file.`,
	Run: func(cmd *cobra.Command, args []string) {
		var C config
		err := viper.Unmarshal(&C)
		if err != nil {
			fmt.Printf("Error parsing configuration %v", err)
			return
		}
		client, err := reviewer.GetClient()
		if err != nil {
			fmt.Printf("Error creating GitHub client %v", err)
			return
		}

		//TODO: validate imput parameters (e.g. Requiered = 0)
		for repoName, repoParams := range C.Repositories {
			if repoParams.Status == false {
				fmt.Printf("- %v/%v Discarded (repo disabled)\n", repoParams.Username, repoName)
				continue
			}
			prInfos, err := reviewer.GetPullRequestInfos(client, repoParams.Username, repoName)
			if err != nil {
				fmt.Printf("Error getting pull request info of repo %v/%v", repoParams.Username, repoName)
				return
			}
			fmt.Printf("+ %v/%v\n", repoParams.Username, repoName)
			for _, prInfo := range prInfos {
				if prInfo.Score >= repoParams.Required {
					fmt.Printf("  + %v MERGE (%v) score %v of %v required\n", prInfo.Number, prInfo.Title, prInfo.Score, repoParams.Required)
					// merge here
				} else {
					fmt.Printf("  - %v NOP   (%v) score %v of %v required\n", prInfo.Number, prInfo.Title, prInfo.Score, repoParams.Required)
				}
			}
		}
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.reviewer.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".reviewer") // name of config file (without extension)
	viper.AddConfigPath("$HOME")     // adding home directory as first search path
	viper.SetEnvPrefix("reviewer")   // so viper.AutomaticEnv will get matching envvars starting with REVIEWER_
	viper.AutomaticEnv()             // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
