package utils

import (
	"bufio"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

const testDir = "testdata"

func createDir(dir string) {
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Fatal().Msgf("%v", err)
	}
}

func TestCopyFile(t *testing.T) {
	createDir(testDir)
	srcFile := filepath.Join(testDir, "test.txt")
	dstFile := filepath.Join(testDir, "text_copy.txt")
	_, err := os.Create(srcFile)
	if err != nil {
		log.Fatal().Msgf("%v", err)
	}

	err = CopyFile(srcFile, dstFile)
	if err != nil {
		assert.FileExistsf(t, dstFile, "dstFile not found")
		_, ferr := os.Open(dstFile)
		if ferr != nil {
			t.Fatal(ferr)
		}
	}
	src, _ := os.ReadFile(srcFile)
	dst, _ := os.ReadFile(dstFile)
	assert.Equal(t, src, dst)
	_ = os.RemoveAll(testDir)
}

func TestCopyByReader(t *testing.T) {
	createDir(testDir)
	srcFile := filepath.Join(testDir, "test.txt")
	dstFile := filepath.Join(testDir, "text_copy.txt")
	err := os.WriteFile(srcFile, []byte{'\n', '\r'}, os.ModePerm)
	if err != nil {
		log.Fatal().Msgf("%v", err)
	}

	src, err := os.Open(srcFile)
	if err != nil {
		log.Fatal().Msgf("%v", err)
	}
	dst, err := os.Create(dstFile)
	if err != nil {
		log.Fatal().Msgf("%v", err)
	}

	err = CopyByReader(src, dst)
	if err != nil {
		assert.FileExistsf(t, dstFile, "dstFile not found")
	}
	src.Close()
	dst.Close()
	srcData, _ := os.ReadFile(srcFile)
	dstData, _ := os.ReadFile(dstFile)
	assert.Equal(t, srcData, dstData)
	_ = os.RemoveAll(testDir)
}

func TestListFiles(t *testing.T) {
	createDir(testDir)
	fileList := []string{"a.csv", "a.txt"}
	for i := range fileList {
		_, err := os.Create(filepath.Join(testDir, fileList[i]))
		if err != nil {
			log.Fatal().Msgf("%v", err)
		}
	}

	testCases := []struct {
		len  int
		exts []string
	}{
		{1, []string{".csv", ".tsv", ".md"}},
		{2, []string{".csv", ".txt"}},
	}

	for _, testCase := range testCases {
		files := ListFiles(testDir, "text/plain; charset=utf-8", testCase.exts...)
		assert.Equal(t, len(files), testCase.len)
		for i := range files {
			assert.FileExists(t, files[i])
		}
	}
	_ = os.RemoveAll(testDir)
}

func TestSkipFirstRow(t *testing.T) {
	createDir(testDir)
	testCases := []struct {
		data []string
	}{
		{[]string{"1", "6", "8"}},
		{[]string{"apple", "banana", "pineapple"}},
		{[]string{"a,$", "b^!", "|c~"}},
	}
	for _, testCase := range testCases {
		file := filepath.Join(testDir, "test.txt")
		err := os.WriteFile(file, []byte(strings.Join(testCase.data, "\n")), os.ModePerm)
		if err != nil {
			log.Fatal().Msgf("%v", err)
		}
		f, err := os.Open(file)
		if err != nil {
			log.Fatal().Msgf("%v", err)
		}
		err = SkipFirstRow(f)
		if err != nil {
			t.Fatal(err)
		}
		b := bufio.NewScanner(f)
		expected := testCase.data[1:]
		i := 0
		for b.Scan() {
			assert.Equal(t, expected[i], b.Text())
			i++
		}
		f.Close()
	}
	_ = os.RemoveAll(testDir)
}

func TestZipAndUnZip(t *testing.T) {
	createDir(testDir)
	srcFile := filepath.Join(testDir, "test.csv")
	err := os.WriteFile(srcFile, []byte{'\n', 1}, os.ModePerm)
	if err != nil {
		log.Fatal().Msgf("%v", err)
	}
	rawSrcData, err := os.ReadFile(srcFile)
	if err != nil {
		log.Fatal().Msgf("%v", err)
	}
	zipFile := strings.Replace(srcFile, filepath.Ext(srcFile), ".zip", 1)
	err = Zip(srcFile, zipFile)
	if err != nil {
		log.Fatal().Msgf("%v", err)
	}
	zipData, err := os.ReadFile(zipFile)
	if err != nil {
		log.Fatal().Msgf("%v", err)
	}
	zipType := http.DetectContentType(zipData)
	assert.Equal(t, "application/zip", zipType)

	err = UnZip(zipFile, testDir)
	if err != nil {
		log.Fatal().Msgf("%v", err)
	}
	newSrcData, err := os.ReadFile(srcFile)
	if err != nil {
		log.Fatal().Msgf("%v", err)
	}
	assert.Equal(t, rawSrcData, newSrcData)
	_ = os.RemoveAll(testDir)
}
