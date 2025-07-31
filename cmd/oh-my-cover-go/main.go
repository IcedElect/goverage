package main

import (
	"fmt"

	"github.com/IcedElect/oh-my-cover-go/internal/profile"
	"github.com/sethvargo/go-githubactions"
	"github.com/spf13/cobra"
)

var (
	profileFile string
	outputDir string
	threshold uint16
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "oh-my-cover-go",
		Short: "A tool for report profiling Go test coverage",
		Run: run,
	}

	rootCmd.PersistentFlags().StringVarP(&profileFile, "profile", "p", "", "coverage profile file")
	rootCmd.PersistentFlags().StringVarP(&outputDir, "output", "o", "", "coverage output directory")
	rootCmd.PersistentFlags().Uint16Var(&threshold, "threshold", 0, "coverage threshold")

	rootCmd.Execute()
}

func run(cmd *cobra.Command, args []string) {
	coveragePercent, err := profile.ProcessProfile(profileFile, outputDir)
	if err != nil {
		fmt.Printf("Error processing profile: %v\n", err)
		return
	}

	githubactions.SetOutput("percent", fmt.Sprintf("%.2f", coveragePercent))

	if coveragePercent < float64(threshold) {
		githubactions.Warningf("Coverage percentage %.2f is below the threshold of %d%%", coveragePercent, threshold)
		return
	}
}
