package main

import (
	"github.com/v2fly/v2ray-core/v4/app/router"
	"net"
	"net/url"
	"strings"
)

const (
	j = '#'
	p = "/"
)

func domainType(uri string) router.Domain_Type {
	if strings.Contains(uri, "domain:") {
		return router.Domain_Domain
	}
	if strings.Contains(uri, "regexp:") {
		return router.Domain_Regex
	}
	if strings.Contains(uri, "full:") {
		return router.Domain_Full
	}
	switch strings.Count(uri, ".") {
	case 1:
		return router.Domain_Domain
	case 2:
		for suffix := range suffixList {
			if strings.Contains(uri, suffix) {
				return router.Domain_Domain
			}
		}
		fallthrough
	default:
		return router.Domain_Full
	}
}

func format(newOrg string, prefix ...string) string {
	for _, s := range prefix {
		newOrg = strings.ReplaceAll(newOrg, s, "")
	}
	return newOrg
}

func parseUrl(raw string) (string, bool) {
	raw = strings.TrimSuffix(raw, "http://")
	raw = strings.TrimSuffix(raw, "https://")
	raw = strings.TrimSuffix(raw, "ftp://")
	raw = strings.TrimSuffix(raw, "ws://")
	raw = strings.TrimSuffix(raw, "wss://")
	
	switch strings.Count(raw, "/") {
	case 0:
		return raw, true
	case 1:
		return raw[:len(raw)-1], true
	default:
		return "", false
	}
}

func ResolveV2Ray(src map[string]struct{}, dst map[string]dT) {
	for k := range src {
		if strings.Contains(k, "#") {
			continue
		}
		ks := strings.Split(k, ":")
		dt := dT{
			Value:  ks[1],
			Format: k,
			Keep:   true,
			Type:   domainType(k),
		}
		switch len(ks) {
		case 2:
			dst[ks[1]] = dt
		case 3:
			k = ks[1]
			dt.Value = k
			dt.Format = ks[0] + ":" + ks[1]
			dt.Type = domainType(dt.Format)
			switch ks[2] {
			case "@cn":
				cnList[k] = dt
			case "@ads":
				blockList[k] = dt
			default:
				dst[k] = dt
			}
		default:
			continue
		}
	}
}

func Resolve(src map[string]struct{}, dst map[string]dT) {
	for k := range src {
		original := k
		// 第一个字符为 # 或 ! 时跳过
		if strings.IndexRune(original, j) == 0 || strings.IndexRune(original, '!') == 0 {
			continue
		}
		// 为空行时跳过
		if strings.TrimSpace(original) == "" {
			continue
		}
		// 中间包含特殊空格的
		if strings.ContainsRune(original, '\t') {
			original = strings.ReplaceAll(original, "\t", " ")
		}
		
		newOrg := strings.ToLower(original)
		
		// 移除前缀为 0.0.0.0 或者 127.0.0.1 (移除第一个空格前的内容)
		index := strings.IndexRune(newOrg, ' ')
		if index > -1 {
			newOrg = strings.ReplaceAll(newOrg, newOrg[:index], "")
		}
		
		// V2Ray
		newOrg = format(newOrg, "domain:", "full:", "regexp:", "keyword:", ":@ads")
		
		// 移除行中的空格
		newOrg = strings.TrimSpace(newOrg)
		// 再一次验证第一个字符为 # 时跳过
		if strings.IndexRune(original, j) == 0 {
			continue
		}
		if strings.ContainsRune(newOrg, j) {
			newOrg = newOrg[:strings.IndexRune(newOrg, j)]
		}
		// dnsmasq-list
		if d := strings.Count(newOrg, p); d == 2 {
			newOrg = newOrg[strings.Index(newOrg, p)+1:]
			newOrg = newOrg[:strings.Index(newOrg, p)]
		}
		
		newOrg = strings.TrimSpace(newOrg)
		// 检测是否有端口号，有则移除端口号
		if i := strings.IndexRune(newOrg, ':'); i != -1 {
			newOrg = newOrg[:i]
		}
		
		// 包含正则符号的
		if strings.ContainsAny(newOrg, "$()*+[?\\^{|") {
			continue
		}
		
		if v, ok := parseUrl(newOrg); ok {
			newOrg = v
		} else {
			continue
		}
		
		urlStr, err := url.Parse(newOrg)
		if err != nil {
			continue
		}
		urlString := urlStr.String()
		// 如果为 IP 则跳过
		if ip := net.ParseIP(urlString); ip != nil {
			continue
		}
		if strings.IndexRune(urlString, '.') == 0 {
			urlString = urlString[1:]
		}
		if urlString == "" {
			continue
		}
		
		dst[urlString] = dT{
			Value: urlString,
			Type:  domainType(urlString),
		}
	}
}
