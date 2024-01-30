package sitepush

import (
	"os"
	"fmt"
	"time"
	"strings"
	"net/http"
	"io/ioutil"
	"encoding/json"

	"github.com/mmcdole/gofeed"
	"github.com/eatmoreapple/openwechat"

	"rabbot/internal/log"
	"rabbot/internal/common"
)

var bot *openwechat.Bot
type urlCacheS struct {
	CommentCache string `json:"commentCache"`
	PostCache string	`json:"postCache"`
}
var urlCacheMap = make(map[string]string)
var urlCache urlCacheS = urlCacheS{"", ""}

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

func parseTime(timeStr string) string {
	// 解析时间字符串
	parsedTime, err := time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", timeStr)
	if err != nil {
		log.RabLog.Errorf("Error parsing time:", err)
		return "获取时间异常"
	}

	// 格式化为指定格式
	formattedTime := parsedTime.Format("2006.01.02 15:04:05")
	return formattedTime
}

func timeAfter(timeStr1, timeStr2 string) bool {
	// 解析时间字符串
	parsedTime1, err := time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", timeStr1)
	if err != nil {
		log.RabLog.Errorf("Error parsing time:", err)
		return false
	}

	parsedTime2, err := time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", timeStr2)
	if err != nil {
		log.RabLog.Errorf("Error parsing time:", err)
		return false
	}

	// 比较两个时间
	if parsedTime1.After(parsedTime2) {
		return true
	} 
	return false
}

func getLastUrl(url, lastTime string) (string, string, error) {
	updateInfo := ""
	method := "GET"

	client := &http.Client {}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		log.RabLog.Error(err)
		return "", "", err
	}
	req.Header.Add("cookie", "{cookie}")

	res, err := client.Do(req)
	if err != nil {
		log.RabLog.Error(err)
		return "", "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.RabLog.Error(err)
		return "", "", err
	}
  	// 使用gofeed解析RSS订阅的内容
	parser := gofeed.NewParser()
	feed, err := parser.Parse(strings.NewReader(string(body)))
	if err != nil {
		log.RabLog.Errorf("解析RSS订阅时发生错误:", err)
		return "", "", err
	}

	updateInfoCross := ""
	for _, item := range feed.Items {
		if (lastTime != "" && timeAfter(lastTime, item.Published)) {
			break
		}
		if (lastTime == item.Published) {
			break
		}
		if (url == "https://{siteurl}/index.php/feed/comments/") {
			if (strings.Contains(item.Link, "cross.html")) {
				updateInfoCross += fmt.Sprintf("%s👀作者：%s\n📝内容：%s\n🔗链接：%s\n📆发布时间：\n%s\n", common.Dilimiter, item.Title, item.Description, item.Link, parseTime(item.Published))
			} else {
				if (updateInfo == "") {
					updateInfo += "🤣小站上有新评论！\n"
				}
				updateInfo += fmt.Sprintf("%s👀作者：%s\n📝内容：%s\n🔗链接：%s\n📆发布时间：\n%s\n", common.Dilimiter, item.Title, item.Description, item.Link, parseTime(item.Published))
			}
		} else {
			updateInfo += fmt.Sprintf("%s👀标题：%s\n📝摘要：%s\n🔗链接：%s\n📆发布时间：\n%s\n", common.Dilimiter, item.Title, item.Description, item.Link, parseTime(item.Published))
		}
	}
	if (updateInfoCross != "") {
		if (updateInfo != "") {
			updateInfo += common.Dilimiter
		}

		updateInfo += "😊时光机上新啦\n" + updateInfoCross
	}

	return feed.Items[0].Published, updateInfo, nil
}

func sendUpdate(post, comment string) {
	self, _ := bot.GetCurrentUser()
	groups, _ := self.Groups()
	group := groups.GetByNickName("咱们仨把日子过好比啥都强")

	responseTxt := ""
	if post != "" {
		responseTxt += "😇小站上有新文章！\n" + post
	}
	if comment != "" {
		responseTxt += comment
	}

	if _, err := group.SendText(responseTxt); err != nil {
		log.RabLog.Errorf("发送信息失败: %v", err)
	}
}

func saveCache() {
	postCache, postUpdata, err := getLastUrl("https://{site_url}/index.php/feed/", urlCache.PostCache)
	if err != nil {
		postUpdata = "获取文章信息失败"
	}
	commentCache, commentUpdate, err := getLastUrl("https://{site_url}/index.php/feed/comments/", urlCache.CommentCache)
	if err != nil {
		commentUpdate = "获取评论信息失败"
	}

	if (urlCache.PostCache == "" || urlCache.CommentCache == "" || timeAfter(postCache, urlCache.PostCache) || timeAfter(commentCache, urlCache.CommentCache)) {
		urlCache.PostCache = postCache
		urlCache.CommentCache = commentCache
		urlCacheMap["postCache"] = urlCache.PostCache
		urlCacheMap["commentCache"] = urlCache.CommentCache
	
		saveJSON(urlCacheMap, common.SPushPath)

		sendUpdate(postUpdata, commentUpdate)
	} else {
		log.RabLog.Debugf("网站没有更新")
	}
}

func getCache() error {
	file, err := os.Open(common.SPushPath)
	if err != nil {
		return  err
	}
	defer file.Close()

	jsonData, err := ioutil.ReadAll(file)
	if err != nil {
		return  err
	}

	err = json.Unmarshal(jsonData, &urlCache)
	if (err != nil) {
		return err
	}
	return nil
}


func SPushEntry(bot_param *openwechat.Bot) {
	bot = bot_param
	if _, err := os.Stat(common.SPushPath); err == nil {
		getCache()
	}
	for {
		log.RabLog.Debugf("开始轮询")
		saveCache()
		time.Sleep(1 * time.Minute)
	}
}