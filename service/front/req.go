package main

import (
	"fmt"
	"net/http"
)

const (
	signupUrl string = "/api/storage/v1/users/signup"
)
const (
	apigw string = "http://localhost:8888"
)

func signupRequest(req *SignupReq) (*Response, error) {
	request := Request{
		Url:    fmt.Sprintf("%s%s", apigw, signupUrl),
		Method: http.MethodPost,
	}
	return ask(request)
}
