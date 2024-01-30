/* 机器人启动入口 */
package entry

import (
	"github.com/eatmoreapple/openwechat"
	"rabbot/internal/handlers"
	"rabbot/internal/sitepush"
	"rabbot/internal/rabmod"
	"rabbot/internal/log"
)

func Run () {
	bot := openwechat.DefaultBot(openwechat.Desktop) // 桌面模式

	// 注册消息处理函数
	bot.MessageHandler = handlers.Handler
	// 注册登陆二维码回调
	bot.UUIDCallback = openwechat.PrintlnQrcodeUrl

	// 创建热存储容器对象
	reloadStorage := openwechat.NewFileHotReloadStorage("storage.json")
	// 执行热登录
	err := bot.HotLogin(reloadStorage)
	if err != nil {
		if err = bot.Login(); err != nil {
			log.RabLog.Errorf("login error: %v", err)
			return
		}
	}

	rabmod.CiyunInit(bot)

	// 轮询网站rss
	go sitepush.SPushEntry(bot)

	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	bot.Block()
}