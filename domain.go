package main

import (
	"bufio"
	"geosite/hosts"
	"io"
	"log"
	"net/http"
	"os"
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

func getBodyFormFile(name string) (dst map[string]struct{}) {
	dst = make(map[string]struct{})
	fi, err := os.Open(name)
	if err != nil {
		log.Println("getBodyFormFile:", err)
		return dst
	}
	body := bufio.NewReader(fi)
	for {
		l, _, e := body.ReadLine()
		if e == io.EOF {
			break
		}
		c := string(l)
		if c == "#" {
			continue
		}
		dst[c] = struct{}{}
	}
	return dst
}

func getSites(path, tag, v2flyTag string, extList ...map[string]struct{}) (int, error) {
	var allow, src map[string]struct{}
	
	if len(extList) == 0 {
		protoList := new(router.GeoSiteList)
		if err := readFiles(filepath.Join(pwd(), v2flySitePathData), protoList); err != nil {
			return 0, err
		}
		for _, i := range protoList.Entry {
			if strings.EqualFold(i.CountryCode, v2flyTag) {
				switch v2flyTag {
				case v2flyBlockTag:
					for _, d := range i.Domain {
						blockList[d.GetValue()] = struct{}{}
					}
					allow = allowList
					src = blockList
				case v2flyDirectTag:
					for _, d := range i.Domain {
						directList[d.GetValue()] = struct{}{}
					}
					src = directList
				}
			}
		}
	} else {
		src = make(map[string]struct{})
		for _, list := range extList {
			for k, v := range list {
				src[k] = v
			}
		}
	}
	
	rules := []string{suffixFull, "", suffixDomain, ""}
	return hosts.WriteFile(filepath.Join(path, tag), src, rules, allow)
}

func init() {
	block, err := hosts.GetUrlsFromTxt("block.txt")
	if err != nil {
		log.Println("read [block.txt] failed, ignore")
	} else {
		log.Println("init ads list ...")
		hosts.Resolve(getBodyFromUrls(block), blockList, allowList)
	}
	
	log.Println("init allow list ...")
	hosts.Resolve(getBodyFromUrls(allowUrls), allowList)
	for _, l := range localList {
		allowList[l] = struct{}{}
	}
	
	log.Println("init cn list ...")
	hosts.Resolve(getBodyFromUrls(directUrls), directList)
	
	log.Println("init ptr list ...")
	hosts.Resolve(getBodyFormFile("ptr.txt"), ptrList)
	
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
