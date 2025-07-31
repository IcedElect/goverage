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

	writeGithubState("coverage_percent", fmt.Sprintf("%.2f", coveragePercent))

	if coveragePercent < float64(threshold) {
		fmt.Printf(
			color.Red("Coverage percent %.2f%% is below the threshold %d%% \n"), 
			coveragePercent, threshold,
		)
		os.Exit(10)
		return
	}
}

func writeGithubState(key, value string) {
	stateFile := os.Getenv("GITHUB_ENV")
    if stateFile == "" {
        fmt.Println("GITHUB_ENV not set")
        os.Exit(1)
    }

    f, err := os.OpenFile(stateFile, os.O_APPEND|os.O_WRONLY, 0600)
    if err != nil {
        fmt.Printf("Failed to open GITHUB_ENV file: %v\n", err)
        os.Exit(1)
    }
    defer f.Close()

    line := fmt.Sprintf("%s=%s\n", key, value)
    if _, err := f.WriteString(line); err != nil {
        fmt.Printf("Failed to write to GITHUB_ENV file: %v\n", err)
        os.Exit(1)
    }
}