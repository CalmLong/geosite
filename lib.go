package main

const (
	localTag     = "local"
	allowTag     = "allow"
	blockTag     = "category-ads-all"
	cnTag        = "cn"
	proxyTag     = "proxy"
	suffixFull   = "full:"
	suffixDomain = "domain:"
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

var blockUrlsForV2Ray = []string{
	"https://raw.githubusercontent.com/CalmLong/domain-list/master/block.txt",
	"https://raw.githubusercontent.com/v2fly/domain-list-community/release/category-ads-all.txt",
}

const (
	domainListCnRaw     = "https://raw.githubusercontent.com/v2fly/domain-list-community/release/cn.txt"
	domainListNotCn     = "https://raw.githubusercontent.com/v2fly/domain-list-community/release/geolocation-!cn.txt"
	suffixListRaw       = "https://raw.githubusercontent.com/publicsuffix/list/master/public_suffix_list.dat"
)

var (
	allowList  = make(map[string]struct{})
	blockList  = make(map[string]struct{})
	cnList     = make(map[string]struct{})
	proxyList  = make(map[string]struct{})
	suffixList = make(map[string]struct{})
	localList  = map[string]struct{}{
		"localhost":             {},
		"ip6-localhost":         {},
		"localhost.localdomain": {},
		"local":                 {},
		"broadcasthost":         {},
		"ip6-loopback":          {},
		"ip6-localnet":          {},
		"ip6-mcastprefix":       {},
		"ip6-allnodes":          {},
		"ip6-allrouters":        {},
		"ip6-allhosts":          {},
	}
)
