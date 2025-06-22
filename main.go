package main

import (
	"yunhudrive/helper"
	"yunhudrive/windows"

	"fyne.io/fyne/v2/app"
)

func main() {

	// 初始化文件和目录
	if !helper.IsExist("user_info.json") {
		helper.SaveUserInfo("", "", "")
	}
	if !helper.IsExist("uploads") {
		helper.CreateFolder("uploads")
	}

	// 初始化应用
	myApp := app.New()

	// 显示窗口
	windows.WindowMain(myApp)
	myApp.Run()
}
