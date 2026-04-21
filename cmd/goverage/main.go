package main

import (
	"fmt"

	"github.com/IcedElect/goverage/internal/cli/ui"
	"github.com/IcedElect/goverage/internal/profile"
	"github.com/sethvargo/go-githubactions"
	"github.com/spf13/cobra"
)

var (
	profileFile string
	outputDir   string
	strategy    string
	threshold   uint16
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "goverage",
		Short: "A fantastic tool for report profiling Go test coverage",
		RunE:  run,
	}

	rootCmd.PersistentFlags().StringVarP(&profileFile, "profile", "p", "", "coverage profile file")
	rootCmd.PersistentFlags().StringVarP(&outputDir, "output", "o", "", "coverage output directory")
	rootCmd.PersistentFlags().
		StringVarP(&strategy, "strategy", "s", "html", "coverage report strategy (html or stdout)")
	rootCmd.PersistentFlags().Uint16Var(&threshold, "threshold", 0, "coverage threshold")

	err := rootCmd.Execute()
	if err != nil {
		ui.Errorlnf("Error executing command: %v", err)
	}
}

func run(cmd *cobra.Command, args []string) error {
	coveragePercent, err := profile.ProcessProfile(profileFile, outputDir)
	if err != nil {
		return fmt.Errorf("error processing profile: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			if coveragePercent < float64(threshold) {
				ui.Errorlnf("Coverage percentage %.2f is below the threshold of %d%%", coveragePercent, threshold)
			}
		}
	}()

	githubactions.SetOutput("percent", fmt.Sprintf("%.2f", coveragePercent))

	if coveragePercent < float64(threshold) {
		githubactions.Warningf("Coverage percentage %.2f is below the threshold of %d%%", coveragePercent, threshold)
	}

	return nil
}
