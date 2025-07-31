package utils

import "runtime/debug"

func GetModulePath() string {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}

	return bi.Main.Path
}