package rabmod

import (
	"fmt"
	"sort"
	"time"
	"errors"
	"strings"
	"math/rand"

	"github.com/eatmoreapple/openwechat"

	"rabbot/internal/log"
	"rabbot/internal/common"
)

func init() {
	common.FuncNameMap["HorseComp"] = HorseComp
	// 把比赛开始也视为一种commond，直接响应赛马的比赛开始
	common.InternalFuncMap["比赛开始"] = "HorseComp"
	// 注册赛马帮助指令
	common.InternalFuncMap["赛马"] = "HorseComp"
}

// 赛马用户结构体
type horseBase struct {
	uname string		// 用户名
	step int			// 每次的步数，覆盖保存
	pos int				// 马当前所处位置
	horseEmoji string	// 马emoji
}

// 赛马全局结构，以群名为key保存在map中
type groupContentBase struct {
	status int
	// status: 群比赛状态
	// 0：群第一次使用赛马功能，需要初始化
	// 1: 还未举办比赛，也未进行比赛
	// 2：举办了比赛，但还不够人数开始
	// 3: 比赛正在进行
	timeStamp int64 // 比赛时间戳，用于判断是否超时
	horses map[string]horseBase  // 赛马用户，以uuid为key
}

var compLength = 20 // 赛道长度
var emojiSet = []string{"🐢", "🎠", "🦓", "🦛", "🐎", "🐴", "🦄", "🚽", "🏃‍♂️", "🏃‍♀️", "🦍", "🦇"} // 选手emoji，随机选取 ！！如果修改的话，至少剩三只马 ！！
var groupEmojiSet map[string][]string = make(map[string][]string)	// 群名为key，emojiset为value，主要用于每个群赛马时保证代表用户的emoji不重复

var groupContent map[string]groupContentBase = make(map[string]groupContentBase)	// 赛马全局结构

// 定时清理对话记录
// 默认超时时间五分钟，五分钟定时检查一次
func CleanOuttimeComp() {
	log.RabLog.Debugf("begin cleanouttimeComp")
	for key, groupDetail := range groupContent {
		if groupContent[key].timeStamp == 0 {
			continue
		}
		timeNow := time.Now().Unix()
		if timeNow > groupDetail.timeStamp && (timeNow - groupDetail.timeStamp) > 60 * 5 {
			log.RabLog.Infof("Group %s horsecomp timeout,clean success", key)
			groupContent[key] = groupContentBase{}
			groupEmojiSet[key] = []string{}
		}
	}
}

// 获取随机的emoji用来代表用户，同一场比赛中每个用户的emoji不会重复
func getHorseEmoji(groupname string) string {
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(groupEmojiSet[groupname]))
	res := groupEmojiSet[groupname][index]
	groupEmojiSet[groupname] = append(groupEmojiSet[groupname][:index], groupEmojiSet[groupname][index + 1:]...)
	return res
}

// 返回当前所有马的位置到群聊中
func replyHorsesPos(msg *openwechat.Message, groupname string) {
	str := ""

	keys := make([]string, 0, len(groupContent[groupname].horses))
	for key := range groupContent[groupname].horses {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		horse := groupContent[groupname].horses[key]
		str += fmt.Sprintf("%s: \n", horse.uname)
		flag := 0
		horseRoad := ""
		for i := 0; i < compLength; i++ {
			if i == horse.pos {
				horseRoad += "H"
				flag = 1
			} else if i == compLength - 1{
				if flag == 0 {
					horseRoad += "H"	// 以H代表马
				} else {
					horseRoad += "_"
				}
			} else {
				horseRoad += "_"
			}
		}

		// 因为Reverse会导致部分emoji格式发生变化，所以先以字母代替然后再替换成emoji
		horseRoad = strings.Replace(common.ReverseString(horseRoad), "H", horse.horseEmoji, -1)
		// 设置起点和终点的emoji
		str += "🚩" + horseRoad + "🏁" + "\n"
	}

	msg.ReplyText(str)
}

// 开始赛马
func beginHorseComp(msg *openwechat.Message, groupname string) {
	msg.ReplyText("比赛开始！")
	// 比赛开始时先响应一次初始位置
	replyHorsesPos(msg, groupname)
	// 间隔两秒，模拟马跑的时间
	time.Sleep(2 * time.Second)

	// 用户信息
	type userInfo struct {
		uname string
		horse string
	}

	// 获胜用户信息，内置只评选前三名
	type winnerStruct struct {
		first userInfo
		second userInfo
		third userInfo
	}

	winner := winnerStruct{userInfo{}, userInfo{}, userInfo{}}
	winnerAll := 0
	horseSpeed := 8

	for {
		for uuid, horse := range groupContent[groupname].horses {
			// 位置为-1的已经到终点了，不需要继续走
			if horse.pos == -1 {
				continue
			}
			rand.Seed(time.Now().UnixNano())
			step := rand.Intn(horseSpeed) + 1 // 每次随即走1到8步
			pos := horse.pos + step
			groupContent[groupname].horses[uuid] = horseBase{
				uname: horse.uname,
				step: step,
				pos: pos,
				horseEmoji: horse.horseEmoji,
			}
		}
		// 马也会累，速度递减，防止最后冲线时候步幅太大
		if horseSpeed != 2 {
			// 速度最少减到2
			horseSpeed = horseSpeed - 1
		}
		
		for uuid, horse := range groupContent[groupname].horses {
			// 按照走完全程的顺序决定名次
			if horse.pos >= compLength {
				groupContent[groupname].horses[uuid] = horseBase{
					uname: horse.uname,
					step: horse.step,
					pos: -1,
					horseEmoji: horse.horseEmoji,
				}
				if winner.first.uname == "" {
					winner = winnerStruct{userInfo{horse.uname, horse.horseEmoji}, winner.second, winner.third}
				} else if winner.second.uname == "" {
					winner = winnerStruct{winner.first, userInfo{horse.uname, horse.horseEmoji}, winner.third}
				} else {
					winner = winnerStruct{winner.first, winner.second, userInfo{horse.uname, horse.horseEmoji}}
					// 决出前三名的话就不用再跑下去了
					winnerAll = 1
				}
			} 
		}
		// 每次跑完响应此时所有马的位置
		replyHorsesPos(msg, groupname)

		if winnerAll == 1 {
			winnerMsg := fmt.Sprintf("比赛结束！\n🥇状元：%s\t%s\n🥈榜眼：%s\t%s\n🥉探花：%s\t%s", winner.first.uname, winner.first.horse,
																								  winner.second.uname, winner.second.horse,
																								  winner.third.uname, winner.third.horse)
			podiumStr := fmt.Sprintf("\n领奖台：\n              __%s__       \n__%s__|             |__%s__\n|____________________|", winner.first.horse, winner.second.horse, winner.third.horse)
			msg.ReplyText(winnerMsg + podiumStr)

			// 比赛结束，重置此群中的赛马结构信息
			groupContent[groupname] = groupContentBase{}
			groupEmojiSet[groupname] = []string{}
			return 
		}
		// 间隔两秒，模拟马跑的时间
		time.Sleep(2 * time.Second)
	}
}

// 赛马入口，赛马使用教程详见rabbot/internal/common中的UseOfHorseComp
// 以群为单位
func HorseComp(requestStruct *common.RequestStruct) (*common.ReplyStruct, error) {
	groupname := requestStruct.Groupname
	// status为0表示需要初始化
	if groupContent[groupname].status == 0 {
		// 初始化赛马结构
		groupContent[groupname] = groupContentBase{
			status: 1,
			timeStamp: time.Now().Unix(),
			horses: make(map[string]horseBase),
		}

		// 深拷贝初始化emojiset
		emojiSetT := make([]string, len(emojiSet))
		copy(emojiSetT, emojiSet)
		groupEmojiSet[groupname] = emojiSetT
	}

	// status为3表示赛马正在进行，不响应指令
	if groupContent[groupname].status == 3 {
		return &common.ReplyStruct{common.MsgTxt, common.HorseCompRunning}, nil
	}

	// 更新赛马结构时间戳
	groupContent[groupname] = groupContentBase{
		status: groupContent[groupname].status,
		timeStamp: time.Now().Unix(),
		horses: groupContent[groupname].horses,
	}

	// 从请求结构中取出需要的字段
	uname, uuid, requestText, msg := requestStruct.Uname, requestStruct.Uuid, requestStruct.RequestTxt, requestStruct.Msg
	horseNum := len(groupContent[groupname].horses)

	if requestStruct.Commond == "比赛开始" {
		requestText = "比赛开始"
	} else if requestStruct.Commond == "赛马" {
		requestText = "玩法"
	}
	switch requestText {
	case "":
		// 如果只有command没有txt，默认视为加入比赛
		if horseNum == 0 {
			// 如果当前没有马，视为创建比赛，并将status改为2
			groupContent[groupname].horses[uuid] = horseBase{uname, 0, 0, getHorseEmoji(groupname)}
			groupContent[groupname] = groupContentBase{
				status: 2,
				timeStamp: time.Now().Unix(),
				horses: groupContent[groupname].horses,
			}			
			return &common.ReplyStruct{common.MsgTxt, common.HorseCompCreateSuccess}, nil
		} else if _, ok := groupContent[groupname].horses[uuid]; ok {
			// 如果该用户已经在map中，说明已经加入了，不允许再加入
			return &common.ReplyStruct{common.MsgTxt, common.HorseCompJoinEd}, nil
		} else {
			// 加入当前用户到map中
			groupContent[groupname].horses[uuid] = horseBase{uname, 0, 0, getHorseEmoji(groupname)}
			horseNum = len(groupContent[groupname].horses)
			tmpTxt := fmt.Sprintf(common.HorseCompJoinSuccess, horseNum)
			if horseNum == len(emojiSet) {
				// 支持的用户最大数由emojiSet的长度决定，如果已经相等了，直接开始比赛
				tmpTxt += "\n" + common.Dilimiter + common.HorseCompTooMuch
				groupContent[groupname] = groupContentBase{
					status: 3,
					timeStamp: time.Now().Unix(),
					horses: groupContent[groupname].horses,
				}
				// 起一个新的goroutine进行比赛
				go beginHorseComp(msg, groupname)
			}
			return &common.ReplyStruct{common.MsgTxt, tmpTxt}, nil
		}
	case "比赛开始":
		// 没有马或者马的数量小于三，不允许开始比赛
		if horseNum == 0 {
			return &common.ReplyStruct{common.MsgTxt, common.HorseCompNotCreated}, nil
		} else if horseNum < 3 {
			return &common.ReplyStruct{common.MsgTxt, common.HorseCompNotEnough}, nil
		}

		groupContent[groupname] = groupContentBase{
			status: 3,
			timeStamp: time.Now().Unix(),
			horses: groupContent[groupname].horses,
		}
		go beginHorseComp(msg, groupname)
		return &common.ReplyStruct{}, errors.New("No need response")
	case "人机对抗":
		//调试用
		if horseNum == 0 {
			groupContent[groupname].horses[uuid] = horseBase{uname, 0, 0, getHorseEmoji(groupname)}
			groupContent[groupname].horses[uuid+"1"] = horseBase{uname+"1", 0, 0, getHorseEmoji(groupname)}
			groupContent[groupname].horses[uuid+"2"] = horseBase{uname+"2", 0, 0, getHorseEmoji(groupname)}
			groupContent[groupname].horses[uuid+"3"] = horseBase{uname+"3", 0, 0, getHorseEmoji(groupname)}
			go beginHorseComp(msg, groupname)
			return &common.ReplyStruct{}, errors.New("No need response")
		}
	default:
		return &common.ReplyStruct{common.MsgTxt, fmt.Sprintf(common.UseOfHorseComp, len(emojiSet), len(emojiSet), emojiSet)}, nil
	}

	return &common.ReplyStruct{common.MsgTxt, fmt.Sprintf(common.UseOfHorseComp, len(emojiSet), len(emojiSet), emojiSet)}, nil
}