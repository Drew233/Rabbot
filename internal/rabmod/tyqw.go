package rabmod

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"errors"
	"encoding/json"

	"rabbot/config"
	"rabbot/internal/log"
	"rabbot/internal/common"
)

var messHistory []common.TyqwMessage

// 导出重置对话接口
func init() {
	common.FuncNameMap["DestroyHistory"] = DestroyHistory
}

func DestroyHistory(uname, uuid string) (*common.ReplyStruct, error) {
	messHistory = messHistory[:0]
	return &common.ReplyStruct{common.MsgTxt, "对话已经重置啦~欢迎继续和我聊天哟"}, nil
}

// 请求通义千问
func GetTyqwReply(content string) (string, error) {

	messHistory = append(messHistory, common.TyqwMessage{Role: "user", Content: content})

	// 检查对话历史长度
	if len(messHistory) > config.RabConfig.TyqwMaxhis {
		// 删除最旧的元素
		messHistory = messHistory[:0]
	}

	// 准备要发送的数据
	data := common.TyqwInput{
		Model:    "qwen-turbo",
		Input: struct {
			Messages []common.TyqwMessage `json:"messages"`
		}{
			Messages: messHistory,
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

	// 提取 text 字段内容并返回
	messHistory = append(messHistory, common.TyqwMessage{Role: "assistant", Content: response.Output.Text})
	return response.Output.Text, nil
}
