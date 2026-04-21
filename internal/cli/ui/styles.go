package ui

import "github.com/pterm/pterm"

var (
	DebugPrinter = pterm.Info.WithPrefix(pterm.Prefix{
		Text:  "🤖",
		Style: &pterm.ThemeDefault.DebugMessageStyle,
	})
	InfoPrinter = pterm.Info.WithPrefix(pterm.Prefix{
		Text:  "ℹ️",
		Style: &pterm.ThemeDefault.InfoMessageStyle,
	})
	WarnPrinter = pterm.Warning.WithPrefix(pterm.Prefix{
		Text:  "⚠️",
		Style: &pterm.ThemeDefault.WarningMessageStyle,
	})
	ErrorPrinter = pterm.Error.WithPrefix(pterm.Prefix{
		Text:  "❗️",
		Style: &pterm.ThemeDefault.ErrorMessageStyle,
	})
	SuccessPrinter = pterm.Success.WithPrefix(pterm.Prefix{
		Text:  "✅",
		Style: &pterm.ThemeDefault.SuccessMessageStyle,
	})
)
