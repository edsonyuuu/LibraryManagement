package model

import (
	"math/rand"
	"time"
)

const (
	// 验证码字符集合
	charSet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	// 验证码长度
	codeLength = 6
)

func GetMessage(userId int64) []SendMsg {
	sql := "select * from send_msg where status=0 and user_id=?"
	var messages []SendMsg
	DB.Raw(sql, userId).Scan(&messages)
	return messages
}

// GenerateCode 生成随机验证码
func GenerateCode() string {
	rand.Seed(time.Now().UnixNano())

	code := make([]byte, codeLength)
	for i := 0; i < codeLength; i++ {
		index := rand.Intn(len(charSet))
		code[i] = charSet[index]
	}
	// 模拟发送验证码的过程
	//time.Sleep(2 * time.Second)
	//fmt.Printf("验证码 %s 已发送", code)
	return string(code)
}
