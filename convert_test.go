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
	_ = f.SetCellStr(sheetName, "A1", "A1")
	_ = f.SetCellInt(sheetName, "B1", 100)
	_ = f.SetCellValue(sheetName, "A2", "Cell")
	_ = f.SetCellValue(sheetName, "B2", 1)
	if err := f.SaveAs(srcFile); err != nil {
		requirement.Error(err)
	}
	assertion.FileExists(srcFile)
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
				if err := ConvertExcelToCSV(srcFile); err != nil {
					requirement.Error(err)
				}
			} else {
				if err := ConvertExcelToTSV(srcFile); err != nil {
					requirement.Error(err)
				}
			}
			csvFile := strings.Replace(srcFile, filepath.Ext(srcFile), "_"+sheetName+testCase.name, 1)
			assertion.FileExists(csvFile)
			got, err := os.ReadFile(csvFile)
			if err != nil {
				requirement.Error(err)
			}
			assertion.Equal(testCase.expected, string(got))
		})
	}

	if err := os.RemoveAll(testDir); err != nil {
		assertion.Error(err)
	}
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
			b := ConvertStringToCharByte(testCase.input)
			assertion.Equal(testCase.input, string(b))
			c := ConvertStringToCharRune(testCase.input)
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
		{"3", []byte{111, 222}, false},
	}
	for _, testCase := range testCases {
		file := filepath.Join(testDir, "test.txt")
		err := os.WriteFile(file, testCase.data, os.ModePerm)
		if err != nil {
			requirement.Error(err)
		}
		err = RemoveNullByteInFile(file)
		if err != nil {
			requirement.Error(err)
		}
		t.Run(testCase.name, func(*testing.T) {
			if testCase.expected {
				assertion.True(HasNullByteInFile(file))
				return
			}
			assertion.False(HasNullByteInFile(file))
		})
	}
	if err := os.RemoveAll(testDir); err != nil {
		requirement.Error(err)
	}
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
		if err != nil {
			requirement.Error(err)
		}
		t.Run(testCase.name, func(*testing.T) {
			data, err := os.Open(file)
			if err != nil {
				requirement.Error(err)
			}
			defer data.Close()
			r, err := RemoveNullByteInReader(data)
			if err != nil {
				requirement.Error(err)
			}
			if testCase.expected {
				assertion.True(HasNullByteInReader(r))
				return
			}
			assertion.False(HasNullByteInReader(r))
		})
	}
	if err := os.RemoveAll(testDir); err != nil {
		requirement.Error(err)
	}
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
		if err != nil {
			requirement.Error(err)
		}
		err = ReplaceDelimiter(file, ",", "|")
		if err != nil {
			requirement.Error(err)
		}
		b, err := os.ReadFile(file)
		if err != nil {
			requirement.Error(err)
		}
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
	if err := os.RemoveAll(testDir); err != nil {
		requirement.Error(err)
	}
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
		if err != nil {
			requirement.Error(err)
		}
		err = ReplaceDosToUnix(file)
		if err != nil {
			requirement.Error(err)
		}
		b, err := os.ReadFile(file)
		if err != nil {
			requirement.Error(err)
		}
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
	if err := os.RemoveAll(testDir); err != nil {
		requirement.Error(err)
	}
}
