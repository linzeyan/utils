package utils

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/rs/zerolog/log"
	"github.com/xuri/excelize/v2"
)

func convertExcel(filePath, ext string, delimiter rune) error {
	excelFile, err := excelize.OpenFile(filePath)
	if err != nil {
		return err
	}
	defer func() {
		if err := excelFile.Close(); err != nil {
			log.Fatal().Msgf("%v", err)
		}
	}()

	for _, sheetName := range excelFile.GetSheetList() {
		csvFile, err := os.Create(strings.Replace(filePath, filepath.Ext(filePath), "_"+sheetName+ext, 1))
		if err != nil {
			return err
		}
		writer := csv.NewWriter(csvFile)
		writer.Comma = delimiter
		rows, err := excelFile.GetRows(sheetName)
		if err != nil {
			return fmt.Errorf("get rows: %w", err)
		}
		for _, row := range rows {
			if err = writer.Write(row); err != nil {
				return fmt.Errorf("write: %w", err)
			}
		}
		writer.Flush()
		csvFile.Close()
		if err = writer.Error(); err != nil {
			return fmt.Errorf("flush: %w", err)
		}
	}
	return nil
}

/* ConvertExcelToCSV converts the Excel file to the CSV format file. */
func ConvertExcelToCSV(filePath string) error {
	return convertExcel(filePath, ".csv", ',')
}

/* ConvertExcelToTSV converts the Excel file to the TSV format file. */
func ConvertExcelToTSV(filePath string) error {
	return convertExcel(filePath, ".tsv", '\t')
}

/* ConvertStringToCharByte converts the given string(char) to a byte slice. */
func ConvertStringToCharByte(s string) []byte {
	r := ConvertStringToCharRune(s)
	b := make([]byte, utf8.RuneLen(r))
	utf8.EncodeRune(b, r)
	return b
}

/* ConvertStringToCharRune converts the given string(char) to a rune. */
func ConvertStringToCharRune(s string) rune {
	r, _, _, err := strconv.UnquoteChar(s, '\'')
	if err != nil {
		return 0
	}
	return r
}

func removeNullByte(data []byte) []byte {
	return bytes.ReplaceAll(data, []byte{0}, []byte{})
}

/* RemoveNullByteInFile removes the ASCII 0 in the file. */
func RemoveNullByteInFile(filePath string) error {
	stat, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	b := removeNullByte(data)
	return os.WriteFile(filePath, b, stat.Mode())
}

/* RemoveNullByte removes the ASCII 0 in the data. */
func RemoveNullByte(data []byte) *bytes.Reader {
	b := removeNullByte(data)
	return bytes.NewReader(b)
}

/* ReplaceDelimiter replaces the old delimiter with the new delimiter in the filePath. */
func ReplaceDelimiter(filePath string, old, new string) error {
	f, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}
	stat, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("stat file: %w", err)
	}
	b := ConvertStringToCharByte(new)
	eol := regexp.MustCompile(old)
	f = eol.ReplaceAllLiteral(f, b)
	return os.WriteFile(filePath, f, stat.Mode())
}

/* ReplaceDosToUnix replaces the Windows end of line(eol) with the Unix eol in the filePath. */
func ReplaceDosToUnix(filePath string) error {
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
