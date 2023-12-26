package rabmod

import (
	"encoding/json"
	"log"
	"os"
	"strings"

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

func updateSteamInfo() {
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
		summary := strings.Replace(e.ChildAttr(".search_review_summary", "data-tooltip-html"), "<br>", "\n", -1)
	
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
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	// 启动爬虫
	err := c.Visit("https://store.steampowered.com/search/?maxprice=free&specials=1")
	if err != nil {
		log.Fatal(err)
	}

	// 将游戏信息保存为JSON文件
	file, err := os.Create("games.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(games)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Games data saved to games.json")
}