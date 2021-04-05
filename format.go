package main

import (
	"net"
	"net/url"
	"strings"
)

const (
	j = '#'
	p = "/"
)

const (
	coverDefault    = 0
	coverOnlyDomain = 1
	coverOnlyFull   = 2
)

func replace(s, old string, new ...string) string {
	s = strings.TrimPrefix(s, old)
	if len(new) >= 2 {
		return new[0] + s + new[1]
	}
	return s
}

func removeSuffix(uri string) (string, bool) {
	if strings.Contains(uri, "regexp:") {
		return strings.ReplaceAll(uri, "regexp:", ""), false
	}
	if strings.Contains(uri, "keyword:") {
		return strings.ReplaceAll(uri, "keyword:", ""), false
	}
	uri = strings.ReplaceAll(uri, suffixFull, "")
	uri = strings.ReplaceAll(uri, suffixDomain, "")
	return uri, true
}

func cover(uri string, action int) (string, bool) {
	if err := net.ParseIP(uri); err != nil {
		return "", false
	}
	uri, _ = removeSuffix(uri)
	switch action {
	case coverOnlyDomain:
		return "domain:" + uri, true
	case coverOnlyFull:
		return "full:" + uri, true
	case coverDefault:
		switch strings.Count(uri, ".") {
		case 1:
			return "domain:" + uri, true
		case 2:
			for suffix := range suffixList {
				if strings.Contains(uri, suffix) {
					return "domain:" + uri, true
				}
			}
			fallthrough
		default:
			return "full:" + uri, true
		}
	default:
		return "", false
	}
}

func format(newOrg string, prefix ...string) string {
	for _, s := range prefix {
		newOrg = strings.ReplaceAll(newOrg, s, "")
	}
	return newOrg
}

func parseUrl(raw string) (string, bool) {
	raw = strings.ReplaceAll(raw, "http://", "")
	raw = strings.ReplaceAll(raw, "https://", "")
	raw = strings.ReplaceAll(raw, "ftp://", "")
	raw = strings.ReplaceAll(raw, "websocket://", "")
	
	switch strings.Count(raw, "/") {
	case 0:
		return raw, true
	case 1:
		return raw[:len(raw)-1], true
	default:
		return "", false
	}
}

func Resolve(src map[string]struct{}, dst map[string]struct{}) {
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
		if uri, ok := cover(urlString, coverDefault); ok {
			dst[uri] = struct{}{}
		}
	}
}
