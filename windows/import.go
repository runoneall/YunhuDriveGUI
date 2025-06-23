package windows

import (
	"path/filepath"
	"yunhudrive/helper"
	"yunhudrive/uiext"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func WindowImport(myApp fyne.App, callback func()) {
	// 初始化窗口
	Window := myApp.NewWindow("云湖网盘GUI - 导入")
	Window.Resize(fyne.NewSize(500, 400))

	// 创建表单
	filePath := widget.NewEntry()
	filePathChoose := widget.NewButton("选择文件", func() {
		uiext.FilePicker(Window, func(path string) {
			filePath.SetText(path)
		})
	})
	filePathWidget := container.NewBorder(
		nil, nil, nil,
		filePathChoose, filePath,
	)

	// 创建导入按钮
	importButton := widget.NewButton("导入", func() {
		filePathText := filePath.Text

		// 尝试解析
		if filePathText == "" {
			dialog.ShowInformation("错误", "请选择文件", Window)
			return
		}
		result, err := helper.LoadUploadResultFromPath(filePathText)
		if err != nil {
			dialog.ShowInformation("错误", "解析文件失败", Window)
			return
		}

		// 验证格式
		if result.Name == "" || len(result.Chunks) == 0 {
			dialog.ShowInformation("错误", "文件格式错误", Window)
			return
		}

		// 导入文件
		savePath := filepath.Join("uploads", result.Name)
		helper.SaveJson(result, savePath)
		dialog.ShowInformation("成功", "文件导入成功", Window)

		// 回调
		callback()

	})

	// 创建页面
	Window.SetContent(container.NewHScroll(
		container.NewVBox(
			&widget.Form{
				Items: []*widget.FormItem{
					{Text: "文件路径", Widget: filePathWidget},
				},
			},
			container.NewCenter(importButton),
		),
	))

	// 显示窗口
	Window.Show()
}
