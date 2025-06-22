package windows

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"yunhudrive/helper"
	"yunhudrive/structs"
	"yunhudrive/uiext"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const maxAllowChunkSize = 40 * 1024 * 1024 // 40MB
const qiniu_token_api = "https://chat-go.jwzhd.com/v1/misc/qiniu-token2"

var UPLOADING_FLAG = false

func WindowUpload(myApp fyne.App, callback func()) {
	// 初始化窗口
	Window := myApp.NewWindow("云湖网盘GUI - 上传")
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

	// 创建日志框
	logArea := widget.NewMultiLineEntry()
	logArea.TextStyle = fyne.TextStyle{Monospace: true}

	// 创建上传按钮
	uploadButton := widget.NewButtonWithIcon("上传", theme.UploadIcon(), func() {
		go func() {

			// 防止重复上传
			if UPLOADING_FLAG {
				return
			}
			UPLOADING_FLAG = true
			defer func() {
				UPLOADING_FLAG = false
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

			// 检查登录状态
			userInfo := helper.LoadUserInfo()
			if userInfo.Token == "" {
				logger("请先登录")
				return
			}

			// 存储MD5列表
			md5List := make([]string, 0)

			// 获取文件路径
			targetPath := filePath.Text
			if targetPath == "" {
				logger("请选择文件")
				return
			}
			logger("上传文件: " + targetPath)

			// 打开文件
			file, err := os.Open(targetPath)
			if err != nil {
				logger("打开文件失败: " + err.Error())
				return
			}
			defer file.Close()

			// 获取文件信息
			fileInfo, err := file.Stat()
			if err != nil {
				logger("获取文件信息失败: " + err.Error())
				return
			}
			fileSize := fileInfo.Size()
			logger(fmt.Sprintf("文件大小: %.2f MB", float64(fileSize)/(1024*1024)))

			// 计算分块数量
			chunkCount := int(fileSize / maxAllowChunkSize)
			if fileSize%maxAllowChunkSize != 0 {
				chunkCount++
			}
			logger(fmt.Sprintf("分块数量: %d", chunkCount))

			// 获取七牛Token
			qiniu_token := helper.GetQiNiuToken(qiniu_token_api, userInfo.Token)
			if qiniu_token == "" {
				logger("获取七牛Token失败")
				return
			}
			logger("获取七牛Token成功: " + qiniu_token)

			// 获取七牛Host
			qiniu_host := helper.GetQiNiuHost(qiniu_token)
			if qiniu_host == "" {
				logger("获取七牛Host失败")
				return
			}
			logger("获取七牛Host成功: " + qiniu_host)

			// 分块读取并上传
			buffer := make([]byte, maxAllowChunkSize)
			for i := 0; i < chunkCount; i++ {
				// 读取分块
				bytesRead, err := file.Read(buffer)
				if err != nil && err != io.EOF {
					logger(fmt.Sprintf("分块 %d 读取失败: %s", i+1, err.Error()))
					return
				}
				chunkData := buffer[:bytesRead]
				logger(fmt.Sprintf("上传分块 %d/%d (%.2f MB)...",
					i+1, chunkCount, float64(bytesRead)/(1024*1024)))

				// 上传分块
				chunkMD5, err := helper.UploadAndGetMD5(qiniu_host, qiniu_token, chunkData)
				if err != nil {
					logger(fmt.Sprintf("分块 %d 上传失败: %s", i+1, err.Error()))
					return
				}
				if chunkMD5 == "" {
					logger(fmt.Sprintf("分块 %d 上传失败: MD5值为空", i+1))
					return
				}
				md5List = append(md5List, chunkMD5)
				logger(fmt.Sprintf("分块 %d 上传成功! MD5: %s", i+1, chunkMD5))
			}
			logger("文件上传完成!")

			// 保存上传结果
			uploadResult := structs.UploadResult{
				Name:   fileInfo.Name(),
				Size:   fileSize,
				Chunks: md5List,
			}
			savePath := filepath.Join("uploads", fileInfo.Name())
			helper.SaveJson(uploadResult, savePath)
			logger("上传结果保存成功: " + savePath)

			// 执行回调
			callback()
			logger("该窗口可以安全关闭")

		}()
	})

	// 创建页面
	Window.SetContent(container.NewBorder(
		container.NewVBox(
			&widget.Form{
				Items: []*widget.FormItem{
					{Text: "文件路径", Widget: filePathWidget},
				},
			},
			container.NewCenter(uploadButton),
		),
		nil, nil, nil,
		container.NewScroll(logArea),
	))

	// 显示窗口
	Window.Show()
}
