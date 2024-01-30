package rabmod

import (
	"os"
	"fmt"
	"time"
	"errors"
	"os/exec"
	"math/rand"

	"rabbot/config"
	"rabbot/internal/log"
	"rabbot/internal/common"

	"github.com/eatmoreapple/openwechat"
)

func init() {
	// 初始化时判断下存放历史记录的路径如果不存在就新建
	if _, err := os.Stat(common.HisDir); os.IsNotExist(err) {
		err := os.MkdirAll(common.HisDir, os.ModePerm)
		if err != nil {
			// 目录创建失败的话不启用此功能
			log.RabLog.Errorf("Make chat history dir failed: %v", err)
			return 
		}
	}
	common.FuncNameMap["Ciyun"] = Ciyun
}

var bot *openwechat.Bot

// 响应指令-词云测试，供测试用
func Ciyun(requestStruct *common.RequestStruct) (*common.ReplyStruct, error) {
	groupname := requestStruct.Groupname
	GenCiyun(groupname, "")
	SendCiyun(groupname)
	return &common.ReplyStruct{}, errors.New("No need response")
}

// 生成词云，如果指定了groupname就生成对应的群聊词云，不指定就生成所有的
func GenCiyun(groupname, time string) {
	for groupName := range config.RabConfig.WbList {
		if groupname != "" && groupname != groupName {
			continue
		}
		sourceFile := common.HisDir + "/history." + groupName
		destFile := common.HisDir + "/" + groupName + ".png"
		
		// 如果词云图片已经存在先删除
		_, err := os.Stat(destFile)
		if err == nil {
			reErr := os.Remove(destFile)
			if reErr != nil {
				log.RabLog.Errorf("%s 删除失败", destFile)
				continue
			}	
		}
	
		command := fmt.Sprintf("python3 ./tool/wc.py %s %s %s", sourceFile, destFile, time)

		_, err = exec.Command("bash", "-c", command).Output()

		if err != nil {
			log.RabLog.Debugf("群聊：%s生成词云失败\n", groupName)
			continue
		}
	}
}

// 发送词云，groupname同生成时一个作用
func SendCiyun(groupname string) {
	for groupName := range config.RabConfig.WbList {
		if config.RabConfig.Features["词云测试"].FeatureGpBlist[groupName] == true {
			log.RabLog.Infof("%s 在黑名单，不发词云", groupName)
			continue
		}

		if groupname != "" && groupname != groupName {
			continue
		}

		if groupName == "" {
			// 凌晨群发的时候，随机延迟1-5分钟，再发词云
			rand.Seed(time.Now().UnixNano())
			randSleep := rand.Intn(5) + 1
			time.Sleep(time.Duration(randSleep) * time.Minute)
		}

		destFile := common.HisDir + "/" + groupName + ".png"

		self, _ := bot.GetCurrentUser()
		groups, _ := self.Groups()
		group := groups.GetByNickName(groupName)

		if group == nil {
			log.RabLog.Infof("群聊%s未添加到通讯录，不发送词云", groupName)
			continue
		}

		img, _ := os.Open(destFile)
		defer img.Close()
	
		if _, err := group.SendImage(img); err != nil {
			log.RabLog.Errorf("发送信息失败: %v", err)
		}
	}
}

func CiyunInit(bot_param *openwechat.Bot) {
	bot = bot_param
}
