// Code generated by MockGen. DO NOT EDIT.
// Source: ./api/proto/gen/account/v1/account_grpc.pb.go

// Package accountmocks is a generated GoMock package.
package accountmocks

import (
	context "context"
	reflect "reflect"

	accountv1 "github.com/NotFound1911/filestore/api/proto/gen/account/v1"
	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockAccountServiceClient is a mock of AccountServiceClient interface.
type MockAccountServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockAccountServiceClientMockRecorder
}

// MockAccountServiceClientMockRecorder is the mock recorder for MockAccountServiceClient.
type MockAccountServiceClientMockRecorder struct {
	mock *MockAccountServiceClient
}

// NewMockAccountServiceClient creates a new mock instance.
func NewMockAccountServiceClient(ctrl *gomock.Controller) *MockAccountServiceClient {
	mock := &MockAccountServiceClient{ctrl: ctrl}
	mock.recorder = &MockAccountServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAccountServiceClient) EXPECT() *MockAccountServiceClientMockRecorder {
	return m.recorder
}

// Login mocks base method.
func (m *MockAccountServiceClient) Login(ctx context.Context, in *accountv1.LoginReq, opts ...grpc.CallOption) (*accountv1.LoginResp, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Login", varargs...)
	ret0, _ := ret[0].(*accountv1.LoginResp)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockAccountServiceClientMockRecorder) Login(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockAccountServiceClient)(nil).Login), varargs...)
}

// Profile mocks base method.
func (m *MockAccountServiceClient) Profile(ctx context.Context, in *accountv1.ProfileReq, opts ...grpc.CallOption) (*accountv1.ProfileResp, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Profile", varargs...)
	ret0, _ := ret[0].(*accountv1.ProfileResp)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Profile indicates an expected call of Profile.
func (mr *MockAccountServiceClientMockRecorder) Profile(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Profile", reflect.TypeOf((*MockAccountServiceClient)(nil).Profile), varargs...)
}

// Signup mocks base method.
func (m *MockAccountServiceClient) Signup(ctx context.Context, in *accountv1.SignupReq, opts ...grpc.CallOption) (*accountv1.SignupResp, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Signup", varargs...)
	ret0, _ := ret[0].(*accountv1.SignupResp)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Signup indicates an expected call of Signup.
func (mr *MockAccountServiceClientMockRecorder) Signup(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Signup", reflect.TypeOf((*MockAccountServiceClient)(nil).Signup), varargs...)
}

// MockAccountServiceServer is a mock of AccountServiceServer interface.
type MockAccountServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockAccountServiceServerMockRecorder
}

// MockAccountServiceServerMockRecorder is the mock recorder for MockAccountServiceServer.
type MockAccountServiceServerMockRecorder struct {
	mock *MockAccountServiceServer
}

// NewMockAccountServiceServer creates a new mock instance.
func NewMockAccountServiceServer(ctrl *gomock.Controller) *MockAccountServiceServer {
	mock := &MockAccountServiceServer{ctrl: ctrl}
	mock.recorder = &MockAccountServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAccountServiceServer) EXPECT() *MockAccountServiceServerMockRecorder {
	return m.recorder
}

// Login mocks base method.
func (m *MockAccountServiceServer) Login(arg0 context.Context, arg1 *accountv1.LoginReq) (*accountv1.LoginResp, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", arg0, arg1)
	ret0, _ := ret[0].(*accountv1.LoginResp)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockAccountServiceServerMockRecorder) Login(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockAccountServiceServer)(nil).Login), arg0, arg1)
}

// Profile mocks base method.
func (m *MockAccountServiceServer) Profile(arg0 context.Context, arg1 *accountv1.ProfileReq) (*accountv1.ProfileResp, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Profile", arg0, arg1)
	ret0, _ := ret[0].(*accountv1.ProfileResp)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Profile indicates an expected call of Profile.
func (mr *MockAccountServiceServerMockRecorder) Profile(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Profile", reflect.TypeOf((*MockAccountServiceServer)(nil).Profile), arg0, arg1)
}

// Signup mocks base method.
func (m *MockAccountServiceServer) Signup(arg0 context.Context, arg1 *accountv1.SignupReq) (*accountv1.SignupResp, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Signup", arg0, arg1)
	ret0, _ := ret[0].(*accountv1.SignupResp)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Signup indicates an expected call of Signup.
func (mr *MockAccountServiceServerMockRecorder) Signup(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Signup", reflect.TypeOf((*MockAccountServiceServer)(nil).Signup), arg0, arg1)
}

// mustEmbedUnimplementedAccountServiceServer mocks base method.
func (m *MockAccountServiceServer) mustEmbedUnimplementedAccountServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedAccountServiceServer")
}

// mustEmbedUnimplementedAccountServiceServer indicates an expected call of mustEmbedUnimplementedAccountServiceServer.
func (mr *MockAccountServiceServerMockRecorder) mustEmbedUnimplementedAccountServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedAccountServiceServer", reflect.TypeOf((*MockAccountServiceServer)(nil).mustEmbedUnimplementedAccountServiceServer))
}

// MockUnsafeAccountServiceServer is a mock of UnsafeAccountServiceServer interface.
type MockUnsafeAccountServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockUnsafeAccountServiceServerMockRecorder
}

// MockUnsafeAccountServiceServerMockRecorder is the mock recorder for MockUnsafeAccountServiceServer.
type MockUnsafeAccountServiceServerMockRecorder struct {
	mock *MockUnsafeAccountServiceServer
}

// NewMockUnsafeAccountServiceServer creates a new mock instance.
func NewMockUnsafeAccountServiceServer(ctrl *gomock.Controller) *MockUnsafeAccountServiceServer {
	mock := &MockUnsafeAccountServiceServer{ctrl: ctrl}
	mock.recorder = &MockUnsafeAccountServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnsafeAccountServiceServer) EXPECT() *MockUnsafeAccountServiceServerMockRecorder {
	return m.recorder
}

// mustEmbedUnimplementedAccountServiceServer mocks base method.
func (m *MockUnsafeAccountServiceServer) mustEmbedUnimplementedAccountServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedAccountServiceServer")
}

// mustEmbedUnimplementedAccountServiceServer indicates an expected call of mustEmbedUnimplementedAccountServiceServer.
func (mr *MockUnsafeAccountServiceServerMockRecorder) mustEmbedUnimplementedAccountServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedAccountServiceServer", reflect.TypeOf((*MockUnsafeAccountServiceServer)(nil).mustEmbedUnimplementedAccountServiceServer))
}
