/*
Copyright © 2021 Parth Shah <parthshah576@gmail.com>

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
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export your wordlist to a text file",
	Long:  `Export your wordlist to a text file in the current working directory, so that it can be shared with others.`,
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		PrintErr(err)
		listFilepath := filepath.Join(cwd, "wordlist.txt")
		file, err := os.Create(listFilepath)
		PrintErr(err)
		defer file.Close()
		wordListTemp := viper.GetStringSlice("wordList")
		for _, line := range wordListTemp {
			file.WriteString(line + "\n")
		}
		fmt.Printf("Wordlist successfully written to %s", listFilepath)
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
}
