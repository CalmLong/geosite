package main

import (
	"bufio"
	"log"
	"net/http"
)

func download(uri string) []string {
	resp, err := http.Get(uri)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	
	name := make([]string, 0)
	
	sc := bufio.NewScanner(resp.Body)
	for sc.Scan() {
		name = append(name, sc.Text())
	}
	
	return name
}

func loadEntry() map[string]*List {
	ref := make(map[string]*List)
	ref[adsTag] = getEntry(adsTag, Resolve(getBodyFromUrls(GetUrlsFromTxt("block.txt")...)))
	ref[cnTag] = getEntry(cnTag, Resolve(getBodyFromUrls(chinaList)))
	ref[proxyTag] = getEntry(proxyTag, Resolve(getBodyFromUrls(proxyList)))
	return ref
}

func getEntry(name string, domain []string) *List {
	list := &List{
		Name: name,
	}
	
	for i := 0; i < len(domain); i++ {
		entry, err := parseEntry(domain[i])
		if err != nil {
			log.Fatalln(err)
		}
		list.Entry = append(list.Entry, entry)
	}

	return list
}