package model

import (
	"archive/zip"
	"bytes"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

type User struct {
	Username string
}

// Zip compresses folder "mail" into archive and IGNORES folders in "mail" - for speed purposes
func (u *User) Zip(mail string, root string, pdf string) ([]byte, error) {
	archive := bytes.NewBuffer(nil)
	w := zip.NewWriter(archive)
	walker := func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if filepath.Base(path) == root {
			return nil
		}
		filename, err := filepath.Rel(mail, path)
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return filepath.SkipDir
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		f, err := w.Create(filename)
		if err != nil {
			return err
		}
		_, err = io.Copy(f, file)
		if err != nil {
			return err
		}
		return nil
	}
	err := filepath.WalkDir(mail, walker)
	if err != nil {
		return nil, err
	}
	err = addPDF(w, pdf)
	if err != nil {
		return nil, err
	}
	w.Close()
	return archive.Bytes(), nil
}

// addPDF adds pdf to archive
func addPDF(w *zip.Writer, pdf string) error {
	pdfFile, err := os.Open(pdf)
	if err != nil {
		return err
	}
	defer pdfFile.Close()
	pdfWriter, err := w.Create(filepath.Base(pdf))
	if err != nil {
		return err
	}
	_, err = io.Copy(pdfWriter, pdfFile)
	if err != nil {
		return err
	}
	return nil
}
