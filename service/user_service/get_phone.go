package user_service

import (
	"LibraryManagementV1/LM_V3/global"
	"strings"
)

var specialPhones = map[string]int{}

//读取白名单

func ReadWhitePhones() {
	WhitePhones := strings.Split(global.Config.SMS.WhiteListedPhoneNumber, ",")
	for _, phone := range WhitePhones {
		specialPhones[phone] = 0
	}
}
