package helper

import (
	"encoding/json"
	"os"
	"yunhudrive/structs"
)

func SaveUserInfo(email string, password string, token string) error {
	return SaveJson(structs.UserInfo{
		Email:    email,
		Password: password,
		Token:    token,
	}, "user_info.json")
}

func LoadUserInfo() structs.UserInfo {
	jsonBytes, err := os.ReadFile("user_info.json")
	if err != nil {
		return structs.UserInfo{}
	}
	var userInfo structs.UserInfo
	err = json.Unmarshal(jsonBytes, &userInfo)
	if err != nil {
		return structs.UserInfo{}
	}
	return userInfo
}
