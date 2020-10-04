package main

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func pwd(r ...string) string {
	str, err := os.Getwd()
	if err != nil {
		log.Panicln(err)
	}
	dir := str
	if len(r) > 0 && strings.TrimSpace(r[0]) != "" {
		dir = filepath.Join(str, r[0])
	}
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Panicln(err)
	}
	return dir
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil || os.IsExist(err) {
		return true
	}
	return false
}

func unzip(fileName string, output ...string) error {
	r, err := zip.OpenReader(fileName)
	if err != nil {
		return err
	}
	var outPath string
	if output != nil {
		outPath = output[0]
	}
	if len(outPath) > 0 {
		if !isExist(outPath) {
			if err := os.MkdirAll(outPath, os.ModePerm); err != nil {
				return err
			}
		}
	}
	for _, z := range r.Reader.File {
		name := z.Name
		if len(outPath) > 0 {
			name = outPath + name
		}
		if z.FileInfo().IsDir() {
			_ = os.MkdirAll(name, os.ModePerm)
			continue
		}
		r, err := z.Open()
		if err != nil {
			return err
		}
		NewFile, err := os.Create(name)
		if err != nil {
			return err
		}
		_, _ = io.Copy(NewFile, r)
		_ = NewFile.Close()
		_ = r.Close()
	}
	return r.Close()
}
