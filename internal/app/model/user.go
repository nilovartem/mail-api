package model

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// User ...
type User struct {
	Mail string
	Link string
}

// NewLink ...
func (u *User) NewLink(ttl time.Duration) { //start & return timer & return User
	u.Link = uuid.New().String()
	//start timer
	go func() {
		ticker := time.NewTicker(ttl)
		for {
			<-ticker.C
			//link expired, "remove" it from list
			u.Link = ""
		}
	}()
}

// relativePath
func relativePath(path string, anchor string) string {
	_, filename, _ := strings.Cut(path, anchor)
	filename = anchor + filename
	return filename
}

// Zip ...
func (u *User) Zip(filename string, mail string) ([]byte, error) {
	archive := bytes.NewBuffer(nil)
	w := zip.NewWriter(archive)
	walker := func(path string, info os.FileInfo, err error) error {
		filename = relativePath(path, mail)
		if err != nil {
			return err
		}
		if info.IsDir() {
			filename = fmt.Sprintf("%s%c", filename, os.PathSeparator)
			_, err = w.Create(filename)
			return err
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
	err := filepath.Walk(filename, walker)
	if err != nil {
		panic(err)
	}
	w.Close()
	fmt.Println(len(archive.Bytes()))
	return archive.Bytes(), nil
}
