package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/xuri/excelize/v2"
)

func TestConvertExcel(t *testing.T) {
	createDir(testDir)
	srcFile := filepath.Join(testDir, "test.xlsx")
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			log.Debug().Msgf("%v", err)
			t.Error(err)
		}
	}()
	sheetName := "Sheet1"
	_ = f.SetCellStr(sheetName, "A1", "A1")
	_ = f.SetCellInt(sheetName, "B1", 100)
	_ = f.SetCellValue(sheetName, "A2", "Cell")
	_ = f.SetCellValue(sheetName, "B2", 1)
	if err := f.SaveAs(srcFile); err != nil {
		log.Fatal().Msgf("%v", err)
		t.Error(err)
	}
	assert.FileExists(t, srcFile)
	testCases := []struct {
		name     string
		expected string
	}{
		{name: ".csv", expected: "A1,100\nCell,1\n"},
		{name: ".tsv", expected: "A1\t100\nCell\t1\n"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if strings.Contains(testCase.name, "csv") {
				if err := ConvertExcelToCSV(srcFile); err != nil {
					t.Fatal(err)
				}
			} else {
				if err := ConvertExcelToTSV(srcFile); err != nil {
					t.Fatal(err)
				}
			}
			csvFile := strings.Replace(srcFile, filepath.Ext(srcFile), "_"+sheetName+testCase.name, 1)
			assert.FileExists(t, csvFile)
			got, err := os.ReadFile(csvFile)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, testCase.expected, string(got))
		})
	}

	if err := os.RemoveAll(testDir); err != nil {
		t.Error(err)
	}
}

func TestRemoveNullByteInFile(t *testing.T) {
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
			log.Fatal().Msgf("%v", err)
		}
		err = RemoveNullByteInFile(file)
		if err != nil {
			t.Fatal(err)
		}
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.expected {
				assert.True(t, HasNullByteInFile(file))
				return
			}
			assert.False(t, HasNullByteInFile(file))
		})
	}
	_ = os.RemoveAll(testDir)
}

func TestRemoveNullByte(t *testing.T) {
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
			log.Fatal().Msgf("%v", err)
		}
		t.Run(testCase.name, func(t *testing.T) {
			data, err := os.ReadFile(file)
			if err != nil {
				log.Fatal().Msgf("%v", err)
			}
			r := RemoveNullByte(data)
			if testCase.expected {
				assert.True(t, HasNullByteInReader(r))
				return
			}
			assert.False(t, HasNullByteInReader(r))
		})
	}
	_ = os.RemoveAll(testDir)
}

func TestReplaceDelimiter(t *testing.T) {
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
			log.Fatal().Msgf("%v", err)
		}
		err = ReplaceDelimiter(file, ",", "|")
		if err != nil {
			t.Fatal(err)
		}
		b, err := os.ReadFile(file)
		if err != nil {
			log.Fatal().Msgf("%v", err)
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
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.expected, counter)
			assert.Equal(t, 0, unexpectd)
		})
	}
	_ = os.RemoveAll(testDir)
}

func TestReplaceDosToUnix(t *testing.T) {
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
			log.Fatal().Msgf("%v", err)
		}
		err = ReplaceDosToUnix(file)
		if err != nil {
			t.Fatal(err)
		}
		b, err := os.ReadFile(file)
		if err != nil {
			log.Fatal().Msgf("%v", err)
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
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.expected, counter)
			assert.Equal(t, 0, unexpectd)
		})
	}
	_ = os.RemoveAll(testDir)
}
