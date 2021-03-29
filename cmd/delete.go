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
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a word from your wordlist",
	Long: `Delete words from your wordlist that may have been added by mistake,
or do not need to be practiced anymore.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please provide the word to be deleted")
		} else if len(args) > 1 {
			fmt.Printf("Expected only 1 argument, received %v arguments\n", len(args))
		} else {
			if cacheSearch := viper.Get(args[0]); cacheSearch != nil {
				err := Unset(args[0])
				PrintErr(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

func Unset(key string) error {
	configMap := viper.AllSettings()
	delete(configMap, key)
	encodedConfig, _ := json.MarshalIndent(configMap, "", " ")
	err := viper.ReadConfig(bytes.NewReader(encodedConfig))
	if err != nil {
		return err
	}
	viper.WriteConfig()
	return nil
}
