package main

import (
	"net"
	"strings"
	"sync"
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
	for uri := range originalMap {
		go func(uri string) {
			defer group.Done()
			if uriStr, ok := removeSuffix(uri); ok {
				for i := 0; i < 2; i++ {
					ip, _ := net.LookupIP(uriStr)
					if len(ip) > 0 {
						return
					}
				}
				deathChan <- uri
			}
		}(uri)
	}
}

func isDeath(originalMap map[string]struct{}) {
	var group sync.WaitGroup
	group.Add(len(originalMap))
	
	deathChan := make(chan string, 1024)
	
	go handle(originalMap, &group, deathChan)
	
	go func() {
		for uri := range deathChan {
			delete(originalMap, uri)
		}
	}()
	
	group.Wait()
	close(deathChan)
}
