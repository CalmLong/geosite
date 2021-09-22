package main

import (
	"fmt"
	"net"
	"net/url"
	"strings"
	
	"github.com/v2fly/v2ray-core/v4/app/router"
)

const (
	j = '#'
	p = "/"
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

func Resolve(src map[string]struct{}) []string {
	name := make([]string, 0)
	for k := range src {
		if HasPattern(k) {
			if strings.Contains(k, "@") {
				continue
			}
			
			name = append(name, k)
			continue
		}
		
		k = strings.TrimSpace(k)
		if k == "" {
			continue
		}
		
		if strings.IndexRune(k, j) == 0 || strings.IndexRune(k, '!') == 0 {
			continue
		}
		
		if strings.ContainsRune(k, '\t') {
			k = strings.ReplaceAll(k, "\t", " ")
		}
		
		newOrg := strings.ToLower(k)
		
		if idx := strings.IndexRune(newOrg, ' '); idx > -1 {
			newOrg = strings.ReplaceAll(newOrg, newOrg[:idx], "")
		}
		
		newOrg = strings.TrimSpace(newOrg)
		if strings.IndexRune(k, j) == 0 {
			continue
		}
		
		if strings.ContainsRune(newOrg, j) {
			newOrg = newOrg[:strings.IndexRune(newOrg, j)]
		}
		if d := strings.Count(newOrg, p); d == 2 {
			newOrg = newOrg[strings.Index(newOrg, p)+1:]
			newOrg = newOrg[:strings.Index(newOrg, p)]
		}
		
		newOrg = strings.TrimSpace(newOrg)
		if i := strings.IndexRune(newOrg, ':'); i != -1 {
			newOrg = newOrg[:i]
		}
		
		if strings.ContainsAny(newOrg, "$()*+[?\\^{|") {
			name = append(name, fmt.Sprintf("%s:%s", router.Domain_Regex, k))
			continue
		}
		
		urlStr, err := url.Parse(newOrg)
		if err != nil {
			continue
		}
		urlString := urlStr.String()
		if ip := net.ParseIP(urlString); ip != nil {
			continue
		}
		if strings.IndexRune(urlString, '.') == 0 {
			urlString = urlString[1:]
		}
		if urlString == "" {
			continue
		}
		
		n, t := Pattern(urlString)
		name = append(name, fmt.Sprintf("%s:%s", t, n))
	}
	
	return name
}
