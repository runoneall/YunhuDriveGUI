package windows

import (
	"fmt"
	"os"
	"path/filepath"
	"yunhudrive/apihttp"
	"yunhudrive/helper"
	"yunhudrive/structs"
	"yunhudrive/uiext"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/adrg/xdg"
)

const download_api = "https://chat-file.jwznb.com"

var DOWNLOADING_FLAG = false

func WindowDownload(myApp fyne.App, uploadResult structs.UploadResult) {
	// 初始化窗口
	Window := myApp.NewWindow("云湖网盘GUI - 下载")
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

	// 默认下载位置
	folderPath.SetText(xdg.UserDirs.Download)

	// 创建日志框
	logArea := widget.NewMultiLineEntry()
	logArea.TextStyle = fyne.TextStyle{Monospace: true}

	// 创建下载按钮
	downloadButton := widget.NewButton("开始下载", func() {
		go func() {

			// 防止重复下载
			if DOWNLOADING_FLAG {
				return
			}
			DOWNLOADING_FLAG = true
			defer func() {
				DOWNLOADING_FLAG = false
			}()

			// 初始化日志输出
			logger := func(msg string) {
				fyne.Do(func() {
					logArea.SetText(logArea.Text + msg + "\n")
				})
			}
			fyne.Do(func() {
				logArea.SetText("")
			})

			// 获取最终下载文件夹
			downloadFolder := folderPath.Text
			if downloadFolder == "" {
				downloadFolder = xdg.UserDirs.Download
			}
			logger("下载到文件夹: " + downloadFolder)

			// 获取要下载的文件名
			fileName := uploadResult.Name
			logger("开始下载: " + fileName)

			// 最终文件路径
			filePath := filepath.Join(downloadFolder, fileName)
			logger("下载到文件: " + filePath)

			// 检查文件路径是否存在
			if helper.IsExist(filePath) {
				logger("不能完成下载, 因为文件已存在")
				return
			}

			// 初始化文件
			file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				logger("不能完成下载, 因为无法创建文件")
				return
			}
			defer file.Close()

			// 获取md5列表
			md5List := uploadResult.Chunks
			chunksTotal := len(md5List)
			logger(fmt.Sprintf("分块数量: %d", chunksTotal))

			// 下载文件
			for index, md5Value := range md5List {
				logger(fmt.Sprintf("正在下载分块: %d/%d", index+1, chunksTotal))

				// 获取分块下载链接
				url := download_api + "/" + md5Value
				logger("获得下载链接: " + url)

				// 获得块二进制数据
				respBytes, err := apihttp.Get(url, map[string]string{
					"referer": "myapp.jwznb.com",
				})
				if err != nil {
					logger("不能完成下载, 因为无法获取分块数据")
					return
				}

				// 写入块数据
				bytesWritten, err := file.Write(respBytes)
				if err != nil {
					logger("不能完成下载, 因为写入分块数据失败")
					return
				}
				if bytesWritten != len(respBytes) {
					logger("不能完成下载, 因为写入分块数据长度不匹配")
					return
				}
				logger(fmt.Sprintf("成功写入块数据: %d", bytesWritten))

			}

			// 下载完成
			logger("下载完成")
			logger("文件保存在: " + filePath)
			logger("该窗口可以被安全关闭")

		}()
	})

	// 创建页面
	Window.SetContent(container.NewBorder(
		container.NewVBox(
			&widget.Form{
				Items: []*widget.FormItem{
					{Text: "下载到", Widget: folderPathWidget},
				},
			},
			container.NewCenter(downloadButton),
		),
		nil, nil, nil,
		container.NewScroll(logArea),
	))

	// 显示窗口
	Window.Show()
}
