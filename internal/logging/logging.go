package logging

import (
	"fmt"
)

// ANSI escape codes for colors
const (
	red    = "\033[31m"
	cyan   = "\033[36m"
	green  = "\033[32m"
	orange = "\033[33m"
	reset  = "\033[0m"
)

func Error(data ...any) {
	fmt.Printf("%s[ERROR] %s%s\n", red, reset, data)
}

func Info(data ...any) {
	fmt.Printf("%s[INFO] %s%s\n", cyan, reset, data)
}

func Success(data ...any) {
	fmt.Printf("%s[SUCCESS] %s%s\n", green, reset, data)
}

func Debug(data ...any) {
	fmt.Printf("%s[DEBUG] %s%s\n", orange, reset, data)
}
