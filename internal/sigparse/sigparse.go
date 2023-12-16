package sigparse

import (
	"os"
	"os/signal"
	"syscall"

	"rabbot/internal/log"
)

// 注册ctrl+c（SIGINT）和SIGTERM的信号处理
func SetupCloseHandler() {
    c := make(chan os.Signal, 2)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-c
        log.RabLog.Infof("\r- Ctrl+C pressed in Terminal")
		// 关闭日志文件句柄
		log.LogFileClose()
        os.Exit(0)
    }()
}
