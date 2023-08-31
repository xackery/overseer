package zip

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Unpack unzips the provided path
func Unpack(srcFile string, dstDir string) error {
	r, err := zip.OpenReader(srcFile)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}
	defer r.Close()

	for _, f := range r.File {
		filePath := filepath.Join(dstDir, f.Name)
		if f.FileInfo().IsDir() {
			err := os.MkdirAll(filePath, os.ModePerm)
			if err != nil {
				return fmt.Errorf("mkdirall: %w", err)
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return fmt.Errorf("mkdirall: %w", err)
		}

		outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return fmt.Errorf("openfile: %w", err)
		}

		rc, err := f.Open()
		if err != nil {
			return fmt.Errorf("open: %w", err)
		}

		_, err = io.Copy(outFile, rc)
		if err != nil {
			return fmt.Errorf("copy: %w", err)
		}

		outFile.Close()
		rc.Close()
	}

	return nil
}
