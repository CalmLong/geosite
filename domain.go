package main

import (
	"bufio"
	"geosite/hosts"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
)

func loadEntry() map[string]*List {
	ref := make(map[string]*List)
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

func getBodyFromUrls(urls []string) (dst map[string]struct{}) {
	dst = make(map[string]struct{})
	for _, u := range urls {
		log.Println(u)
		resp, err := http.Get(u)
		if err != nil {
			log.Println(err)
			return dst
		}
		body := bufio.NewReader(resp.Body)
		for {
			l, _, e := body.ReadLine()
			if e == io.EOF {
				break
			}
			dst[string(l)] = struct{}{}
		}
	}
	return dst
}

func init() {
	block, err := hosts.GetUrlsFromTxt("block.txt")
	if err != nil {
		log.Println("read [block.txt] failed, ignore")
	} else {
		log.Println("init ads list ...")
		hosts.Resolve(getBodyFromUrls(block), blockList)
		hosts.Resolve(getBodyFromUrls([]string{domainListAdsAllRaw}), blockList)
	}
	
	log.Println("init allow list ...")
	hosts.Resolve(getBodyFromUrls(allowUrls), allowList)
	for _, l := range localList {
		allowList[l] = struct{}{}
	}
	
	log.Println("init cn list ...")
	hosts.Resolve(getBodyFromUrls(directUrls), cnList)
	hosts.Resolve(getBodyFromUrls([]string{domainListCnRaw}), cnList)
}
