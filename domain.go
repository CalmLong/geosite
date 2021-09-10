package main

import (
	"log"
	"strings"
	
	"github.com/v2fly/v2ray-core/v4/app/router"
)

func loadEntry() map[string]*List {
	ref := make(map[string]*List)
	ref[adsTag] = getEntry(adsTag, blockList)
	ref[cnTag] = getEntry(cnTag, cnList)
	ref[proxyTag] = getEntry(proxyTag, proxyList)
	return ref
}

func getEntry(name string, value map[string]router.Domain_Type) *List {
	full, domain, plain, regex := getDomain(value)

	lines := make([]string, 0)
	lines = append(append(full, domain...), regex...)
	lines = append(lines, plain...)

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

func getDomain(value map[string]router.Domain_Type) ([]string, []string, []string, []string) {
	full := make([]string, 0)
	domain := make([]string, 0)
	plain := make([]string, 0)
	regex := make([]string, 0)
	
	for k, v := range value {
		vd := strings.ToLower(v.String()) + ":" + k
		
		switch v {
		case router.Domain_Full:
			full = append(full, vd)
		case router.Domain_Domain:
			domain = append(domain, vd)
		case router.Domain_Plain:
			plain = append(plain, vd)
		case router.Domain_Regex:
			regex = append(regex, vd)
		}
	}
	
	return full, domain, plain, regex
}

func init() {
	block, err := GetUrlsFromTxt("block.txt")
	if err != nil {
		log.Println("read [block.txt] failed, ignore")
	} else {
		log.Printf("Prepare %s list ...\n", adsTag)
		Resolve(getBodyFromUrls(block), blockList, router.Domain_Full)
	}

	log.Printf("Prepare %s list ...\n", cnTag)
	Resolve(getBodyFromUrls(directUrls), cnList, router.Domain_Domain)

	log.Printf("Prepare %s list ...\n", proxyTag)
	ResolveV2Ray(getBodyFromUrls([]string{v2rayNotCn}), proxyList)
}
