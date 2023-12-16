/* 调用摸鱼日历接口，返回摸鱼日历 */

package fishcal

import (
	"os"

	"rabbot/internal/log"
	"rabbot/internal/common"
	"rabbot/internal/rabhttp"
)

// 下载摸鱼日历并返回
func GetFishCal(uname, uuid string) (*common.ReplyStruct, error) {
	var cal_data common.Cal_data
	if err := rabhttp.RabHttpGetJson(common.CalenderUrl, &cal_data); err != nil {
		log.RabLog.Errorf("Get fish calender api failed, %v", err)
		return nil, err
	}

	tmpFilePath := common.GenTmpFilePath()

	if _, err := os.Stat(tmpFilePath); err == nil {
		log.RabLog.Debug("File cache exist, send it directly")
		return &common.ReplyStruct{common.MsgPic, tmpFilePath}, nil
	} else {
		log.RabLog.Errorf("Check tmpFile failed, %v", err)
	}

	if err:= rabhttp.RabHttpGetPic(cal_data.Url, tmpFilePath); err != nil {
		log.RabLog.Errorf("Download pic from %s failed\n, %v", cal_data.Url, err)
		return nil, err
	}

	return &common.ReplyStruct{common.MsgPic, tmpFilePath}, nil
}