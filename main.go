package main

import (
	"github.com/golang/protobuf/proto"
	"github.com/v2fly/v2ray-core/v4/app/router"
	"io/ioutil"
	"log"
	"time"
)

func main() {
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
