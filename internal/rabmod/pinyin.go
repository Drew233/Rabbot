package rabmod

import (
	"github.com/mozillazg/go-pinyin"

	"rabbot/internal/common"
)

func init() {
	common.FuncNameMap["TransPinyin"] = TransPinyin
}

func trans(str string) (string, bool) {
	a := pinyin.NewArgs()
	// 包含声调
	a.Style = pinyin.Tone
	pinyinArrays := pinyin.Pinyin(str, a)
	if len(pinyinArrays) == 0 || len(pinyinArrays[0]) == 0 {
		return str, false
	}

	return pinyinArrays[0][0], true
}

func TransPinyin(requestStruct *common.RequestStruct)(*common.ReplyStruct, error) {
	requestTxt := requestStruct.RequestTxt
	replyTxt := ""
	for _, s := range requestTxt {
		transed, ifTrans := trans(string(s))
		if replyTxt != "" && ifTrans == true {
			replyTxt += " "
		}
		replyTxt += transed
	}

	return &common.ReplyStruct{common.MsgTxt, replyTxt}, nil
}
