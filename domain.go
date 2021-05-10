package main

import (
	"bufio"
	"github.com/v2fly/v2ray-core/v4/app/router"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
)

func loadEntry() map[string]*List {
	ref := make(map[string]*List)
	ref[adsTag] = getEntry(adsTag, blockList)
	ref[cnTag] = getEntry(cnTag, cnList)
	ref[proxyTag] = getEntry(proxyTag, proxyList)
	return ref
}

func getEntry(name string, value map[string]dT) *List {
	full := make([]string, 0)
	domain := make([]string, 0)
	regex := make([]string, 0)
	for _, vv := range value {
		if vv.Keep {
			switch vv.Type {
			case router.Domain_Regex:
				regex = append(regex, vv.Format)
			case router.Domain_Domain:
				domain = append(domain, vv.Format)
			case router.Domain_Full:
				full = append(full, vv.Format)
			}
		} else {
			switch vv.Type {
			case router.Domain_Regex:
				regex = append(regex, strings.ToLower(router.Domain_Regex.String())+":"+vv.Value)
			case router.Domain_Domain:
				domain = append(domain, strings.ToLower(router.Domain_Domain.String())+":"+vv.Value)
			case router.Domain_Full:
				full = append(full, strings.ToLower(router.Domain_Full.String())+":"+vv.Value)
			}
		}
	}
	
	sort.Strings(full)
	sort.Strings(domain)
	
	lines := make([]string, 0)
	lines = append(append(full, domain...), regex...)
	
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
		log.Printf("init %s list ...\n", adsTag)
		Resolve(getBodyFromUrls(block), blockList)
		ResolveV2Ray(getBodyFromUrls(blockUrlsForV2Ray), blockList)
	}
	
	log.Println("init suffix list ...")
	initSuffix(suffixListRaw)
	
	log.Println("init allow list ...")
	Resolve(getBodyFromUrls(allowUrls), allowList)
	for ak, av := range allowList {
		if bv, ok := blockList[ak]; ok {
			if bv.Type == router.Domain_Domain && av.Type == router.Domain_Domain || bv.Type == router.Domain_Full && av.Type == router.Domain_Full {
				delete(blockList, ak)
			}
		}
	}
	
	log.Printf("init %s list ...\n", cnTag)
	Resolve(getBodyFromUrls(directUrls), cnList)
	ResolveV2Ray(getBodyFromUrls(directUrlsForV2Ray), cnList)
	
	log.Printf("init %s list ...\n", proxyTag)
	ResolveV2Ray(getBodyFromUrls([]string{v2rayNotCn}), proxyList)
}
