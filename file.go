package utils

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	"github.com/duke-git/lancet/v2/fileutil"
)

/* CopyByReader copies the src Reader to the dst Writer. */
func CopyByReader(src io.Reader, dst io.Writer, buffer ...[]byte) error {
	if len(buffer) == 0 {
		buffer = append(buffer, make([]byte, 1024*4))
	}
	buf := buffer[0]
	if buf == nil {
		buf = make([]byte, 1024*4)
	}
	for {
		n, err := src.Read(buf)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("read: %w", err)
		}
		if _, err := dst.Write(buf[:n]); err != nil {
			return fmt.Errorf("write: %w", err)
		}
	}
}

/* CopyFile copies the src file to the dst file. */
func CopyFile(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !sourceFileStat.Mode().IsRegular() {
		return err
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	destination, err := os.Create(dst)
	if err != nil {
		return err
	}

	defer func() {
		source.Close()
		destination.Close()
	}()

	buf := make([]byte, 1024*4)
	return CopyByReader(source, destination, buf)
}

/* ListFiles returns a slice containing file paths for a given content type and file extension. */
func ListFiles(dir, contentType string, fileExteinsion ...string) []string {
	var files []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		fileBytes, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read file: %w, %s", err, path)
		}
		if http.DetectContentType(fileBytes) != contentType {
			return nil
		}
		var fileExt []bool
		for _, v := range fileExteinsion {
			if filepath.Ext(path) != v {
				fileExt = append(fileExt, true)
			}
		}
		if len(fileExteinsion) == len(fileExt) {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		return nil
	}
	return files
}

/* SkipFirstRow skips the first row in f. */
func SkipFirstRow(f *os.File) error {
	row1, err := bufio.NewReader(f).ReadSlice('\n')
	if err != nil {
		return fmt.Errorf("bufio read slice: %w", err)
	}
	_, err = f.Seek(int64(len(row1)), io.SeekStart)
	if err != nil {
		return fmt.Errorf("seek file: %w", err)
	}
	return nil
}

/* UnZip unzips the file and save it to dst directory. */
func UnZip(zip, dst string) error {
	return fileutil.UnZip(zip, dst)
}

/* Zip creates a zip file, src could be a single file or a directory. */
func Zip(src, zip string) error {
	return fileutil.Zip(src, zip)
}
