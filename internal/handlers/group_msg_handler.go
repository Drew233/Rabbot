package handlers

import (
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

	log.RabLog.Info(msg.IsPaiYiPai(), msg.IsTickled(), msg.IsTickledMe())
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
		return err
	}

	reply, err = HandleRequestText(common.GenRequestStruct(groupSender.DisplayName,
														   groupSender.UserName,
														   groupName,
														   requestText))
	if err != nil {
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