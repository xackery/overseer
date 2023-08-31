package zip

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Unpack unzips or gunzips the provided path
func Unpack(srcFile string, dstDir string) error {
	ext := filepath.Ext(srcFile)
	if ext == ".gz" {
		err := gunzip(srcFile, dstDir)
		if err != nil {
			return fmt.Errorf("gunzip: %w", err)
		}
	}
	if ext != ".zip" {
		return fmt.Errorf("invalid extension: %s", ext)
	}
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

func gunzip(srcFile string, dstDir string) error {
	r, err := os.Open(srcFile)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}
	defer r.Close()

	uncompressedStream, err := gzip.NewReader(r)
	if err != nil {
		return fmt.Errorf("newreader: %w", err)
	}

	tarReader := tar.NewReader(uncompressedStream)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("next: %w", err)
		}

		// Construct the path for the extracted file
		extractedFilePath := filepath.Join(dstDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			// Create directories
			if err := os.MkdirAll(extractedFilePath, os.ModePerm); err != nil {
				return fmt.Errorf("mkdirall: %w", err)
			}
		case tar.TypeReg:
			// Create regular files
			outFile, err := os.Create(extractedFilePath)
			if err != nil {
				return fmt.Errorf("create: %w", err)
			}
			defer outFile.Close()

			if _, err := io.Copy(outFile, tarReader); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported file type: %v", header.Typeflag)
		}
	}
	return nil
}
