package terminal

import (
	"fmt"
	"strings"
)

func ansiCombine(str, color string) string {
	return "\033[" + color + "m" + str + "\033[0m"
}

const (
	Bold          = ";1"
	Italic        = ";3"
	Underline     = ";4"
	ReverseBg     = ";7"
	Strikethrough = ";9"

	Reset         = "0"
	Black         = "30"
	Red           = "31"
	Green         = "32"
	Yellow        = "33"
	Blue          = "34"
	Magenta       = "35"
	Cyan          = "36"
	White         = "37"
	BrightBlack   = "90"
	BrightRed     = "91"
	BrightGreen   = "92"
	BrightYellow  = "93"
	BrightBlue    = "94"
	BrightMagenta = "95"
	BrightCyan    = "96"
	BrightWhite   = "97"

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

func PrintColored(str string, actions ...string) {
	fmt.Print(SprintColored(str, actions...))
}

func PrintlnColored(str string, actions ...string) {
	PrintColored(str+"\n", actions...)
}

func SprintColored(str string, actions ...string) string {
	return fmt.Sprint(ansiCombine(str, strings.Join(actions, "")))
}

func Sprintf(format string, a ...any) string {
	return fmt.Sprintf(format, a...)
}
