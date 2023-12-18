# RabBot

[![standard-readme compliant](https://img.shields.io/badge/openwechat-Rabbot-brightgreen.svg?style=flat-square)](https://github.com/Drew233/rabbot)

RabBot是基于[openwechat](https://github.com/eatMoreApple/openwechat)实现的微信机器人框架

框架包含
1. 通过[logrus](https://github.com/sirupsen/logrus)库,实现全局日志管理
2. 支持日记级别热更新（Info->Debug）
3. 统一消息处理接口
4. 配置化，可以通过仅修改配置达到部分自定义效果
    4.1 群聊黑名单
    4.2 功能是否启用（最小粒度为群聊）

## 内容列表

- [RabBot](#rabbot)
	- [内容列表](#内容列表)
	- [背景](#背景)
	- [使用需知](#使用需知)
	- [安装](#安装)
		- [下载代码到本地](#下载代码到本地)
		- [安装go和git](#安装go和git)
		- [安装框架依赖的模块](#安装框架依赖的模块)
	- [使用说明](#使用说明)
		- [配置文件说明](#配置文件说明)
		- [go run](#go-run)
		- [go build](#go-build)
		- [调试](#调试)
	- [徽章](#徽章)
	- [新增插件](#新增插件)
	- [相关仓库](#相关仓库)
	- [维护者](#维护者)
	- [如何贡献](#如何贡献)
	- [使用许可](#使用许可)

## 背景

很早之前就想做个微信机器人玩玩，当时只看到了主流的`wechaty`，但是登陆限制太多，又完全联系不上开源激励计划，所以一直没有机会

偶然间看到`openwechat`的打破登录限制，于是就上头了，Go是边做边学的，代码问题可能会有不少，欢迎提建议

## 使用需知
> [!WARNING]
> 1. 采用的openwechat提供的SDK，桌面微信协议
> 2. 使用机器人需要一个微信号，扫码登录且手机不可退出登录（有需要的话可以断网退出，具体操作可以在网上搜到）
> 3. 作为机器人的微信号需要实名认证，在实名认证是会让绑银行卡
> 4. 会有概率被封号

## 安装

### 下载代码到本地
通过git clone
```bash
$ git clone https://github.com/Drew233/Rabbot.git
```
或者
Download ZIP下载源码

或者前往Releases下载源码或可执行文件（推荐源码）

### 安装go和git
运行依赖于 [go](https://golang.google.cn/) 和 [git](https://npmjs.com)

### 安装框架依赖的模块
通过`go mod tidy`添加缺失的模块以及移除不需要的模块

```bash
$ go mod tidy
```

## 使用说明

### 配置文件说明
配置文件默认路径为`./rabbot/config/rabbot.config`，启动之前一定要把配置文件中需要自己修改的部分完成，否则会导致功能失效

这里必须要修改的只有：botName和groupWhiteList，其他字段可以根据自己需要修改

整体为一个json格式的字符串，下面会解释每个字段的含义

```txt
{
	"botName": "小兔子",
    // botName 机器人微信昵称，因为现在暂时取不到登陆账号自己在群聊中的备注，所以直接采用微信昵称。主要用于判断是否是@自己的消息。注意：如果在群聊中改了自己的昵称，会导致功能失效
	"datadir": "./rabdata",    
    // datadir 数据存放目录，主要用于存放日志以及一些临时文件
	"defaultMsg": {
		"dullMsg": "我只是一只小兔子，小兔子怎么会想到现在该说什么呢",
        // dullMsg 收到未知指令时的回复
		"errorMsg": "谁把我做成麻辣兔头了？！"
        // errorMsg 出现内部错误时候的回复
	},
    // defaultMsg 默认回复消息
	"groupWhiteList" : [
		"一群单身狗",
		"OW灌水群"
	],
    // groupWhiteList 群聊白名单，机器人只会响应群聊白名单里面的消息
	"rablog": {
		"maxsize": 5,
        // maxsize 日志文件的最大大小，单位为MB。当日志文件大小达到这个值时，会触发日志切割。
		"maxbackups": 3,
        // maxbackups 最多保留的旧日志文件数量，超过这个数量的旧日志文件会被删除。
		"maxage": 30,
        // maxage 旧日志文件的最大保留天数，超过这个天数的旧日志文件会被删除。
		"compress": true
        // compress 是否压缩旧日志文件，设置为 true 表示会对旧日志文件进行压缩。
	},
    // rablog 日志相关参数设置，日志路径固定为 ${datadir}/log/rabbot.log
	"cron": {
		"tmpCleanCron": "0 0 * * *"
        // tmpCleanCron 每天凌晨执行一次，这里主要是定时清理tmp目录下的文件
    },
    // 定时任务的cron表达式
	"feature": {
		"摸鱼日历": {
        // 摸鱼日历 功能名称，这里是直接匹配的指令名
			"enable": true,
            // enable 功能是否开启
			"entry": "GetFishCal",
			// entry 调用插件接口
			"groupBlackList": {
				"一群单身狗": true,
				"洋佬带下俺8": true
			}
            // groupBlackList 功能开启的情况下群聊黑名单，如果群聊名在黑名单中为true，那么此功能即使开启也不会在此群聊生效
		},
		"抽签": {
			"enable": true,
			"entry": "DrawLots",
			"groupBlackList": {
				"OW灌水群": true
			}
		}
	},
    // feature 任务开关，这里两个任务是当前内置的，后续扩展需要把功能添加进来，不在列表内的不会生效
	"mysql": {
		"username": "root",
        // username 用户名
		"password": "root",
        // password 密码
		"ip": "127.0.0.1",
        // ip 数据库ip地址
		"port": 3306,
        // ip 数据库端口
		"dbname": "chouqian",
        // dbname 数据库名
		"charset": "utf8"
        // charset 内容编码格式
	}
    // mysql 数据库配置，如果不需要数据库相关的功能，可以不配
}
```

配置文件处理好之后就可以运行了，有两种方式

### go run
直接通过go run执行main.go，在linux上面可以通过screen来放在后台执行

```bash
go run ./robbot/cmd/main.go
```

### go build
也可以通过go build构造可执行程序执行
```bash
go build ./robbot/cmd/main.go -o rabbot
chmod +x rabbot
./rabbot
```

通过build执行时有两个可选参数
```bash
Usage of rabbot.exe:
  -cfg-path string
        config file path (default "./config/rabbot.config")
  -data-path string
        program data file path (default "./data")
```

> [!NOTE]
> 运行起来之后会给出一个链接扫码登陆即可，但是手机上的微信退出的话机器人会一起退出。
> 解决方法：
> 1. 两个手机
> 2. 手机断网后退出登录，然后登陆其他账号再联网（退出登录时候需要多等一会）

### 调试
在`./rabdata/tmp/`目录下新建一个文件`RabDbg`，源码中Debug级别的日志就会输出

## 徽章

[![standard-readme compliant](https://img.shields.io/badge/readme%20style-standard-brightgreen.svg?style=flat-square)](https://github.com/RichardLitt/standard-readme)

## 新增插件

本框架下的插件业务代码统一放在rabmod包中（当前仅支持全匹配）

在本框架下添加插件需要以下操作（以抽签插件为例）
1. 将业务代码放到/rabbot/internal/rabmod/下（新增一个drlots.go）
2. 业务代码中添加init，将函数入口信息放在全局的map中（这里map的key要能和配置文件中的entry对应上）
```go
func init() {
	common.FuncNameMap["DrawLots"] = DrawLots
}
```
3. 业务代码coding...
> [!IMPORTANT]
> 插件入口函数的参数和返回值需和框架中保持一致
```go
func DrawLots(uname, uuid string) (*common.ReplyStruct, error)
```
1. 修改配置文件，在features下添加对应插件的信息即可
```json

		"抽签": {
			"enable": true,
			"entry": "DrawLots",
			"groupBlackList": {
				"OW灌水群": true
			}
		}
```

## 相关仓库

- [Openwechat](https://github.com/eatMoreApple/openwechat) - golang微信SDK
- [Wechatbot](https://github.com/djun/wechatbot) — 为个人微信接入ChatGPT
- [Standard-readme](https://github.com/RichardLitt/standard-readme) - A standard style for README files

## 维护者

[@Drew233](https://github.com/Drew233)。

## 如何贡献

非常欢迎你的加入！[提一个 Issue](https://github.com/RichardLitt/standard-readme/issues/new) 或者提交一个 Pull Request。

标准 Readme 遵循 [Contributor Covenant](http://contributor-covenant.org/version/1/3/0/) 行为规范。

## 使用许可

[MIT](LICENSE) © Richard Littauer