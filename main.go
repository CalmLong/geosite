package main

import (
	"flag"
	"github.com/golang/protobuf/proto"
	"github.com/v2fly/v2ray-core/v4/app/router"
	"io/ioutil"
	"log"
	"time"
)

var otherPut bool

func command() {
	onlyDomain := flag.Bool("d", false, "only output domain: domains")
	onlyFull := flag.Bool("f", false, "only output full: domains")
	death := flag.Bool("D", false, "detect and remove invalid domain names")
	
	outCn := flag.Bool("cn", false, "if you use -F, output in the format")
	format := flag.String("F", "", "format output")
	
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
		go func() {
			for {
				time.Sleep(1 * time.Minute)
				log.Println("processing ...")
			}
		}()
		isDeathList(allowList, blockList, cnList)
		log.Printf("done. %.2fm", time.Now().Sub(t).Minutes())
	}
	
	if *outCn {
		otherPut = true
		load2Format(*format, cnList)
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
	// always output full format
	output(coverOnlyFull, localList)
	output(coverOnlyFull, allowList)
	
	command()
	

	if !otherPut {
		v2rayGeoSite()
		return
	}
}
