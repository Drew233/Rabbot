package handlers

import (
	"fmt"
	"reflect"
	// "strings"
	"rabbot/internal/rabmod"
	"rabbot/internal/common"
	"rabbot/internal/log"
	"rabbot/config"
)


func HandleRequestText(reqStruct *common.RequestStruct) (*common.ReplyStruct, error) {
	var reply *common.ReplyStruct
	var err error

	// 优先匹配默认回复
	if common.DefaultReply[reqStruct.RequestTxt] != "" {
		return &common.ReplyStruct{common.MsgTxt, common.DefaultReply[reqStruct.RequestTxt]}, nil
	}

	if !rabmod.FuncMap[reqStruct.RequestTxt].IsValid() {
		// 不是内置指令，请求文心一言 TODO
		return &common.ReplyStruct{common.MsgTxt, config.RabConfig.DefaultMsg.DullMsg}, nil
	}

	if config.RabConfig.Features[reqStruct.RequestTxt].Enable != true || config.RabConfig.Features[reqStruct.RequestTxt].FeatureGpBlist[reqStruct.Groupname] == true {
		// 功能未启用
		log.RabLog.Infof("%s not enable in group %s, now feature config is %v", reqStruct.RequestTxt, reqStruct.Groupname, config.RabConfig.Features)
		return &common.ReplyStruct{common.MsgTxt, fmt.Sprintf(common.FeatureDisabled, reqStruct.RequestTxt)}, nil
	}

	args := []reflect.Value{
		reflect.ValueOf(reqStruct.Uname),
		reflect.ValueOf(reqStruct.Uuid),
	}
	result := rabmod.FuncMap[reqStruct.RequestTxt].Call(args)

	reply = result[0].Interface().(*common.ReplyStruct)
	errData := result[1].Interface()
	if errData != nil {
		err = result[1].Interface().(error)
		return &common.ReplyStruct{}, err
	}

	return reply, nil
}