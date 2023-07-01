package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xuri/excelize/v2"
)

func TestConvertExcel(t *testing.T) {
	assertion := assert.New(t)
	requirement := require.New(t)
	createDir(testDir)
	srcFile := filepath.Join(testDir, "test.xlsx")
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			requirement.Error(err)
		}
	}()
	sheetName := "Sheet1"
	err := f.SetCellStr(sheetName, "A1", "A1")
	requirement.Nil(err)
	err = f.SetCellInt(sheetName, "B1", 100)
	requirement.Nil(err)
	err = f.SetCellValue(sheetName, "A2", "Cell")
	requirement.Nil(err)
	err = f.SetCellValue(sheetName, "B2", 1)
	requirement.Nil(err)
	err = f.SaveAs(srcFile)
	requirement.Nil(err)
	requirement.FileExists(srcFile)

	testCases := []struct {
		name     string
		expected string
	}{
		{name: ".csv", expected: "A1,100\nCell,1\n"},
		{name: ".tsv", expected: "A1\t100\nCell\t1\n"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(*testing.T) {
			if strings.Contains(testCase.name, "csv") {
				err := ConvertExcelToCSV(srcFile)
				requirement.Nil(err)
			} else {
				err := ConvertExcelToTSV(srcFile)
				requirement.Nil(err)
			}
			csvFile := strings.Replace(srcFile, filepath.Ext(srcFile), "_"+sheetName+testCase.name, 1)
			requirement.FileExists(csvFile)
			got, err := os.ReadFile(csvFile)
			requirement.Nil(err)
			assertion.Equal(testCase.expected, string(got))
		})
	}

	requirement.Nil(os.RemoveAll(testDir))
}

func TestConvertStringToChar(t *testing.T) {
	assertion := assert.New(t)
	testCases := []struct {
		input    string
		expected rune
	}{
		{input: "\t", expected: 9},
		{input: "\n", expected: 10},
		{input: "\r", expected: 13},
		{input: ",", expected: 44},
		{input: "|", expected: 124},
	}

	for _, testCase := range testCases {
		t.Run(testCase.input, func(*testing.T) {
			b, err := ConvertStringToCharByte(testCase.input)
			assertion.Nil(err)
			assertion.Equal(testCase.input, string(b))
			c, err := ConvertStringToCharRune(testCase.input)
			assertion.Nil(err)
			assertion.Equal(testCase.expected, c)
			assertion.Equal(testCase.input, fmt.Sprintf("%c", c))
		})
	}
}

func TestRemoveNullByteInFile(t *testing.T) {
	assertion := assert.New(t)
	requirement := require.New(t)
	createDir(testDir)

	testCases := []struct {
		name     string
		data     []byte
		expected bool
	}{
		{"1", []byte{0, 1, 2, 3}, false},
		{"2", []byte{44, 55, 00, 77, 88}, false},
	}

	for _, testCase := range testCases {
		file := filepath.Join(testDir, "test.txt")
		err := os.WriteFile(file, testCase.data, os.ModePerm)
		requirement.Nil(err)

		t.Run(testCase.name, func(*testing.T) {
			b0, err := HasNullByteInFile(file)
			assertion.Nil(err)
			err = RemoveNullByteInFile(file)
			assertion.Nil(err)
			b, err := HasNullByteInFile(file)
			assertion.Nil(err)
			assertion.NotEqual(testCase.expected, b0)
			assertion.Equal(testCase.expected, b)
		})
	}
	requirement.Nil(os.RemoveAll(testDir))
}

func TestRemoveNullByte(t *testing.T) {
	assertion := assert.New(t)
	requirement := require.New(t)
	createDir(testDir)

	testCases := []struct {
		name     string
		data     []byte
		expected bool
	}{
		{"1", []byte{0, 1, 2, 3}, false},
		{"2", []byte{44, 55, 00, 77, 0}, false},
	}

	for _, testCase := range testCases {
		file := filepath.Join(testDir, "test.txt")
		err := os.WriteFile(file, testCase.data, os.ModePerm)
		requirement.Nil(err)
		data, err := os.Open(file)
		requirement.Nil(err)
		defer data.Close()
		t.Run(testCase.name, func(*testing.T) {
			b0, err := HasNullByteInReader(data)
			requirement.Nil(err)
			r, err := RemoveNullByteInReader(data)
			assertion.Nil(err)
			b, err := HasNullByteInReader(r)
			requirement.Nil(err)
			assertion.NotEqual(testCase.expected, b0)
			assertion.Equal(testCase.expected, b)
		})
	}
	requirement.Nil(os.RemoveAll(testDir))
}

func TestReplaceDelimiter(t *testing.T) {
	assertion := assert.New(t)
	requirement := require.New(t)
	createDir(testDir)
	testCases := []struct {
		name     string
		data     []string
		expected int
	}{
		{"1", []string{"123", "456", "abc", "def"}, 3},
		{"2", []string{"dog", "cat"}, 1},
	}
	for _, testCase := range testCases {
		file := filepath.Join(testDir, "test.txt")
		err := os.WriteFile(file, []byte(strings.Join(testCase.data, ",")), os.ModePerm)
		requirement.Nil(err)
		err = ReplaceDelimiter(file, ",", "|")
		requirement.Nil(err)
		b, err := os.ReadFile(file)
		requirement.Nil(err)
		var counter, unexpectd int
		for _, c := range b {
			if c == '|' {
				counter++
			}
			if c == ',' {
				unexpectd++
			}
		}
		t.Run(testCase.name, func(*testing.T) {
			assertion.Equal(testCase.expected, counter)
			assertion.Equal(0, unexpectd)
		})
	}
	requirement.Nil(os.RemoveAll(testDir))
}

func TestReplaceDosToUnix(t *testing.T) {
	assertion := assert.New(t)
	requirement := require.New(t)
	createDir(testDir)
	testCases := []struct {
		name     string
		data     []string
		expected int
	}{
		{"1", []string{"123", "456", "abc", "def"}, 3},
		{"2", []string{"dog", "cat"}, 1},
	}
	for _, testCase := range testCases {
		file := filepath.Join(testDir, "test.txt")
		err := os.WriteFile(file, []byte(strings.Join(testCase.data, "\r\n")), os.ModePerm)
		requirement.Nil(err)
		err = ReplaceDosToUnix(file)
		requirement.Nil(err)
		b, err := os.ReadFile(file)
		requirement.Nil(err)
		var counter, unexpectd int
		for _, c := range b {
			if c == '\n' {
				counter++
			}
			if c == '\r' {
				unexpectd++
			}
		}
		t.Run(testCase.name, func(*testing.T) {
			assertion.Equal(testCase.expected, counter)
			assertion.Equal(0, unexpectd)
		})
	}
	requirement.Nil(os.RemoveAll(testDir))
}
