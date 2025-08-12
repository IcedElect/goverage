package utils

import (
	"fmt"
	"runtime/debug"
)

func GetModulePath() string {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		fmt.Println("Failed to read build info")
		return "unknown"
	}

	return bi.Main.Path
}