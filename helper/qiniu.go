package helper

import (
	"encoding/json"
	"fmt"
	"strings"
	"yunhudrive/apihttp"
	"yunhudrive/structs"
)

func GetQiNiuToken(url string, token string) string {
	respBytes, err := apihttp.Get(url, map[string]string{
		"token": token,
	})
	if err != nil {
		return ""
	}
	var qiNiuTokenResp structs.YunhuLoginResponse
	err = json.Unmarshal(respBytes, &qiNiuTokenResp)
	if err != nil {
		return ""
	}
	if qiNiuTokenResp.Msg != "success" {
		return ""
	}
	return qiNiuTokenResp.Data.Token
}

func GetQiNiuHost(token string) string {
	ak := strings.Split(token, ":")[0]
	bucket := "chat68-file"
	hostUrl := fmt.Sprintf("https://api.qiniu.com/v4/query?ak=%s&bucket=%s", ak, bucket)
	respBytes, err := apihttp.Get(hostUrl, map[string]string{})
	if err != nil {
		return ""
	}
	var qiNiuHostResp structs.QiNiuHostResponse
	err = json.Unmarshal(respBytes, &qiNiuHostResp)
	if err != nil {
		return ""
	}
	if len(qiNiuHostResp.Hosts) == 0 {
		return ""
	}
	host := qiNiuHostResp.Hosts[0].Up.Domains[0]
	if !strings.HasPrefix(host, "https://") {
		host = "https://" + host
	}
	return host
}
