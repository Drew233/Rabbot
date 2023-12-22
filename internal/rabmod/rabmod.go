package rabmod

import (
	"reflect"

	"rabbot/config"
	"rabbot/internal/log"
	"rabbot/internal/common"
)

var FuncMap map[string]reflect.Value = make(map[string]reflect.Value)

func ModInit() {
	// 把内置的指令加载到配置中
	for commond, entry := range common.InternalFuncMap {
		config.RabConfig.Features[commond] = config.FeatureStruct{
			Enable: true,
			Entry: entry,
			FeatureGpBlist: map[string]bool{},
		}
	}
	// 加载配置文件中的模块
	for feaName, feaStruct := range config.RabConfig.Features {
		feaFunc := reflect.ValueOf(common.FuncNameMap[feaStruct.Entry])

		if feaFunc.IsValid() && feaFunc.Kind() == reflect.Func {
			FuncMap[feaName] = feaFunc
		} else {
			log.RabLog.Errorf("Feature %s load failed of feature function %s invalid", feaName, feaStruct.Entry)
		}
	}
	log.RabLog.Infof("Rabmod load finish, now funcmap is %v", FuncMap)
}