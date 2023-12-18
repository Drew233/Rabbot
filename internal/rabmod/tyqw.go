package rabmod

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"errors"
	"encoding/json"

	"rabbot/config"
	"rabbot/internal/log"
	"rabbot/internal/common"
)

// 请求通义千问
func GetTyqwReply(content string) (string, error) {
	// 准备要发送的数据
	data := []byte(fmt.Sprintf(`{
		"model": "qwen-turbo",
		"input": {
			"messages": [
				{
					"role": "user",
					"content": "%s"
				}
			]
		}
	}`, content))

	// 创建一个请求
	req, err := http.NewRequest("POST", common.TyqwApiUrl, bytes.NewBuffer(data))
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
	return response.Output.Text, nil
}
