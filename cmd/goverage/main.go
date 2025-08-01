package main

import (
	"fmt"

	"github.com/IcedElect/goverage/internal/profile"
	"github.com/labstack/gommon/color"
	"github.com/sethvargo/go-githubactions"
	"github.com/spf13/cobra"
)

var (
	profileFile string
	outputDir string
	strategy string
	threshold uint16
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "goverage",
		Short: "A fantastic tool for report profiling Go test coverage",
		Run: run,
	}

	rootCmd.PersistentFlags().StringVarP(&profileFile, "profile", "p", "", "coverage profile file")
	rootCmd.PersistentFlags().StringVarP(&outputDir, "output", "o", "", "coverage output directory")
	rootCmd.PersistentFlags().StringVarP(&strategy, "strategy", "s", "html", "coverage report strategy (html or stdout)")
	rootCmd.PersistentFlags().Uint16Var(&threshold, "threshold", 0, "coverage threshold")

	rootCmd.Execute()
}

func run(cmd *cobra.Command, args []string) {
	coveragePercent, err := profile.ProcessProfile(profileFile, outputDir)
	if err != nil {
		fmt.Printf("Error processing profile: %v\n", err)
		return
	}

	defer func() {
        if r := recover(); r != nil {
			if coveragePercent < float64(threshold) {
			fmt.Printf(color.Red("Coverage percentage %.2f is below the threshold of %d%%"), coveragePercent, threshold)
			}
        }
    }()

	githubactions.SetOutput("percent", fmt.Sprintf("%.2f", coveragePercent))

	if coveragePercent < float64(threshold) {
		githubactions.Warningf("Coverage percentage %.2f is below the threshold of %d%%", coveragePercent, threshold)
		return
	}
}
