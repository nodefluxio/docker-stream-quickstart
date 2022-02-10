package zipper

import (
	"archive/zip"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type zipperStruct struct{}

// NewZipper will create an object that represent the zipper interface
func NewZipper() Zipper {
	return &zipperStruct{}
}

func (z *zipperStruct) Create(ctx context.Context, source, target string) (*os.File, error) {
	zipfile, err := os.Create(target)
	if err != nil {
		return nil, err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Stat(source)
	if err != nil {
		return nil, err
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})

	return zipfile, err
}
