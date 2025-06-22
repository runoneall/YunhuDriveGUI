package windows

import (
	"fmt"
	"image/color"
	"strings"
	"yunhudrive/helper"
	"yunhudrive/structs"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var REFRESH_FLAG = true

func WindowMain(myApp fyne.App) {

	// 初始化窗口
	Window := myApp.NewWindow("云湖网盘GUI")
	Window.Resize(fyne.NewSize(800, 600))

	// 上传按钮
	uploadButton := widget.NewButtonWithIcon("上传", theme.UploadIcon(), func() {
		WindowUpload(myApp, func() {
			REFRESH_FLAG = true
		})
	})

	// 刷新按钮
	refreshButton := widget.NewButtonWithIcon("刷新", theme.ViewRefreshIcon(), func() {
		REFRESH_FLAG = true
	})

	// 登录按钮
	loginButton := widget.NewButtonWithIcon("登录", theme.LoginIcon(), func() {
		WindowLogin(myApp)
	})

	// 登出按钮
	logoutButton := widget.NewButtonWithIcon("登出", theme.LogoutIcon(), func() {
		helper.SaveUserInfo("", "", "")
		dialog.ShowInformation("提示", "登出成功", Window)
	})

	// 退出按钮
	exitButton := widget.NewButtonWithIcon("退出", theme.CancelIcon(), func() {
		myApp.Quit()
	})

	// 顶部按钮区
	topButtons := container.NewVBox(
		container.NewBorder(
			nil,
			canvas.NewRectangle(color.Gray{0x99}),
			nil, nil,
			container.NewHBox(
				uploadButton,
				refreshButton,
				loginButton,
				logoutButton,
				exitButton,
			),
		),
	)

	// 右部文件信息
	fileShowArea := container.NewVBox(
		widget.NewLabel("选择一个文件以使用"),
	)
	fileControlButtons := func(uploadResult structs.UploadResult) *fyne.Container {

		// 返回按钮区容器
		return container.NewHBox(

			// 下载按钮
			widget.NewButtonWithIcon("下载", theme.DownloadIcon(), func() {
				fmt.Println("下载文件：" + uploadResult.Name)
			}),

			// 删除按钮
			widget.NewButtonWithIcon("删除", theme.DeleteIcon(), func() {

				// 删除上传结果文件
				uploadResult.Delete()

				// 恢复文件信息区
				fileShowArea.RemoveAll()
				fileShowArea.Add(widget.NewLabel("选择一个文件以使用"))

				// 刷新文件列表
				REFRESH_FLAG = true
			}),

			// 分享按钮
			widget.NewButtonWithIcon("分享", theme.UploadIcon(), func() {
				fmt.Println("分享文件：" + uploadResult.Name)
			}),
		)
	}
	drawFileShowArea := func(file string) {
		fileShowArea.RemoveAll()

		// 获取上传结果
		uploadResult, _ := helper.LoadUploadResult(file)

		// 添加文件操作按钮
		fileShowArea.Add(fileControlButtons(uploadResult))

		// 添加配置预览
		uploadResultContent := uploadResult.FormatJson()
		uploadResultShowArea := widget.NewMultiLineEntry()
		uploadResultShowArea.SetText(uploadResultContent)
		uploadResultShowArea.TextStyle = fyne.TextStyle{Monospace: true}
		uploadResultShowArea.SetMinRowsVisible(
			strings.Count(uploadResultContent, "\n") + 1,
		)
		fileShowArea.Add(uploadResultShowArea)
	}

	// 左部文件列表
	fileList := container.NewVBox(
		widget.NewLabel("文件列表"),
	)

	// 绘制文件列表
	drawFileList := func() {
		fileList.RemoveAll()

		// 加载文件列表
		files := helper.ListFiles("uploads")
		if len(files) == 0 {
			fileList.Add(widget.NewLabel("暂无文件"))
		}
		if len(files) > 0 {
			for _, file := range files {
				fileList.Add(widget.NewButton(file, func() {
					drawFileShowArea(file)
				}))
			}
		}
	}
	go func() {
		for {
			if REFRESH_FLAG {
				fyne.Do(drawFileList)
				REFRESH_FLAG = false
			}
		}
	}()

	// 创建一个水平分割容器
	mainContainer := container.NewHSplit(
		container.NewScroll(fileList),
		container.NewVScroll(fileShowArea),
	)
	mainContainer.SetOffset(1.0 / 3.0)

	// 创建页面
	Window.SetContent(container.NewBorder(
		topButtons,
		nil, nil, nil,
		mainContainer,
	))

	// 显示窗口
	Window.Show()
}
