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
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"

	"github.com/cheggaaa/pb/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import words into your wordlist",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("1 argument required, none provided")
		} else if len(args) > 1 {
			fmt.Printf("Only 1 argument required, %v were provided", len(args))
		} else {
			lines, err := readLines(args[0])
			PrintErr(err)
			bar := pb.Default.Start(len(lines))
			newWordCounter := 0
			errLog := ""
			for _, word := range lines {
				if cacheSearch := viper.Get(word); cacheSearch != nil {
					//def found in cache, so pass
					bar.Increment()
					continue
				} else {
					//send req and get word def
					var reqURL = "https://api.dictionaryapi.dev/api/v2/entries/en_US/" + word
					res, err := http.Get(reqURL)
					PrintErr(err)
					defer res.Body.Close()
					body, err := ioutil.ReadAll(res.Body)
					PrintErr(err)
					if res.StatusCode == 200 {
						//fetch from api and save in cache
						fetchFromAPI(word, body)
						newWordCounter += 1
						bar.Increment()
					} else if res.StatusCode == 429 {
						//api rate limit reached, need to try again after some time
						fmt.Println("API Rate limit reached. Please try again after 5 minutes.")
						break
					} else if res.StatusCode == 404 {
						//no word def found
						wordNotFound := fmt.Sprintf("No definition found for %s\n", word)
						errLog += wordNotFound
						bar.Increment()
					}
				}
			}
			bar.Finish()
			fmt.Printf("Imported %v new words.\n", newWordCounter)
			fmt.Printf("%s", errLog)
		}
	},
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func init() {
	rootCmd.AddCommand(importCmd)
}

func fetchFromAPI(word string, body []byte) {
	var result interface{}
	err := json.Unmarshal(body, &result)
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
			Reverse(keys)
			pos := ""
			for _, key := range keys {
				if key == "partOfSpeech" {
					posTemp := defs[key].(string)
					//fmt.Printf("%s%s%s\n", colors.Cyan, strings.ToUpper(posTemp), colors.Reset)
					pos = posTemp
				} else if key == "definitions" {
					//fmt.Printf("%s\n\n", defs[key].([]interface{})[0].(map[string]interface{})["definition"])
					for _, j := range defs[key].([]interface{}) {
						defTemp := j.(map[string]interface{})["definition"]
						example, exampleExists := j.(map[string]interface{})["example"]
						//fmt.Printf("%s\u279C%s%s%s\n", colors.Blue, " ", colors.Reset, defTemp)
						cacheSave[pos] = append(cacheSave[pos], defTemp)
						if exampleExists {
							//fmt.Printf("%s\u2605%s%s%s\n", colors.Yellow, " ", colors.Reset, example)
							cacheSave[pos] = append(cacheSave[pos], "68f3fde1-8c1a-49eb-9f27-8d951b049142"+example.(string))
						}
					}
					//fmt.Println()
				}
			}
		}
	}
	wordListTemp := viper.GetStringSlice("wordList")
	wordListTemp = append(wordListTemp, word)
	viper.Set("wordList", wordListTemp)
	viper.Set(word, cacheSave)
	viper.WriteConfig()

}
