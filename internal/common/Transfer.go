package common

import (
	
	"github.com/eatmoreapple/openwechat"
)

/*
	@Struct 请求处理函数所用结构体
*/
type RequestStruct struct {
	Uname string // 昵称
	Uuid  string // 用户唯一标识
	Groupname string // 群昵称，取不到或没有的时候为""
	RequestTxt string // 请求字符串
	Commond string // 请求指令
	Msg *openwechat.Message // openwechat Message对象
}

func GenRequestStruct (uname, uuid, gname, requesttxt, commond string, msg *openwechat.Message) *RequestStruct {
	return &RequestStruct{
		Uname: uname,
		Uuid: uuid,
		Groupname: gname,
		RequestTxt: requesttxt,
		Commond: commond,
		Msg: msg,
	}
}

/* 
	@Struct 消息处理返回数据结构

	@Member reType: 返回数据类型
			1: 文本消息
			2: 图片消息,此时reText是文件路径
*/
type ReplyStruct struct {
	ReType int
	ReText string
}


const (
	MsgTxt = 1
	MsgPic = 2
)