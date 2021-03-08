package hosts

import (
	"bufio"
	"io"
	"net/url"
	"os"
	"sort"
	"strings"
)

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
