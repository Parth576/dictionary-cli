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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"github.com/Parth576/gowords/colors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search definition for any word using an online dictionary, Example: gowords search hello",
	Run: func(cmd *cobra.Command, args []string) {
		if cacheSearch := viper.Get(args[0]); cacheSearch != nil {
			//fmt.Println("Fetching definition from cache...")
			for k, v := range cacheSearch.(map[string]interface{}) {
				fmt.Printf("%s%s%s\n", colors.Cyan, strings.ToUpper(k), colors.Reset)
				for _, def := range v.([]interface{}) {
					if strings.HasPrefix(def.(string), "68f3fde1-8c1a-49eb-9f27-8d951b049142") {
						fmt.Printf("%s\u2605%s%s%s\n", colors.Yellow, " ", colors.Reset, def.(string)[36:])
					} else {
						fmt.Printf("%s\u279C%s%s%s\n", colors.Blue, " ", colors.Reset, def)
					}
				}
				fmt.Println()
			}
		} else {
			//fmt.Println("Fetching definition from API...")
			searchInDict(args[0])
		}
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
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
		cacheSave := make(map[string][]interface{})

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
				pos := ""
				for _, key := range keys {
					if key == "partOfSpeech" {
						posTemp := defs[key].(string)
						fmt.Printf("%s%s%s\n", colors.Cyan, strings.ToUpper(posTemp), colors.Reset)
						pos = posTemp
					} else if key == "definitions" {
						//fmt.Printf("%s\n\n", defs[key].([]interface{})[0].(map[string]interface{})["definition"])
						for _, j := range defs[key].([]interface{}) {
							defTemp := j.(map[string]interface{})["definition"]
							example, exampleExists := j.(map[string]interface{})["example"]
							fmt.Printf("%s\u279C%s%s%s\n", colors.Blue, " ", colors.Reset, defTemp)
							cacheSave[pos] = append(cacheSave[pos], defTemp)
							if exampleExists {
								fmt.Printf("%s\u2605%s%s%s\n", colors.Yellow, " ", colors.Reset, example)
								cacheSave[pos] = append(cacheSave[pos], "68f3fde1-8c1a-49eb-9f27-8d951b049142"+example.(string))
							}
						}
						fmt.Println()
					}
				}
			}
		}
		wordListTemp := viper.GetStringSlice("wordList")
		wordListTemp = append(wordListTemp, word)
		viper.Set("wordList", wordListTemp)
		viper.Set(word, cacheSave)
		viper.WriteConfig()

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
