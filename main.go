package main

import (
	"flag"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/v2fly/v2ray-core/v4/app/router"
	"io/ioutil"
	"log"
	"time"
)

func command() {
	death := flag.Bool("D", false, "detect and remove invalid domain names")
	format := flag.String("F", "v2ray", "")
	flag.Parse()
	if *death {
		t := time.Now()
		log.Printf("clear invalid domain names ...")
		go func() {
			for {
				time.Sleep(1 * time.Minute)
				log.Println("processing ...")
			}
		}()
		isDeathList(blockList)
		log.Printf("done. %.2fm", time.Now().Sub(t).Minutes())
	}
	switch *format {
	case "domain":
		domainGeoSite()
	case "agh":
		aghGeoSite()
	case "v2ray":
		v2rayGeoSite()
	default:
		v2rayGeoSite()
	}
}

func aghGeoSite() {
	full, domain, _ := getDomain(cnList)
	f := make([]string, len(full))
	for i, s := range full {
		f[i] = fmt.Sprintf("[/%s/]223.5.5.5", trimDomain(s))
	}
	d := make([]string, len(domain))
	for i, s := range domain {
		d[i] = fmt.Sprintf("[/%s/]223.5.5.5", trimDomain(s))
	}
	if err := writer2File("agh-cn.txt", f, d); err != nil {
		log.Fatalln(err)
	}
}

func domainGeoSite() {
	full, domain, _ := getDomain(blockList)
	if err := writer2File("domain-block.txt", full, domain); err != nil {
		log.Fatalln(err)
	}
	full, domain, _ = getDomain(cnList)
	if err := writer2File("domain-cn.txt", full, domain); err != nil {
		log.Fatalln(err)
	}
	full, domain, _ = getDomain(proxyList)
	if err := writer2File("domain-proxy.txt", full, domain); err != nil {
		log.Fatalln(err)
	}
}

func v2rayGeoSite() {
	t := time.Now()
	log.Printf("creating ...")
	ref := loadEntry()
	protoList := new(router.GeoSiteList)
	for _, list := range ref {
		pl, err := ParseList(list, ref)
		if err != nil {
			log.Fatalln("Failed: ", err)
		}
		site, err := pl.toProto()
		if err != nil {
			log.Fatalln("Failed: ", err)
		}
		protoList.Entry = append(protoList.Entry, site)
	}
	protoBytes, err := proto.Marshal(protoList)
	if err != nil {
		log.Fatalln("Failed: ", err)
	}
	if err := ioutil.WriteFile("geosite.dat", protoBytes, 0644); err != nil {
		log.Fatalln("Failed: ", err)
	} else {
		log.Println("geosite.dat has been generated successfully.")
	}
	log.Printf("created. %ds", int64(time.Now().Sub(t).Seconds()))
}

func main() {
	command()
}
