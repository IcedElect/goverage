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
	var outputDir string
	var hosted bool

	var rootCmd = &cobra.Command{
		Use:   "oh-my-cover-go",
		Short: "A tool for report profiling Go test coverage",
		Run: run,
	}

	rootCmd.PersistentFlags().StringVarP(&profileFile, "profile", "p", "", "coverage profile file")
	rootCmd.PersistentFlags().StringVarP(&outputDir, "output", "o", "", "coverage output directory")
	rootCmd.PersistentFlags().BoolVar(&hosted, "hosted", false, "coverage output directory")

	rootCmd.Execute()
}

func run(cmd *cobra.Command, args []string) {
	coverProfileFileName, err := cmd.Flags().GetString("profile")
	if err != nil {
		fmt.Printf("Error getting profile flag: %v\n", err)
		return
	}

	coverOutputDir, err := cmd.Flags().GetString("output")
	if err != nil {
		fmt.Printf("Error getting output flag: %v\n", err)
		return
	}

	hosted, err := cmd.Flags().GetBool("hosted")
	if err != nil {
		fmt.Printf("Error getting hosted flag: %v\n", err)
		return
	}

	err = profile.ProcessProfile(coverProfileFileName, coverOutputDir, hosted)
	if err != nil {
		fmt.Printf("Error processing profile: %v\n", err)
		return
	}
}