package hosts

import (
	"bufio"
	"net"
	"net/url"
	"sort"
	"strings"
)

var domainSuffix = []string{".com.cn", ".net.cn", ".org.cn", ".gov.cn", ".ah.cn", ".bj.cn", ".cq.cn", ".fj.cn",
	".gd.cn", ".gs.cn", ".gx.cn", ".gz.cn", ".ha.cn", ".hb.cn", ".he.cn", ".hi.cn", ".hk.cn", ".hn.cn", ".jl.cn",
	".js.cn", ".jx.cn", ".ln.cn", ".mo.cn", ".nm.cn", ".nx.cn", ".qh.cn", ".sc.cn", ".sd.cn", ".sh.cn", ".sn.cn",
	".sx.cn", ".tj.cn", ".tw.cn", ".xj.cn", ".yn.cn", ".zj.cn",
}

const (
	j = '#'
	
	dnsmasqIndex = "server=/"
	dnsmasqLast  = "/114.114.114.114"
)

func Classify(dst map[string]struct{}, writer *bufio.Writer, params []string) int {
	domains := make([]string, 0)
	fulls := make([]string, 0)
	for k := range dst {
		if err := net.ParseIP(k); err != nil {
			continue
		}
		switch strings.Count(k, ".") {
		case 1:
			domains = append(domains, k)
		case 2:
			var is bool
			for _, suffix := range domainSuffix {
				if strings.Contains(k, suffix) {
					is = true
					domains = append(domains, k)
					break
				}
			}
			if is {
				continue
			}
			fallthrough
		default:
			fulls = append(fulls, k)
		}
	}
	sort.Strings(fulls)
	sort.Strings(domains)
	for _, f := range fulls {
		_, _ = writer.WriteString(params[0] + f + params[1] + "\n")
	}
	for _, d := range domains {
		_, _ = writer.WriteString(params[2] + d + params[3] + "\n")
	}
	return len(domains) + len(fulls)
}

func format(newOrg string, prefix []string) string {
	for _, s := range prefix {
		newOrg = strings.ReplaceAll(newOrg, s, "")
	}
	return newOrg
}

func isNotExit(original string, allow ...map[string]struct{}) bool {
	if len(allow) > 0 && len(allow[0]) > 0 {
		if _, ok := allow[0][original]; ok {
			return false
		}
	}
	return true
}

func parseUrl(raw string) string {
	raw = strings.ReplaceAll(raw, "http://", "")
	raw = strings.ReplaceAll(raw, "https://", "")
	raw = strings.ReplaceAll(raw, "ftp://", "")
	raw = strings.ReplaceAll(raw, "websocket://", "")
	if i := strings.IndexRune(raw, '/'); i != -1 {
		return raw[:i]
	}
	return raw
}

func Resolve(src map[string]struct{}, dst map[string]struct{}, allow ...map[string]struct{}) {
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
		newOrg := original
		// 移除前缀为 0.0.0.0 或者 127.0.0.1 (移除第一个空格前的内容)
		index := strings.IndexRune(original, ' ')
		if index > -1 {
			newOrg = strings.ReplaceAll(original, original[:index], "")
		}
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
		newOrg = format(newOrg, []string{dnsmasqIndex, dnsmasqLast})
		// adblock
		if strings.ContainsRune(newOrg, '^') {
			// 子域名包含 * 的不会被解析
			if strings.ContainsRune(newOrg, '*') {
				continue
			}
			// 表达式不会被解析
			if strings.Contains(newOrg, "/^") {
				continue
			}
			// 基础白名单规则会被一同解析
			newOrg = format(newOrg, []string{"||", "^", "@@"})
		}
		newOrg = strings.TrimSpace(newOrg)
		// 检测是否有端口号，有则移除端口号
		if strings.ContainsRune(newOrg, ':') {
			newOrg = newOrg[:strings.IndexRune(newOrg, ':')]
		}
		newOrg = parseUrl(newOrg)
		if !isNotExit(newOrg, allow...) {
			continue
		}
		urlStr, err := url.Parse(newOrg)
		if err != nil {
			continue
		}
		// 如果为 IP 则跳过
		if err := net.ParseIP(urlStr.String()); err != nil {
			continue
		}
		dst[urlStr.String()] = struct{}{}
	}
}
