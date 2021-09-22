package main

import (
	"bufio"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func GetUrlsFromTxt(name string) []string {
	fi, err := os.Open(name)
	if err != nil {
		log.Println(err)
		return []string{}
	}
	defer func() {
		_ = fi.Close()
	}()
	
	urls := make([]string, 0)
	
	sc := bufio.NewScanner(fi)
	for sc.Scan() {
		urlStr := sc.Text()
		if strings.TrimSpace(urlStr) == "" {
			continue
		}
		if strings.IndexRune(urlStr, j) == 0 {
			continue
		}
		u, err := url.Parse(urlStr)
		if err != nil {
			log.Fatalln(err)
		}
		urls = append(urls, u.String())
	}
	
	return urls
}

func getBodyFromUrls(uri ...string) (dst map[string]struct{}) {
	dst = make(map[string]struct{})
	for _, u := range uri {
		log.Println(u)
		resp, err := http.Get(u)
		if err != nil {
			log.Println(err)
			return dst
		}
		sc := bufio.NewScanner(resp.Body)
		for sc.Scan() {
			dst[sc.Text()] = struct{}{}
		}
		_ = resp.Body.Close()
	}
	return dst
}
