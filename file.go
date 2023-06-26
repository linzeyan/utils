package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func SkipFirstRow(f *os.File) error {
	row1, err := bufio.NewReader(f).ReadSlice('\n')
	if err != nil {
		return fmt.Errorf("bufio read slice: %w", err)
	}
	_, err = f.Seek(int64(len(row1)), io.SeekStart)
	if err != nil {
		return fmt.Errorf("seek file: %w", err)
	}
	return nil
}
