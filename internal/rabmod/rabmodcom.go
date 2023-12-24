package rabmod

import (
	"io/ioutil"
	"encoding/json"

	"rabbot/internal/log"
)

func saveJSON(data interface{}, filename string) error {
	// 将数据序列化为JSON格式
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		log.RabLog.Errorf("Error encoding JSON:", err)
		return err
	}

	// 将JSON数据写入文件
	err = ioutil.WriteFile(filename, jsonData, 0644)
	if err != nil {
		log.RabLog.Errorf("Error writing to file:", err)
		return err
	}
	return nil
}