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

	"github.com/bytedance/sonic"
	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/parquet"
	"github.com/xitongsys/parquet-go/source"
	"github.com/xitongsys/parquet-go/writer"
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
			printError(err)
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
			return wrapError(err)
		}
		for _, row := range rows {
			if err = writer.Write(row); err != nil {
				return wrapError(err)
			}
		}
		writer.Flush()
		csvFile.Close()
		if err = writer.Error(); err != nil {
			return wrapError(err)
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
		return 0, wrapError(err)
	}
	return r, nil
}

/* JSONMarshal returns the JSON encoding bytes of v. */
func JSONMarshal(v any) ([]byte, error) {
	return sonic.Marshal(v)
}

/* MarshalString returns the JSON encoding string of v. */
func JSONMarshalString(v any) (string, error) {
	return sonic.MarshalString(v)
}

/* JSONUnmarshal parses the JSON-encoded data and stores the result in the value pointed to by v. */
func JSONUnmarshal(data []byte, v any) error {
	return sonic.Unmarshal(data, v)
}

/* JSONUnmarshalString is like JSONUnmarshal, except buf is a string. */
func JSONUnmarshalString(data string, v any) error {
	return sonic.UnmarshalString(data, v)
}

/*
ParquetWriter creates the file and a ParquetWriter capable of writing data in parquet format,
obj is a object with tags or JSON schema string.
*/
func ParquetWriter(dstFile string, obj any) (source.ParquetFile, *writer.ParquetWriter, error) {
	fw, err := local.NewLocalFileWriter(dstFile)
	if err != nil {
		return nil, nil, fmt.Errorf("writer: %w", err)
	}
	pw, err := writer.NewParquetWriter(fw, obj, 4)
	if err != nil {
		return nil, nil, fmt.Errorf("parquet handler: %w", err)
	}
	pw.RowGroupSize = 128 * 1024 * 1024 //128M
	pw.PageSize = 8 * 1024              //8K
	pw.CompressionType = parquet.CompressionCodec_ZSTD
	return fw, pw, nil
}

/* RemoveNullByteInFile removes the ASCII 0 in the file. */
func RemoveNullByteInFile(filePath string) error {
	stat, err := os.Stat(filePath)
	if err != nil {
		return wrapError(err)
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		return wrapError(err)
	}
	b := RemoveNullByte(data)
	return os.WriteFile(filePath, b, stat.Mode())
}

/* RemoveNullByteInReader removes the ASCII 0 in reader. */
func RemoveNullByteInReader(reader io.Reader) (io.Reader, error) {
	r, err := io.ReadAll(reader)
	if err != nil {
		return nil, wrapError(err)
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
		return wrapError(err)
	}
	stat, err := os.Stat(filePath)
	if err != nil {
		return wrapError(err)
	}
	b, err := ConvertStringToCharByte(new)
	if err != nil {
		return wrapError(err)
	}
	eol := regexp.MustCompile(old)
	f = eol.ReplaceAllLiteral(f, b)
	return os.WriteFile(filePath, f, stat.Mode())
}

/* ReplaceDosToUnix replaces the Windows end of line(eol) with the Unix eol in the filePath. */
func ReplaceDosToUnix(filePath string) error {
	f, err := os.ReadFile(filePath)
	if err != nil {
		return wrapError(err)
	}
	stat, err := os.Stat(filePath)
	if err != nil {
		return wrapError(err)
	}
	eol := regexp.MustCompile(`\r\n`)
	f = eol.ReplaceAllLiteral(f, []byte{'\n'})
	return os.WriteFile(filePath, f, stat.Mode())
}
