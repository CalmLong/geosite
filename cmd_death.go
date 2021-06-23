package main

import (
	"bufio"
	"encoding/json"
	"github.com/v2fly/v2ray-core/v4/app/router"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	NOERROR  = 0
	NXDOMAIN = 3
)

type Answer struct {
	Status int32 `json:"Status"`
}

func doRequest(uri string) int {
	uri = "https://dns.google/resolve?name=" + uri
	resp, err := http.Get(uri)
	if err != nil {
		return NOERROR
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var ret Answer
	if err := json.Unmarshal(body, &ret); err != nil {
		return NOERROR
	}
	if ret.Status == 3 {
		return NXDOMAIN
	}
	return NOERROR
}

func handle(valueMap map[string]dT, deathChan chan dT, wait chan int) {
	limit := make(chan struct{}, 1000)
	
	for _, ov := range valueMap {
		
		limit <- struct{}{}
		
		go func(uri dT, limit chan struct{}) {
			if uri.Type == router.Domain_Full {
				if doRequest(uri.Value) == NXDOMAIN {
					deathChan <- uri
				}
			}
			<-limit
		}(ov, limit)
	}
	
	wait <- 1
}

func rwCache(valueMap map[string]dT, write bool) {
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

func readChan(deathChan chan dT) map[string]dT {
	v := make(map[string]dT)
	timer := time.NewTimer(2 * time.Second)
	for {
		select {
		case <-timer.C:
			close(deathChan)
			return v
		case uri := <-deathChan:
			v[uri.Value] = uri
			timer.Reset(2 * time.Second)
		}
	}
}

func isDeath(valueMap map[string]dT) {
	rwCache(valueMap, false)
	
	wait := make(chan int, 1)
	
	deathChan := make(chan dT, len(valueMap))
	handle(valueMap, deathChan, wait)
	
	<-wait
	
	deathMap := readChan(deathChan)
	rwCache(deathMap, true)
}

func isDeathList(valueMaps ...map[string]dT) {
	for i := 0; i < len(valueMaps); i++ {
		isDeath(valueMaps[i])
	}
}
