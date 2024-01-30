// 定时任务
package scheduler

import (
	"time"
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
		// 获取当前时间的前一天作为词云标题，发送词云给每个开启功能的群聊
		currentTime := time.Now()
		previousDay := currentTime.AddDate(0, 0, -1)
		formattedDate := previousDay.Format("2006.1.2")
		rabmod.GenCiyun("", formattedDate)
		rabmod.SendCiyun("")
		//发送完成后删除掉这个目录下的文件，这里会把凌晨N（1~5*开启功能的群聊）分钟内的聊天记录统计不上，问题不大
		deleteTmpFiles(common.HisDir)
	})
	if err != nil {
		log.RabLog.Errorf("Add cron job failed, %v", err)
		return
	}
	// 一小时执行一次的任务
	err = CronClean.AddFunc(config.RabConfig.Cron.CronPerH, func(){
		// 因为通过率一直在变，所以每小时清空一次每日一题的数据
		deleteFile(common.LCDailyFile)
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
		// 清理超时的赛马比赛缓存
		rabmod.CleanOuttimeComp()
	})
	
	if err != nil {
		log.RabLog.Errorf("Add cron job failed, %v", err)
		return
	}
	CronClean.Start()
}