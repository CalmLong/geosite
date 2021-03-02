package main

import "path/filepath"

const (
	v2flyNoneTag   = ""
	v2flyBlockTag  = "CATEGORY-ADS-ALL"
	v2flyDirectTag = "CN"
	allowTag       = "allow"
	adsTag         = "ads"
	cnTag          = "cn"
	suffixFull     = "full:"
	suffixDomain   = "domain:"
)

const (
	v2flySites    = "https://github.com/v2fly/domain-list-community/archive/master.zip"
	v2flySitePath = "domain-list-community-master"
)

const geoSitePath = "geodata"

var v2flySitePathData = filepath.Join("domain-list-community-master", "data")

var allowUrls = []string{
	"https://raw.githubusercontent.com/CalmLong/allow-list/master/allow.txt",
	"https://raw.githubusercontent.com/privacy-protection-tools/dead-horse/master/anti-ad-white-list.txt",
	"https://raw.githubusercontent.com/anudeepND/whitelist/master/domains/whitelist.txt",
	"https://raw.githubusercontent.com/notracking/hosts-blocklists-scripts/master/hostnames.whitelist.txt",
}

var directUrls = []string{
	"https://raw.githubusercontent.com/felixonmars/dnsmasq-china-list/master/accelerated-domains.china.conf",
	"https://raw.githubusercontent.com/felixonmars/dnsmasq-china-list/master/apple.china.conf",
	"https://raw.githubusercontent.com/felixonmars/dnsmasq-china-list/master/google.china.conf",
}

var allowList = make(map[string]struct{}, 0)
var blockList = make(map[string]struct{}, 0)
var directList = make(map[string]struct{}, 0)

var localList = []string{
	"localhost",
	"ip6-localhost",
	"localhost.localdomain",
	"local",
	"broadcasthost",
	"ip6-loopback",
	"ip6-localnet",
	"ip6-mcastprefix",
	"ip6-allnodes",
	"ip6-allrouters",
	"ip6-allhosts",
}
