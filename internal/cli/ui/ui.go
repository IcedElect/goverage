package ui

import "github.com/pterm/pterm"

func Printlnf(format string, args ...interface{}) {
	pterm.Printfln(format, args...)
}

func Debuglnf(format string, args ...interface{}) {
	DebugPrinter.Printfln(format, args...)
}

func Infolnf(format string, args ...interface{}) {
	InfoPrinter.Printfln(format, args...)
}

func Warnlnf(format string, args ...interface{}) {
	WarnPrinter.Printfln(format, args...)
}

func Errorlnf(format string, args ...interface{}) {
	ErrorPrinter.Printfln(format, args...)
}

func Successlnf(format string, args ...interface{}) {
	SuccessPrinter.Printfln(format, args...)
}
