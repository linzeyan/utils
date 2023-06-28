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

	if err := os.RemoveAll(testDir); err != nil {
		requirement.Error(err)
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
	}

	for _, testCase := range testCases {
		file := filepath.Join(testDir, "test.txt")
		err := os.WriteFile(file, testCase.data, os.ModePerm)
		if err != nil {
			requirement.Error(err)
		}

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
		data, err := os.Open(file)
		if err != nil {
			requirement.Error(err)
		}
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

func TestToInt64(t *testing.T) {
	assertion := assert.New(t)
	testCases := []struct {
		data     any
		expected int64
	}{
		{data: "123", expected: 123},
		{data: 456, expected: 456},
		{data: byte(255), expected: 255},
		{data: rune(111), expected: 111},
		{data: 077, expected: 63},
		{data: 0x77, expected: 119},
		{data: 0b1111, expected: 15},
		{data: map[string]string{}, expected: 0},
		{data: [2]int{1}, expected: 0},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%v", testCase.data), func(*testing.T) {
			i, err := ToInt64(testCase.data)
			if err != nil {
				assertion.Error(err)
			}
			assertion.Equal(testCase.expected, i)
		})
	}
}
