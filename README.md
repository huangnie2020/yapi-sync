1. 安装 
```
go get -u github.com/huangnie2020/yapi-sync

```
2. 配置
- 创建 yapi-sync.json
```json
{
  "type": "swagger",
  "token": "yourtoken",
  "file": "read.swagger.json,write.swagger.json",
  "merge": "good",
  "server": "http://yourip:3000"
}
```