package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	success = 0
	noRet   = 3
)

type aa struct {
	Status int32 `json:"Status"`
}

func request(domain string) int {
	domain = fmt.Sprintf("https://dns.google/resolve?name=%s", domain)
	resp, err := http.Get(domain)
	if err != nil {
		return success
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var ret aa
	if err := json.Unmarshal(body, &ret); err != nil {
		return success
	}
	if ret.Status == 3 {
		return noRet
	}
	return success
}

func handle(originalMap map[string]struct{}, deathChan chan string) {
	var group sync.WaitGroup
	limit := make(chan struct{}, 1000)
	
	for uri := range originalMap {
		
		limit <- struct{}{}
		
		group.Add(1)
		
		go func(uri string, limit chan struct{}) {
			if uriStr, ok := removeSuffix(uri); ok {
				if request(uriStr) == noRet {
					deathChan <- uri
				}
			}
			group.Done()
			<-limit
		}(uri, limit)
	}
	
	group.Wait()
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

func readChan(deathChan chan string) map[string]struct{} {
	v := make(map[string]struct{})
	timer := time.NewTimer(2 * time.Second)
	for {
		select {
		case <-timer.C:
			close(deathChan)
			return v
		case uri := <-deathChan:
			v[uri] = struct{}{}
			timer.Reset(2 * time.Second)
		}
	}
}

func isDeath(originalMap map[string]struct{}) {
	rwCache(originalMap, false)
	
	deathChan := make(chan string, len(originalMap))
	handle(originalMap, deathChan)
	
	deathMap := readChan(deathChan)
	rwCache(deathMap, true)
}

func isDeathList(originalMaps ...map[string]struct{}) {
	for i := 0; i < len(originalMaps); i++ {
		isDeath(originalMaps[i])
	}
}
