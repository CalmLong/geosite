package main

import (
	"github.com/v2fly/v2ray-core/v4/app/router"
	"time"
)

const (
	adsTag   = "category-ads-all"
	cnTag    = "cn"
	proxyTag = "geolocation-!cn"
)

var allowUrls = []string{
	"https://raw.githubusercontent.com/CalmLong/allow-list/master/allow.txt",
	"https://raw.githubusercontent.com/privacy-protection-tools/dead-horse/master/anti-ad-white-list.txt",
}

var directUrls = []string{
	"https://raw.githubusercontent.com/felixonmars/dnsmasq-china-list/master/accelerated-domains.china.conf",
	"https://raw.githubusercontent.com/felixonmars/dnsmasq-china-list/master/apple.china.conf",
	"https://raw.githubusercontent.com/felixonmars/dnsmasq-china-list/master/google.china.conf",
}

const (
	v2rayNotCn    = "https://raw.githubusercontent.com/v2fly/domain-list-community/release/geolocation-!cn.txt"
	suffixListRaw = "https://raw.githubusercontent.com/publicsuffix/list/master/public_suffix_list.dat"
)

type dT struct {
	Value  string
	Format string
	Keep   bool
	Type   router.Domain_Type
}

var (
	allowList  = make(map[string]dT)
	blockList  = make(map[string]dT)
	cnList     = make(map[string]dT)
	proxyList  = make(map[string]dT)
	suffixList = make(map[string]struct{})
)

var nowTime string

func init() {
	nowTime = "# TIME: " + time.Now().In(time.FixedZone("CST", 8*3600)).Format("2006-01-02 15:04:05")
}
