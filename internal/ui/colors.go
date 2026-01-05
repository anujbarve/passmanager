// internal/ui/colors.go
package ui

import "fmt"

// ANSI color codes
const (
	Reset      = "\033[0m"
	Bold       = "\033[1m"
	Dim        = "\033[2m"
	Italic     = "\033[3m"
	Underline  = "\033[4m"
	
	Black      = "\033[30m"
	Red        = "\033[31m"
	Green      = "\033[32m"
	Yellow     = "\033[33m"
	Blue       = "\033[34m"
	Magenta    = "\033[35m"
	Cyan       = "\033[36m"
	White      = "\033[37m"
	
	BgBlack    = "\033[40m"
	BgRed      = "\033[41m"
	BgGreen    = "\033[42m"
	BgYellow   = "\033[43m"
	BgBlue     = "\033[44m"
	BgMagenta  = "\033[45m"
	BgCyan     = "\033[46m"
	BgWhite    = "\033[47m"
)

func Success(msg string) string {
	return fmt.Sprintf("%s%s✓ %s%s", Bold, Green, msg, Reset)
}

func Error(msg string) string {
	return fmt.Sprintf("%s%s✗ %s%s", Bold, Red, msg, Reset)
}

func Warning(msg string) string {
	return fmt.Sprintf("%s%s⚠ %s%s", Bold, Yellow, msg, Reset)
}

func Info(msg string) string {
	return fmt.Sprintf("%s%s→ %s%s", Bold, Cyan, msg, Reset)
}

func Title(msg string) string {
	return fmt.Sprintf("%s%s%s%s", Bold, Magenta, msg, Reset)
}

func Subtle(msg string) string {
	return fmt.Sprintf("%s%s%s", Dim, msg, Reset)
}

func Highlight(msg string) string {
	return fmt.Sprintf("%s%s%s%s", Bold, Yellow, msg, Reset)
}