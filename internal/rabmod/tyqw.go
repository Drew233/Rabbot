package rabmod

import (
	"bytes"
	"time"
	"io/ioutil"
	"net/http"
	"errors"
	"encoding/json"

	"rabbot/config"
	"rabbot/internal/log"
	"rabbot/internal/common"
)

var messHistoryMap map[string][]common.TyqwMessage = make(map[string][]common.TyqwMessage)
var messTimeStamp map[string] int64 = make(map[string] int64)

// 导出重置对话接口
func init() {
	common.FuncNameMap["DestroyHistory"] = DestroyHistory
}

// 定时清理对话记录
// 默认超时时间五分钟，五分钟定时检查一次
func CleanOuttimeHistory() {
	log.RabLog.Debugf("begin cleanouttimeHistory")
	for key, _ := range messHistoryMap {
		timeNow := time.Now().Unix()
		if timeNow > messTimeStamp[key] && (timeNow - messTimeStamp[key]) > 60 * 5 {
			log.RabLog.Infof("user %s timeout destory history success", key)
			delete(messHistoryMap, key)
			delete(messTimeStamp, key)
		}
	}
}

// 重置对话接口
func DestroyHistory(uname, uuid string) (*common.ReplyStruct, error) {
	messHistoryMap[uuid] = messHistoryMap[uuid][:0]
	log.RabLog.Infof("user %s destory history success", uname)
	return &common.ReplyStruct{common.MsgTxt, "对话已经重置啦~欢迎继续和我聊天哟"}, nil
}

// 请求通义千问
func GetTyqwReply(content, uuid string) (string, error) {
	// 每次收到请求时更新时间戳
	messTimeStamp[uuid] = time.Now().Unix()

	if _, ok := messHistoryMap[uuid]; !ok {
		messHistoryMap[uuid] = []common.TyqwMessage{}
	}
	messHistoryMap[uuid] = append(messHistoryMap[uuid], common.TyqwMessage{Role: "user", Content: content})

	log.RabLog.Debugf("user %s message history len is %d", uuid, len(messHistoryMap[uuid]))
	// 检查对话历史长度
	if len(messHistoryMap[uuid]) > config.RabConfig.TyqwMaxhis {
		// 删除最旧的元素
		messHistoryMap[uuid] = messHistoryMap[uuid][:0]
	}

	// 准备要发送的数据
	data := common.TyqwInput{
		Model:    "qwen-turbo",
		Input: struct {
			Messages []common.TyqwMessage `json:"messages"`
		}{
			Messages: messHistoryMap[uuid],
		},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.RabLog.Errorf("创建请求失败: %v", err)
		return "", err
	}

	// 创建一个请求
	req, err := http.NewRequest("POST", common.TyqwApiUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		log.RabLog.Errorf("创建请求失败: %v", err)
		return "", err
	}

	// 设置请求头
	req.Header.Set("Authorization", "Bearer " + config.RabConfig.TyqwToken)
	req.Header.Set("Content-Type", "application/json")

	// 发送请求并获取响应
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.RabLog.Errorf("发送请求失败: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应内容
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.RabLog.Errorf("读取响应失败: %v", err)
		return "", err
	}

	// 如果是token错误，视为没有配置，返回默认回复
	if resp.StatusCode == 401 {
		log.RabLog.Infof("Token错误")
		return "", errors.New("invalid token")
	}

	// 如果响应码不是200，打日志并返回error
	if resp.StatusCode != 200 {
		log.RabLog.Errorf("请求失败, 响应码: %d, 响应内容: %s", resp.StatusCode, string(respBody))
		return "", errors.New("请求成功，但响应失败")
	}

	var response common.TyqwResponse
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		log.RabLog.Errorf("解析 JSON 失败: %v", err)
		return "", err
	}

	hintRes := ""
	// 提取 text 字段内容并返回
	if len(messHistoryMap[uuid]) == 1 {
		hintRes = "欢迎来和兔兔聊天，不过呢兔兔脑子就那么点，如果五分钟你都没有接着说话我会忘了你的，并且只能记住十次对话哦。如果你想聊一个新的话题，记得和我说“重置对话”。\n" + common.Dilimiter
	}

	messHistoryMap[uuid] = append(messHistoryMap[uuid], common.TyqwMessage{Role: "assistant", Content: response.Output.Text})
	return hintRes + response.Output.Text, nil
}
