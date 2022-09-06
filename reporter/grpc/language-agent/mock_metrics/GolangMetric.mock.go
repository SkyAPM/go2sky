package mock_v3

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	grpc0 "google.golang.org/grpc"
	v3 "skywalking.apache.org/repo/goapi/collect/common/v3"
	agentv3 "skywalking.apache.org/repo/goapi/collect/language/agent/v3"
)

// MockGolangMetricReportServiceClient is a mock of GolangMetricReportServiceClient interface.
type MockGolangMetricReportServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockGolangMetricReportServiceClientMockRecorder
}

// MockGolangMetricReportServiceClientMockRecorder is the mock recorder for MockGolangMetricReportServiceClient.
type MockGolangMetricReportServiceClientMockRecorder struct {
	mock *MockGolangMetricReportServiceClient
}

// NewMockGolangMetricReportServiceClient creates a new mock instance.
func NewMockGolangMetricReportServiceClient(ctrl *gomock.Controller) *MockGolangMetricReportServiceClient {
	mock := &MockGolangMetricReportServiceClient{ctrl: ctrl}
	mock.recorder = &MockGolangMetricReportServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGolangMetricReportServiceClient) EXPECT() *MockGolangMetricReportServiceClientMockRecorder {
	return m.recorder
}

// Collect mocks base method.
func (m *MockGolangMetricReportServiceClient) Collect(ctx context.Context, in *agentv3.GolangMetricCollection, opts ...grpc0.CallOption) (*v3.Commands, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Collect", varargs...)
	ret0, _ := ret[0].(*v3.Commands)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Collect indicates an expected call of Collect.
func (mr *MockGolangMetricReportServiceClientMockRecorder) Collect(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Collect", reflect.TypeOf((*MockGolangMetricReportServiceClient)(nil).Collect), varargs...)
}

// MockGolangMetricReportServiceServer is a mock of GolangMetricReportServiceServer interface.
type MockGolangMetricReportServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockGolangMetricReportServiceServerMockRecorder
}

// MockGolangMetricReportServiceServerMockRecorder is the mock recorder for MockGolangMetricReportServiceServer.
type MockGolangMetricReportServiceServerMockRecorder struct {
	mock *MockGolangMetricReportServiceServer
}

// NewMockGolangMetricReportServiceServer creates a new mock instance.
func NewMockGolangMetricReportServiceServer(ctrl *gomock.Controller) *MockGolangMetricReportServiceServer {
	mock := &MockGolangMetricReportServiceServer{ctrl: ctrl}
	mock.recorder = &MockGolangMetricReportServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGolangMetricReportServiceServer) EXPECT() *MockGolangMetricReportServiceServerMockRecorder {
	return m.recorder
}

// Collect mocks base method.
func (m *MockGolangMetricReportServiceServer) Collect(arg0 context.Context, arg1 *agentv3.GolangMetricCollection) (*v3.Commands, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Collect", arg0, arg1)
	ret0, _ := ret[0].(*v3.Commands)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Collect indicates an expected call of Collect.
func (mr *MockGolangMetricReportServiceServerMockRecorder) Collect(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Collect", reflect.TypeOf((*MockGolangMetricReportServiceServer)(nil).Collect), arg0, arg1)
}

// mustEmbedUnimplementedGolangMetricReportServiceServer mocks base method.
func (m *MockGolangMetricReportServiceServer) mustEmbedUnimplementedGolangMetricReportServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedGolangMetricReportServiceServer")
}

// mustEmbedUnimplementedGolangMetricReportServiceServer indicates an expected call of mustEmbedUnimplementedGolangMetricReportServiceServer.
func (mr *MockGolangMetricReportServiceServerMockRecorder) mustEmbedUnimplementedGolangMetricReportServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedGolangMetricReportServiceServer", reflect.TypeOf((*MockGolangMetricReportServiceServer)(nil).mustEmbedUnimplementedGolangMetricReportServiceServer))
}

// MockUnsafeGolangMetricReportServiceServer is a mock of UnsafeGolangMetricReportServiceServer interface.
type MockUnsafeGolangMetricReportServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockUnsafeGolangMetricReportServiceServerMockRecorder
}

// MockUnsafeGolangMetricReportServiceServerMockRecorder is the mock recorder for MockUnsafeGolangMetricReportServiceServer.
type MockUnsafeGolangMetricReportServiceServerMockRecorder struct {
	mock *MockUnsafeGolangMetricReportServiceServer
}

// NewMockUnsafeGolangMetricReportServiceServer creates a new mock instance.
func NewMockUnsafeGolangMetricReportServiceServer(ctrl *gomock.Controller) *MockUnsafeGolangMetricReportServiceServer {
	mock := &MockUnsafeGolangMetricReportServiceServer{ctrl: ctrl}
	mock.recorder = &MockUnsafeGolangMetricReportServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnsafeGolangMetricReportServiceServer) EXPECT() *MockUnsafeGolangMetricReportServiceServerMockRecorder {
	return m.recorder
}

// mustEmbedUnimplementedGolangMetricReportServiceServer mocks base method.
func (m *MockUnsafeGolangMetricReportServiceServer) mustEmbedUnimplementedGolangMetricReportServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedGolangMetricReportServiceServer")
}

// mustEmbedUnimplementedGolangMetricReportServiceServer indicates an expected call of mustEmbedUnimplementedGolangMetricReportServiceServer.
func (mr *MockUnsafeGolangMetricReportServiceServerMockRecorder) mustEmbedUnimplementedGolangMetricReportServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedGolangMetricReportServiceServer", reflect.TypeOf((*MockUnsafeGolangMetricReportServiceServer)(nil).mustEmbedUnimplementedGolangMetricReportServiceServer))
}
