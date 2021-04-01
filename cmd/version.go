package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get the version of gowords",
	Run: func(cmd *cobra.Command, args []string) {
		os := runtime.GOOS
		arch := runtime.GOARCH
		fmt.Printf("v0.0.8 %s/%s\n", os, arch)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
