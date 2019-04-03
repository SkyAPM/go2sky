// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/tetratelabs/go2sky/reporter/grpc/register (interfaces: RegisterClient)

// Package mock_register is a generated GoMock package.
package mock_register

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	common "github.com/tetratelabs/go2sky/reporter/grpc/common"
	register "github.com/tetratelabs/go2sky/reporter/grpc/register"
	grpc "google.golang.org/grpc"
	reflect "reflect"
)

// MockRegisterClient is a mock of RegisterClient interface
type MockRegisterClient struct {
	ctrl     *gomock.Controller
	recorder *MockRegisterClientMockRecorder
}

// MockRegisterClientMockRecorder is the mock recorder for MockRegisterClient
type MockRegisterClientMockRecorder struct {
	mock *MockRegisterClient
}

// NewMockRegisterClient creates a new mock instance
func NewMockRegisterClient(ctrl *gomock.Controller) *MockRegisterClient {
	mock := &MockRegisterClient{ctrl: ctrl}
	mock.recorder = &MockRegisterClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRegisterClient) EXPECT() *MockRegisterClientMockRecorder {
	return m.recorder
}

// DoEndpointRegister mocks base method
func (m *MockRegisterClient) DoEndpointRegister(arg0 context.Context, arg1 *register.Enpoints, arg2 ...grpc.CallOption) (*register.EndpointMapping, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DoEndpointRegister", varargs...)
	ret0, _ := ret[0].(*register.EndpointMapping)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DoEndpointRegister indicates an expected call of DoEndpointRegister
func (mr *MockRegisterClientMockRecorder) DoEndpointRegister(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DoEndpointRegister", reflect.TypeOf((*MockRegisterClient)(nil).DoEndpointRegister), varargs...)
}

// DoNetworkAddressRegister mocks base method
func (m *MockRegisterClient) DoNetworkAddressRegister(arg0 context.Context, arg1 *register.NetAddresses, arg2 ...grpc.CallOption) (*register.NetAddressMapping, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DoNetworkAddressRegister", varargs...)
	ret0, _ := ret[0].(*register.NetAddressMapping)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DoNetworkAddressRegister indicates an expected call of DoNetworkAddressRegister
func (mr *MockRegisterClientMockRecorder) DoNetworkAddressRegister(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DoNetworkAddressRegister", reflect.TypeOf((*MockRegisterClient)(nil).DoNetworkAddressRegister), varargs...)
}

// DoServiceAndNetworkAddressMappingRegister mocks base method
func (m *MockRegisterClient) DoServiceAndNetworkAddressMappingRegister(arg0 context.Context, arg1 *register.ServiceAndNetworkAddressMappings, arg2 ...grpc.CallOption) (*common.Commands, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DoServiceAndNetworkAddressMappingRegister", varargs...)
	ret0, _ := ret[0].(*common.Commands)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DoServiceAndNetworkAddressMappingRegister indicates an expected call of DoServiceAndNetworkAddressMappingRegister
func (mr *MockRegisterClientMockRecorder) DoServiceAndNetworkAddressMappingRegister(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DoServiceAndNetworkAddressMappingRegister", reflect.TypeOf((*MockRegisterClient)(nil).DoServiceAndNetworkAddressMappingRegister), varargs...)
}

// DoServiceInstanceRegister mocks base method
func (m *MockRegisterClient) DoServiceInstanceRegister(arg0 context.Context, arg1 *register.ServiceInstances, arg2 ...grpc.CallOption) (*register.ServiceInstanceRegisterMapping, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DoServiceInstanceRegister", varargs...)
	ret0, _ := ret[0].(*register.ServiceInstanceRegisterMapping)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DoServiceInstanceRegister indicates an expected call of DoServiceInstanceRegister
func (mr *MockRegisterClientMockRecorder) DoServiceInstanceRegister(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DoServiceInstanceRegister", reflect.TypeOf((*MockRegisterClient)(nil).DoServiceInstanceRegister), varargs...)
}

// DoServiceRegister mocks base method
func (m *MockRegisterClient) DoServiceRegister(arg0 context.Context, arg1 *register.Services, arg2 ...grpc.CallOption) (*register.ServiceRegisterMapping, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DoServiceRegister", varargs...)
	ret0, _ := ret[0].(*register.ServiceRegisterMapping)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DoServiceRegister indicates an expected call of DoServiceRegister
func (mr *MockRegisterClientMockRecorder) DoServiceRegister(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DoServiceRegister", reflect.TypeOf((*MockRegisterClient)(nil).DoServiceRegister), varargs...)
}
