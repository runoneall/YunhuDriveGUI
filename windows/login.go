package windows

import (
	"time"
	"yunhudrive/helper"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const login_api = "https://chat-go.jwzhd.com/v1/user/email-login"

func WindowLogin(myApp fyne.App) {
	// 初始化窗口
	Window := myApp.NewWindow("云湖网盘GUI - 登录")
	Window.Resize(fyne.NewSize(400, 150))

	// 创建表单
	email := widget.NewEntry()
	password := widget.NewPasswordEntry()

	// 如果已经登录
	userInfo := helper.LoadUserInfo()
	if userInfo.Token != "" {
		email.SetText(userInfo.Email)
		password.SetText(userInfo.Password)
	}

	// 创建登录按钮
	loginButton := widget.NewButtonWithIcon("登录", theme.LoginIcon(), func() {
		token := helper.LoginAndGetToken(login_api, email.Text, password.Text)

		// 登录失败
		if token == "" {
			dialog.ShowInformation("错误", "登录失败", Window)
		}

		// 登录成功
		if token != "" {
			helper.SaveUserInfo(email.Text, password.Text, token)
			dialog.ShowInformation("提示", "登录成功", Window)

			// 延时关闭窗口
			fyne.Do(func() {
				time.Sleep(2 * time.Second)
				Window.Close()
			})
		}
	})

	// 创建页面
	Window.SetContent(container.NewHScroll(
		container.NewVBox(
			&widget.Form{
				Items: []*widget.FormItem{
					{Text: "邮箱", Widget: email},
					{Text: "密码", Widget: password},
				},
			},
			container.NewCenter(loginButton),
		),
	))

	// 显示窗口
	Window.Show()
}
