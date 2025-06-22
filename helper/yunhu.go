package helper

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"yunhudrive/apihttp"
	"yunhudrive/structs"
)

func LoginAndGetToken(url string, email string, password string) string {
	jsonBytes, err := json.Marshal(structs.YunhuLogin{
		Email:    email,
		Password: password,
		DeviceID: apihttp.RandomUserAgent(),
		Platform: "YunhuDriveGUI",
	})
	if err != nil {
		return ""
	}
	respBytes, err := apihttp.Post(url, "application/json", jsonBytes)
	if err != nil {
		return ""
	}
	var loginResp structs.YunhuLoginResponse
	err = json.Unmarshal(respBytes, &loginResp)
	if err != nil {
		return ""
	}
	if loginResp.Msg != "success" {
		return ""
	}
	return loginResp.Data.Token
}

func UploadAndGetMD5(url string, qiniu_token string, chunkData []byte) (string, error) {

	// 计算块的MD5
	hash := md5.Sum(chunkData)
	md5Value := hex.EncodeToString(hash[:])

	// 准备multipart表单
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 添加token和key字段
	_ = writer.WriteField("token", qiniu_token)
	_ = writer.WriteField("key", md5Value)

	// 添加chunk数据
	chunkPart, err := writer.CreateFormFile("file", md5Value)
	if err != nil {
		return "", err
	}
	_, err = chunkPart.Write(chunkData)
	if err != nil {
		return "", err
	}

	// 构建表单
	err = writer.Close()
	if err != nil {
		return "", err
	}

	// 发送请求
	respBytes, err := apihttp.PostMultipartForm(
		url, map[string]string{
			"Content-Type": writer.FormDataContentType(),
		}, body,
	)
	if err != nil {
		return "", err
	}

	// 解析响应
	var qiniuUploadResp structs.QiNiuUploadResponse
	err = json.Unmarshal(respBytes, &qiniuUploadResp)
	if err != nil {
		return "", err
	}
	if qiniuUploadResp.Key != md5Value {
		return "", fmt.Errorf("upload failed, key not match")
	}
	return md5Value, nil

}
