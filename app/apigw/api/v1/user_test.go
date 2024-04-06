package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	userv1 "github.com/NotFound1911/filestore/app/account/api/proto/gen/user/v1"
	usermocks "github.com/NotFound1911/filestore/app/account/api/proto/gen/user/v1/mocks"
	"github.com/NotFound1911/filestore/app/account/service"
	"github.com/NotFound1911/filestore/errs"
	"github.com/NotFound1911/filestore/pkg/server"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHandler_SignUp(t *testing.T) {
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) userv1.UserServiceClient

		reqBuilder func(t *testing.T) *http.Request

		wantCode int
		wantBody server.Result
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) userv1.UserServiceClient {
				userSvc := usermocks.NewMockUserServiceClient(ctrl)
				userSvc.EXPECT().Signup(gomock.Any(), &userv1.SignupReq{
					User: &userv1.User{
						Email:    "123@qq.com",
						Password: "hello#world123",
					},
				}).Return(&userv1.SignupResp{}, nil)
				return userSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost,
					"/api/storage/v1/users/signup", bytes.NewReader([]byte(`{
"email": "123@qq.com",
"password": "hello#world123",
"confirm_password": "hello#world123"
}`)))
				assert.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantCode: http.StatusOK,
			wantBody: server.Result{
				Msg:  "注册成功",
				Code: 2000,
			},
		},
		{
			name: "Bind出错",
			mock: func(ctrl *gomock.Controller) userv1.UserServiceClient {
				userSvc := usermocks.NewMockUserServiceClient(ctrl)
				return userSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost,
					"/api/storage/v1/users/signup", bytes.NewReader([]byte(`{
		"email": "123@qq.com",
		"password": "hello#world"
		}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},

			wantCode: http.StatusOK,
			wantBody: server.Result{
				Code: errs.UserInvalidInput,
				Msg:  "两次输入的密码不相等",
			},
		},
		{
			name: "邮箱格式不对",
			mock: func(ctrl *gomock.Controller) userv1.UserServiceClient {
				userSvc := usermocks.NewMockUserServiceClient(ctrl)
				return userSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost,
					"/api/storage/v1/users/signup", bytes.NewReader([]byte(`{
		"email": "123@",
		"password": "hello#world123",
		"confirm_password": "hello#world123"
		}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},

			wantCode: http.StatusOK,
			wantBody: server.Result{
				Code: errs.UserInvalidInput,
				Msg:  "非法邮箱格式",
			},
		},
		{
			name: "两次密码输入不同",
			mock: func(ctrl *gomock.Controller) userv1.UserServiceClient {
				userSvc := usermocks.NewMockUserServiceClient(ctrl)
				return userSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost,
					"/api/storage/v1/users/signup", bytes.NewReader([]byte(`{
		"email": "123@qq.com",
		"password": "hello#world123455",
		"confirm_password": "hello#world123"
		}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},

			wantCode: http.StatusOK,
			wantBody: server.Result{
				Code: errs.UserInvalidInput,
				Msg:  "两次输入的密码不相等",
			},
		},
		{
			name: "密码格式不对",
			mock: func(ctrl *gomock.Controller) userv1.UserServiceClient {
				userSvc := usermocks.NewMockUserServiceClient(ctrl)
				return userSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost,
					"/api/storage/v1/users/signup", bytes.NewReader([]byte(`{
		"email": "123@qq.com",
		"password": "hello",
		"confirm_password": "hello"
		}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},

			wantCode: http.StatusOK,
			wantBody: server.Result{
				Code: errs.UserInvalidInput,
				Msg:  "密码必须包含字母、数字、特殊字符,并且不少于八位",
			},
		},
		{
			name: "系统错误",
			mock: func(ctrl *gomock.Controller) userv1.UserServiceClient {
				userSvc := usermocks.NewMockUserServiceClient(ctrl)
				userSvc.EXPECT().Signup(gomock.Any(), &userv1.SignupReq{
					User: &userv1.User{
						Email:    "123@qq.com",
						Password: "hello#world123",
					},
				}).Return(&userv1.SignupResp{}, errors.New("db错误"))
				return userSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost,
					"/api/storage/v1/users/signup", bytes.NewReader([]byte(`{
		"email": "123@qq.com",
		"password": "hello#world123",
		"confirm_password": "hello#world123"
		}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},
			wantBody: server.Result{
				Code: errs.UserInternalServerError,
				Msg:  "注册失败",
				Data: "db错误",
			},
			wantCode: http.StatusOK,
		},
		{
			name: "邮箱冲突",
			mock: func(ctrl *gomock.Controller) userv1.UserServiceClient {
				userSvc := usermocks.NewMockUserServiceClient(ctrl)
				userSvc.EXPECT().Signup(gomock.Any(), &userv1.SignupReq{
					User: &userv1.User{
						Email:    "123@qq.com",
						Password: "hello#world123",
					},
				}).Return(&userv1.SignupResp{}, service.ErrDuplicateEmail)
				return userSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost,
					"/api/storage/v1/users/signup", bytes.NewReader([]byte(`{
		"email": "123@qq.com",
		"password": "hello#world123",
		"confirm_password": "hello#world123"
		}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},
			wantBody: server.Result{
				Code: errs.UserInternalServerError,
				Msg:  "注册失败",
				Data: "邮箱或电话冲突",
			},
			wantCode: http.StatusOK,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			userSvc := tc.mock(ctrl)
			hdl := NewUserHandler(userSvc, nil)

			gin.SetMode(gin.TestMode)
			core := gin.Default()
			hdl.RegisterUserRoutes(core)
			// 准备Req和记录的 recorder
			req := tc.reqBuilder(t)
			recorder := httptest.NewRecorder()

			// 执行
			core.ServeHTTP(recorder, req)

			// 断言结果
			assert.Equal(t, tc.wantCode, recorder.Code)
			var res server.Result
			err := json.Unmarshal(recorder.Body.Bytes(), &res)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantBody, res)
		})
	}
}
