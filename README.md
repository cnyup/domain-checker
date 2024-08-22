这是一个用于检测证书文件是否到期以及还有多久到期的简单脚本


## 结果显示
目前支持结果显示到：
1. 控制台显示 
2. 发送到钉钉



## 如何安装
### 源码安装
1. git clone
2. make build
3. 运行二进制文件查看参数


## 如何使用
直接运行二进制文件可查看相关命令
```shell
./domain-checker
check domain cert and send result to somewhere

Usage:
  domain-checker [command]

Available Commands:
  check       check cert
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  version     Print the version number of domain-checker

Flags:
  -h, --help   help for domain-checker

Use "domain-checker [command] --help" for more information about a command.
```
查看check子命令相关：
```shell
./domain-checker check
check cert

Usage:
  domain-checker check [command]

Available Commands:
  ding        Print the result to dingTalk
  stdout      Print the result of stdout

Flags:
      --days int        Alarm threshold days default 15 (default 15)
  -d, --dir string      cert file dir
  -h, --help            help for check
  -f, --path string     cert file path
      --suffix string   file suffix match

Use "domain-checker check [command] --help" for more information about a command.

```

只检测单个文件并将结果显示在终端控制台
```shell
./domain-checker check  stdout  -f  "your cert file path" 
```
只检测单个文件并将结果发送到钉钉
```shell
./domain-checker check  dingtalk  -f  "your cert file path"   --token  "your token"
```

检测目录下后缀为crt的证书文件并将结果显示在控制台
```shell
./domain-checker check  dingtalk  -f  "your cert file path"  --token  "your token"
```

检测目录下后缀为crt的证书文件并将结果发送到钉钉
建议使用检测目录时,指定suffix参数 否则结果可能不是预期的
```shell
./domain-checker check  dingtalk  -d  "your cert dir " --suffix ".crt"  --token  "your token"
```


### TODO
1. 结果显示支持企业微信,飞书
2. 支持命令补全
3. 支持参数从文件中获取
