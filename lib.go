package main

const (
	localTag     = "local"
	allowTag     = "allow"
	blockTag     = "block"
	cnTag        = "cn"
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

const (
	domainListCnRaw     = "https://raw.githubusercontent.com/v2fly/domain-list-community/release/cn.txt"
	domainListAdsAllRaw = "https://raw.githubusercontent.com/v2fly/domain-list-community/release/category-ads-all.txt"
	suffixListRaw       = "https://raw.githubusercontent.com/publicsuffix/list/master/public_suffix_list.dat"
)

var (
	allowList  = make(map[string]struct{})
	blockList  = make(map[string]struct{})
	cnList     = make(map[string]struct{})
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
