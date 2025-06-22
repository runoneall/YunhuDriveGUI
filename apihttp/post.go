package apihttp

import (
	"bytes"
	"io"
	"net/http"
)

func Post(url string, data_type string, data []byte) ([]byte, error) {
	resp, err := http.Post(url, data_type, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func PostMultipartForm(url string, headers map[string]string, data *bytes.Buffer) ([]byte, error) {
	req, err := http.NewRequest("POST", url, data)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
