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
	"math/rand"
	"strings"
	"time"

	"github.com/Parth576/gowords/colors"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var numberOfDefs int

type dataStruct struct {
	Word    string
	Content string
}

// randomCmd represents the random command
var randomCmd = &cobra.Command{
	Use:   "random",
	Short: "Print random words from the cache with their definition. Useful for memorizing words.",
	Run: func(cmd *cobra.Command, args []string) {
		wordList := viper.GetStringSlice("wordList")
		listLength := len(wordList)
		data := []dataStruct{}

		if numberOfDefs > listLength {
			fmt.Printf("Cache only has %v words but number supplied was %v. \nPlease try again with a number which is in range.\n", listLength, numberOfDefs)
		} else {
			rand.Seed(time.Now().Unix())
			completedWords := make(map[int]struct{}, listLength)
			i := 0
			for i < numberOfDefs {
				randomIndex := rand.Intn(listLength)
				if _, exists := completedWords[randomIndex]; exists {
					continue
				} else {
					i++
					completedWords[randomIndex] = struct{}{}
					content := getPageContent(randomIndex, wordList)
					//fmt.Printf("\n~ %s%s%s ~\n\n%s\n\n", colors.Yellow, strings.ToUpper(wordList[randomIndex]), colors.Reset, content)
					data = append(data, dataStruct{strings.ToUpper(wordList[randomIndex]), content})
				}
			}
		}

		templates := &promptui.SelectTemplates{
			Label:    "{{ . }}",
			Active:   "\u279C {{ .Word | yellow }}",
			Inactive: "{{ .Word | yellow }}",
			Selected: " ",
			Details:  "\n{{ .Content }}",
		}

		prompt := promptui.Select{
			Label:     "Press ENTER to exit",
			Items:     data,
			Templates: templates,
			Size:      10,
		}
		_, _, err := prompt.Run()
		PrintErr(err)
	},
}

func getPageContent(index int, wordList []string) string {
	definition := viper.Get(wordList[index])
	resultString := ""
	if definition != nil {
		for k, v := range definition.(map[string]interface{}) {
			resultString += fmt.Sprintf("%s%s%s\n", colors.Cyan, strings.ToUpper(k), colors.Reset)
			for _, def := range v.([]interface{}) {
				if strings.HasPrefix(def.(string), "68f3fde1-8c1a-49eb-9f27-8d951b049142") {
					resultString += fmt.Sprintf("%s\u2605%s%s%s\n", colors.Yellow, " ", colors.Reset, def.(string)[36:])
				} else {
					resultString += fmt.Sprintf("%s\u279C%s%s%s\n", colors.Blue, " ", colors.Reset, def)
				}
			}
			resultString += "\n"
		}
	}
	return resultString
}

func init() {
	rootCmd.AddCommand(randomCmd)
	randomCmd.Flags().IntVarP(&numberOfDefs, "number", "n", 1, "Number of random definitions to print")
}
