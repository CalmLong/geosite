package main

import (
	"net"
	"net/url"
	"strings"

	"github.com/v2fly/v2ray-core/v4/app/router"
)

const (
	j = '#'
	p = "/"
)

func domainType(uri string) router.Domain_Type {
	if strings.HasPrefix(uri, "domain:") {
		return router.Domain_Domain
	}
	if strings.HasPrefix(uri, "regexp:") {
		return router.Domain_Regex
	}
	if strings.HasPrefix(uri, "full:") {
		return router.Domain_Full
	}
	return router.Domain_Plain
}

func ResolveV2Ray(src map[string]struct{}, dst map[string]router.Domain_Type) {
	for k := range src {
		if strings.Contains(k, "#") {
			continue
		}
		// domain:github.com:@noCN
		og := strings.Split(k, ":")
		
		switch len(og) {
		case 2:
			dst[og[1]] = domainType(k)
		case 3:
			v := og[1] // github.com
			t := domainType(k)
			switch og[2] {
			case "@cn":
				cnList[v] = t
			case "@ads":
				blockList[v] = t
			default:
				dst[v] = t
			}
		default:
			continue
		}
	}
}

func Resolve(src map[string]struct{}, dst map[string]router.Domain_Type, rt router.Domain_Type) {
	for k := range src {
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
			dst[newOrg] = router.Domain_Regex
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

		dst[urlString] = rt
	}

	src = nil
}
