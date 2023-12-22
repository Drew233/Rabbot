package handlers

import (
	"fmt"
	"reflect"
	"strings"

	"rabbot/internal/rabmod"
	"rabbot/internal/common"
	"rabbot/internal/log"
	"rabbot/config"
)


func HandleRequestText(reqStruct *common.RequestStruct) (*common.ReplyStruct, error) {
	var reply *common.ReplyStruct
	var err error
	var commond string = ""
	var requestTxt string = ""

	// 响应消息格式：@{robot} commond requestTxt
	// requestTxt可选
	parts := strings.SplitN(reqStruct.RequestTxt, " ", 2)
	if len(parts) == 2 {
		commond, requestTxt = parts[0], parts[1]
	} else {
		commond = reqStruct.RequestTxt
	}

	log.RabLog.Infof("Received a msg, commond: %s, request txt: %s", commond, requestTxt)

	// 优先匹配默认回复
	if common.DefaultReply[commond] != "" {
		return &common.ReplyStruct{common.MsgTxt, common.DefaultReply[commond]}, nil
	}

	if !rabmod.FuncMap[commond].IsValid() {
		// 不是内置指令，请求通义千问
		responseTxt, err := rabmod.GetTyqwReply(commond, reqStruct.Uuid)
		if err != nil {
			if err.Error() == "请求成功，但响应失败" {
				// 如果是接口调用失败，不用返回“麻辣秃头”
				return &common.ReplyStruct{common.MsgTxt, common.UnknownReply}, nil
			} else if err.Error() == "invalid token"{ 
				return &common.ReplyStruct{common.MsgTxt, config.RabConfig.DefaultMsg.DullMsg}, nil
			}
			return &common.ReplyStruct{}, err
		}
		return &common.ReplyStruct{common.MsgTxt, responseTxt}, nil
	}

	if config.RabConfig.Features[commond].Enable != true || config.RabConfig.Features[commond].FeatureGpBlist[reqStruct.Groupname] == true {
		// 功能未启用
		log.RabLog.Infof("%s not enable in group %s, now feature config is %v", commond, reqStruct.Groupname, config.RabConfig.Features)
		return &common.ReplyStruct{common.MsgTxt, fmt.Sprintf(common.FeatureDisabled, commond)}, nil
	}

	reqStruct = &common.RequestStruct{
		Uname: reqStruct.Uname,
		Uuid: reqStruct.Uuid,
		Groupname: reqStruct.Groupname,
		RequestTxt: requestTxt,
		Commond: commond,
		Msg: reqStruct.Msg,
	}

	args := []reflect.Value{
		reflect.ValueOf(reqStruct),
	}
	result := rabmod.FuncMap[commond].Call(args)

	reply = result[0].Interface().(*common.ReplyStruct)
	errData := result[1].Interface()
	if errData != nil {
		err = result[1].Interface().(error)
		return &common.ReplyStruct{}, err
	}
	return reply, nil
}