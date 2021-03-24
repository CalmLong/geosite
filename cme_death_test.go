package main

import (
	"bufio"
	"fmt"
	"github.com/AdguardTeam/dnsproxy/upstream"
	"github.com/miekg/dns"
	"io"
	"net"
	"os"
	"sync"
	"testing"
)

func TestLoad(t *testing.T) {
	address := "tls://8.8.8.8:853"
	u, err := upstream.AddressToUpstream(address, upstream.Options{})
	
	originalMap := make(map[string]struct{})
	
	fi, err := os.Open("geoSiteDeathCacheList")
	if err != nil {
		t.Error(err)
		return
	}
	
	buff := bufio.NewReader(fi)
	
	for {
		s, _, e := buff.ReadLine()
		if e == io.EOF {
			break
		}
		originalMap[string(s)] = struct{}{}
	}
	
	var group sync.WaitGroup
	limit := make(chan struct{}, 1000)
	
	for uri := range originalMap {
		
		limit <- struct{}{}
		
		group.Add(1)
		
		go func(uri string, limit chan struct{}) {
			if uriStr, ok := removeSuffix(uri); ok {
				msg, err := u.Exchange(createHostTestMessage(uriStr))
				if err == nil && len(msg.Answer) > 0 {
					if a, ok := msg.Answer[0].(*dns.A); ok {
						if !net.IPv4(0, 0, 0, 0).Equal(a.A) {
							fmt.Println(uriStr)
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
