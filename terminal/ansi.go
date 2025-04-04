package terminalutils

import (
	"fmt"
	"strings"
)

// ansiCombine wraps a string with ANSI color codes.
// Optimization: Simple string concatenation with minimal allocations.
func ansiCombine(str, color string) string {
	return "\033[" + color + "m" + str + "\033[0m"
}

const (
	Bold            = ";1"
	Italic          = ";3"
	Underline       = ";4"
	ReverseBg       = ";7"
	Strikethrough   = ";9"
	Reset           = "0"
	Black           = "30"
	Red             = "31"
	Green           = "32"
	Yellow          = "33"
	Blue            = "34"
	Magenta         = "35"
	Cyan            = "36"
	White           = "37"
	BrightBlack     = "90"
	BrightRed       = "91"
	BrightGreen     = "92"
	BrightYellow    = "93"
	BrightBlue      = "94"
	BrightMagenta   = "95"
	BrightCyan      = "96"
	BrightWhite     = "97"
	BgBlack         = "40"
	BgRed           = "41"
	BgGreen         = "42"
	BgYellow        = "43"
	BgBlue          = "44"
	BgMagenta       = "45"
	BgCyan          = "46"
	BgWhite         = "47"
	BgBrightBlack   = "100"
	BgBrightRed     = "101"
	BgBrightGreen   = "102"
	BgBrightYellow  = "103"
	BgBrightBlue    = "104"
	BgBrightMagenta = "105"
	BgBrightCyan    = "106"
	BgBrightWhite   = "107"
)

// PrintColored prints a string with ANSI color codes.
// Optimization: Direct output with minimal overhead.
func PrintColored(str string, actions ...string) {
	fmt.Print(SprintColored(str, actions...))
}

// PrintlnColored prints a string with ANSI color codes followed by a newline.
// Optimization: Reuses SprintColored for efficiency.
func PrintlnColored(str string, actions ...string) {
	PrintColored(str+"\n", actions...)
}

// SprintColored returns a string with ANSI color codes applied.
// Optimization: Efficient string joining with minimal allocations.
func SprintColored(str string, actions ...string) string {
	return fmt.Sprint(ansiCombine(str, strings.Join(actions, "")))
}

// Sprintf wraps fmt.Sprintf for formatted string output.
// Optimization: Direct passthrough to standard library function.
func Sprintf(format string, a ...any) string {
	return fmt.Sprintf(format, a...)
}
