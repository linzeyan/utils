package utils

import (
	"fmt"
	"io"
	"log"
	"os"
)

func CopyByReader(src io.Reader, dst io.Writer) error {
	buf := make([]byte, 1024*1024)
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
		if err := source.Close(); err != nil {
			log.Fatalln(err)
		}
		if err := destination.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	buf := make([]byte, 1024)
	for {
		n, err := source.Read(buf)
		if err != nil {
			if err != io.EOF {
				return nil
			}
			return err
		}
		if _, err := destination.Write(buf[:n]); err != nil {
			return err
		}
	}
}
