package utils

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
)

func logFatal(e error) {
	if e != nil {
		log.Fatal().Msgf("%v", e)
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
