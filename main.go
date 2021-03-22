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
	onlyDomain := flag.Bool("d", false, "only output domain: domains")
	onlyFull := flag.Bool("f", false, "only output full: domains")
	death := flag.Bool("D", false, "detect and remove invalid domain names")
	flag.Parse()
	if *onlyDomain {
		log.Printf("only output domain")
		output(coverOnlyDomain, blockList)
	}
	if *onlyFull {
		log.Printf("only output full")
		output(coverOnlyFull, blockList)
	}
	if *death {
		t := time.Now()
		log.Printf("clear invalid domain names ...")
		isDeathList(allowList, blockList, cnList)
		log.Printf("done. %.2fm", time.Now().Sub(t).Minutes())
	}
}

func main() {
	// allowList always output full format
	output(coverOnlyFull, allowList)
	
	command()
	
	log.Printf("creating ...")
	
	t := time.Now()
	
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
