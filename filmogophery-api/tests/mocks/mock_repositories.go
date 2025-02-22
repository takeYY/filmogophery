// Code generated by MockGen. DO NOT EDIT.
// Source: filmogophery/internal/app/repositories (interfaces: IGenreRepository,IImpressionRepository,IMediaRepository,IMovieRepository,IRecordRepository)
//
// Generated by this command:
//
//	mockgen -package=mocks -destination=tests/mocks/mock_repositories.go filmogophery/internal/app/repositories IGenreRepository,IImpressionRepository,IMediaRepository,IMovieRepository,IRecordRepository
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	repositories "filmogophery/internal/app/repositories"
	model "filmogophery/internal/pkg/gen/model"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockIGenreRepository is a mock of IGenreRepository interface.
type MockIGenreRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIGenreRepositoryMockRecorder
	isgomock struct{}
}

// MockIGenreRepositoryMockRecorder is the mock recorder for MockIGenreRepository.
type MockIGenreRepositoryMockRecorder struct {
	mock *MockIGenreRepository
}

// NewMockIGenreRepository creates a new mock instance.
func NewMockIGenreRepository(ctrl *gomock.Controller) *MockIGenreRepository {
	mock := &MockIGenreRepository{ctrl: ctrl}
	mock.recorder = &MockIGenreRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIGenreRepository) EXPECT() *MockIGenreRepositoryMockRecorder {
	return m.recorder
}

// FindByNames mocks base method.
func (m *MockIGenreRepository) FindByNames(ctx context.Context, names []string) ([]*model.Genre, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByNames", ctx, names)
	ret0, _ := ret[0].([]*model.Genre)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByNames indicates an expected call of FindByNames.
func (mr *MockIGenreRepositoryMockRecorder) FindByNames(ctx, names any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByNames", reflect.TypeOf((*MockIGenreRepository)(nil).FindByNames), ctx, names)
}

// MockIImpressionRepository is a mock of IImpressionRepository interface.
type MockIImpressionRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIImpressionRepositoryMockRecorder
	isgomock struct{}
}

// MockIImpressionRepositoryMockRecorder is the mock recorder for MockIImpressionRepository.
type MockIImpressionRepositoryMockRecorder struct {
	mock *MockIImpressionRepository
}

// NewMockIImpressionRepository creates a new mock instance.
func NewMockIImpressionRepository(ctrl *gomock.Controller) *MockIImpressionRepository {
	mock := &MockIImpressionRepository{ctrl: ctrl}
	mock.recorder = &MockIImpressionRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIImpressionRepository) EXPECT() *MockIImpressionRepositoryMockRecorder {
	return m.recorder
}

// FindAll mocks base method.
func (m *MockIImpressionRepository) FindAll(ctx context.Context) ([]*model.MovieImpression, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAll", ctx)
	ret0, _ := ret[0].([]*model.MovieImpression)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAll indicates an expected call of FindAll.
func (mr *MockIImpressionRepositoryMockRecorder) FindAll(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*MockIImpressionRepository)(nil).FindAll), ctx)
}

// Save mocks base method.
func (m *MockIImpressionRepository) Save(ctx context.Context, input repositories.SaveImpressionInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, input)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockIImpressionRepositoryMockRecorder) Save(ctx, input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockIImpressionRepository)(nil).Save), ctx, input)
}

// Update mocks base method.
func (m *MockIImpressionRepository) Update(ctx context.Context, input repositories.UpdateImpressionInput) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, input)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockIImpressionRepositoryMockRecorder) Update(ctx, input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockIImpressionRepository)(nil).Update), ctx, input)
}

// MockIMediaRepository is a mock of IMediaRepository interface.
type MockIMediaRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIMediaRepositoryMockRecorder
	isgomock struct{}
}

// MockIMediaRepositoryMockRecorder is the mock recorder for MockIMediaRepository.
type MockIMediaRepositoryMockRecorder struct {
	mock *MockIMediaRepository
}

// NewMockIMediaRepository creates a new mock instance.
func NewMockIMediaRepository(ctrl *gomock.Controller) *MockIMediaRepository {
	mock := &MockIMediaRepository{ctrl: ctrl}
	mock.recorder = &MockIMediaRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIMediaRepository) EXPECT() *MockIMediaRepositoryMockRecorder {
	return m.recorder
}

// FindAll mocks base method.
func (m *MockIMediaRepository) FindAll(ctx context.Context) ([]*model.WatchMedia, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAll", ctx)
	ret0, _ := ret[0].([]*model.WatchMedia)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAll indicates an expected call of FindAll.
func (mr *MockIMediaRepositoryMockRecorder) FindAll(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*MockIMediaRepository)(nil).FindAll), ctx)
}

// FindByCode mocks base method.
func (m *MockIMediaRepository) FindByCode(ctx context.Context, code *string) (*model.WatchMedia, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByCode", ctx, code)
	ret0, _ := ret[0].(*model.WatchMedia)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByCode indicates an expected call of FindByCode.
func (mr *MockIMediaRepositoryMockRecorder) FindByCode(ctx, code any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByCode", reflect.TypeOf((*MockIMediaRepository)(nil).FindByCode), ctx, code)
}

// MockIMovieRepository is a mock of IMovieRepository interface.
type MockIMovieRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIMovieRepositoryMockRecorder
	isgomock struct{}
}

// MockIMovieRepositoryMockRecorder is the mock recorder for MockIMovieRepository.
type MockIMovieRepositoryMockRecorder struct {
	mock *MockIMovieRepository
}

// NewMockIMovieRepository creates a new mock instance.
func NewMockIMovieRepository(ctrl *gomock.Controller) *MockIMovieRepository {
	mock := &MockIMovieRepository{ctrl: ctrl}
	mock.recorder = &MockIMovieRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIMovieRepository) EXPECT() *MockIMovieRepositoryMockRecorder {
	return m.recorder
}

// FindAll mocks base method.
func (m *MockIMovieRepository) FindAll(ctx context.Context) ([]*model.Movie, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAll", ctx)
	ret0, _ := ret[0].([]*model.Movie)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAll indicates an expected call of FindAll.
func (mr *MockIMovieRepositoryMockRecorder) FindAll(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*MockIMovieRepository)(nil).FindAll), ctx)
}

// FindByID mocks base method.
func (m *MockIMovieRepository) FindByID(ctx context.Context, id *int32) (*model.Movie, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", ctx, id)
	ret0, _ := ret[0].(*model.Movie)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID.
func (mr *MockIMovieRepositoryMockRecorder) FindByID(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockIMovieRepository)(nil).FindByID), ctx, id)
}

// Save mocks base method.
func (m *MockIMovieRepository) Save(ctx context.Context, input repositories.SaveMovieInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, input)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockIMovieRepositoryMockRecorder) Save(ctx, input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockIMovieRepository)(nil).Save), ctx, input)
}

// MockIRecordRepository is a mock of IRecordRepository interface.
type MockIRecordRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIRecordRepositoryMockRecorder
	isgomock struct{}
}

// MockIRecordRepositoryMockRecorder is the mock recorder for MockIRecordRepository.
type MockIRecordRepositoryMockRecorder struct {
	mock *MockIRecordRepository
}

// NewMockIRecordRepository creates a new mock instance.
func NewMockIRecordRepository(ctrl *gomock.Controller) *MockIRecordRepository {
	mock := &MockIRecordRepository{ctrl: ctrl}
	mock.recorder = &MockIRecordRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIRecordRepository) EXPECT() *MockIRecordRepositoryMockRecorder {
	return m.recorder
}

// FindAll mocks base method.
func (m *MockIRecordRepository) FindAll(ctx context.Context) ([]*model.MovieWatchRecord, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAll", ctx)
	ret0, _ := ret[0].([]*model.MovieWatchRecord)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAll indicates an expected call of FindAll.
func (mr *MockIRecordRepositoryMockRecorder) FindAll(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*MockIRecordRepository)(nil).FindAll), ctx)
}

// FindByImpressionID mocks base method.
func (m *MockIRecordRepository) FindByImpressionID(ctx context.Context, id *int32) ([]*model.MovieWatchRecord, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByImpressionID", ctx, id)
	ret0, _ := ret[0].([]*model.MovieWatchRecord)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByImpressionID indicates an expected call of FindByImpressionID.
func (mr *MockIRecordRepositoryMockRecorder) FindByImpressionID(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByImpressionID", reflect.TypeOf((*MockIRecordRepository)(nil).FindByImpressionID), ctx, id)
}

// Save mocks base method.
func (m *MockIRecordRepository) Save(ctx context.Context, input repositories.SaveRecordInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, input)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockIRecordRepositoryMockRecorder) Save(ctx, input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockIRecordRepository)(nil).Save), ctx, input)
}
