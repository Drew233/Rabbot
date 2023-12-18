/* 通用变量 */
package common

var (
	Dilimiter = "---------\n"	// 分隔符
	DataDir = "./rabdata"		// 数据目录
	TmpDir = DataDir + "/tmp"   // 临时目录
	LogDir = DataDir + "/log"	// 日志目录
	PicDir = DataDir + "/pic"	// 图片目录
	LogFilename = LogDir + "/rabbot.log"  // 日志文件名
	DebugFlag = TmpDir + "/RabDbg"  // 调试标记
	FeatureDisabled = "嘿，您猜怎么着，我有%s的功能，但就是不给你用"  // 功能未启用提示语
	UnknownReply = "诶呀呀，你这是什么问题？我才不要告诉你哼"
	DefaultReply = map[string]string {
		"我吃柠檬": "兔兔那么可爱，怎么可以吃兔兔",
		"": "你没事吧？没事的话我建议你去玩会游戏",
	}							// 默认对话
	FuncNameMap map[string]interface{} // 模块函数
)

func init() {
	FuncNameMap = make(map[string]interface{})
}