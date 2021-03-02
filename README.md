# GeoSite

广告以及国内域名的合集

## 运行

```bash
./geosite
```

稍等片刻输出 `geosite.dat`

* 可识别 `https_proxy` 环境变量

### 配置示例

```json
      {
        "type": "field",
        "domain": [
          "geosite:allow"
        ],
        "outboundTag": "direct"
      },      
      {
        "type": "field",
        "domain": [
          "geosite:ads"
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

## 其他

程序运行时加载同级目录中名为 `block.txt` 的文件，文件内容为包含域名的 Raw Url；请确保该文件存在

`block.txt` 内的 URL 可以自由添加或删除，其中的域名用于 `ads` 标签

## 引用以下项目

源码：

* [github.com/v2fly/domain-list-community](https://github.com/v2fly/domain-list-community)

域名数据来源：

* [github.com/CalmLong/whitelist](https://github.com/CalmLong/whitelist) `allow`
* [github.com/privacy-protection-tools/dead-horse](https://github.com/privacy-protection-tools/dead-horse) `allow`
* [github.com/neodevpro/neodevhost](https://github.com/neodevpro/neodevhost) `allow`
* [github.com/anudeepND/whitelist](https://github.com/anudeepND/whitelist) `allow`
* [github.com/felixonmars/dnsmasq-china-list](https://github.com/felixonmars/dnsmasq-china-list) `direct`

仅用于测试：

* [https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts](https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts)
* [https://raw.githubusercontent.com/privacy-protection-tools/anti-AD/master/anti-ad-domains.txt](https://raw.githubusercontent.com/privacy-protection-tools/anti-AD/master/anti-ad-domains.txt)




