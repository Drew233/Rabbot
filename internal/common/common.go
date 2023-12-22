/* 通用变量 */
package common

var (
	Dilimiter = "---------\n"	// 分隔符
	DataDir = "./rabdata"		// 数据目录
	TmpDir = DataDir + "/tmp"   // 临时目录
	LogDir = DataDir + "/log"	// 日志目录
	PicDir = DataDir + "/pic"	// 图片目录
	LogFilename = LogDir + "/rabbot.log"  // 日志文件名
	XiJsonFile = TmpDir + "/xi.json"	// 喜加一数据缓存文件名
	DebugFlag = TmpDir + "/RabDbg"  // 调试标记
	FeatureDisabled = "嘿，您猜怎么着，我有%s的功能，但就是不给你用"  // 功能未启用提示语
	UnknownReply = "诶呀呀，你这是什么问题？我才不要告诉你哼"
	DefaultReply = map[string]string {
		"我吃柠檬": "兔兔那么可爱，怎么可以吃兔兔",
		"": "你没事吧？没事的话我建议你去玩会游戏",
	}							// 默认对话
	FuncNameMap map[string]interface{} // 模块函数
	UseOfHorseComp = "想玩赛马小游戏吗？\n加入指令：@{bot} 加入赛马\n开始指令：@{bot} 比赛开始\n注意：\n1. 每场比赛的参与人数最少三人\n2. 同时只能存在一场赛马\n3. 五分钟内没有赛马的信息交互，会自动删除比赛\n4. 当前马厩中有%d匹“马”，每场比赛最多%d人参加\n马厩：%v"
	HorseCompNotCreated = "当前还没有创建比赛，请@我并输入加入比赛来加入赛马"
	HorseCompCreateSuccess = "创建赛马比赛成功，当前参与人数: 1"
	HorseCompJoinSuccess = "加入比赛成功，当前参与人数: %d"
	HorseCompRunning = "赛马正在进行，请耐心等待比赛结束"
	HorseCompNotEnough = "一场比赛至少有三人参加"
	HorseCompTooMuch = "本场比赛参与人数已达上限，比赛马上开始"
	HorseCompJoinEd = "您已经加入了本场比赛，请勿重复参加"
	XiGameStr = "🕹游戏名：%s\n🕛️开始时间：%s\n🕛️结束时间：%s\n🔗领取链接：%s\n"
)

func init() {
	FuncNameMap = make(map[string]interface{})
}