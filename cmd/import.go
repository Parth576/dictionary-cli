package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

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
				if viper.IsSet(word) {
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
						newWordCounter++
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
	var result []Body
	err := json.Unmarshal(body, &result)
	PrintErr(err)
	var cacheSave []Meaning
	for _, val := range result {
		cacheSave = val.Meanings
	}
	wordListTemp := viper.GetStringSlice("wordList")
	wordListTemp = append(wordListTemp, word)
	viper.Set("wordList", wordListTemp)
	viper.Set(word, cacheSave)
	viper.WriteConfig()
}
