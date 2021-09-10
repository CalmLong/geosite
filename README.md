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

## 运行

```bash
./geosite
```

稍等片刻输出 `geosite.dat`

* 支持 HTTP 代理环境变量

> 程序运行时自动加载同级目录中名为 `block.txt` 的文件，内容为域名列表的 URL (参考 [block.txt](block.txt)) 用于 `category-ads-all` 标签
 
## 引用以下项目

* [github.com/v2fly/domain-list-community](https://github.com/v2fly/domain-list-community)
* [github.com/CalmLong/domain-list](https://github.com/CalmLong/domain-list)
* [github.com/felixonmars/dnsmasq-china-list](https://github.com/felixonmars/dnsmasq-china-list)
