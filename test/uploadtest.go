package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	v1 "github.com/NotFound1911/filestore/api/rest/upload/v1"
	serv "github.com/NotFound1911/filestore/pkg/server"
	account "github.com/NotFound1911/filestore/service/account/run"
	apigw "github.com/NotFound1911/filestore/service/apigw/run"
	file_manager "github.com/NotFound1911/filestore/service/file_manager/run"
	transfer "github.com/NotFound1911/filestore/service/transfer/run"
	upload "github.com/NotFound1911/filestore/service/upload/run"
	"github.com/NotFound1911/filestore/util"
	"io"
	"io/ioutil"
	"math"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// 文件上传测试
func main() {
	initServ()
	time.Sleep(time.Second * 5)
	u := user{
		Email:    "123@qq.com",
		Password: "hello#world123",
	}
	reqAccountSignup(&u)
	token := reqAccountLogin(&u)
	file := "./tmp/test.txt"
	reqUploadSingleFile(file, token)
	file = "./tmp/test.jpg"
	id, err := reqInitMultiUploadFile(file, token)
	if err != nil {
		fmt.Println("初始化任务失败：", err)
		return
	}
	startMultiUploadFile(file, token, id)
	reqMultiUploadFileMerge(file, token, id)
	select {}
}

// 服务初始化
func initServ() {
	go upload.Run()
	go apigw.Run()
	go account.Run()
	go file_manager.Run()
	go transfer.Run()
	time.Sleep(time.Second * 5)
}

type user struct {
	Email    string
	Password string
}

// 账户注册请求
func reqAccountSignup(u *user) {
	// 定义请求体数据
	requestData := map[string]string{
		"email":            u.Email,
		"password":         u.Password,
		"confirm_password": u.Password,
	}

	// 将请求体数据编码为 JSON 格式
	requestBody, err := json.Marshal(requestData)
	if err != nil {
		fmt.Println("Failed to marshal JSON:", err)
		return
	}

	// 创建 POST 请求
	resp, err := http.Post("http://localhost:8888/api/storage/v1/users/signup",
		"application/json",
		bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Failed to make POST request:", err)
		return
	}
	defer resp.Body.Close()
}

// 账户登录
func reqAccountLogin(u *user) string {
	// 定义请求体数据
	requestData := map[string]string{
		"email":    u.Email,
		"password": u.Password,
	}

	// 将请求体数据编码为 JSON 格式
	requestBody, err := json.Marshal(requestData)
	if err != nil {
		fmt.Println("Failed to marshal JSON:", err)
		return ""
	}

	// 创建 POST 请求
	resp, err := http.Post("http://localhost:8888/api/storage/v1/users/login", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Failed to make POST request:", err)
		return ""
	}
	defer resp.Body.Close()
	token := fmt.Sprintf("%s %s", resp.Header.Get("X-Refresh-Token"), resp.Header.Get("X-Jwt-Token"))
	return token
}

// 单文件上传
func reqUploadSingleFile(filePath, token string) {
	url := "http://localhost:8889/api/storage/v1/upload/upload-file"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	defer file.Close()
	part1, err := writer.CreateFormFile("file", filepath.Base(filePath))
	_, err = io.Copy(part1, file)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = writer.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", token)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

// 分块文件上传初始化
func reqInitMultiUploadFile(filePath, token string) (int64, error) {
	url := "http://localhost:8889/api/storage/v1/upload/init-multi-upload-file"
	method := "POST"
	size, err := util.GetFileSize(filePath)
	if err != nil {
		return -1, err
	}
	sha1, err := util.GetFileSha1(filePath)
	if err != nil {
		return -1, err
	}
	// 准备要发送的 JSON 数据
	jsonData := v1.InitMultiUploadFileReq{
		FileName: filepath.Base(filePath),
		FileSize: size,
		FileSha1: sha1,
	}
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return -1, err
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return -1, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	// 创建 HTTP 客户端并发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return -1, err
	}
	defer resp.Body.Close()

	// 处理响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return -1, err
	}
	fmt.Println(string(body))
	res := serv.Result{}
	if err := json.Unmarshal(body, &res); err != nil {
		return -1, err
	}
	data, ok := res.Data.(map[string]interface{})
	if !ok {
		return -1, fmt.Errorf("%v is not type map[string]interface{}", res.Data)
	}
	id, ok := data["id"].(float64)
	if !ok {
		return -1, fmt.Errorf("%v is not type float64", data["id"])
	}
	return int64(id), nil
}

// 分块文件开始上传
func startMultiUploadFile(filePath, token string, uploadId int64) error {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Failed to open file:", err)
		return err
	}
	defer file.Close()
	fileSize, err := util.GetFileSize(filePath)
	if err != nil {
		return err
	}
	//chunkSize := 1024.0
	chunkSize := 1024.0 * 1024
	currentChunk := int64(1)
	totalChunk := int64(math.Ceil(float64(fileSize) / chunkSize))
	var wg sync.WaitGroup
	ch := make(chan struct{}, 2)
	var multiErr error
	for currentChunk <= totalChunk {
		start := (currentChunk - 1) * int64(chunkSize)
		end := minH(fileSize, start+int64(chunkSize))
		buffer := make([]byte, end-start)
		// 循环读取，会自动偏移
		_, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			fmt.Println("读取文件长度失败", err)
			break
		}
		sha1 := util.Sha1(buffer)
		// 多协程上传
		ch <- struct{}{}
		wg.Add(1)
		fp := filePart{
			UploadId: uploadId,
			ChunkId:  currentChunk,
			Name:     filepath.Base(filePath),
			Content:  buffer,
			Count:    totalChunk,
			Sha1:     sha1,
		}
		go func(fp filePart, wg *sync.WaitGroup, token string) {
			defer wg.Done()
			if err = reqMultiUploadFilePart(fp, token); err != nil {
				multiErr = err
				fmt.Printf("传输文件分块失败:%v\n", err)
			}
			<-ch
		}(fp, &wg, token)
		currentChunk += 1
	}
	wg.Wait()
	return multiErr
}

// 分块文件合并
func reqMultiUploadFileMerge(filePath, token string, uploadId int64) error {
	url := "http://localhost:8889/api/storage/v1/upload/multi-upload-file-merge"
	method := "POST"
	size, err := util.GetFileSize(filePath)
	if err != nil {
		return err
	}
	sha1, err := util.GetFileSha1(filePath)
	if err != nil {
		return err
	}
	// 定义请求体数据
	requestData := v1.MultiUploadFileMergeReq{
		UploadId: uploadId,
		FileName: filepath.Base(filePath),
		FileSize: size,
		FileSha1: sha1,
	}

	// 将请求体数据编码为 JSON 格式
	requestBody, err := json.Marshal(requestData)
	if err != nil {
		fmt.Println("Failed to marshal JSON:", err)
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))

	if err != nil {
		fmt.Println(err)
		return err
	}
	req.Header.Add("Authorization", token)
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(string(body))
	return nil
}

// 开始分块传输
type filePart struct {
	UploadId int64
	ChunkId  int64
	Name     string
	Content  []byte
	Count    int64
	Sha1     string
}

func reqMultiUploadFilePart(fp filePart, token string) error {
	baseURL := "http://localhost:8889/api/storage/v1/upload/multi-upload-file-part"
	method := "POST"

	// 准备要发送的 params 数据
	// 创建一个 url.Values 对象来存储参数
	params := url.Values{}
	// 添加参数
	params.Add("chunk_id", fmt.Sprintf("%d", fp.ChunkId))
	params.Add("upload_id", fmt.Sprintf("%d", fp.UploadId))
	params.Add("file_name", fp.Name)
	params.Add("chunk_sha1", fp.Sha1)
	params.Add("chunk_size", fmt.Sprintf("%d", len(fp.Content)))
	params.Add("count", fmt.Sprintf("%d", fp.Count))
	urlWithParams := baseURL + "?" + params.Encode()

	// 创建一个 buffer 用于存储 multipart 请求体
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// 添加一个文件表单项
	fieldName := "file"

	// 使用 CreateFormFile 创建文件表单项
	part, err := writer.CreateFormFile(fieldName, fp.Name)
	if err != nil {
		panic(err)
	}

	// 将文件内容写入到表单项中
	_, err = io.Copy(part, bytes.NewReader(fp.Content))
	if err != nil {
		panic(err)
	}

	// 关闭 writer 完成 multipart 请求体的构建
	writer.Close()

	// 现在 body 包含了完整的 multipart 请求体

	// 创建 HTTP 请求
	req, err := http.NewRequest(method, urlWithParams, io.NopCloser(&body))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", token)

	// 创建 HTTP 客户端并发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return err
	}
	defer resp.Body.Close()

	// 处理响应
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(string(respBody))
	return nil
}
func minH(a, b int64) int64 {
	if a <= b {
		return a
	} else {
		return b
	}
}
