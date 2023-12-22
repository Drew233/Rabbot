/* 调用色图接口，返回色图 */

package rabmod

import (
	"os"

	"rabbot/config"
	"rabbot/internal/log"
	"rabbot/internal/common"
	"rabbot/internal/rabhttp"
)
func init() {
	common.FuncNameMap["GetSetu"] = GetSetu
}
// 下载色图并返回
func GetSetu(requestStruct *common.RequestStruct) (*common.ReplyStruct, error) {
	var stData common.SetuData
	if err := rabhttp.RabHttpGetJson(common.SetuUrl, &stData); err != nil {
		log.RabLog.Errorf("Get fish calender api failed, %v", err)
		return nil, err
	}

	tmpFilePath := common.GenPicFilePath()

	if _, err := os.Stat(tmpFilePath); err == nil {
		log.RabLog.Debug("File cache exist, send it directly")
		return &common.ReplyStruct{common.MsgPic, tmpFilePath}, nil
	} else {
		log.RabLog.Infof("Check tmpFile failed, %v", err)
	}

	if len(stData.Pics) <= 0 {
		log.RabLog.Error("Setu api return pic nums less than 0")
		return &common.ReplyStruct{common.MsgTxt, config.RabConfig.DefaultMsg.ErrMsg}, nil
	}

	if err:= rabhttp.RabHttpGetPic(stData.Pics[0], tmpFilePath); err != nil {
		log.RabLog.Errorf("Download pic from %s failed\n, %v", stData.Pics[0], err)
		return nil, err
	}

	return &common.ReplyStruct{common.MsgPic, tmpFilePath}, nil
}