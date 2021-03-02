package main

import (
	"bufio"
	"geosite/hosts"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"v2ray.com/core/app/router"
)

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

func getSites(path, tag string) (int, error) {
	var src map[string]struct{}
	protoList := new(router.GeoSiteList)
	if err := readFiles(filepath.Join(pwd(), v2flySitePathData), protoList); err != nil {
		return 0, err
	}
	switch tag {
	case adsTag:
		for _, i := range protoList.Entry {
			if strings.EqualFold(i.CountryCode, v2flyBlockTag) {
				for _, d := range i.Domain {
					blockList[d.GetValue()] = struct{}{}
				}
				src = blockList
			}
		}
	case cnTag:
		for _, i := range protoList.Entry {
			if strings.EqualFold(i.CountryCode, v2flyDirectTag) {
				for _, d := range i.Domain {
					directList[d.GetValue()] = struct{}{}
				}
				src = directList
			}
		}
	case allowTag:
		src = allowList
	}
	rules := []string{suffixFull, "", suffixDomain, ""}
	return hosts.WriteFile(filepath.Join(path, tag), src, rules, false)
}

func init() {
	block, err := hosts.GetUrlsFromTxt("block.txt")
	if err != nil {
		log.Println("read [block.txt] failed, ignore")
	} else {
		log.Println("init ads list ...")
		hosts.Resolve(getBodyFromUrls(block), blockList)
	}
	
	log.Println("init allow list ...")
	hosts.Resolve(getBodyFromUrls(allowUrls), allowList)
	for _, l := range localList {
		allowList[l] = struct{}{}
	}
	
	log.Println("init cn list ...")
	hosts.Resolve(getBodyFromUrls(directUrls), directList)
	
	log.Println(v2flySites)
	name := filepath.Base(v2flySites)
	if err := getFile(v2flySites, name); err != nil {
		log.Fatalln(err)
	}
	log.Printf("unzip: %s", v2flySitePath)
	if err := unzip(name); err != nil {
		log.Fatalln(err)
	}
}
