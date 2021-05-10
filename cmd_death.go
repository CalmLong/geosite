package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/v2fly/v2ray-core/v4/app/router"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	success  = 0
	nxdomain = 3
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
		return nxdomain
	}
	return success
}

func handle(originalMap map[string]dT, deathChan chan dT) {
	var group sync.WaitGroup
	limit := make(chan struct{}, 1000)

	for _, ov := range originalMap {

		limit <- struct{}{}

		group.Add(1)

		go func(uri dT, limit chan struct{}) {
			if uri.Type == router.Domain_Full {
				if request(uri.Value) == nxdomain {
					deathChan <- uri
				}
			}
			group.Done()
			<-limit
		}(ov, limit)
	}

	group.Wait()
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

func isDeath(originalMap map[string]dT) {
	rwCache(originalMap, false)

	deathChan := make(chan dT, len(originalMap))
	handle(originalMap, deathChan)

	deathMap := readChan(deathChan)
	rwCache(deathMap, true)
}

func isDeathList(originalMaps ...map[string]dT) {
	for i := 0; i < len(originalMaps); i++ {
		isDeath(originalMaps[i])
	}
}
