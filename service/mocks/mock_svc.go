// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mocks is a generated GoMock package.
package mocks

import (
	cache "IM-Backend/cache"
	dao "IM-Backend/dao"
	table "IM-Backend/model/table"
	context "context"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	gorm "gorm.io/gorm"
)

// MockGormWriter is a mock of GormWriter interface.
type MockGormWriter struct {
	ctrl     *gomock.Controller
	recorder *MockGormWriterMockRecorder
}

// MockGormWriterMockRecorder is the mock recorder for MockGormWriter.
type MockGormWriterMockRecorder struct {
	mock *MockGormWriter
}

// NewMockGormWriter creates a new mock instance.
func NewMockGormWriter(ctrl *gomock.Controller) *MockGormWriter {
	mock := &MockGormWriter{ctrl: ctrl}
	mock.recorder = &MockGormWriterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGormWriter) EXPECT() *MockGormWriterMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockGormWriter) Create(ctx context.Context, svc string, t dao.Table) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, svc, t)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockGormWriterMockRecorder) Create(ctx, svc, t interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockGormWriter)(nil).Create), ctx, svc, t)
}

// Delete mocks base method.
func (m *MockGormWriter) Delete(ctx context.Context, svc string, t dao.Table, where map[string]interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, svc, t, where)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockGormWriterMockRecorder) Delete(ctx, svc, t, where interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockGormWriter)(nil).Delete), ctx, svc, t, where)
}

// GetGormTx mocks base method.
func (m *MockGormWriter) GetGormTx(ctx context.Context) *gorm.DB {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGormTx", ctx)
	ret0, _ := ret[0].(*gorm.DB)
	return ret0
}

// GetGormTx indicates an expected call of GetGormTx.
func (mr *MockGormWriterMockRecorder) GetGormTx(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGormTx", reflect.TypeOf((*MockGormWriter)(nil).GetGormTx), ctx)
}

// InTx mocks base method.
func (m *MockGormWriter) InTx(ctx context.Context, f func(context.Context) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InTx", ctx, f)
	ret0, _ := ret[0].(error)
	return ret0
}

// InTx indicates an expected call of InTx.
func (mr *MockGormWriterMockRecorder) InTx(ctx, f interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InTx", reflect.TypeOf((*MockGormWriter)(nil).InTx), ctx, f)
}

// Update mocks base method.
func (m *MockGormWriter) Update(ctx context.Context, svc string, t dao.Table) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, svc, t)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockGormWriterMockRecorder) Update(ctx, svc, t interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockGormWriter)(nil).Update), ctx, svc, t)
}

// MockGormPostReader is a mock of GormPostReader interface.
type MockGormPostReader struct {
	ctrl     *gomock.Controller
	recorder *MockGormPostReaderMockRecorder
}

// MockGormPostReaderMockRecorder is the mock recorder for MockGormPostReader.
type MockGormPostReaderMockRecorder struct {
	mock *MockGormPostReader
}

// NewMockGormPostReader creates a new mock instance.
func NewMockGormPostReader(ctrl *gomock.Controller) *MockGormPostReader {
	mock := &MockGormPostReader{ctrl: ctrl}
	mock.recorder = &MockGormPostReaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGormPostReader) EXPECT() *MockGormPostReaderMockRecorder {
	return m.recorder
}

// CheckPostExist mocks base method.
func (m *MockGormPostReader) CheckPostExist(ctx context.Context, svc string, id uint64) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckPostExist", ctx, svc, id)
	ret0, _ := ret[0].(bool)
	return ret0
}

// CheckPostExist indicates an expected call of CheckPostExist.
func (mr *MockGormPostReaderMockRecorder) CheckPostExist(ctx, svc, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckPostExist", reflect.TypeOf((*MockGormPostReader)(nil).CheckPostExist), ctx, svc, id)
}

// GetPostCommentIds mocks base method.
func (m *MockGormPostReader) GetPostCommentIds(ctx context.Context, svc string, id uint64) ([]uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPostCommentIds", ctx, svc, id)
	ret0, _ := ret[0].([]uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPostCommentIds indicates an expected call of GetPostCommentIds.
func (mr *MockGormPostReaderMockRecorder) GetPostCommentIds(ctx, svc, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPostCommentIds", reflect.TypeOf((*MockGormPostReader)(nil).GetPostCommentIds), ctx, svc, id)
}

// GetPostInfos mocks base method.
func (m *MockGormPostReader) GetPostInfos(ctx context.Context, svc string, ids ...uint64) ([]table.PostInfo, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, svc}
	for _, a := range ids {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetPostInfos", varargs...)
	ret0, _ := ret[0].([]table.PostInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPostInfos indicates an expected call of GetPostInfos.
func (mr *MockGormPostReaderMockRecorder) GetPostInfos(ctx, svc interface{}, ids ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, svc}, ids...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPostInfos", reflect.TypeOf((*MockGormPostReader)(nil).GetPostInfos), varargs...)
}

// GetPostInfosByExtra mocks base method.
func (m *MockGormPostReader) GetPostInfosByExtra(ctx context.Context, svc, key string, value interface{}) ([]table.PostInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPostInfosByExtra", ctx, svc, key, value)
	ret0, _ := ret[0].([]table.PostInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPostInfosByExtra indicates an expected call of GetPostInfosByExtra.
func (mr *MockGormPostReaderMockRecorder) GetPostInfosByExtra(ctx, svc, key, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPostInfosByExtra", reflect.TypeOf((*MockGormPostReader)(nil).GetPostInfosByExtra), ctx, svc, key, value)
}

// GetPostLike mocks base method.
func (m *MockGormPostReader) GetPostLike(ctx context.Context, svc string, id uint64) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPostLike", ctx, svc, id)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPostLike indicates an expected call of GetPostLike.
func (mr *MockGormPostReaderMockRecorder) GetPostLike(ctx, svc, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPostLike", reflect.TypeOf((*MockGormPostReader)(nil).GetPostLike), ctx, svc, id)
}

// MockGormCommentReader is a mock of GormCommentReader interface.
type MockGormCommentReader struct {
	ctrl     *gomock.Controller
	recorder *MockGormCommentReaderMockRecorder
}

// MockGormCommentReaderMockRecorder is the mock recorder for MockGormCommentReader.
type MockGormCommentReaderMockRecorder struct {
	mock *MockGormCommentReader
}

// NewMockGormCommentReader creates a new mock instance.
func NewMockGormCommentReader(ctrl *gomock.Controller) *MockGormCommentReader {
	mock := &MockGormCommentReader{ctrl: ctrl}
	mock.recorder = &MockGormCommentReaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGormCommentReader) EXPECT() *MockGormCommentReaderMockRecorder {
	return m.recorder
}

// CheckCommentExist mocks base method.
func (m *MockGormCommentReader) CheckCommentExist(ctx context.Context, svc string, commentID ...uint64) map[uint64]bool {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, svc}
	for _, a := range commentID {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CheckCommentExist", varargs...)
	ret0, _ := ret[0].(map[uint64]bool)
	return ret0
}

// CheckCommentExist indicates an expected call of CheckCommentExist.
func (mr *MockGormCommentReaderMockRecorder) CheckCommentExist(ctx, svc interface{}, commentID ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, svc}, commentID...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckCommentExist", reflect.TypeOf((*MockGormCommentReader)(nil).CheckCommentExist), varargs...)
}

// GetChildCommentCnt mocks base method.
func (m *MockGormCommentReader) GetChildCommentCnt(ctx context.Context, svc string, commentID ...uint64) (map[uint64]int, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, svc}
	for _, a := range commentID {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetChildCommentCnt", varargs...)
	ret0, _ := ret[0].(map[uint64]int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetChildCommentCnt indicates an expected call of GetChildCommentCnt.
func (mr *MockGormCommentReaderMockRecorder) GetChildCommentCnt(ctx, svc interface{}, commentID ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, svc}, commentID...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChildCommentCnt", reflect.TypeOf((*MockGormCommentReader)(nil).GetChildCommentCnt), varargs...)
}

// GetChildCommentIDAfterCursor mocks base method.
func (m *MockGormCommentReader) GetChildCommentIDAfterCursor(ctx context.Context, svc string, fatherID uint64, cursor time.Time, limit uint) ([]uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChildCommentIDAfterCursor", ctx, svc, fatherID, cursor, limit)
	ret0, _ := ret[0].([]uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetChildCommentIDAfterCursor indicates an expected call of GetChildCommentIDAfterCursor.
func (mr *MockGormCommentReaderMockRecorder) GetChildCommentIDAfterCursor(ctx, svc, fatherID, cursor, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChildCommentIDAfterCursor", reflect.TypeOf((*MockGormCommentReader)(nil).GetChildCommentIDAfterCursor), ctx, svc, fatherID, cursor, limit)
}

// GetCommentInfosByID mocks base method.
func (m *MockGormCommentReader) GetCommentInfosByID(ctx context.Context, svc string, ids ...uint64) ([]table.PostCommentInfo, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, svc}
	for _, a := range ids {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetCommentInfosByID", varargs...)
	ret0, _ := ret[0].([]table.PostCommentInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCommentInfosByID indicates an expected call of GetCommentInfosByID.
func (mr *MockGormCommentReaderMockRecorder) GetCommentInfosByID(ctx, svc interface{}, ids ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, svc}, ids...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCommentInfosByID", reflect.TypeOf((*MockGormCommentReader)(nil).GetCommentInfosByID), varargs...)
}

// GetCommentLike mocks base method.
func (m *MockGormCommentReader) GetCommentLike(ctx context.Context, svc string, ids ...uint64) (map[uint64]int64, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, svc}
	for _, a := range ids {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetCommentLike", varargs...)
	ret0, _ := ret[0].(map[uint64]int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCommentLike indicates an expected call of GetCommentLike.
func (mr *MockGormCommentReaderMockRecorder) GetCommentLike(ctx, svc interface{}, ids ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, svc}, ids...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCommentLike", reflect.TypeOf((*MockGormCommentReader)(nil).GetCommentLike), varargs...)
}

// GetUserIDByCommentID mocks base method.
func (m *MockGormCommentReader) GetUserIDByCommentID(ctx context.Context, svc string, commentID uint64) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserIDByCommentID", ctx, svc, commentID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserIDByCommentID indicates an expected call of GetUserIDByCommentID.
func (mr *MockGormCommentReaderMockRecorder) GetUserIDByCommentID(ctx, svc, commentID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserIDByCommentID", reflect.TypeOf((*MockGormCommentReader)(nil).GetUserIDByCommentID), ctx, svc, commentID)
}

// MockWriteCache is a mock of WriteCache interface.
type MockWriteCache struct {
	ctrl     *gomock.Controller
	recorder *MockWriteCacheMockRecorder
}

// MockWriteCacheMockRecorder is the mock recorder for MockWriteCache.
type MockWriteCacheMockRecorder struct {
	mock *MockWriteCache
}

// NewMockWriteCache creates a new mock instance.
func NewMockWriteCache(ctrl *gomock.Controller) *MockWriteCache {
	mock := &MockWriteCache{ctrl: ctrl}
	mock.recorder = &MockWriteCacheMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWriteCache) EXPECT() *MockWriteCacheMockRecorder {
	return m.recorder
}

// AddKVToSet mocks base method.
func (m *MockWriteCache) AddKVToSet(ctx context.Context, expire time.Duration, kv ...cache.KV) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, expire}
	for _, a := range kv {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AddKVToSet", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddKVToSet indicates an expected call of AddKVToSet.
func (mr *MockWriteCacheMockRecorder) AddKVToSet(ctx, expire interface{}, kv ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, expire}, kv...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddKVToSet", reflect.TypeOf((*MockWriteCache)(nil).AddKVToSet), varargs...)
}

// DelKV mocks base method.
func (m *MockWriteCache) DelKV(ctx context.Context, kv cache.KV) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DelKV", ctx, kv)
	ret0, _ := ret[0].(error)
	return ret0
}

// DelKV indicates an expected call of DelKV.
func (mr *MockWriteCacheMockRecorder) DelKV(ctx, kv interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DelKV", reflect.TypeOf((*MockWriteCache)(nil).DelKV), ctx, kv)
}

// SetKV mocks base method.
func (m *MockWriteCache) SetKV(ctx context.Context, expire time.Duration, kv ...cache.KV) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, expire}
	for _, a := range kv {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SetKV", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetKV indicates an expected call of SetKV.
func (mr *MockWriteCacheMockRecorder) SetKV(ctx, expire interface{}, kv ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, expire}, kv...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetKV", reflect.TypeOf((*MockWriteCache)(nil).SetKV), varargs...)
}

// MockReadCache is a mock of ReadCache interface.
type MockReadCache struct {
	ctrl     *gomock.Controller
	recorder *MockReadCacheMockRecorder
}

// MockReadCacheMockRecorder is the mock recorder for MockReadCache.
type MockReadCacheMockRecorder struct {
	mock *MockReadCache
}

// NewMockReadCache creates a new mock instance.
func NewMockReadCache(ctrl *gomock.Controller) *MockReadCache {
	mock := &MockReadCache{ctrl: ctrl}
	mock.recorder = &MockReadCacheMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReadCache) EXPECT() *MockReadCacheMockRecorder {
	return m.recorder
}

// GetKV mocks base method.
func (m *MockReadCache) GetKV(ctx context.Context, kv cache.KV) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetKV", ctx, kv)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetKV indicates an expected call of GetKV.
func (mr *MockReadCacheMockRecorder) GetKV(ctx, kv interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetKV", reflect.TypeOf((*MockReadCache)(nil).GetKV), ctx, kv)
}

// GetValFromSet mocks base method.
func (m *MockReadCache) GetValFromSet(ctx context.Context, kv cache.KV) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetValFromSet", ctx, kv)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetValFromSet indicates an expected call of GetValFromSet.
func (mr *MockReadCacheMockRecorder) GetValFromSet(ctx, kv interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetValFromSet", reflect.TypeOf((*MockReadCache)(nil).GetValFromSet), ctx, kv)
}

// MGetKV mocks base method.
func (m *MockReadCache) MGetKV(ctx context.Context, kv ...cache.KV) []bool {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range kv {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "MGetKV", varargs...)
	ret0, _ := ret[0].([]bool)
	return ret0
}

// MGetKV indicates an expected call of MGetKV.
func (mr *MockReadCacheMockRecorder) MGetKV(ctx interface{}, kv ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, kv...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MGetKV", reflect.TypeOf((*MockReadCache)(nil).MGetKV), varargs...)
}

// MockTrashFinder is a mock of TrashFinder interface.
type MockTrashFinder struct {
	ctrl     *gomock.Controller
	recorder *MockTrashFinderMockRecorder
}

// MockTrashFinderMockRecorder is the mock recorder for MockTrashFinder.
type MockTrashFinderMockRecorder struct {
	mock *MockTrashFinder
}

// NewMockTrashFinder creates a new mock instance.
func NewMockTrashFinder(ctrl *gomock.Controller) *MockTrashFinder {
	mock := &MockTrashFinder{ctrl: ctrl}
	mock.recorder = &MockTrashFinderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTrashFinder) EXPECT() *MockTrashFinderMockRecorder {
	return m.recorder
}

// FindTrashCommentID mocks base method.
func (m *MockTrashFinder) FindTrashCommentID(ctx context.Context, svc string) []uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindTrashCommentID", ctx, svc)
	ret0, _ := ret[0].([]uint64)
	return ret0
}

// FindTrashCommentID indicates an expected call of FindTrashCommentID.
func (mr *MockTrashFinderMockRecorder) FindTrashCommentID(ctx, svc interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindTrashCommentID", reflect.TypeOf((*MockTrashFinder)(nil).FindTrashCommentID), ctx, svc)
}

// FindTrashCommentIDByPostID mocks base method.
func (m *MockTrashFinder) FindTrashCommentIDByPostID(ctx context.Context, svc string, postID uint64) []uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindTrashCommentIDByPostID", ctx, svc, postID)
	ret0, _ := ret[0].([]uint64)
	return ret0
}

// FindTrashCommentIDByPostID indicates an expected call of FindTrashCommentIDByPostID.
func (mr *MockTrashFinderMockRecorder) FindTrashCommentIDByPostID(ctx, svc, postID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindTrashCommentIDByPostID", reflect.TypeOf((*MockTrashFinder)(nil).FindTrashCommentIDByPostID), ctx, svc, postID)
}

// FindTrashCommentLikeByCommentID mocks base method.
func (m *MockTrashFinder) FindTrashCommentLikeByCommentID(ctx context.Context, svc string, commentID uint64) []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindTrashCommentLikeByCommentID", ctx, svc, commentID)
	ret0, _ := ret[0].([]string)
	return ret0
}

// FindTrashCommentLikeByCommentID indicates an expected call of FindTrashCommentLikeByCommentID.
func (mr *MockTrashFinderMockRecorder) FindTrashCommentLikeByCommentID(ctx, svc, commentID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindTrashCommentLikeByCommentID", reflect.TypeOf((*MockTrashFinder)(nil).FindTrashCommentLikeByCommentID), ctx, svc, commentID)
}

// FindTrashCommentLikeByPostID mocks base method.
func (m *MockTrashFinder) FindTrashCommentLikeByPostID(ctx context.Context, svc string, postID uint64) map[uint64][]string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindTrashCommentLikeByPostID", ctx, svc, postID)
	ret0, _ := ret[0].(map[uint64][]string)
	return ret0
}

// FindTrashCommentLikeByPostID indicates an expected call of FindTrashCommentLikeByPostID.
func (mr *MockTrashFinderMockRecorder) FindTrashCommentLikeByPostID(ctx, svc, postID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindTrashCommentLikeByPostID", reflect.TypeOf((*MockTrashFinder)(nil).FindTrashCommentLikeByPostID), ctx, svc, postID)
}

// FindTrashPostID mocks base method.
func (m *MockTrashFinder) FindTrashPostID(ctx context.Context, svc string) []uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindTrashPostID", ctx, svc)
	ret0, _ := ret[0].([]uint64)
	return ret0
}

// FindTrashPostID indicates an expected call of FindTrashPostID.
func (mr *MockTrashFinderMockRecorder) FindTrashPostID(ctx, svc interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindTrashPostID", reflect.TypeOf((*MockTrashFinder)(nil).FindTrashPostID), ctx, svc)
}

// FindTrashPostLikeByPostID mocks base method.
func (m *MockTrashFinder) FindTrashPostLikeByPostID(ctx context.Context, svc string, postID uint64) []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindTrashPostLikeByPostID", ctx, svc, postID)
	ret0, _ := ret[0].([]string)
	return ret0
}

// FindTrashPostLikeByPostID indicates an expected call of FindTrashPostLikeByPostID.
func (mr *MockTrashFinderMockRecorder) FindTrashPostLikeByPostID(ctx, svc, postID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindTrashPostLikeByPostID", reflect.TypeOf((*MockTrashFinder)(nil).FindTrashPostLikeByPostID), ctx, svc, postID)
}

// MockSvcHandler is a mock of SvcHandler interface.
type MockSvcHandler struct {
	ctrl     *gomock.Controller
	recorder *MockSvcHandlerMockRecorder
}

// MockSvcHandlerMockRecorder is the mock recorder for MockSvcHandler.
type MockSvcHandlerMockRecorder struct {
	mock *MockSvcHandler
}

// NewMockSvcHandler creates a new mock instance.
func NewMockSvcHandler(ctrl *gomock.Controller) *MockSvcHandler {
	mock := &MockSvcHandler{ctrl: ctrl}
	mock.recorder = &MockSvcHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSvcHandler) EXPECT() *MockSvcHandlerMockRecorder {
	return m.recorder
}

// GetAllServices mocks base method.
func (m *MockSvcHandler) GetAllServices() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllServices")
	ret0, _ := ret[0].([]string)
	return ret0
}

// GetAllServices indicates an expected call of GetAllServices.
func (mr *MockSvcHandlerMockRecorder) GetAllServices() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllServices", reflect.TypeOf((*MockSvcHandler)(nil).GetAllServices))
}

// GetSecretByName mocks base method.
func (m *MockSvcHandler) GetSecretByName(name string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSecretByName", name)
	ret0, _ := ret[0].(string)
	return ret0
}

// GetSecretByName indicates an expected call of GetSecretByName.
func (mr *MockSvcHandlerMockRecorder) GetSecretByName(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSecretByName", reflect.TypeOf((*MockSvcHandler)(nil).GetSecretByName), name)
}

// MockTrashCleaner is a mock of TrashCleaner interface.
type MockTrashCleaner struct {
	ctrl     *gomock.Controller
	recorder *MockTrashCleanerMockRecorder
}

// MockTrashCleanerMockRecorder is the mock recorder for MockTrashCleaner.
type MockTrashCleanerMockRecorder struct {
	mock *MockTrashCleaner
}

// NewMockTrashCleaner creates a new mock instance.
func NewMockTrashCleaner(ctrl *gomock.Controller) *MockTrashCleaner {
	mock := &MockTrashCleaner{ctrl: ctrl}
	mock.recorder = &MockTrashCleanerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTrashCleaner) EXPECT() *MockTrashCleanerMockRecorder {
	return m.recorder
}

// DeleteComment mocks base method.
func (m *MockTrashCleaner) DeleteComment(ctx context.Context, svc string, commentID ...uint64) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, svc}
	for _, a := range commentID {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteComment", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteComment indicates an expected call of DeleteComment.
func (mr *MockTrashCleanerMockRecorder) DeleteComment(ctx, svc interface{}, commentID ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, svc}, commentID...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteComment", reflect.TypeOf((*MockTrashCleaner)(nil).DeleteComment), varargs...)
}

// DeleteCommentLike mocks base method.
func (m *MockTrashCleaner) DeleteCommentLike(ctx context.Context, svc string, commentID uint64, userIDs ...string) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, svc, commentID}
	for _, a := range userIDs {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteCommentLike", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCommentLike indicates an expected call of DeleteCommentLike.
func (mr *MockTrashCleanerMockRecorder) DeleteCommentLike(ctx, svc, commentID interface{}, userIDs ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, svc, commentID}, userIDs...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCommentLike", reflect.TypeOf((*MockTrashCleaner)(nil).DeleteCommentLike), varargs...)
}

// DeletePostLike mocks base method.
func (m *MockTrashCleaner) DeletePostLike(ctx context.Context, svc string, postID uint64, userIDs ...string) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, svc, postID}
	for _, a := range userIDs {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeletePostLike", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePostLike indicates an expected call of DeletePostLike.
func (mr *MockTrashCleanerMockRecorder) DeletePostLike(ctx, svc, postID interface{}, userIDs ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, svc, postID}, userIDs...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePostLike", reflect.TypeOf((*MockTrashCleaner)(nil).DeletePostLike), varargs...)
}
