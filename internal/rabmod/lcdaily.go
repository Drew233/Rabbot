package rabmod

import (
	"os"
	"fmt"
	"strings"
	"io/ioutil"
	"net/http"
	"encoding/json"

	"rabbot/internal/log"
	"rabbot/internal/common"
)

// 导出每日一题接口
func init() {
	common.FuncNameMap["LCDaily"] = LCDaily
}

type tagInfo struct {
	Name string `json:"name"`
	NameCN string `json:"nameTranslated"`
}

type quesInfo struct {
	Title string `json:"title"`
	TitleCn string `json:"titleCn"`
	TitleSlug string `json:"titleSlug"`
	AcRate float64 `json:"acRate"`
	Diff string `json:"difficulty"`
	Tags []tagInfo `json:"topicTags"`
}

type leetcodeProblem struct {
	Data struct {
		TodayRecord []struct {
			Question quesInfo `json:"question"`
		}`json:"todayRecord"`
	}`json:"data"`
}

func getRateHint(rate float64) string {
	if rate < 0.1 {
		return fmt.Sprintf("😨本题历史通过率仅有%.2f%%，尽力而为吧骚年", rate * 100)
	}
	if rate < 0.3 {
		return fmt.Sprintf("🫢本题历史通过率为%.2f%%，也不是很难对吧？", rate * 100)
	}
	if rate < 0.6 {
		return fmt.Sprintf("😄本题历史通过率为%.2f%%，想必以你的实力一定是洒洒水吧", rate * 100)
	}

	return fmt.Sprintf("🌼本题历史通过率高达%.2f%%，还不赶紧去秒了？", rate * 100)
}

func getTagStr(tags []tagInfo) string {
	str := ""
	for _, tag := range tags {
		if str != "" {
			str += "，"
		}
		if tag.NameCN == "" {
			str += tag.Name
		} else {
			str += tag.NameCN
		}
	}
	return str
}

func updateLCDInfo() error {
	url := "https://leetcode.cn/graphql/"
	method := "POST"

	payload := strings.NewReader("{\"query\":\"\\n    query questionOfToday {\\n  todayRecord {\\n    question {\\n      difficulty\\n      title\\n      titleCn: translatedTitle\\n      titleSlug\\n      acRate\\n      topicTags {\\n        name\\n        nameTranslated: translatedName\\n      }\\n    }\\n  }\\n}\\n    \",\"variables\":{}}")

	client := &http.Client {}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return err
	}

	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Language", "en,zh-CN;q=0.9,zh;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0")
	req.Header.Add("content-type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}

	var data leetcodeProblem
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return saveJSON(data, common.LCDailyFile)
}

func LCDaily(requestStruct *common.RequestStruct) (*common.ReplyStruct, error) {
	if _, err := os.Stat(common.LCDailyFile); err == nil {
		log.RabLog.Debug("File cache exist")
	} else {
		if err := updateLCDInfo(); err != nil {
			log.RabLog.Errorf("UpdateLCDaily failed, %v", err)
			return nil, err
		}
	}

	// 读取力扣每日一题json文件
	file, err := os.Open(common.LCDailyFile)
	if err != nil {
		log.RabLog.Errorf("Open LCDaily json file failed, %v", err)
		return nil, err
	}
	defer file.Close()

	jsonData, err := ioutil.ReadAll(file)
	if err != nil {
		log.RabLog.Errorf("Read LCDaily json file failed, %v", err)
		return nil, err
	}

	var data leetcodeProblem
	err = json.Unmarshal(jsonData, &data)
	if err != nil || len(data.Data.TodayRecord) == 0 {
		log.RabLog.Errorf("Translate LCDaily json file content to json failed, %v", err)
		return nil, err
	}

	replyStr := ""
	str := "力扣每日一题：\n题目：%s\n难度：%s\n标签：%s\n题目链接：%s\n"
	for _, record := range data.Data.TodayRecord {
		ques := record.Question
		replyStr = (fmt.Sprintf(str, ques.TitleCn, ques.Diff, getTagStr(ques.Tags), "https://leetcode.cn/problems/" + ques.TitleSlug) + getRateHint(ques.AcRate))
	}

	return &common.ReplyStruct{common.MsgTxt, replyStr}, nil
}