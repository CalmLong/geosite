package main

import (
	"bufio"
	"errors"
	"geosite/hosts"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"v2ray.com/core/app/router"
)

func getBodyFromUrls(urls []string) ([]io.Reader, error) {
	body := make([]io.Reader, 0)
	for _, u := range urls {
		log.Println(u)
		resp, err := http.Get(u)
		if err != nil {
			return body, err
		}
		body = append(body, resp.Body)
	}
	return body, nil
}

func getSites(path, suffix, tag string) (int, error) {
	protoList := new(router.GeoSiteList)
	if err := readFiles(pwd()+v2flySitePathData, protoList); err != nil {
		return 0, err
	}
	
	orgTag, err := os.Create(suffix)
	if err != nil {
		return 0, err
	}
	
	buff := bufio.NewWriter(orgTag)
	allow := map[string]struct{}{}
	
	for _, i := range protoList.Entry {
		if strings.EqualFold(i.CountryCode, tag) {
			for _, d := range i.Domain {
				_, err := buff.WriteString(d.GetValue() + "\n")
				if err != nil {
					return 0, err
				}
			}
			switch tag {
			case v2flyBlockTag:
				for k := range blockList {
					_, err := buff.WriteString(k + "\n")
					if err != nil {
						return 0, err
					}
				}
				if err := buff.Flush(); err != nil {
					return 0, err
				}
				allow = allowList
			case v2flyDirectTag:
				for k := range directList {
					_, err := buff.WriteString(k + "\n")
					if err != nil {
						return 0, err
					}
				}
				if err := buff.Flush(); err != nil {
					return 0, err
				}
			default:
				return 0, errors.New("unsupported tag")
			}
		}
	}
	
	_ = orgTag.Close()
	
	reader, err := os.Open(orgTag.Name())
	if err != nil {
		return 0, err
	}
	
	rules := []string{suffixFull, "", suffixDomain, ""}
	total, err := hosts.WriteFile(path+"/"+suffix, []io.Reader{reader}, rules, allow)
	if err != nil {
		return 0, err
	}
	
	_ = reader.Close()
	
	_ = os.RemoveAll(orgTag.Name())
	_ = os.RemoveAll(reader.Name())
	return total, nil
}

func init() {
	block, err := hosts.GetUrlsFromTxt("block.txt")
	if err != nil {
		log.Fatalln(err)
	}
	
	log.Println("init allow list ...")
	body, err := getBodyFromUrls(allowUrls)
	if err != nil {
		log.Fatalln(err)
	}
	hosts.Resolve(body, allowList)
	
	log.Println("init block list ...")
	body, err = getBodyFromUrls(block)
	if err != nil {
		log.Fatalln(err)
	}
	hosts.Resolve(body, blockList, allowList)
	
	log.Println("init direct list ...")
	body, err = getBodyFromUrls(directUrls)
	if err != nil {
		log.Fatalln(err)
	}
	hosts.Resolve(body, directList)
	hosts.AppendLocal(directList)
	
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
