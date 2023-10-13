package printer

import "fmt"

const (
	errorColor   string = "\033[1;31m%s\033[0m"
	warningColor        = "\033[1;33m%s\033[0m"
)

// PrintError formats an error according to a format specifier and writes it to standard output.
// It returns the number of bytes written and any write error encountered.
func PrintError(format string, a ...any) (int, error) {
	return fmt.Printf(errorColor, fmt.Sprintf(format, a...))
}

// PrintWarning formats a warning according to a format specifier and writes it to standard output.
// It returns the number of bytes written and any write error encountered.
func PrintWarning(format string, a ...any) (int, error) {
	return fmt.Printf(warningColor, fmt.Sprintf(format, a...))
}

// Print formats according to a format specifier and writes it to standard output.
// It returns the number of bytes written and any write error encountered.
func Print(format string, a ...any) (int, error) {
	return fmt.Printf(fmt.Sprintf(format, a...))
}
