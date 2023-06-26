package utils

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"unicode/utf8"
)

func Dos2Unix(filePath string) error {
	f, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}
	stat, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("stat file: %w", err)
	}
	eol := regexp.MustCompile(`\r\n`)
	f = eol.ReplaceAllLiteral(f, []byte{'\n'})
	return os.WriteFile(filePath, f, stat.Mode())
}

func RemoveNullByte(data []byte) *bytes.Reader {
	b := bytes.ReplaceAll(data, []byte{0}, []byte{})
	reader := bytes.NewReader(b)
	return reader
}

func ReplaceDelimiter(filename string, old, new string) error {
	f, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}
	stat, err := os.Stat(filename)
	if err != nil {
		return fmt.Errorf("stat file: %w", err)
	}
	r, _, _, err := strconv.UnquoteChar(new, '\'')
	if err != nil {
		return fmt.Errorf("convert char: %w", err)
	}
	b := make([]byte, utf8.RuneLen(r))
	utf8.EncodeRune(b, r)
	eol := regexp.MustCompile(old)
	f = eol.ReplaceAllLiteral(f, b)
	return os.WriteFile(filename, f, stat.Mode())
}
