package printer

import "fmt"

const (
	errorColor   string = "\033[1;31m%s\033[0m"
	warningColor        = "\033[1;33m%s\033[0m"
)

func PrintError(format string, a ...any) {
	fmt.Printf(errorColor, fmt.Sprintf(format, a...))
}

func PrintWarning(format string, a ...any) {
	fmt.Printf(warningColor, fmt.Sprintf(format, a...))
}

func Print(format string, a ...any) {
	fmt.Printf(fmt.Sprintf(format, a...))
}
