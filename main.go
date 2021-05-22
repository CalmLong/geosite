package main

import (
	"flag"
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
	full, domain, _ := getDomain(cnList, false)
	fd := formatDomain("[/%s/]223.5.5.5", full, domain)
	upstream := []string{"tls://8.8.8.8", "tls://8.8.4.4"}
	if err := writer2File(upstream, "agh-cn.txt", fd); err != nil {
		log.Fatalln(err)
	}
}

func domainGeoSite() {
	full, domain, _ := getDomain(blockList, false)
	if err := writer2File(nil, "domain-block.txt", full, domain); err != nil {
		log.Fatalln(err)
	}
	full, domain, _ = getDomain(cnList, false)
	if err := writer2File(nil, "domain-cn.txt", full, domain); err != nil {
		log.Fatalln(err)
	}
	full, domain, _ = getDomain(proxyList, false)
	if err := writer2File(nil, "domain-proxy.txt", full, domain); err != nil {
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
