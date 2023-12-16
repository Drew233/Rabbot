package log

import (
	"os"
	_log "log"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"rabbot/config"
	"rabbot/internal/common"
)

var RabLog = logrus.New()
var logFile *lumberjack.Logger

func LogFileClose() {
	_log.Println("close log file")
	logFile.Close()
	return
}

func RabLogInit() {
	// 如果日志目录不存在就新建
	if _, err := os.Stat(common.LogDir); os.IsNotExist(err) {
		err := os.MkdirAll(common.LogDir, 0755)
		if err != nil {
			_log.Fatalln("Create log dir {%s} failed, %v, now exist\n", common.LogDir, err)
		} 
	}

	// 循环写日志
	logFile = &lumberjack.Logger{
		Filename: common.LogFilename,
		MaxSize: config.RabConfig.RabLogConfig.Maxsize, // 日志最大占用，单位MB
		MaxBackups: config.RabConfig.RabLogConfig.Maxbackups, //最多切片数量
		MaxAge: config.RabConfig.RabLogConfig.Maxage, //日志最多存活时间
		Compress: config.RabConfig.RabLogConfig.Compress, //是否压缩
	}

	// 日志默认级别设置为Info
	RabLog.SetLevel(logrus.InfoLevel)
	// 记录日志时自动包含调用函数的信息
	RabLog.SetReportCaller(true)
	// 日志输出到文件
	RabLog.SetOutput(logFile)

	RabLog.Infof("Rablog init finished")
}