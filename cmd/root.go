/*
Copyright © 2021 Parth Shah parthshah576@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// PrintErr is used to print the errors to error log
func PrintErr(e error) {
	if e != nil {
		log.Fatalln(e)
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gowords",
	Short: "Get definitions for words or save new words you learn",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

func initConfig() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	configPath := filepath.Join(usr.HomeDir, "gowords.json")
	_, err = os.Stat(configPath)
	if os.IsNotExist(err) {
		var file, err = os.Create(configPath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
	}
	viper.SetConfigFile(configPath)
	//viper.SetDefault("hello", map[string][]interface{}{"EXCLAMATION": []interface{}{"Used as a greeting or to begin a phone conversation."}, "NOUN": []interface{}{"An utterance of \"hello\"; a greeting."}, "INTRANSITIVE VERB": []interface{}{"Say or shout \"hello\"; greet someone."}})
	viper.SetDefault("wordList", []string{})
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("could not find config file")
		} else {
			//fmt.Println(err)
		}
	}
	viper.WriteConfig()
	//viper.WatchConfig()
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	initConfig()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
