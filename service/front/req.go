package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	signupUrl string = "/api/storage/v1/users/signup"
	loginUrl  string = "/api/storage/v1/users/login"
)
const (
	apigw string = "http://localhost:8888"
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
