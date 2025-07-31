package main

import (
	"fmt"

	"github.com/IcedElect/oh-my-cover-go/internal/profile"
	"github.com/spf13/cobra"
)

func init() {

}

func main() {
	var profileFile string

	var rootCmd = &cobra.Command{
		Use:   "oh-my-cover-go",
		Short: "A tool for report profiling Go test coverage",
		Run: run,
	}

	rootCmd.PersistentFlags().StringVar(&profileFile, "profile", "", "coverage profile file")

	rootCmd.Execute()
}

func run(cmd *cobra.Command, args []string) {
	coverProfileFileName, err := cmd.Flags().GetString("profile")
	if err != nil {
		fmt.Printf("Error getting profile flag: %v\n", err)
		return
	}

	err = profile.ProcessProfile(coverProfileFileName, "coverage")
	if err != nil {
		fmt.Printf("Error processing profile: %v\n", err)
		return
	}
}