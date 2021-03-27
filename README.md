# GeoSite

广告以及国内域名的合集

## 运行

```bash
./geosite
```

稍等片刻输出 `geosite.dat`

* 可识别 `https_proxy` 环境变量

程序运行时自动加载同级目录中名为 `block.txt` 的文件，内容为域名列表的 URL (参考 [block.txt](block.txt)) 用于 `block` 标签
，请确保该文件存在

### 命令参数

所有参数默认为关闭状态

* `-f` 强制外置域名输出 `full:` 格式
* `-d` 强制外置域名输出 `domain:` 格式

```bash
# -d 和 -f 同为 true 时 -f 生效
# 两者未指定时或都为 false 时则根据域名自动处理格式
./geosite -f=true
```

* `-D` 自动检测并移除无效域名

使用 [Google DoH](https://dns.google) 查询所有域名，状态码为 `3` 时即为无效域名；
此操作会增大创建时间(内置域名数据≤1分钟)

```bash
./geosite -D=true
```

### 配置示例

```json
      {
        "type": "field",
        "domain": [
          "geosite:local",
          "geosite:allow"
        ],
        "outboundTag": "direct"
      },
      {
        "type": "field",
        "domain": [
          "geosite:block"
        ],
        "outboundTag": "block"
      },
      {
        "type": "field",
        "domain": [
          "geosite:cn"
        ],
        "outboundTag": "direct"
      }
```

## 引用以下项目

* [github.com/v2fly/domain-list-community](https://github.com/v2fly/domain-list-community)

内置域名数据：

* [github.com/CalmLong/whitelist](https://github.com/CalmLong/whitelist) `allow`
* [github.com/privacy-protection-tools/dead-horse](https://github.com/privacy-protection-tools/dead-horse) `allow`
* [github.com/felixonmars/dnsmasq-china-list](https://github.com/felixonmars/dnsmasq-china-list) `cn`
* [github.com/publicsuffix/list](https://github.com/publicsuffix/list)

仅用于测试：

* [https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts](https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts)
* [https://raw.githubusercontent.com/privacy-protection-tools/anti-AD/master/anti-ad-domains.txt](https://raw.githubusercontent.com/privacy-protection-tools/anti-AD/master/anti-ad-domains.txt)




