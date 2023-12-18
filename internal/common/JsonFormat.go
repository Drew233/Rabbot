package common
// 该文件中定义整个项目中用到的json格式

// 抽签API
var DrlotUrl = "https://api.jcwle.com/api/chouqian?&type=%s&apiKey=0452574d625d61135116a710c0eea03e"
// 不同类型的抽签，根据随机数选择随机类型
var DrlotTypes = [...](string){"guanyin", "yuelao", "mazu", "lvzu"}
/*
返回数据格式
{"code":200,"msg":"查询成功","data":[{"name":"您抽出观音灵签第","text":" 55 签"},{"name":"签曰：","text":"周武王登位"},{"name":"吉凶：","text":"中签"},{"name":"宫位：","text":"丑宫"},{"name":"签诗：","text":"父贤传子子传孙，衣食丰隆只靠天；堂上椿萱人快乐，饥饭渴饮困时眠。"},{"name":"签语：","text":"此卦接竹引泉之象，凡事谋望大吉也。"},{"name":"解签：","text":"接竹引泉，流传不绝，君子谋望，无不欣悦。"},{"name":"仙机：","text":"家宅→欠安　自身→还愿　求财→利　交易→成　婚姻→合　六甲→平安　行人→动　田蚕→吉　六畜→损　寻人→难　公讼→有贵人　移徙→吉　失物→见　疾病→禳星　山坟→安"},{"name":"详解：","text":"世代传承子孙优秀贤能，衣食无缺富贵在天;金玉满堂无忧无虑自然舒爽，饿了就吃，倦了就眠的生活必然适意。祖宗积德，福禄后昆，荣禄并耀，光裕门庭。此签接竹引泉之象，凡事着谋吉利。本签示于弟子曰。子孙贤。衣禄丰盈富在天。君尔之命。即是富贵天定者。大吉大利之签者。功名。交易。婚姻。求财皆如意者。而且吃饭困时眠之运。逢此大吉利之时。宜多修德。积善。更可发扬光大者。此签有”失控惹祸”之意。提醒当事人，遇事不宜意气用事。有时人与人之间难免意见不合而发生争执口角。须知冲动于事无补，无论对方是否有错在先，但如果因此而忍不下一口气、非要争个理出来，却可能导致撕破脸不相往来、见面又尴尬的局面，甚至拖累到不相干的人替你收拾残局。这都是因为欠缺理智且没有顾虑到后果所造成的影响。俗云：”冤家宜解不宜结”。与其吵得面红耳赤闹得不可开交，还不如冷静下来，想个好方法，既不会伤害到对方，又可以避免让事情重蹈覆辙。心平气和解决问题，而非制造问题。"},{"name":"","text":"1.测算结果若是理想，固然是一件可喜可贺的事情，如果测算结果不理想，缘主也不必灰心。缘主的命运，即便算得再准，也还是需要缘主自己去把握。算命的目的是为了趋吉避凶，顺势而行。"},{"name":"","text":"2.正所谓一命二运三风水，四积阴德五读书，六名七相八敬神，九交贵人十修身。平时积善行德，心存善念，必有善行，善念善行，天必佑之。"}]}
这里只关注一部分，详见下面的NameMap
*/
var DrlotNameMap = map[string]string{
	"签曰：": "",
	"签诗：": "",
	"签语：": "",
	"诗曰：": "",
	"解签：": Dilimiter,
	"解曰：": Dilimiter,
}
// 抽签API接口返回数据格式
type DrlotsData struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		Name string `json:"name"`
		Text string `json:"text"`
	} `json:"data"`
}

// 摸鱼日历API
var CalenderUrl = "https://api.vvhan.com/api/moyu?type=json"
/*
返回数据格式
{"success":true,"url":"https:\/\/web-static.4ce.cn\/storage\/bucket\/v1\/2532bc63e3266c483d6ecf52be96175a.jpg"}
*/
// 摸鱼日历API接口返回数据格式
type Cal_data struct {
	Success bool `json:"success`
	Url string `json:"url`
}

// 色图API
// https://img.jitsu.top/#/
var SetuUrl = "https://moe.jitsu.top/img/?sort=setu&type=json"
/*
返回数据格式
{
    "code": 200,
    "alert": "别整天搁那儿爬来爬去，小孩子吗？API不是做给你爬的，我希望API能发挥它本身的作用，爬点涩图对你有什么好处？",
    "pics": [
        "https://pic.rmb.bdstatic.com/bjh/497b235dec5329bbf3eee43cfe539975.jpeg"
    ]
}
*/
type SetuData struct {
	Code int `json:"code"`
	Pics []string `json:"pics"`
}