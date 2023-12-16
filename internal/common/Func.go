/* 通用函数 */
package common

import (
	"os"
	"fmt"
	"time"

	"github.com/eatmoreapple/openwechat"
)

/*
	@Func 获取群聊名称
	@param msg openwechat框架中的msg结构
	@return string 返回获取到的群聊名称
*/
func GetGroupName(msg *openwechat.Message) string {
	sender, _ := msg.Sender()
	group,_ := sender.AsGroup()		// 将sender转为group类型
	group.Detail()					// 获取详细信息，保证信息及时更新
	return group.NickName
}

/*
	@Func 生成临时文件名，格式为./data/tmp/{timestamp}
	@return string 临时文件的路径
*/
func GenTmpFilePath() string {
	if _, err := os.Stat(TmpDir); os.IsNotExist(err) {
		err := os.MkdirAll(TmpDir, 0755)
		if err != nil {
			// log.Rabbot.Errorf("Create dir %s failed, %v", TmpDir, err)
		} 
	}
	currentTime := time.Now().Format("20060102")

	return fmt.Sprintf("%s/%s", TmpDir, currentTime)
}

/*
	@Func 发送图片
	@param filePath 文件路径
	@param msg openwechat框架中的msg结构
	@return error 出错的话返回error
*/
func ReplayPic(filePath string, msg *openwechat.Message) error {
	filePic, err := os.Open(filePath)
	if err != nil {
		// log.Rabbot.Errorf("Open file %s filed, %v", filePath, err)
		return nil
	}
	defer filePic.Close()

	msg.ReplyImage(filePic)

	return nil
}

/*
	@Func 发送消息，仅对本来的接口加了一个异常处理
	@param replayTxt 消息内容
	@param msg openwechat框架中的msg结构
	@return error 出错的话返回error
*/
func ReplyTxt(replayTxt string, msg *openwechat.Message) error {
	if _, err := msg.ReplyText(replayTxt); err != nil {
		// log.Rabbot.Errorf("Reply %s Failed, %v", replayTxt, err)
		return err
	}

	return nil
}

/*
	@Func 判断是否需要开启debug级别日志
	@return bool 需要开启返回true，否则返回false
*/
func IsDbgMode() bool {
	if _, err := os.Stat(DebugFlag); err == nil {
		return true
	}

	return false
}

/*
	@Func 更新相关文件路径
*/
func DirUpdate(dirData string) {
	DataDir = dirData		// 数据目录
	TmpDir = DataDir + "/tmp"   // 临时目录
	LogDir = DataDir + "/log"	// 日志目录
	LogFilename = LogDir + "/rabbot.log"  // 日志文件名
	DebugFlag = TmpDir + "/RabDbg"  // 调试标记
}