package main

import (
	"bufio"
	"github.com/AdguardTeam/dnsproxy/upstream"
	"github.com/miekg/dns"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

func createHostTestMessage(host string) *dns.Msg {
	req := dns.Msg{}
	req.Id = dns.Id()
	req.RecursionDesired = true
	name := host + "."
	req.Question = []dns.Question{
		{Name: name, Qtype: dns.TypeA, Qclass: dns.ClassINET},
	}
	return &req
}

func handle(originalMap map[string]struct{}, deathChan chan string) {
	address := "tls://8.8.8.8:853"
	doq, err := upstream.AddressToUpstream(address, upstream.Options{})
	
	if err != nil {
		panic(err)
	}
	
	var group sync.WaitGroup
	limit := make(chan struct{}, 1000)
	
	for uri := range originalMap {
		
		limit <- struct{}{}
		
		group.Add(1)
		
		go func(uri string, limit chan struct{}) {
			if uriStr, ok := removeSuffix(uri); ok {
				
				msg, err := doq.Exchange(createHostTestMessage(uriStr))
				if err == nil {
					if msg.Answer == nil || len(msg.Answer) <= 0 {
						deathChan <- uri
					} else {
						if a, ok := msg.Answer[0].(*dns.A); ok {
							if net.IPv4(0, 0, 0, 0).Equal(a.A) {
								deathChan <- uri
							}
						}
					}
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
	
	deathFirst := make(chan string, len(originalMap))
	
	handle(originalMap, deathFirst)
	
	deathFirstMap := readChan(deathFirst)
	
	deathSecond := make(chan string, len(deathFirstMap))
	
	handle(deathFirstMap, deathSecond)
	
	deathMap := readChan(deathSecond)
	
	rwCache(deathMap, true)
}

func isDeathList(originalMaps ...map[string]struct{}) {
	for i := 0; i < len(originalMaps); i++ {
		isDeath(originalMaps[i])
	}
}
