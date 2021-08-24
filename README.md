1. 安装 
```
git clone https://github.com/huangnie2020/yapi-sync
cd yaoi-sync
go build
./yapi-sync -path .

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
