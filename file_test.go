package utils

import (
	"bufio"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testDir = "testdata"

func createDir(dir string, mode ...fs.FileMode) {
	m := os.ModePerm
	if len(mode) != 0 {
		m = mode[0]
	}
	err := os.MkdirAll(dir, m)
	if err != nil {
		panic(err)
	}
}

func TestFileMode(t *testing.T) {
	assertion := assert.New(t)
	requirement := require.New(t)
	for i, v := range []fs.FileMode{ModeRead, ModeReadOnly, ModeReadExec} {
		dir := fmt.Sprintf("%s%d", testDir, i)
		createDir(dir, v)
		defer os.RemoveAll(dir)
		stat, err := os.Stat(dir)
		requirement.Nil(err)
		requirement.DirExists(dir)
		assertion.Equal(v.Perm(), stat.Mode().Perm())
	}
}

func TestCopyFile(t *testing.T) {
	assertion := assert.New(t)
	requirement := require.New(t)
	createDir(testDir)
	srcFile := filepath.Join(testDir, "test.txt")
	dstFile := filepath.Join(testDir, "text_copy.txt")

	err := CopyFile(srcFile, dstFile)
	requirement.Error(err)

	_, err = os.Create(srcFile)
	requirement.Nil(err)
	err = CopyFile(srcFile, dstFile)
	requirement.Nil(err)
	src, err := os.ReadFile(srcFile)
	requirement.Nil(err)
	dst, err := os.ReadFile(dstFile)
	requirement.Nil(err)
	assertion.Equal(src, dst)
	requirement.Nil(os.RemoveAll(testDir))
}

func TestCopyByReader(t *testing.T) {
	assertion := assert.New(t)
	requirement := require.New(t)
	createDir(testDir)
	srcFile := filepath.Join(testDir, "test.txt")
	dstFile := filepath.Join(testDir, "text_copy.txt")
	err := os.WriteFile(srcFile, []byte{'\n', '\r'}, os.ModePerm)
	requirement.Nil(err)

	src, err := os.Open(srcFile)
	requirement.Nil(err)
	dst, err := os.Create(dstFile)
	requirement.Nil(err)
	err = CopyByReader(src, dst)
	requirement.Nil(err)
	defer src.Close()
	defer dst.Close()
	srcData, err := os.ReadFile(srcFile)
	requirement.Nil(err)
	dstData, err := os.ReadFile(dstFile)
	requirement.Nil(err)
	assertion.Equal(srcData, dstData)
	requirement.Nil(os.RemoveAll(testDir))
}

func TestListFiles(t *testing.T) {
	assertion := assert.New(t)
	requirement := require.New(t)
	createDir(testDir)
	fileList := []string{"a.csv", "a.txt"}
	for i := range fileList {
		_, err := os.Create(filepath.Join(testDir, fileList[i]))
		requirement.Nil(err)
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
		assertion.Equal(len(files), testCase.len)
		for i := range files {
			assertion.FileExists(files[i])
		}
	}
	requirement.Nil(os.RemoveAll(testDir))
}

func TestSkipFirstRow(t *testing.T) {
	assertion := assert.New(t)
	requirement := require.New(t)
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
		requirement.Nil(err)
		f, err := os.Open(file)
		requirement.Nil(err)
		err = SkipFirstRow(f)
		assertion.Nil(err)
		b := bufio.NewScanner(f)
		expected := testCase.data[1:]
		i := 0
		for b.Scan() {
			assertion.Equal(expected[i], b.Text())
			i++
		}
		f.Close()
	}
	requirement.Nil(os.RemoveAll(testDir))
}

func TestZipAndUnZip(t *testing.T) {
	assertion := assert.New(t)
	requirement := require.New(t)
	createDir(testDir)

	srcFile := filepath.Join(testDir, "test.csv")
	err := os.WriteFile(srcFile, []byte{'\n', 1}, os.ModePerm)
	requirement.Nil(err)

	rawSrcData, err := os.ReadFile(srcFile)
	requirement.Nil(err)
	zipFile := strings.Replace(srcFile, filepath.Ext(srcFile), ".zip", 1)

	err = Zip(srcFile, zipFile)
	assertion.Nil(err)
	zipData, err := os.ReadFile(zipFile)
	requirement.Nil(err)
	zipType := http.DetectContentType(zipData)
	assertion.Equal("application/zip", zipType)

	err = UnZip(zipFile, testDir)
	assertion.Nil(err)
	newSrcData, err := os.ReadFile(srcFile)
	requirement.Nil(err)
	assertion.Equal(rawSrcData, newSrcData)
	requirement.Nil(os.RemoveAll(testDir))
}
