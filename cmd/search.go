package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Parth576/gowords/colors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Body struct {
	Word     string    `json:"word"`
	Meanings []Meaning `json:"meanings"`
}

type Meaning struct {
	POS         string `json:"partOfSpeech" mapstructure:"partOfSpeech"`
	Definitions []Def  `json:"definitions" mapstructure:"definitions"`
}

type Def struct {
	Definition string `json:"definition" mapstructure:"definition"`
	//Synonyms   []string `json:"synonyms"`
	Example string `json:"example" mapstructure:"example"`
}

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search definition for any word using an online dictionary, Example: gowords search hello",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please provide word to be searched")
		} else if len(args) > 1 {
			fmt.Printf("Only 1 argument needed but %v were provided\n", len(args))
		} else {
			if viper.IsSet(args[0]) {
				// Fetching definition from cache
				var cache []Meaning
				err := viper.UnmarshalKey(args[0], &cache)
				PrintErr(err)
				for _, v := range cache {
					fmt.Printf("%s%s%s\n", colors.Cyan, strings.ToUpper(v.POS), colors.Reset)
					for _, def := range v.Definitions {
						fmt.Printf("%s\u279C%s%s%s\n", colors.Blue, " ", colors.Reset, def.Definition)
						if def.Example != "" {
							fmt.Printf("%s\u2605%s%s%s\n", colors.Yellow, " ", colors.Reset, def.Example)
						}
					}
					fmt.Println()
				}
			} else {
				fetch(args[0])
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}

func fetch(word string) {
	var reqURL = "https://api.dictionaryapi.dev/api/v2/entries/en_US/" + word
	res, err := http.Get(reqURL)
	PrintErr(err)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	PrintErr(err)
	if res.StatusCode == 200 {
		var result []Body
		err = json.Unmarshal(body, &result)
		PrintErr(err)
		var cacheSave []Meaning
		for _, val := range result {
			cacheSave = val.Meanings
			for _, meaning := range val.Meanings {
				fmt.Printf("%s%s%s\n", colors.Cyan, strings.ToUpper(meaning.POS), colors.Reset)
				for _, defs := range meaning.Definitions {
					fmt.Printf("%s\u279C%s%s%s\n", colors.Blue, " ", colors.Reset, defs.Definition)
					if defs.Example != "" {
						fmt.Printf("%s\u2605%s%s%s\n", colors.Yellow, " ", colors.Reset, defs.Example)
					}
				}
				fmt.Println()
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
