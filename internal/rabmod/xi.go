package rabmod

// import (
// 	"os"
// 	"fmt"
// 	"strings"
// 	"io/ioutil"
// 	"encoding/json"

// 	"github.com/gocolly/colly"
	
// 	"rabbot/internal/log"
// 	"rabbot/internal/common"
// )

// // 网页中提取到的表格数据
// // 数据来源：https://steamstats.cn/xi
// type TableData struct {
//     Headers []string
//     Rows    [][]string
// 	Urlmap     map[int]string // 主要用于保存游戏的URL，因为在游戏名和操作都会有url，逻辑中是覆盖逻辑，操作下面的url是最终取到的
// }

// var gbMap map[string][]map[string]string = make(map[string][]map[string]string)

// func getList() error {
// 	// 创建一个 Collector 对象
// 	c := colly.NewCollector()

// 	// 注册回调函数，处理表头
// 	c.OnHTML("div.container[data-v-5071b21d]", func(tables *colly.HTMLElement) {
// 		tables.ForEach("table", func(_ int, table *colly.HTMLElement) {	
// 			headTxt := ""
// 			tableData := TableData{}
// 			tableData.Urlmap = make(map[int]string)
// 			// headTxt是大表头，当前仅有两个：1. 当前限免游戏 2. 即将限免游戏
// 			table.ForEach("thead tr th h2", func(_ int, th *colly.HTMLElement) {
// 				headTxt = strings.TrimSpace(th.Text)
// 			})

// 			table.ForEach("thead tr th", func(_ int, th *colly.HTMLElement) {
// 				headTxtNow := strings.TrimSpace(th.Text)
// 				if headTxtNow != headTxt {
// 					tableData.Headers = append(tableData.Headers, strings.TrimSpace(th.Text))
// 				}
// 			})
	
// 			// 提取表格内容
// 			table.ForEach("tbody tr", func(_ int, tr *colly.HTMLElement) {
// 				var row []string
// 				tr.ForEach("td", func(_ int, td *colly.HTMLElement) {
// 					td.ForEach("a", func(_ int, a *colly.HTMLElement) {
// 						if (a.Attr("href") != "") {
// 							tableData.Urlmap[len(tableData.Rows)] = a.Attr("href")
// 						}
// 					})
// 					row = append(row, strings.TrimSpace(td.Text))
// 				})
// 				tableData.Rows = append(tableData.Rows, row)
// 			})

// 			// 按照下表把表头和内容绑定在一起后放到全局的map中，后续转成json用
// 			for indexN, table := range tableData.Rows {
// 				var tmpMap map[string]string = make(map[string]string)
// 				for index, head := range table {
// 					tmpMap[tableData.Headers[index]] = head
// 				}
// 				tmpMap["url"] = tableData.Urlmap[indexN]
// 				gbMap[headTxt] = append(gbMap[headTxt], tmpMap)
// 			}
			
// 		})
// 	})

// 	// 访问目标网页
// 	err := c.Visit("https://steamstats.cn/xi")
// 	if err != nil {
// 		log.RabLog.Errorf("Visit steamstats failed: %v", err)
// 		return err
// 	}
// 	return nil
// }

// func UpdateXiInfo() error {
// 	err := getList()
// 	if err != nil {
// 		return err
// 	}

// 	jsonData, err := json.MarshalIndent(gbMap, "", "  ")
// 	if err != nil {
// 		log.RabLog.Errorf("JSON encoding error: %v", err)
// 		return err
// 	}

// 	// 爬取到的json数据，保存到本地tmp目录下，每天凌晨定时清理
// 	err = ioutil.WriteFile(common.XiJsonFile, jsonData, 0644)
// 	if err != nil {
// 		log.RabLog.Errorf("Error writing JSON file: %v", err)
// 		return err
// 	}
	
// 	log.RabLog.Infof("Update xi info success: %v", err)
// 	return nil
// }

// // 导出请求喜加一数据接口
// func init() {
// 	common.FuncNameMap["XiPlusOne"] = XiPlusOne
// }

// // 游戏信息，格式参考/rabdata/tmp/xi.json
// type gameInfo struct {
// 	StartTime string `json:"start time"`
// 	EndTime string `json:"end time"`
// 	Name string `json:"name"`
// 	Url string `json:"url"`
// }

// // 所有游戏信息
// type allGmaeInfo struct {
// 	Currently    []gameInfo `json:"Currently live promotions"`
// 	Upcoming     []gameInfo `json:"Upcoming promotions"`
// }

// func XiPlusOne(requestStruct *common.RequestStruct) (*common.ReplyStruct, error) {
// 	if _, err := os.Stat(common.XiJsonFile); err == nil {
// 		log.RabLog.Debug("File cache exist")
// 	} else {
// 		if err := UpdateXiInfo(); err != nil {
// 			log.RabLog.Errorf("UpdateXiInfo failed, %v", err)
// 			return nil, err
// 		}
// 	}

// 	// 读取喜加一json文件
// 	file, err := os.Open(common.XiJsonFile)
// 	if err != nil {
// 		log.RabLog.Errorf("Open xi json file failed, %v", err)
// 		return nil, err
// 	}
// 	defer file.Close()

// 	jsonData, err := ioutil.ReadAll(file)
// 	if err != nil {
// 		log.RabLog.Errorf("Read xi json file failed, %v", err)
// 		return nil, err
// 	}

// 	var gmaeData allGmaeInfo
// 	err = json.Unmarshal(jsonData, &gmaeData)
// 	if (err != nil) {
// 		log.RabLog.Errorf("Translate xi json file content to json failed, %v", err)
// 		return nil, err
// 	}

// 	xiTxt := "早买早享受，晚买有折扣，不买🆓免费送\n当前限免🎮：\n"
// 	xiTxt += common.Dilimiter
// 	for _, game := range gmaeData.Currently {
// 		xiTxt += fmt.Sprintf(common.XiGameStr, game.Name, game.StartTime, game.EndTime, game.Url)
// 		xiTxt += common.Dilimiter
// 	}

// 	xiTxt += "即将限免🎮：\n"
// 	xiTxt += common.Dilimiter
// 	for _, game := range gmaeData.Upcoming {
// 		xiTxt += fmt.Sprintf(common.XiGameStr, game.Name, game.StartTime, game.EndTime, game.Url)
// 		xiTxt += common.Dilimiter
// 	}

// 	return &common.ReplyStruct{common.MsgTxt, xiTxt}, nil
// }