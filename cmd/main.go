package main

import (
	"flag"

	"rabbot/config"
	"rabbot/internal/log"
	"rabbot/internal/entry"
	"rabbot/internal/mysql"
	"rabbot/internal/common"
	"rabbot/internal/rabmod"
	"rabbot/internal/sigparse"
	"rabbot/internal/hotdebug"
	"rabbot/internal/scheduler"
)


func main() {

	// 配置文件路径，默认./config/rabbot.config
	cfgPath := flag.String("cfg-path", config.ConfigPath, "config file path")
	// 数据存储路径，默认./data
	dataPath := flag.String("data-path", "./rabdata", "program data file path")
	// 解析命令行参数并更新相关文件路径
	flag.Parse()
	common.DirUpdate(*dataPath)

	// 加载配置
	config.LoadConfig(*cfgPath)

	// 日志初始化
	log.RabLogInit()

	// 功能模块加载
	rabmod.ModInit()

	// 注册SIGINT, SIGTERM信号处理函数
	sigparse.SetupCloseHandler()

	// 数据库连接初始化
	mysql.InitDB()

	// 启动一个goroutine做调试级别热更新
	go hotdebug.HotDebugInit()

	// 启动一个goroutine(并发)开启定时任务
	// 1. 每天凌晨删除/data/tmp/目录下的文件
	go scheduler.RunSheduler()

	log.RabLog.Infof("Init task finished, begin run rabbot, now config is %v", config.RabConfig)

	// 启动机器人
	entry.Run()

}