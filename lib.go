package main

import "path/filepath"

const (
	v2flyBlockTag  = "CATEGORY-ADS-ALL"
	v2flyDirectTag = "CN"
	adsTag         = "ads"
	cnTag          = "cn"
	suffixFull     = "full:"
	suffixDomain   = "domain:"
)

const (
	v2flySites    = "https://github.com/v2fly/domain-list-community/archive/master.zip"
	v2flySitePath = "domain-list-community-master"
	
	allowSites1 = "https://raw.githubusercontent.com/CalmLong/allow-list/master/allow.txt"
	allowSites2 = "https://raw.githubusercontent.com/privacy-protection-tools/dead-horse/master/anti-ad-white-list.txt"
	allowSites3 = "https://raw.githubusercontent.com/neodevpro/neodevhost/master/allow"
	allowSites4 = "https://raw.githubusercontent.com/anudeepND/whitelist/master/domains/whitelist.txt"
	
	directSite1 = "https://raw.githubusercontent.com/felixonmars/dnsmasq-china-list/master/accelerated-domains.china.conf"
	directSite2 = "https://raw.githubusercontent.com/felixonmars/dnsmasq-china-list/master/apple.china.conf"
	directSite3 = "https://raw.githubusercontent.com/felixonmars/dnsmasq-china-list/master/google.china.conf"
)

const geoSitePath = "geodata"

var v2flySitePathData = filepath.Join("domain-list-community-master", "data")

var allowUrls = []string{allowSites1, allowSites2, allowSites3, allowSites4}
var directUrls = []string{directSite1, directSite2, directSite3}

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
