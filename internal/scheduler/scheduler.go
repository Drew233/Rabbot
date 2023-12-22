// 定时任务
package scheduler

import (
	"github.com/robfig/cron"

	"rabbot/internal/log"
	"rabbot/internal/rabmod"
	"rabbot/internal/common"
	"rabbot/config"
)

var CronClean *cron.Cron

func RunSheduler() {
	CronClean := cron.New()
	// 每天凌晨的任务
	err := CronClean.AddFunc(config.RabConfig.Cron.CronDaily, func(){
		deleteTmpFiles(common.TmpDir)
	})
	if err != nil {
		log.RabLog.Errorf("Add cron job failed, %v", err)
		return
	}
	// 五分钟执行一次的任务
	err = CronClean.AddFunc(config.RabConfig.Cron.CronPerFM, func(){
		// 清理图片文件夹
		deletePicFiles(common.PicDir)
		// 清理通义千问缓存
		rabmod.CleanOuttimeHistory()
	})
	
	if err != nil {
		log.RabLog.Errorf("Add cron job failed, %v", err)
		return
	}
	CronClean.Start()
}