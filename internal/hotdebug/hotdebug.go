package hotdebug

import (
	"os"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/fsnotify/fsnotify"
	"rabbot/internal/common"
	"rabbot/internal/log"
)

func HotDebugInit()	error {
	// 先看下调试标记文件有没有
	if common.IsDbgMode() == true {
		log.RabLog.Infof("Debug flag exist, change log level to DebugLevel")
		log.RabLog.SetLevel(logrus.DebugLevel)
	}

	// 如果TMP目录不存在就新建
	if _, err := os.Stat(common.TmpDir); os.IsNotExist(err) {
		err := os.MkdirAll(common.TmpDir, 0755)
		if err != nil {
			fmt.Printf("Create dir {%s} failed, %v\n", common.TmpDir, err)
		} 
	}

	watch, err := fsnotify.NewWatcher()
	if err != nil {
		log.RabLog.Errorf("Create new fsnotify failed, err info: %v", err)
		return err
	}
	defer watch.Close()

	if err := watch.Add(common.TmpDir); err != nil {
		log.RabLog.Errorf("Watch debug flag add failed, err info: %v", err)
		return err
	}

	go func() {
		for {
			select {
			case ev := <-watch.Events:
				{
					fmt.Println(common.IsDbgMode())
					if (ev.Op&fsnotify.Create == fsnotify.Create && (common.IsDbgMode() == true)) {
						log.RabLog.Info("Debug flag created, change log level to DebugLevel")
						log.RabLog.SetLevel(logrus.DebugLevel)
					}
					if ((ev.Op&fsnotify.Remove == fsnotify.Remove || ev.Op&fsnotify.Rename == fsnotify.Rename) && common.IsDbgMode() == false) {
						log.RabLog.Info("Debug flag lossed, change log level to InfoLevel")
						log.RabLog.SetLevel(logrus.InfoLevel)
					}
				}
			case err := <-watch.Errors:
				{
					log.RabLog.Errorf("Select watch error happened, errinfo: %v", err)
					return ;
				}
			}
		}
	}()

	select{}
}