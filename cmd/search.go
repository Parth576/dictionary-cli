/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"github.com/Parth576/gowords/colors"
	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "search word using an online dictionary",
	Run: func(cmd *cobra.Command, args []string) {
		searchInDict(args[0])
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// searchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// searchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func searchInDict(word string) {
	var reqURL = "https://api.dictionaryapi.dev/api/v2/entries/en_US/" + word
	res, err := http.Get(reqURL)
	PrintErr(err)
	defer res.Body.Close()
	var result interface{}
	body, err := ioutil.ReadAll(res.Body)
	PrintErr(err)
	if res.StatusCode == 200 {
		err = json.Unmarshal(body, &result)
		PrintErr(err)
		var r = result.([]interface{})[0].(map[string]interface{})
		var meaning = r["meanings"]

		switch meaning := meaning.(type) {
		case []interface{}:
			for _, v := range meaning {
				defs := v.(map[string]interface{})
				keys := make([]string, 0)
				for k, _ := range defs {
					keys = append(keys, k)
				}
				sort.Strings(keys)
				reverse(keys)
				for _, key := range keys {
					if key == "partOfSpeech" {
						fmt.Printf("%s%s%s\n", colors.Cyan, strings.ToUpper(defs[key].(string)), colors.Reset)
					} else if key == "definitions" {
						//fmt.Printf("%s\n\n", defs[key].([]interface{})[0].(map[string]interface{})["definition"])
						for _, j := range defs[key].([]interface{}) {
							fmt.Printf("%s\u279C%s%s%s\n", colors.Blue, " ", colors.Reset, j.(map[string]interface{})["definition"])
						}
						fmt.Println()
					}
				}
			}
		}

	} else if res.StatusCode == 429 {
		fmt.Println("API Rate Limit reached. Please try again after some time.")
	} else if res.StatusCode == 404 {
		fmt.Println("No definition found for " + word)
		//Give option to manually enter definition
	}

}

func reverse(ss []string) {
	last := len(ss) - 1
	for i := 0; i < len(ss)/2; i++ {
		ss[i], ss[last-i] = ss[last-i], ss[i]
	}
}
