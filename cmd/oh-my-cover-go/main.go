package main

import (
	"fmt"
	"os"

	"github.com/IcedElect/oh-my-cover-go/internal/profile"
	"github.com/labstack/gommon/color"
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

	fmt.Printf(`::set-output name=coverage_percent::%.2f`, coveragePercent)
	fmt.Print("\n")

	if coveragePercent < float64(threshold) {
		fmt.Printf(
			color.Red("Coverage percent %.2f%% is below the threshold %d%% \n"), 
			coveragePercent, threshold,
		)
		os.Exit(10)
		return
	}
}
