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
				wordList := viper.GetStringSlice("wordList")
				wordList = remove(wordList, args[0])
				viper.Set("wordList", wordList)
				viper.WriteConfig()
			} else {
				fmt.Printf("%s not found in the wordlist.", args[0])
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

// https://github.com/spf13/viper/issues/632#issuecomment-581597492
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

// https://stackoverflow.com/a/34070691/12370518
func remove(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}
