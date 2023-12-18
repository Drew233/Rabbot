package rabmod

import (
	"reflect"

	"rabbot/config"
	"rabbot/internal/log"
	"rabbot/internal/common"
)

var FuncMap map[string]reflect.Value = make(map[string]reflect.Value)

func ModInit() {
	for feaName, feaStruct := range config.RabConfig.Features {
		log.RabLog.Info(feaStruct.Entry)
		
		feaFunc := reflect.ValueOf(common.FuncNameMap[feaStruct.Entry])

		if feaFunc.IsValid() && feaFunc.Kind() == reflect.Func {
			FuncMap[feaName] = feaFunc
		} else {
			log.RabLog.Errorf("Feature %s load failed of feature function %s invalid", feaName, feaStruct.Entry)
		}
	}
	log.RabLog.Infof("Rabmod load finish, now funcmap is %v", FuncMap)
}