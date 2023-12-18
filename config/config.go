package config

import (
	"encoding/json"
	"os"
	_log "log"
)

type featureStruct struct {
	Enable bool `json:"enable"`
	Entry string `json:"entry"`
	FeatureGpBlist map[string]bool `json:"groupBlackList"`
}

// 定义配置文件格式
type RabConfigStruct struct {
	BotName string `json:"botName"`
	Datadir string `json:"Datadir"`
	DefaultMsg struct {
		DullMsg string `json:"dullMsg"`
		ErrMsg string `json:"ErrorMsg"`
	} `json:"defaultMsg"`
	GroupWhiteList []string `json:"groupWhiteList"`
	Cron struct {
		TmpCleanCron string `json:"tmpCleanCron"`
	}
	RabLogConfig struct {
		Maxsize int `json:"maxsize"`
		Maxbackups int `json:"maxbackups"`
		Maxage int `json:"maxage"`
		Compress bool `json:"compress"`
	} `json:"rablog"`
	// Features map[string]bool `json:"features"`
	Features map[string]featureStruct `json:"feature"`
	MysqlConfig struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Ip string `json:"ip"`
		Port int `json:"port"`
		Dbname string `json:"chouqian"`
		Charset string `json:"charset"`
	} `json:"mysql"`
	WbList map[string]bool
}

var RabConfig *RabConfigStruct
var ConfigPath = "./config/rabbot.config"

// 加载全局配置到RabConfig
func LoadConfig(cfgPath string) {
	configCont, err := os.ReadFile(cfgPath);
	if err != nil {
		_log.Fatalln("Read config file error, file_path is " + cfgPath)
	}

	var _rabConfig = new(RabConfigStruct)
	if err = json.Unmarshal(configCont, _rabConfig); err != nil {
		_log.Fatalln("Load config file error, file_path is " + cfgPath, err)
	}

	RabConfig = _rabConfig
	RabConfig.WbList = map[string]bool{}

	for _, groupName := range _rabConfig.GroupWhiteList {
		RabConfig.WbList[groupName] = true
	}
}