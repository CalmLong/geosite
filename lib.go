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

var ptrList = map[string]struct{}{
	"127.in-addr.arpa": {},
	"1.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.ip6.arpa": {},
	"onion": {},
	"test": {},
	"invalid": {},
	"10.in-addr.arpa": {},
	"16.172.in-addr.arpa": {},
	"17.172.in-addr.arpa": {},
	"18.172.in-addr.arpa": {},
	"19.172.in-addr.arpa": {},
	"20.172.in-addr.arpa": {},
	"21.172.in-addr.arpa": {},
	"22.172.in-addr.arpa": {},
	"23.172.in-addr.arpa": {},
	"24.172.in-addr.arpa": {},
	"25.172.in-addr.arpa": {},
	"26.172.in-addr.arpa": {},
	"27.172.in-addr.arpa": {},
	"28.172.in-addr.arpa": {},
	"29.172.in-addr.arpa": {},
	"30.172.in-addr.arpa": {},
	"31.172.in-addr.arpa": {},
	"168.192.in-addr.arpa": {},
	"0.in-addr.arpa": {},
	"254.169.in-addr.arpa": {},
	"2.0.192.in-addr.arpa": {},
	"100.51.198.in-addr.arpa": {},
	"113.0.203.in-addr.arpa": {},
	"255.255.255.255.in-addr.arpa": {},
	"0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.ip6.arpa": {},
	"d.f.ip6.arpa": {},
	"8.e.f.ip6.arpa": {},
	"9.e.f.ip6.arpa": {},
	"a.e.f.ip6.arpa": {},
	"b.e.f.ip6.arpa": {},
	"8.b.d.0.1.0.0.2.ip6.arpa": {},
	"64.100.in-addr.arpa": {},
	"65.100.in-addr.arpa": {},
	"66.100.in-addr.arpa": {},
	"67.100.in-addr.arpa": {},
	"68.100.in-addr.arpa": {},
	"69.100.in-addr.arpa": {},
	"70.100.in-addr.arpa": {},
	"71.100.in-addr.arpa": {},
	"72.100.in-addr.arpa": {},
	"73.100.in-addr.arpa": {},
	"74.100.in-addr.arpa": {},
	"75.100.in-addr.arpa": {},
	"76.100.in-addr.arpa": {},
	"77.100.in-addr.arpa": {},
	"78.100.in-addr.arpa": {},
	"79.100.in-addr.arpa": {},
	"80.100.in-addr.arpa": {},
	"81.100.in-addr.arpa": {},
	"82.100.in-addr.arpa": {},
	"83.100.in-addr.arpa": {},
	"84.100.in-addr.arpa": {},
	"85.100.in-addr.arpa": {},
	"86.100.in-addr.arpa": {},
	"87.100.in-addr.arpa": {},
	"88.100.in-addr.arpa": {},
	"89.100.in-addr.arpa": {},
	"90.100.in-addr.arpa": {},
	"91.100.in-addr.arpa": {},
	"92.100.in-addr.arpa": {},
	"93.100.in-addr.arpa": {},
	"94.100.in-addr.arpa": {},
	"95.100.in-addr.arpa": {},
	"96.100.in-addr.arpa": {},
	"97.100.in-addr.arpa": {},
	"98.100.in-addr.arpa": {},
	"99.100.in-addr.arpa": {},
	"100.100.in-addr.arpa": {},
	"101.100.in-addr.arpa": {},
	"102.100.in-addr.arpa": {},
	"103.100.in-addr.arpa": {},
	"104.100.in-addr.arpa": {},
	"105.100.in-addr.arpa": {},
	"106.100.in-addr.arpa": {},
	"107.100.in-addr.arpa": {},
	"108.100.in-addr.arpa": {},
	"109.100.in-addr.arpa": {},
	"110.100.in-addr.arpa": {},
	"111.100.in-addr.arpa": {},
	"112.100.in-addr.arpa": {},
	"113.100.in-addr.arpa": {},
	"114.100.in-addr.arpa": {},
	"115.100.in-addr.arpa": {},
	"116.100.in-addr.arpa": {},
	"117.100.in-addr.arpa": {},
	"118.100.in-addr.arpa": {},
	"119.100.in-addr.arpa": {},
	"120.100.in-addr.arpa": {},
	"121.100.in-addr.arpa": {},
	"122.100.in-addr.arpa": {},
	"123.100.in-addr.arpa": {},
	"124.100.in-addr.arpa": {},
	"125.100.in-addr.arpa": {},
	"126.100.in-addr.arpa": {},
	"127.100.in-addr.arpa": {},
}

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
