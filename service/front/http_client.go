package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	serv "github.com/NotFound1911/filestore/pkg/server"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"
)

type Response struct {
	serv.Result
	RespStatusCode int
	RespHeader     http.Header
}

var (
	Client *http.Client //HTTPClient
)

type Request struct {
	Url       string
	Body      io.Reader
	HeaderSet map[string]string
	Method    string
	Params    map[string]string
}

func init() {
	Client = &http.Client{
		Timeout: 300 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
			DisableKeepAlives: true,
			Proxy:             http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second, // tcp连接超时时间
				KeepAlive: 60 * time.Second, // 保持长连接的时间
				DualStack: true,
			}).DialContext, // 设置连接的参数
			MaxIdleConns:          100, // 最大空闲连接
			MaxConnsPerHost:       100,
			MaxIdleConnsPerHost:   100,              // 每个host保持的空闲连接数
			ExpectContinueTimeout: 30 * time.Second, // 等待服务第一响应的超时时间
			IdleConnTimeout:       60 * time.Second, // 空闲连接的超时时间
		},
	}
}

// ask 建立http请求，返回header信息
func ask(requester Request) (*Response, error) {
	fmt.Println("url:", requester.Url)
	request, err := http.NewRequest(requester.Method, requester.Url, requester.Body)
	if err != nil {
		return &Response{RespStatusCode: http.StatusBadRequest}, err
	}
	// header 添加字段,包含token
	if requester.HeaderSet != nil {
		for k, v := range requester.HeaderSet {
			request.Header.Set(k, v)
		}
	}
	// query params
	if requester.Params != nil {
		params := make(url.Values)
		for k, v := range requester.Params {
			params.Add(k, v)
		}
		request.URL.RawQuery = params.Encode()
	}

	resp, err := Client.Do(request)
	if err != nil {
		return &Response{RespStatusCode: http.StatusBadRequest}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(resp.Body)

	return checkRespStatus(resp)
}

// checkRespStatus 状态检查
func checkRespStatus(resp *http.Response) (*Response, error) {
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	respRes := Response{}
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		if err := json.Unmarshal(bodyBytes, &respRes); err != nil {
			return nil, err
		}
		respRes.RespHeader = resp.Header
		respRes.RespStatusCode = resp.StatusCode
		return &respRes, nil
	}
	return nil, errors.New(string(bodyBytes))
}
