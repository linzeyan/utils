package utils

import (
	"fmt"
	"os"
)

func printError(e error) {
	if e != nil {
		fmt.Fprintln(os.Stderr, wrapError(e))
	}
}

func wrapError(e error) error {
	return fmt.Errorf("utils: %w", e)
}
