// 定时任务
package scheduler

import (
	"os"
	"path/filepath"

	"github.com/robfig/cron"

	"rabbot/internal/log"
	"rabbot/internal/common"
	"rabbot/config"
)

var CronClean *cron.Cron

func RunDailyFileCleanup() {
	CronClean := cron.New()
	err := CronClean.AddFunc(config.RabConfig.Cron.TmpCleanCron, func(){
		deleteFiles(common.TmpDir)
	})
	if err != nil {
		log.RabLog.Errorf("Add cron job failed, %v", err)
		return
	}
	CronClean.Start()
}

// 定时删除./data/tmp目录下的文件，默认一天删除一次
func deleteFiles(dirPath string) {
	log.RabLog.Info("Daily file cleanup begin")
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.RabLog.Errorf("Walk file %s failed, %v", path, err)
			return err
		}
		if !info.IsDir() {
			err := os.Remove(path)
			if err != nil {
				log.RabLog.Errorf("Remove file %s failed, %v", dirPath, err)
			} else {
				log.RabLog.Infof("Remove file %s successed", dirPath)
			}
		}

		return nil
	})
	if err != nil {
		log.RabLog.Errorf("Walk dirpath %s failed, %v", dirPath, err)
		return
	}
}