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
		} else if numberOfDefs <= 0 {
			fmt.Println("Please enter a valid number.")
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
					data = append(data, dataStruct{strings.ToUpper(wordList[randomIndex]), content})
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
		}

	},
}

func getPageContent(index int, wordList []string) string {
	var cache []Meaning
	err := viper.UnmarshalKey(wordList[index], &cache)
	PrintErr(err)
	resultString := ""
	for _, v := range cache {
		resultString += fmt.Sprintf("%s%s%s\n", colors.Cyan, strings.ToUpper(v.POS), colors.Reset)
		for _, def := range v.Definitions {
			resultString += fmt.Sprintf("%s\u279C%s%s%s\n", colors.Blue, " ", colors.Reset, def.Definition)
			if def.Example != "" {
				resultString += fmt.Sprintf("%s\u2605%s%s%s\n", colors.Yellow, " ", colors.Reset, def.Example)
			}
		}
		resultString += "\n"
	}
	return resultString
}

func init() {
	rootCmd.AddCommand(randomCmd)
	randomCmd.Flags().IntVarP(&numberOfDefs, "number", "n", 1, "Number of random definitions to print")
}
