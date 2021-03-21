package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

func removeSuffix(uri string) (string, bool) {
	for _, l := range localList {
		if strings.EqualFold(l, uri) {
			return "", false
		}
	}
	if strings.Contains(uri, "regexp:") {
		return "", false
	}
	if strings.Contains(uri, "keyword:") {
		return "", false
	}
	uri = strings.ReplaceAll(uri, suffixFull, "")
	uri = strings.ReplaceAll(uri, suffixDomain, "")
	return uri, true
}

func handle(originalMap map[string]struct{}, group *sync.WaitGroup, deathChan chan string) {
	var mutex sync.Mutex
	for uri := range originalMap {
		mutex.Lock()
		go func(uri string) {
			defer group.Done()
			if uriStr, ok := removeSuffix(uri); ok {
				ip, _ := net.LookupIP(uriStr)
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
	log.Println("reading death list...")

	timer := time.NewTimer(2 * time.Second)
	for {
		select {
		case <-timer.C:
			close(deathChan)
			log.Println("writing death list...")
			rwCache(deathMap, true)
			return
		case uri := <-deathChan:
			deathMap[uri] = struct{}{}
			delete(originalMap, uri)
			timer.Reset(2 * time.Second)
		}
	}
}
