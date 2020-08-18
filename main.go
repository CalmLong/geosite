package main

import (
	"fmt"
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
		fmt.Println("Failed: ", err)
		return
	}
	log.Printf("block sties: %d", total)
	
	total, err = getSites(dir, directTag, v2flyDirectTag)
	if err != nil {
		fmt.Println("Failed: ", err)
		return
	}
	log.Printf("direct sties: %d", total)
	
	protoList := new(router.GeoSiteList)
	if err := readFiles(dir, protoList); err != nil {
		log.Printf("protoList err: %s", err.Error())
		return
	}
	protoBytes, err := proto.Marshal(protoList)
	if err != nil {
		fmt.Println("Failed:", err)
		return
	}
	if err := ioutil.WriteFile("geosite.dat", protoBytes, 0644); err != nil {
		fmt.Println("Failed: ", err)
	}
	
	log.Println("deleted tmp files")
	_ = os.RemoveAll(dir)
	_ = os.RemoveAll(pwd() + "master.zip")
	_ = os.RemoveAll(pwd() + v2flySitePath)
	
	log.Printf("created. %ds", int64(time.Now().Sub(t).Seconds()))
}
