package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func writer2File(head []string, name string, list ...[]string) error {
	buff := bytes.NewBuffer([]byte{})
	for _, item := range list {
		for _, h := range head {
			buff.WriteString(h + "\n")
		}
		for _, s := range item {
			s = strings.TrimPrefix(s, "full:")
			s = strings.TrimPrefix(s, "domain:")
			buff.WriteString(s + "\n")
		}
	}
	return ioutil.WriteFile(name, buff.Bytes(), os.ModePerm)
}

func trimDomain(s string) string {
	s = strings.TrimPrefix(s, "full:")
	s = strings.TrimPrefix(s, "domain:")
	return s
}

func formatDomain(v string, d ...[]string) []string {
	a := make([]string, 0)
	for _, s := range d {
		for _, sv := range s {
			a = append(a, fmt.Sprintf(v, sv))
		}
	}
	return a
}
