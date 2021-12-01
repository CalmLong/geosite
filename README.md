# GeoSite

广告以及国内域名的合集

## 获取

Github Action 每日自动构建

* **geosite.dat**: [https://raw.githubusercontent.com/CalmLong/geosite/release/v2ray/geosite.dat](https://raw.githubusercontent.com/CalmLong/geosite/release/v2ray/geosite.dat)

### 配置示例

```json
{
  "rules": [
    {
      "type": "field",
      "domain": [
        "geosite:category-ads-all"
      ],
      "outboundTag": "block"
    },
    {
      "type": "field",
      "domain": [
        "geosite:cn"
      ],
      "outboundTag": "direct"
    },
    {
      "type": "field",
      "domain": [
        "geosite:geolocation-!cn"
      ],
      "outboundTag": "proxy"
    }
  ]
}
```
 
## 引用以下项目

* [github.com/v2fly/domain-list-community](https://github.com/v2fly/domain-list-community)
* [github.com/CalmLong/domain-list](https://github.com/CalmLong/domain-list)
* [github.com/felixonmars/dnsmasq-china-list](https://github.com/felixonmars/dnsmasq-china-list)
