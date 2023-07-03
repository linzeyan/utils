package utils

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/rs/zerolog/log"
)

func createDir(dir string, mode ...fs.FileMode) {
	m := os.ModePerm
	if len(mode) != 0 {
		m = mode[0]
	}
	err := os.MkdirAll(dir, m)
	logPanic(err)
}

func logPanic(e error) {
	if e != nil {
		log.Panic().Caller().Msgf("%v", e)
	}
}

func printError(e error) {
	if e != nil {
		fmt.Fprintln(os.Stderr, wrapError(e))
	}
}

func wrapError(e error) error {
	return fmt.Errorf("utils: %w", e)
}
