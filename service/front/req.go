package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	signupUrl   string = "/api/storage/v1/users/signup"
	loginUrl    string = "/api/storage/v1/users/login"
	profileUrl  string = "/api/storage/v1/users/profile"
	listUrl     string = "/api/storage/v1/users/file-list"
	getDownload string = "/api/storage/v1/download/download-url"
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
		Url:    fmt.Sprintf("%s%s", download, getDownload),
		Method: http.MethodPost,
		Body:   body,
		HeaderSet: map[string]string{
			"Authorization": fmt.Sprintf("%s %s", X_Refresh_Token, X_Jwt_Token),
		},
	}
	return ask(request)
}
