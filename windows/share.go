package windows

import (
	"fmt"
	"path/filepath"
	"time"
	"yunhudrive/helper"
	"yunhudrive/structs"
	"yunhudrive/uiext"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/adrg/xdg"
)

func WindowShare(myApp fyne.App, uploadResult structs.UploadResult) {
	// 初始化窗口
	Window := myApp.NewWindow("云湖网盘GUI - 分享")
	Window.Resize(fyne.NewSize(500, 400))

	// 创建表单
	folderPath := widget.NewEntry()
	folderPathChoose := widget.NewButton("选择文件夹", func() {
		uiext.FolderPicker(Window, func(path string) {
			folderPath.SetText(path)
		})
	})
	folderPathWidget := container.NewBorder(
		nil, nil, nil,
		folderPathChoose, folderPath,
	)

	// 默认分享位置
	folderPath.SetText(xdg.UserDirs.Download)

	// 创建分享按钮
	shareButton := widget.NewButton("导出源信息文件", func() {

		// 确保文件夹路径存在
		if !helper.IsExist(folderPath.Text) {
			dialog.ShowInformation("错误", "文件夹路径不存在", Window)
			return
		}

		// 拼接导出文件路径
		exportFilePath := filepath.Join(
			folderPath.Text,
			fmt.Sprintf("share-%d.dat", time.Now().UnixMilli()),
		)

		// 保存到文件
		helper.SaveJson(uploadResult, exportFilePath)
		dialog.ShowInformation("提示", "分享源信息文件成功\n保存至: "+exportFilePath, Window)

	})

	// 创建页面
	Window.SetContent(container.NewHScroll(
		container.NewVBox(
			widget.NewLabel("导出源信息文件到指定目录以分享"),
			widget.NewLabel("任何拥有该文件的人都可以下载源文件"),
			widget.NewLabel("请放心, 该文件不会泄露您的登录信息"),
			&widget.Form{
				Items: []*widget.FormItem{
					{Text: "保存到", Widget: folderPathWidget},
				},
			},
			container.NewCenter(shareButton),
		),
	))

	// 显示窗口
	Window.Show()
}
