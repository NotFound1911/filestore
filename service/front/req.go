package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	signupUrl      string = "/api/storage/v1/users/signup"
	loginUrl       string = "/api/storage/v1/users/login"
	profileUrl     string = "/api/storage/v1/users/profile"
	listUrl        string = "/api/storage/v1/users/file-list"
	getDownloadUrl string = "/api/storage/v1/download/download-url"
	downloadUrl    string = "/api/storage/v1/download/download"
)
const (
	apigw    string = "http://localhost:8888"
	download string = "http://localhost:8890"
)

func signupRequest(req *SignupReq) (*Response, error) {
	// 将结构体实例转换为 JSON 字符串
	jsonData, err := json.Marshal(req)
	if err != nil {
		fmt.Println("JSON marshal error:", err)
		return nil, err
	}

	// 创建一个 bytes.Buffer，并将 JSON 数据写入其中
	body := bytes.NewBuffer(jsonData)
	request := Request{
		Url:    fmt.Sprintf("%s%s", apigw, signupUrl),
		Method: http.MethodPost,
		Body:   body,
	}
	return ask(request)
}

func loginRequest(req *LoginReq) (*Response, error) {
	// 将结构体实例转换为 JSON 字符串
	jsonData, err := json.Marshal(req)
	if err != nil {
		fmt.Println("JSON marshal error:", err)
		return nil, err
	}

	// 创建一个 bytes.Buffer，并将 JSON 数据写入其中
	body := bytes.NewBuffer(jsonData)
	request := Request{
		Url:    fmt.Sprintf("%s%s", apigw, loginUrl),
		Method: http.MethodPost,
		Body:   body,
	}
	return ask(request)
}

func profileRequest() (*Response, error) {

	request := Request{
		Url:    fmt.Sprintf("%s%s", apigw, profileUrl),
		Method: http.MethodGet,
		HeaderSet: map[string]string{
			"Authorization": fmt.Sprintf("%s %s", X_Refresh_Token, X_Jwt_Token),
		},
	}
	return ask(request)
}

func listRequest() (*Response, error) {
	request := Request{
		Url:    fmt.Sprintf("%s%s", apigw, listUrl),
		Method: http.MethodGet,
		HeaderSet: map[string]string{
			"Authorization": fmt.Sprintf("%s %s", X_Refresh_Token, X_Jwt_Token),
		},
	}
	return ask(request)
}

func downloadUrlRequest(req *DownloadReq) (*Response, error) {
	// 将结构体实例转换为 JSON 字符串
	jsonData, err := json.Marshal(req)
	if err != nil {
		fmt.Println("JSON marshal error:", err)
		return nil, err
	}

	// 创建一个 bytes.Buffer，并将 JSON 数据写入其中
	body := bytes.NewBuffer(jsonData)
	request := Request{
		Url:    fmt.Sprintf("%s%s", download, getDownloadUrl),
		Method: http.MethodPost,
		Body:   body,
		HeaderSet: map[string]string{
			"Authorization": fmt.Sprintf("%s %s", X_Refresh_Token, X_Jwt_Token),
			"Content-Type":  "application/json",
		},
	}
	return ask(request)
}

func downloadRequest(info map[string]string) (*http.Response, error) {
	request := Request{
		Url:    fmt.Sprintf("%s%s", download, downloadUrl),
		Method: http.MethodPost,
		HeaderSet: map[string]string{
			"Authorization": fmt.Sprintf("%s %s", X_Refresh_Token, X_Jwt_Token),
			"Content-Type":  "application/json",
		},
		Params: map[string]string{
			"filename": info["filename"],
			"bucket":   info["bucket"],
			"name":     info["name"],
		},
	}
	req, err := http.NewRequest(request.Method, request.Url, request.Body)
	if err != nil {
		return nil, err
	}
	// header 添加字段,包含token
	if request.HeaderSet != nil {
		for k, v := range request.HeaderSet {
			req.Header.Set(k, v)
		}
	}
	// query params
	if request.Params != nil {
		params := make(url.Values)
		for k, v := range request.Params {
			params.Add(k, v)
		}
		req.URL.RawQuery = params.Encode()
	}

	resp, err := Client.Do(req)
	return resp, err
}
