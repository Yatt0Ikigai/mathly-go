// Code generated by MockGen. DO NOT EDIT.
// Source: game_library.go
//
// Generated by this command:
//
//	mockgen -source=game_library.go -package games -destination=game_library_mock.go
//

// Package games is a generated GoMock package.
package games

import (
	common_games "mathly/internal/sockets/games/common"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockGameLibrary is a mock of GameLibrary interface.
type MockGameLibrary struct {
	ctrl     *gomock.Controller
	recorder *MockGameLibraryMockRecorder
	isgomock struct{}
}

// MockGameLibraryMockRecorder is the mock recorder for MockGameLibrary.
type MockGameLibraryMockRecorder struct {
	mock *MockGameLibrary
}

// NewMockGameLibrary creates a new mock instance.
func NewMockGameLibrary(ctrl *gomock.Controller) *MockGameLibrary {
	mock := &MockGameLibrary{ctrl: ctrl}
	mock.recorder = &MockGameLibraryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGameLibrary) EXPECT() *MockGameLibraryMockRecorder {
	return m.recorder
}

// StartNewGame mocks base method.
func (m *MockGameLibrary) StartNewGame(game AvailableGames, c common_games.GameConfig) common_games.Game {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StartNewGame", game, c)
	ret0, _ := ret[0].(common_games.Game)
	return ret0
}

// StartNewGame indicates an expected call of StartNewGame.
func (mr *MockGameLibraryMockRecorder) StartNewGame(game, c any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartNewGame", reflect.TypeOf((*MockGameLibrary)(nil).StartNewGame), game, c)
}
