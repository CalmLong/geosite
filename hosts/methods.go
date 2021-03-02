package hosts

import (
	"bufio"
	"io"
	"net/url"
	"os"
	"sort"
	"strings"
)

func WriteFile(fileName string, src map[string]struct{}, prefix []string, force bool) (int, error) {
	file, err := os.Create(fileName)
	if err != nil {
		return 0, err
	}
	dst := make(map[string]struct{}, 0)
	Resolve(src, dst)
	buff := bufio.NewWriter(file)
	total := Classify(dst, buff, prefix, force)
	_ = buff.Flush()
	return total, file.Close()
}

func GetUrlsFromTxt(name string) ([]string, error) {
	tmpUrls := make(map[string]struct{}, 0)
	fi, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	buff := bufio.NewReader(fi)
	for {
		b, _, e := buff.ReadLine()
		if e == io.EOF {
			break
		}
		urlStr := string(b)
		if strings.TrimSpace(urlStr) == "" {
			continue
		}
		if strings.IndexRune(urlStr, j) == 0 {
			continue
		}
		u, err := url.Parse(urlStr)
		if err != nil {
			return nil, err
		}
		tmpUrls[u.String()] = struct{}{}
	}
	urls := make([]string, 0)
	for k, _ := range tmpUrls {
		urls = append(urls, k)
	}
	sort.Strings(urls)
	return urls, fi.Close()
}
