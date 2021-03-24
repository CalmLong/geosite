package main

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
)

func loadEntry() map[string]*List {
	ref := make(map[string]*List)
	ref[localTag] = getEntry(localTag, localList)
	ref[allowTag] = getEntry(allowTag, allowList)
	ref[blockTag] = getEntry(blockTag, blockList)
	ref[cnTag] = getEntry(cnTag, cnList)
	return ref
}

func getEntry(name string, value map[string]struct{}) *List {
	full := make([]string, 0)
	domain := make([]string, 0)
	other := make([]string, 0)
	
	for k := range value {
		if strings.Contains(k, suffixFull) {
			full = append(full, k)
			continue
		}
		if strings.Contains(k, suffixDomain) {
			domain = append(domain, k)
			continue
		}
		other = append(other, k)
	}
	
	sort.Strings(full)
	sort.Strings(domain)
	sort.Strings(other)
	
	lines := make([]string, 0)
	lines = append(append(full, domain...), other...)
	
	list := &List{
		Name: name,
	}
	
	for _, line := range lines {
		entry, err := parseEntry(line)
		if err != nil {
			log.Fatalln(err)
		}
		list.Entry = append(list.Entry, entry)
	}
	return list
}

func initSuffix(uri string) {
	log.Println(uri)
	resp, err := http.Get(uri)
	if err != nil {
		log.Fatalln(err)
	}
	body := bufio.NewReader(resp.Body)
	for {
		l, _, e := body.ReadLine()
		if e == io.EOF {
			break
		}
		line := string(l)
		if strings.Contains(line, "//") {
			continue
		}
		if strings.Count(line, ".") != 1 {
			continue
		}
		if idx := strings.IndexRune(line, '*'); idx != -1 {
			line = line[idx+1:]
		}
		line = "." + line
		suffixList[line] = struct{}{}
	}
}

func init() {
	block, err := GetUrlsFromTxt("block.txt")
	if err != nil {
		log.Println("read [block.txt] failed, ignore")
	} else {
		log.Println("init ads list ...")
		Resolve(getBodyFromUrls(block), blockList)
		Resolve(getBodyFromUrls([]string{domainListAdsAllRaw}), blockList)
	}
	
	log.Println("init suffix list ...")
	initSuffix(suffixListRaw)
	
	log.Println("init allow list ...")
	Resolve(getBodyFromUrls(allowUrls), allowList)
	
	log.Println("init cn list ...")
	Resolve(getBodyFromUrls(directUrls), cnList)
	Resolve(getBodyFromUrls([]string{domainListCnRaw}), cnList)
}
