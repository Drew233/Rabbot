package rabmod

import (
	"encoding/json"
	"fmt"
	"time"
	"io/ioutil"
	"net/http"
	"os"
	
	"rabbot/internal/log"
	"rabbot/internal/common"
)

type dataStr struct {
	Data struct {
		Catalog struct {
			SearchStore struct {
				Elements []gameInfo `json:"elements"`
			} `json:"searchStore"`
		} `json:"Catalog"`
	} `json:data`
}

type PromotionsItem struct {
	PromotionalOffers []struct {
		StartDate string `json:"startDate"`
		EndDate string `json:"endDate"`
	} `json: "promotionalOffers"`
}

type Mapping struct {
	PageSlug string `json:"pageSlug"`
	PageType string `json:"pageType"`
}

type CatalogNs struct {
	Mappings []Mapping `json:"mappings"`
}

type gameInfo struct {
	ProductSlug string `json:"productSlug"`
	Promotions struct {
		PromotionalOffers []PromotionsItem `json:"promotionalOffers"`
		UpcomingPromotionalOffers []PromotionsItem `json:"upcomingPromotionalOffers"`
	}
	Title string `json:"title"`
	Price struct {
		TotalPrice struct {
			DiscountPrice int `json:"discountPrice"`
		} `json: totalPrice`
	} `json: "price`
	Url string
	CatalogNs CatalogNs `json:"catalogNs"`
}

// 导出请求喜加一数据接口
func init() {
	common.FuncNameMap["XiPlusOne"] = XiPlusOne
}

func getTime(timeStr string) string {
	inputTime, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return ""
	}

	return inputTime.Format("2006-01-02")
}

func updateXiInfo() error {
	url := "https://store-site-backend-static.ak.epicgames.com/freeGamesPromotions"

	// 发送HTTP请求
	response, err := http.Get(url)
	if err != nil {
		log.RabLog.Errorf("Error making HTTP request:", err)
		return err
	}
	defer response.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.RabLog.Errorf("Error reading response body:", err)
		return err
	}

	// 解析JSON
	var data dataStr
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.RabLog.Errorf("Error decoding JSON:", err)
		return err
	}
	
	return saveJSON(data, common.XiJsonFile)
}

func XiPlusOne(requestStruct *common.RequestStruct) (*common.ReplyStruct, error) {

	if _, err := os.Stat(common.XiJsonFile); err == nil {
		log.RabLog.Debug("File cache exist")
	} else {
		if err := updateXiInfo(); err != nil {
			log.RabLog.Errorf("UpdateXiInfo failed, %v", err)
			return nil, err
		}
	}

	// 读取喜加一json文件
	file, err := os.Open(common.XiJsonFile)
	if err != nil {
		log.RabLog.Errorf("Open xi json file failed, %v", err)
		return nil, err
	}
	defer file.Close()

	jsonData, err := ioutil.ReadAll(file)
	if err != nil {
		log.RabLog.Errorf("Read xi json file failed, %v", err)
		return nil, err
	}

	var data dataStr
	err = json.Unmarshal(jsonData, &data)
	if (err != nil) {
		log.RabLog.Errorf("Translate xi json file content to json failed, %v", err)
		return nil, err
	}

	xiStr, upXiStr := "", ""
	for _, gameInfo := range data.Data.Catalog.SearchStore.Elements {
		var proItem PromotionsItem
		var gameUrl = ""
		if (gameInfo.ProductSlug != "") {
			gameUrl = "https://epicgames.com/store/product/" + gameInfo.ProductSlug
		} else if (len(gameInfo.CatalogNs.Mappings) > 0 && gameInfo.CatalogNs.Mappings[0].PageSlug != "") {
			gameUrl = "https://store.epicgames.com/zh-CN/p/" + gameInfo.CatalogNs.Mappings[0].PageSlug
		} else {
			gameUrl = "阿哦，小兔子找不到链接，自己上去看看呢？"
		}
		if len(gameInfo.Promotions.PromotionalOffers) > 0 && len(gameInfo.Promotions.PromotionalOffers[0].PromotionalOffers) > 0 {
			proItem = gameInfo.Promotions.PromotionalOffers[0]
			xiStr += fmt.Sprintf(common.XiGameStr, gameInfo.Title, getTime(proItem.PromotionalOffers[0].StartDate), getTime(proItem.PromotionalOffers[0].EndDate), gameUrl) + common.Dilimiter
		} else if len(gameInfo.Promotions.UpcomingPromotionalOffers) > 0 && len(gameInfo.Promotions.UpcomingPromotionalOffers[0].PromotionalOffers) > 0 {
			proItem = gameInfo.Promotions.UpcomingPromotionalOffers[0]
			upXiStr += fmt.Sprintf(common.XiGameStr, gameInfo.Title, getTime(proItem.PromotionalOffers[0].StartDate), getTime(proItem.PromotionalOffers[0].EndDate), gameUrl) + common.Dilimiter
		} else {
			continue
		}
	}

	if (upXiStr == "") {
		upXiStr = "啊哦，小兔子也找不到有什么免费游戏了，再等等咯"
	}
	
	replyStr := "早买早享受，晚买有折扣，不买🆓免费送\nEpic当前限免🎮：\n" + common.Dilimiter + xiStr + "Epic即将限免🎮：\n" + common.Dilimiter + upXiStr

	steamStr, err := GetSXiInfo()
	if err == nil {
		replyStr += steamStr
	}

	return &common.ReplyStruct{common.MsgTxt, replyStr}, nil
}