package main

import (
	"github.com/v2fly/v2ray-core/v4/app/router"
)

const (
	adsTag   = "category-ads-all"
	cnTag    = "cn"
	proxyTag = "geolocation-!cn"
)

var directUrls = []string{
	"https://raw.githubusercontent.com/felixonmars/dnsmasq-china-list/master/accelerated-domains.china.conf",
	"https://raw.githubusercontent.com/felixonmars/dnsmasq-china-list/master/apple.china.conf",
	"https://raw.githubusercontent.com/felixonmars/dnsmasq-china-list/master/google.china.conf",
}

const (
	v2rayNotCn    = "https://raw.githubusercontent.com/v2fly/domain-list-community/release/geolocation-!cn.txt"
)

var (
	blockList  = make(map[string]router.Domain_Type)
	cnList     = make(map[string]router.Domain_Type)
	proxyList  = make(map[string]router.Domain_Type)
)