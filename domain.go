package main

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

func dict2Slice(itemMap map[string]struct{}) []string {
	item := make([]string, 0)
	for t := range itemMap {
		item = append(item, t)
	}
	sort.Strings(item)
	return item
}

func pickWriter(header []string, name string, items ...[]string) error {
	buff := bytes.NewBuffer([]byte{})
	for _, h := range header {
		buff.WriteString(h)
	}
	for _, item := range items {
		for _, s := range item {
			buff.WriteString(s)
		}
	}
	return ioutil.WriteFile(name, buff.Bytes(), os.ModePerm)
}

func autoPick(f string) {
	t := time.Now().Format("2006-01-02 15:04:05")
	switch f {
	case "clash":
		const rule = "  - DOMAIN,;;  - DOMAIN-SUFFIX,;"

		blockFull, blockDomain := load2Format(rule, blockList)
		cnFull, cnDomain := load2Format(rule, cnList)
		proxyFull, proxyDomain := load2Format(rule, proxyList)

		header := []string{"# TIME: ", t, "\n", "payload:", "\n"}

		if err := pickWriter(header, "block.yaml", blockFull, blockDomain); err != nil {
			log.Fatalln(err)
		}
		if err := pickWriter(header, "cn.yaml", cnFull, cnDomain); err != nil {
			log.Fatalln(err)
		}
		if err := pickWriter(header, "proxy.yaml", proxyFull, proxyDomain); err != nil {
			log.Fatalln(err)
		}
		log.Println("clash.yaml has been generated successfully.")
	default:
		log.Fatalf("unsupported %s", f)
	}
}

func load2Format(format string, domainList map[string]struct{}) ([]string, []string) {
	f := strings.Split(format, ";")
	if len(f) != 4 {
		log.Fatalln("format err: ", format)
	}

	full := make(map[string]struct{})
	domain := make(map[string]struct{})
	for k := range domainList {
		if strings.Contains(k, suffixFull) {
			delete(domainList, k)
			k = replace(k, suffixFull, f[0], f[1]) + "\n"
			full[k] = struct{}{}
			continue
		}
		if strings.Contains(k, suffixDomain) {
			delete(domainList, k)
			k = replace(k, suffixDomain, f[2], f[3]) + "\n"
			domain[k] = struct{}{}
			continue
		}
	}

	return dict2Slice(full), dict2Slice(domain)
}

func loadEntry() map[string]*List {
	ref := make(map[string]*List)
	ref[localTag] = getEntry(localTag, localList)
	ref[allowTag] = getEntry(allowTag, allowList)
	ref[blockTag] = getEntry(blockTag, blockList)
	ref[cnTag] = getEntry(cnTag, cnList)
	ref[proxyTag] = getEntry(proxyTag, proxyList)
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
		ResolveV2Ray(getBodyFromUrls(blockUrlsForV2Ray), blockList)
	}

	log.Println("init suffix list ...")
	initSuffix(suffixListRaw)

	log.Println("init allow list ...")
	Resolve(getBodyFromUrls(allowUrls), allowList)

	log.Println("init cn list ...")
	Resolve(getBodyFromUrls(directUrls), cnList)
	ResolveV2Ray(getBodyFromUrls(directUrlsForV2Ray), cnList)

	log.Println("init proxy list ...")
	ResolveV2Ray(getBodyFromUrls([]string{v2rayNotCn}), proxyList)
}
