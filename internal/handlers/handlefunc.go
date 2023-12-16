package handlers

import (
	"fmt"
	// "strings"
	"rabbot/internal/mods/drlots"
	"rabbot/internal/mods/fishcal"
	"rabbot/internal/common"
	"rabbot/internal/log"
	"rabbot/config"
)

var funcMap = map[string] func (uname, uuid string) (*common.ReplyStruct, error) {
	"抽签": drlots.DrawLots,
	"摸鱼日历": fishcal.GetFishCal,
}

func HandleRequestText(reqStruct *common.RequestStruct) (*common.ReplyStruct, error) {
	var reply *common.ReplyStruct
	var err error

	// 优先匹配默认回复
	if common.DefaultReply[reqStruct.RequestTxt] != "" {
		return &common.ReplyStruct{common.MsgTxt, common.DefaultReply[reqStruct.RequestTxt]}, nil
	}

	if funcMap[reqStruct.RequestTxt] == nil {
		// 不是内置指令，请求文心一言 TODO
		return &common.ReplyStruct{common.MsgTxt, config.RabConfig.DefaultMsg.DullMsg}, nil
	}

	if config.RabConfig.Features[reqStruct.RequestTxt].Enable != true || config.RabConfig.Features[reqStruct.RequestTxt].FeatureGpBlist[reqStruct.Groupname] == true {
		// 功能未启用
		log.RabLog.Infof("%s not enable in group %s, now feature config is %v", reqStruct.RequestTxt, reqStruct.Groupname, config.RabConfig.Features)
		return &common.ReplyStruct{common.MsgTxt, fmt.Sprintf(common.FeatureDisabled, reqStruct.RequestTxt)}, nil
	}

	if reply, err = funcMap[reqStruct.RequestTxt](reqStruct.Uname, reqStruct.Uuid); err != nil {
		return &common.ReplyStruct{}, err
	}

	return reply, nil
}