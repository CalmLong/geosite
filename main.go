package main

import (
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"log"
	"os"
	"time"
	"v2ray.com/core/app/router"
)

var dir = pwd(geoSitePath)

func main() {
	log.Printf("creating ...")
	
	t := time.Now()
	
	total, err := getSites(dir, blockTag, v2flyBlockTag)
	if err != nil {
		log.Fatalln("Failed: ", err)
	}
	log.Printf("block sties: %d", total)
	
	total, err = getSites(dir, directTag, v2flyDirectTag)
	if err != nil {
		log.Fatalln("Failed: ", err)
	}
	log.Printf("direct sties: %d", total)
	
	protoList := new(router.GeoSiteList)
	if err := readFiles(dir, protoList); err != nil {
		log.Fatalf("protoList err: %s", err.Error())
	}
	protoBytes, err := proto.Marshal(protoList)
	if err != nil {
		log.Fatalln("Failed: ", err)
	}
	if err := ioutil.WriteFile("geosite.dat", protoBytes, 0644); err != nil {
		log.Fatalln("Failed: ", err)
	}
	
	log.Println("deleted tmp files")
	_ = os.RemoveAll(dir)
	_ = os.RemoveAll(pwd() + "master.zip")
	_ = os.RemoveAll(pwd() + v2flySitePath)
	
	log.Printf("created. %ds", int64(time.Now().Sub(t).Seconds()))
}
