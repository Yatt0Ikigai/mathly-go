// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/go-co-op/gocron/v2 (interfaces: Scheduler,Job)
//
// Generated by this command:
//
//	mockgen -destination=scheduler_mock.go -package=mocks github.com/go-co-op/gocron/v2 Scheduler,Job
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"
	time "time"

	gocron "github.com/go-co-op/gocron/v2"
	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockScheduler is a mock of Scheduler interface.
type MockScheduler struct {
	ctrl     *gomock.Controller
	recorder *MockSchedulerMockRecorder
	isgomock struct{}
}

// MockSchedulerMockRecorder is the mock recorder for MockScheduler.
type MockSchedulerMockRecorder struct {
	mock *MockScheduler
}

// NewMockScheduler creates a new mock instance.
func NewMockScheduler(ctrl *gomock.Controller) *MockScheduler {
	mock := &MockScheduler{ctrl: ctrl}
	mock.recorder = &MockSchedulerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockScheduler) EXPECT() *MockSchedulerMockRecorder {
	return m.recorder
}

// Jobs mocks base method.
func (m *MockScheduler) Jobs() []gocron.Job {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Jobs")
	ret0, _ := ret[0].([]gocron.Job)
	return ret0
}

// Jobs indicates an expected call of Jobs.
func (mr *MockSchedulerMockRecorder) Jobs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Jobs", reflect.TypeOf((*MockScheduler)(nil).Jobs))
}

// JobsWaitingInQueue mocks base method.
func (m *MockScheduler) JobsWaitingInQueue() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "JobsWaitingInQueue")
	ret0, _ := ret[0].(int)
	return ret0
}

// JobsWaitingInQueue indicates an expected call of JobsWaitingInQueue.
func (mr *MockSchedulerMockRecorder) JobsWaitingInQueue() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "JobsWaitingInQueue", reflect.TypeOf((*MockScheduler)(nil).JobsWaitingInQueue))
}

// NewJob mocks base method.
func (m *MockScheduler) NewJob(arg0 gocron.JobDefinition, arg1 gocron.Task, arg2 ...gocron.JobOption) (gocron.Job, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "NewJob", varargs...)
	ret0, _ := ret[0].(gocron.Job)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewJob indicates an expected call of NewJob.
func (mr *MockSchedulerMockRecorder) NewJob(arg0, arg1 any, arg2 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewJob", reflect.TypeOf((*MockScheduler)(nil).NewJob), varargs...)
}

// RemoveByTags mocks base method.
func (m *MockScheduler) RemoveByTags(arg0 ...string) {
	m.ctrl.T.Helper()
	varargs := []any{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "RemoveByTags", varargs...)
}

// RemoveByTags indicates an expected call of RemoveByTags.
func (mr *MockSchedulerMockRecorder) RemoveByTags(arg0 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveByTags", reflect.TypeOf((*MockScheduler)(nil).RemoveByTags), arg0...)
}

// RemoveJob mocks base method.
func (m *MockScheduler) RemoveJob(arg0 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveJob", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveJob indicates an expected call of RemoveJob.
func (mr *MockSchedulerMockRecorder) RemoveJob(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveJob", reflect.TypeOf((*MockScheduler)(nil).RemoveJob), arg0)
}

// Shutdown mocks base method.
func (m *MockScheduler) Shutdown() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Shutdown")
	ret0, _ := ret[0].(error)
	return ret0
}

// Shutdown indicates an expected call of Shutdown.
func (mr *MockSchedulerMockRecorder) Shutdown() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Shutdown", reflect.TypeOf((*MockScheduler)(nil).Shutdown))
}

// Start mocks base method.
func (m *MockScheduler) Start() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Start")
}

// Start indicates an expected call of Start.
func (mr *MockSchedulerMockRecorder) Start() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockScheduler)(nil).Start))
}

// StopJobs mocks base method.
func (m *MockScheduler) StopJobs() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StopJobs")
	ret0, _ := ret[0].(error)
	return ret0
}

// StopJobs indicates an expected call of StopJobs.
func (mr *MockSchedulerMockRecorder) StopJobs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StopJobs", reflect.TypeOf((*MockScheduler)(nil).StopJobs))
}

// Update mocks base method.
func (m *MockScheduler) Update(arg0 uuid.UUID, arg1 gocron.JobDefinition, arg2 gocron.Task, arg3 ...gocron.JobOption) (gocron.Job, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1, arg2}
	for _, a := range arg3 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Update", varargs...)
	ret0, _ := ret[0].(gocron.Job)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockSchedulerMockRecorder) Update(arg0, arg1, arg2 any, arg3 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1, arg2}, arg3...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockScheduler)(nil).Update), varargs...)
}

// MockJob is a mock of Job interface.
type MockJob struct {
	ctrl     *gomock.Controller
	recorder *MockJobMockRecorder
	isgomock struct{}
}

// MockJobMockRecorder is the mock recorder for MockJob.
type MockJobMockRecorder struct {
	mock *MockJob
}

// NewMockJob creates a new mock instance.
func NewMockJob(ctrl *gomock.Controller) *MockJob {
	mock := &MockJob{ctrl: ctrl}
	mock.recorder = &MockJobMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockJob) EXPECT() *MockJobMockRecorder {
	return m.recorder
}

// ID mocks base method.
func (m *MockJob) ID() uuid.UUID {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ID")
	ret0, _ := ret[0].(uuid.UUID)
	return ret0
}

// ID indicates an expected call of ID.
func (mr *MockJobMockRecorder) ID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ID", reflect.TypeOf((*MockJob)(nil).ID))
}

// LastRun mocks base method.
func (m *MockJob) LastRun() (time.Time, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LastRun")
	ret0, _ := ret[0].(time.Time)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LastRun indicates an expected call of LastRun.
func (mr *MockJobMockRecorder) LastRun() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LastRun", reflect.TypeOf((*MockJob)(nil).LastRun))
}

// Name mocks base method.
func (m *MockJob) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockJobMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockJob)(nil).Name))
}

// NextRun mocks base method.
func (m *MockJob) NextRun() (time.Time, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NextRun")
	ret0, _ := ret[0].(time.Time)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NextRun indicates an expected call of NextRun.
func (mr *MockJobMockRecorder) NextRun() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NextRun", reflect.TypeOf((*MockJob)(nil).NextRun))
}

// NextRuns mocks base method.
func (m *MockJob) NextRuns(arg0 int) ([]time.Time, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NextRuns", arg0)
	ret0, _ := ret[0].([]time.Time)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NextRuns indicates an expected call of NextRuns.
func (mr *MockJobMockRecorder) NextRuns(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NextRuns", reflect.TypeOf((*MockJob)(nil).NextRuns), arg0)
}

// RunNow mocks base method.
func (m *MockJob) RunNow() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RunNow")
	ret0, _ := ret[0].(error)
	return ret0
}

// RunNow indicates an expected call of RunNow.
func (mr *MockJobMockRecorder) RunNow() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunNow", reflect.TypeOf((*MockJob)(nil).RunNow))
}

// Tags mocks base method.
func (m *MockJob) Tags() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Tags")
	ret0, _ := ret[0].([]string)
	return ret0
}

// Tags indicates an expected call of Tags.
func (mr *MockJobMockRecorder) Tags() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Tags", reflect.TypeOf((*MockJob)(nil).Tags))
}
