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
		for index, line := range wordListTemp {
			if index == len(wordListTemp)-1 {
				file.WriteString(line)
			} else {
				file.WriteString(line + "\n")
			}
		}
		fmt.Printf("Wordlist successfully written to %s", listFilepath)
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
}
