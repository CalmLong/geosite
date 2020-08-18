package main

const (
	v2flyBlockTag  = "CATEGORY-ADS-ALL"
	v2flyDirectTag = "CN"
	blockTag       = "ads"
	directTag      = "cn"
	suffixFull     = "full:"
	suffixDomain   = "domain:"
)

const (
	v2flySites    = "https://github.com/v2fly/domain-list-community/archive/master.zip"
	v2flySitePath = "domain-list-community-master"
	
	v2flySitePathData = "domain-list-community-master/data"
	
	allowSites1 = "https://raw.githubusercontent.com/CalmLong/whitelist/master/white.txt"
	allowSites2 = "https://raw.githubusercontent.com/privacy-protection-tools/dead-horse/master/anti-ad-white-list.txt"
	
	directSite = "https://raw.githubusercontent.com/felixonmars/dnsmasq-china-list/master/accelerated-domains.china.conf"
)

const geoSitePath = "geodata"

var allowUrls = []string{allowSites1, allowSites2}
var directUrls = []string{directSite}

var allowList = make(map[string]struct{}, 0)
var blockList = make(map[string]struct{}, 0)
var directList = make(map[string]struct{}, 0)
