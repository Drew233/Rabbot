package rabmod

import (
	"encoding/json"
	"os"
	"strings"
	"fmt"
	"io/ioutil"

	"github.com/gocolly/colly"
	
	"rabbot/internal/log"
	"rabbot/internal/common"
)

type Game struct {
	Title       string `json:"title"`
	Price       string `json:"price"`
	Link        string `json:"link"`
	Summary     string `json:"summary"`
}

func updateSteamInfo() error {
	// 创建Colly收集器
	c := colly.NewCollector()

	// 设置User-Agent
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
		r.Headers.Set("Accept-Language", "zh-CN,zh;")
	})

	// 创建一个切片，用于保存游戏信息
	var games []Game
	// 注册回调函数，提取游戏信息
	c.OnHTML(".search_result_row", func(e *colly.HTMLElement) {

		title := e.ChildText(".search_name")
		price := e.ChildText(".discount_final_price")
		link := strings.Split(e.Attr("href"), "?")[0]
		summary := strings.Replace(e.ChildAttr(".search_review_summary", "data-tooltip-html"), "<br>", ",", -1)
	
		// 创建游戏结构体并添加到切片中
		game := Game{
			Title:       title,
			Price:       price,
			Link:        link,
			Summary: summary,
		}
		games = append(games, game)
	})

	// 错误处理回调函数
	c.OnError(func(r *colly.Response, err error) {
		log.RabLog.Errorf("Request steam free game failed, err: %v", err)
		return 
	})

	// 启动爬虫
	err := c.Visit("https://store.steampowered.com/search/?maxprice=free&specials=1")
	if err != nil {
		log.RabLog.Errorf("Request steam free game failed, err: %v", err)
		return err
	}

	// 将游戏信息保存为JSON文件
	file, err := os.Create(common.XiSJsonFile)
	if err != nil {
		log.RabLog.Errorf("Open file %s failed, err: %v", common.XiSJsonFile, err)
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(games)
	if err != nil {
		log.RabLog.Errorf("Encoder json failed, err: %v", common.XiSJsonFile, err)
		return err
	}

	return nil
}

func GetSXiInfo() (string, error) {
	if _, err := os.Stat(common.XiSJsonFile); err == nil {
		log.RabLog.Debug("File cache exist")
	} else {
		if err := updateSteamInfo(); err != nil {
			log.RabLog.Errorf("UpdateXiSteamInfo failed, %v", err)
			return "", err
		}
	}

	file, err := os.Open(common.XiSJsonFile)
	if err != nil {
		log.RabLog.Errorf("Open steam xi json file failed, %v", err)
		return "", err
	}
	defer file.Close()
	jsonData, err := ioutil.ReadAll(file)
	if err != nil {
		log.RabLog.Errorf("Read steam xi json file failed, %v", err)
		return "", err
	}

	var data []Game
	err = json.Unmarshal(jsonData, &data)
	if (err != nil) {
		log.RabLog.Errorf("Translate steam xi json file content to json failed, %v", err)
		return "", err
	}

	str := "Steam当前限免🎮："
	for _, value := range data {
		str += "\n" + common.Dilimiter + fmt.Sprintf("🕹游戏名：%s\n💰参考价格：%s\n🗣️历史评价：%s\n🔗领取链接：%s", value.Title, value.Price, value.Summary, value.Link)
	}
	return str, nil
}