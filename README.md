# GeoSite

广告以及国内域名的合集

## 运行

```bash
./geosite
```

稍等片刻输出 `geosite.dat`

* 支持 HTTP 代理环境变量

> 程序运行时自动加载同级目录中名为 `block.txt` 的文件，内容为域名列表的 URL (参考 [block.txt](block.txt)) 用于 `block` 标签
> ，请确保该文件存在

[点击这里](https://raw.githubusercontent.com/CalmLong/geosite/release/v2ray/geosite.dat)下载已构建的数据包，由 Github Action 每日 UTC+08:00 2 点自动构建

### 命令参数

* `-D` 自动检测并移除无效域名

使用 [Google DoH](https://dns.google) 查询所有域名，DNS Response Code 等于 `3` 时即为无效域名；
此操作会增大创建时间(内置域名数据≤1分钟)

```bash
./geosite -D=true
```

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
* [github.com/privacy-protection-tools/dead-horse](https://github.com/privacy-protection-tools/dead-horse)
* [github.com/felixonmars/dnsmasq-china-list](https://github.com/felixonmars/dnsmasq-china-list)
* [github.com/publicsuffix/list](https://github.com/publicsuffix/list)

