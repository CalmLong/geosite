package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"strings"
	
	"github.com/v2fly/v2ray-core/v4/app/router"
)

const (
	adsTag   = "category-ads-all"
	cnTag    = "cn"
	proxyTag = "geolocation-!cn"
)

const (
	adList    = "https://raw.githubusercontent.com/carrnot/adblock-list/release/domain.txt"
	chinaList = "https://raw.githubusercontent.com/carrnot/china-domain-list/release/domain.txt"
	proxyList = "https://raw.githubusercontent.com/v2fly/domain-list-community/release/geolocation-!cn.txt"
)

const (
	_Full    = "full:"
	_Domain  = "domain:"
	_Keyword = "keyword:"
	_Regexp  = "regexp:"
)

const (
	_LenFull    = len(_Full)
	_LenDomain  = len(_Domain)
	_LenKeyword = len(_Keyword)
	_LenRegexp  = len(_Regexp)
)

func Pattern(s string) (string, router.Domain_Type) {
	if i := strings.Index(s, "@"); i != -1 {
		s = s[:i-1]
	}
	if strings.HasPrefix(s, _Full) {
		return s[_LenFull:], router.Domain_Full
	}
	if strings.HasPrefix(s, _Domain) {
		return s[_LenDomain:], router.Domain_Domain
	}
	if strings.HasPrefix(s, _Keyword) {
		return s[_LenKeyword:], router.Domain_Plain
	}
	if strings.HasPrefix(s, _Regexp) {
		return s[_LenRegexp:], router.Domain_Regex
	}
	return s, router.Domain_Full
}

func HasPattern(s string) bool {
	if strings.HasPrefix(s, _Full) {
		return true
	}
	if strings.HasPrefix(s, _Domain) {
		return true
	}
	if strings.HasPrefix(s, _Keyword) {
		return true
	}
	if strings.HasPrefix(s, _Regexp) {
		return true
	}
	return false
}

func Resolve(src []string) []string {
	name := make([]string, 0, len(src))
	for i := 0; i < len(src); i++ {
		if !HasPattern(src[i]) || strings.Contains(src[i], "@") {
			continue
		}
		n, t := Pattern(src[i])
		name = append(name, fmt.Sprintf("%s:%s", t, n))
	}
	return name
}

func getDomain(uri ...string) []string {
	tmp := make(map[string]struct{})
	for _, u := range uri {
		log.Println(u)
		resp, err := http.Get(u)
		if err != nil {
			log.Panic(err)
		}
		sc := bufio.NewScanner(resp.Body)
		for sc.Scan() {
			tmp[sc.Text()] = struct{}{}
		}
		_ = resp.Body.Close()
	}
	domain := make([]string, 0, len(tmp))
	for d := range tmp {
		domain = append(domain, d)
	}
	return domain
}

func loadEntry() map[string]*List {
	ref := make(map[string]*List)
	ref[adsTag] = getEntry(adsTag, Resolve(getDomain(adList)))
	ref[cnTag] = getEntry(cnTag, Resolve(getDomain(chinaList)))
	ref[proxyTag] = getEntry(proxyTag, Resolve(getDomain(proxyList)))
	return ref
}

func getEntry(name string, domain []string) *List {
	list := &List{
		Name: name,
	}
	
	for i := 0; i < len(domain); i++ {
		entry, err := parseEntry(domain[i])
		if err != nil {
			log.Fatalln(err)
		}
		list.Entry = append(list.Entry, entry)
	}
	
	return list
}
