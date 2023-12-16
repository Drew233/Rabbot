/*
	调用抽签api接口，实现抽签解签
*/
package drlots

import (
	"fmt"
	"math/rand"
	"time"

	"rabbot/internal/log"
	"rabbot/internal/common"
	"rabbot/internal/rabhttp"
)

func genRandom(max int) int {
	// 设置随机数种子
	rand.Seed(time.Now().UnixNano())

	// 生成随机数
	randomNumber := rand.Intn(max) // 生成 0 到 100 之间的随机整数

	return randomNumber
}

// 拼接随机抽签url
func genUrl() string {
	drowLotType := common.DrlotTypes[genRandom(len(common.DrlotTypes))]
	url := fmt.Sprintf(common.DrlotUrl, drowLotType)

	return url
}

// 抽签
func DrawLots(uname, uuid string) (*common.ReplyStruct, error) {

	url := genUrl()

	var data common.DrlotsData
	if err := rabhttp.RabHttpGetJson(url, &data); err != nil {
		return nil, err
	}

	finalRes := "\n" + data.Data[0].Name + data.Data[0].Text + "\n" + common.Dilimiter
	for _, item := range data.Data {
		if prefix, ok := common.DrlotNameMap[item.Name]; ok {
			finalRes += prefix + item.Name + item.Text + "\n"
		}
	}

	log.RabLog.Debug("DrawLots res is " + finalRes)

	return &common.ReplyStruct{common.MsgTxt, finalRes}, nil
}
