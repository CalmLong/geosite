package main

import (
	"bufio"
	"context"
	"github.com/v2fly/v2ray-core/v4/app/dns"
	dns_feature "github.com/v2fly/v2ray-core/v4/features/dns"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

func handle(originalMap map[string]struct{}, group *sync.WaitGroup, deathChan chan string) {
	dnsLookup := dns.NewLocalNameServer()
	var mutex sync.Mutex
	for uri := range originalMap {
		mutex.Lock()
		go func(uri string) {
			defer group.Done()
			if uriStr, ok := removeSuffix(uri); ok {
				ip, _ := dnsLookup.QueryIP(context.Background(), uriStr, net.IPv4zero, dns_feature.IPOption{
					IPv4Enable: true,
					IPv6Enable: true,
				}, true)
				if len(ip) > 0 {
					return
				}
				deathChan <- uri
			}
		}(uri)
		mutex.Unlock()
	}
}

func rwCache(valueMap map[string]struct{}, write bool) {
	if len(valueMap) == 0 {
		return
	}
	const cacheName = "geoSiteDeathCacheList"
	
	fi, err := os.OpenFile(cacheName, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}
	
	if write {
		for k := range valueMap {
			_, err := fi.WriteString(k + "\n")
			if err != nil {
				log.Fatalln(err)
			}
		}
	} else {
		reader := bufio.NewReader(fi)
		for {
			d, _, e := reader.ReadLine()
			if e == io.EOF {
				break
			}
			delete(valueMap, string(d))
		}
	}
	
	_ = fi.Close()
}

func isDeath(originalMap map[string]struct{}) {
	rwCache(originalMap, false)
	
	nowNum := len(originalMap)
	
	var group sync.WaitGroup
	group.Add(nowNum)
	
	deathChan := make(chan string, nowNum)
	
	go handle(originalMap, &group, deathChan)
	
	group.Wait()
	
	deathMap := make(map[string]struct{})
	
	timer := time.NewTimer(2 * time.Second)
	for {
		select {
		case <-timer.C:
			close(deathChan)
			rwCache(deathMap, true)
			return
		case uri := <-deathChan:
			deathMap[uri] = struct{}{}
			delete(originalMap, uri)
			timer.Reset(2 * time.Second)
		}
	}
}

func isDeathList(originalMaps ...map[string]struct{}) {
	for i := 0; i < len(originalMaps); i++ {
		isDeath(originalMaps[i])
	}
}
