// 封装http为rabbot可用的接口
package rabhttp

import (
	"os"
	"io"
	"time"
	"net/http"
	"encoding/json"

	"rabbot/internal/log"
)

/* 
	@Func 向指定url发起get请求，并解析返回结果中的json
	@param url 发起请求的url
	@param jsonStruct 用于储存json的变量，用interface{}是为了通用，可以传入任意的参数，参数有效性由调用者决定
					  同时也是输出，转换后的json会赋值给jsonStruct
	@return error 如果出错的话返回错误信息
*/
func RabHttpGetJson(url string, jsonStruct interface{}) error {
	// 默认三秒超时
	cli := http.Client{Timeout: 3 * time.Second}
	resp, err := cli.Get(url)
	if err != nil {
		log.RabLog.Errorf("Get %s failed, %v", url, err)
		return err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&jsonStruct)
	if err != nil {
		log.RabLog.Errorf("Decode json string %s failed, %v", resp.Body, err)
		return err
	}

	_,err = json.Marshal(&jsonStruct)
	if err != nil {
		log.RabLog.Errorf("Parse json string %s failed, %v", resp.Body, err)
		return err
	}

	return nil
}

/*
	@Func 下载指定url的图片
	@param url 发起请求的url
	@param picPath 下载下来的图片保存的路径
*/
func RabHttpGetPic(url, picPath string) error {
	// 默认三秒超时
	cli := http.Client{Timeout: 3 * time.Second}
	resp, err := cli.Get(url)
	if err != nil {
		log.RabLog.Errorf("Download {%s} picture failed", url)
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(picPath)
	if err != nil {
		log.RabLog.Errorf("Create file {%s} filed", picPath)
		return err
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.RabLog.Errorf("Copy response to file error, %v", err)
		return err
	}

	return nil
}