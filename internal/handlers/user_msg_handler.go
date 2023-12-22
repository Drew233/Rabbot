package handlers

import (
	"strings"

	"github.com/eatmoreapple/openwechat"
	"rabbot/internal/common"
	"rabbot/internal/log"
	"rabbot/config"
)

var _ MessageHandlerInterface = (*UserMessageHandler)(nil)

// UserMessageHandler 私聊消息处理
type UserMessageHandler struct {
}

// handle 处理消息
func (g *UserMessageHandler) handle(msg *openwechat.Message) error {
	if msg.IsText() && msg.Content != "" {
		return g.ReplyText(msg)
	}
	return nil
}

// NewUserMessageHandler 创建私聊处理器
func NewUserMessageHandler() MessageHandlerInterface {
	return &UserMessageHandler{}
}

// ReplyText 发送文本消息给朋友
func (g *UserMessageHandler) ReplyText(msg *openwechat.Message) error {
	var err error
	var reply *common.ReplyStruct

	// 接收私聊消息
	sender, err := msg.Sender()
	if err != nil {
		log.RabLog.Errorf("gpt sender error: %v", err)
		msg.ReplyText(config.RabConfig.DefaultMsg.ErrMsg)
		return err
	}

	requestText := strings.Trim(msg.Content, "\n")
	if sender.DisplayName == "" {
		sender.DisplayName = sender.NickName
	}
	log.RabLog.Debugf("Received User %v Text Msg : %v", sender.DisplayName, msg.Content)

	reply, err = HandleRequestText(common.GenRequestStruct(sender.DisplayName,
														   sender.UserName,
														   "",
														   requestText))
	if err != nil {
		log.RabLog.Errorf("gpt request error: %v", err)
		msg.ReplyText(config.RabConfig.DefaultMsg.ErrMsg)
		return err
	}

	// 回复用户
	switch reply.ReType {
	case common.MsgTxt:
		return common.ReplyTxt(strings.Trim(strings.TrimSpace(reply.ReText), "\n"), msg)
	case common.MsgPic:
		return common.ReplayPic(reply.ReText, msg)
	default:
		// 出错时候返回认错txt
		return common.ReplyTxt(config.RabConfig.DefaultMsg.DullMsg, msg)
	}
}