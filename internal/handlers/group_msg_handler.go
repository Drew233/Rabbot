package handlers

import (
	"os"
	"bufio"
	"strings"

	"github.com/eatmoreapple/openwechat"
	"rabbot/internal/common"
	"rabbot/internal/log"
	"rabbot/config"
)

var _ MessageHandlerInterface = (*GroupMessageHandler)(nil)
var paiCnt = 0

// GroupMessageHandler 群消息处理
type GroupMessageHandler struct {
}

// handle 处理消息
func (g *GroupMessageHandler) handle(msg *openwechat.Message) error {
	groupName, terr := common.GetGroupName(msg)
	if terr != "" {
		log.RabLog.Errorf("get group name failed at %s", terr)
		return nil
	}

	if config.RabConfig.WbList[groupName] != true {
		log.RabLog.Debugf("Group => {%s} not support, no response", groupName)
		return nil
	}

	if msg.IsText() && msg.Content != "" {
		// 保存每个群的聊天记录，用于生成词云
		filePath := common.HisDir + "/history." + groupName
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			file, _ := os.Create(filePath)
			defer file.Close()
		}
		file, _ := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		defer file.Close()

		writer := bufio.NewWriter(file)
		defer writer.Flush()

		data := []byte(msg.Content)
		_, _ = writer.Write(data)

		return g.ReplyText(msg)
	}
	return nil
}

// NewGroupMessageHandler 创建群消息处理器
func NewGroupMessageHandler() MessageHandlerInterface {
	return &GroupMessageHandler{}
}

// ReplyText 发送文本消息到群
func (g *GroupMessageHandler) ReplyText(msg *openwechat.Message) error {
	var err error
	var reply *common.ReplyStruct

	// 接收群消息
	groupName, terr := common.GetGroupName(msg)
	if terr != "" {
		log.RabLog.Errorf("get group name failed at %s", terr)
		return nil
	}

	log.RabLog.Debugf("Received Group %v Text Msg : %v", groupName, msg.Content)

	if msg.IsPaiYiPai() || msg.IsTickled() || msg.IsTickledMe() {
		if paiCnt + 1 < len(config.RabConfig.DefaultMsg.PaiMsg) {
			paiCnt++
		}
		msg.ReplyText(config.RabConfig.DefaultMsg.PaiMsg[paiCnt])
		return nil
	}
	if paiCnt - 1 >= 0 {
		paiCnt--
	}

	// 不是@的不处理
	if !msg.IsAt() {
		return nil
	}

	// 获取@我的用户
	groupSender, err := msg.SenderInGroup()
	if groupSender.DisplayName == "" {
		groupSender.DisplayName = groupSender.NickName
	}

	// 替换掉@文本
	replaceText := "@" + config.RabConfig.BotName
	requestText := strings.TrimSpace(strings.ReplaceAll(msg.Content, replaceText, ""))
	
	if err != nil {
		log.RabLog.Errorf("get sender in group error :%v", err)
		msg.ReplyText(config.RabConfig.DefaultMsg.ErrMsg)
		return err
	}

	reply, err = HandleRequestText(common.GenRequestStruct(groupSender.DisplayName,
														   groupSender.UserName,
														   groupName,
														   requestText,
														   "",
														   msg))
	if err != nil {
		// 如果返回的错误信息是No need reponse，不需要做处理
		// 这里设计的不太好，后面有空再改吧
		if err.Error() == "No need response" {
			return nil
		}
		log.RabLog.Errorf("gpt request error: %v", err)
		msg.ReplyText(config.RabConfig.DefaultMsg.ErrMsg)
		return err
	}

	// 构造@发起请求用户的字符串，虽然现在@成员没用
	atText := "@" + groupSender.DisplayName + " "

	switch reply.ReType {
	case common.MsgTxt:
		replyText := atText + strings.Trim(strings.TrimSpace(reply.ReText), "\n")
		return common.ReplyTxt(replyText, msg)
	case common.MsgPic:
		return common.ReplayPic(reply.ReText, msg)
	default:
		// 出错时候返回认错txt
		replyText := atText + config.RabConfig.DefaultMsg.DullMsg
		return common.ReplyTxt(replyText, msg)
	}
}