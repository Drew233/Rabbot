/* 调用摸鱼日历接口，返回摸鱼日历 */

package rabmod

import (
	"os"

	"rabbot/internal/log"
	"rabbot/internal/common"
	"rabbot/internal/rabhttp"
)

func init() {
	common.FuncNameMap["GetFishCal"] = GetFishCal
}

// 下载摸鱼日历并返回
func GetFishCal(requestStruct *common.RequestStruct) (*common.ReplyStruct, error) {
	// var cal_data common.Cal_data
	// if err := rabhttp.RabHttpGetJson(common.CalenderUrl, &cal_data); err != nil {
	// 	log.RabLog.Errorf("Get fish calender api failed, %v", err)
	// 	return nil, err
	// }

	tmpFilePath := common.GenTmpFilePath()

	if _, err := os.Stat(tmpFilePath); err == nil {
		log.RabLog.Debug("File cache exist, send it directly")
		return &common.ReplyStruct{common.MsgPic, tmpFilePath}, nil
	} else {
		log.RabLog.Errorf("Check tmpFile failed, %v", err)
	}

	if err:= rabhttp.RabHttpGetPic(common.CalenderUrl, tmpFilePath); err != nil {
		log.RabLog.Errorf("Download pic from %s failed\n, %v", common.CalenderUrl, err)
		return nil, err
	}

	return &common.ReplyStruct{common.MsgPic, tmpFilePath}, nil
}