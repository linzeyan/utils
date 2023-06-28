package utils

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/duke-git/lancet/v2/convertor"
	"github.com/rs/zerolog/log"
	"github.com/xuri/excelize/v2"
)

/* ConvertExcelToTSV converts the Excel file to the text file by delimiter. */
func ConvertExcel(filePath, ext string, delimiter rune) error {
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
		fileName := strings.Replace(filePath, filepath.Ext(filePath), "_"+sheetName+ext, 1)
		csvFile, err := os.Create(fileName)
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
	return ConvertExcel(filePath, ".csv", ',')
}

/* ConvertExcelToTSV converts the Excel file to the TSV format file. */
func ConvertExcelToTSV(filePath string) error {
	return ConvertExcel(filePath, ".tsv", '\t')
}

/* ConvertStringToCharByte converts the given string(char) to a byte slice, if error returns nil. */
func ConvertStringToCharByte(s string) ([]byte, error) {
	r, err := ConvertStringToCharRune(s)
	if err != nil {
		return nil, err
	}
	b := make([]byte, utf8.RuneLen(r))
	utf8.EncodeRune(b, r)
	return b, nil
}

/* ConvertStringToCharRune converts the given string(char) to rune, if error returns 0 and error. */
func ConvertStringToCharRune(s string) (rune, error) {
	r, _, _, err := strconv.UnquoteChar(s, '\'')
	if err != nil {
		return 0, err
	}
	return r, nil
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
	b := RemoveNullByte(data)
	return os.WriteFile(filePath, b, stat.Mode())
}

/* RemoveNullByteInReader removes the ASCII 0 in reader. */
func RemoveNullByteInReader(reader io.Reader) (io.Reader, error) {
	r, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	b := RemoveNullByte(r)
	buf := bytes.NewBuffer(b)
	return buf, nil
}

/* RemoveNullByte removes the ASCII 0 in the data. */
func RemoveNullByte(data []byte) []byte {
	return bytes.ReplaceAll(data, []byte{0}, []byte{})
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
	b, err := ConvertStringToCharByte(new)
	if err != nil {
		return fmt.Errorf("converts: %w", err)
	}
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

/* ToInt64 converts v to int64 value, if input is not numerical, return 0 and error. */
func ToInt64(v any) (int64, error) {
	return convertor.ToInt(v)
}
