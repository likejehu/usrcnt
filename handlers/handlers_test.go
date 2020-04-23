package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/likejehu/usrcnt/handlers/mocks"
	"github.com/likejehu/usrcnt/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	testToken = "testToken"
	notSetErr = session.ErrorNotSet
	error400  = session.Error400
	error500  = errors.New("shit happens")
)

const notSet = "cookie is not set"

func TestHello(t *testing.T) {

	t.Run("new user success", func(t *testing.T) {
		//setup
		mockStore := &mocks.Store{}
		mockSession := &mocks.SessionManager{}
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		mockSession.On("ReadCookie", rec, req).Return(notSet, notSetErr)
		mockSession.On("NewST").Return(testToken)
		mockStore.On("Set", testToken, testToken).Return(nil)
		mockSession.On("SetCookie", mock.Anything, testToken).Return()
		mockStore.On("SETNXToZero", usrCountKey).Return(nil)
		mockStore.On("Increment", usrCountKey).Return(13, nil)
		mockStore.On("Get", usrCountKey).Return(13, nil)
		handler := &Handler{mockStore, mockSession}
		handler.Hello(rec, req)
		mockSession.AssertExpectations(t)
		mockStore.AssertExpectations(t)
		//assertions
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("new user fail on set token", func(t *testing.T) {
		//setup
		mockStore := &mocks.Store{}
		mockSession := &mocks.SessionManager{}
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		mockSession.On("ReadCookie", rec, req).Return(notSet, notSetErr)
		mockSession.On("NewST").Return(testToken)
		mockStore.On("Set", testToken, testToken).Return(error500)
		handler := &Handler{mockStore, mockSession}
		handler.Hello(rec, req)
		mockSession.AssertExpectations(t)
		mockStore.AssertExpectations(t)
		//assertions
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

	})

	t.Run("new user fail on SETNXToZero", func(t *testing.T) {
		//setup
		mockStore := &mocks.Store{}
		mockSession := &mocks.SessionManager{}
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		mockSession.On("ReadCookie", rec, req).Return(notSet, notSetErr)
		mockSession.On("NewST").Return(testToken)
		mockStore.On("Set", testToken, testToken).Return(nil)
		mockSession.On("SetCookie", mock.Anything, testToken).Return()
		mockStore.On("SETNXToZero", usrCountKey).Return(error500)
		handler := &Handler{mockStore, mockSession}
		handler.Hello(rec, req)
		mockSession.AssertExpectations(t)
		mockStore.AssertExpectations(t)
		//assertions
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("new user fail on Increment", func(t *testing.T) {
		//setup
		mockStore := &mocks.Store{}
		mockSession := &mocks.SessionManager{}
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		mockSession.On("ReadCookie", rec, req).Return(notSet, notSetErr)
		mockSession.On("NewST").Return(testToken)
		mockStore.On("Set", testToken, testToken).Return(nil)
		mockSession.On("SetCookie", mock.Anything, testToken).Return()
		mockStore.On("SETNXToZero", usrCountKey).Return(nil)
		mockStore.On("Increment", usrCountKey).Return(0, error500)
		handler := &Handler{mockStore, mockSession}
		handler.Hello(rec, req)
		mockSession.AssertExpectations(t)
		mockStore.AssertExpectations(t)
		//assertions
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

	})

	t.Run("new user fail on reading the cookies", func(t *testing.T) {
		//setup
		mockStore := &mocks.Store{}
		mockSession := &mocks.SessionManager{}
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		mockSession.On("ReadCookie", rec, req).Return("bad req", error400)
		handler := &Handler{mockStore, mockSession}
		handler.Hello(rec, req)
		mockSession.AssertExpectations(t)
		//assertions
		assert.Equal(t, http.StatusBadRequest, rec.Code)

	})

	t.Run("old user success", func(t *testing.T) {
		//setup
		mockStore := &mocks.Store{}
		mockSession := &mocks.SessionManager{}
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		mockSession.On("ReadCookie", rec, req).Return(testToken, nil)
		mockStore.On("Exists", testToken).Return(1, nil)
		mockStore.On("Get", usrCountKey).Return(13, nil)
		handler := &Handler{mockStore, mockSession}
		handler.Hello(rec, req)
		mockSession.AssertExpectations(t)
		mockStore.AssertExpectations(t)
		//assertions
		assert.Equal(t, http.StatusOK, rec.Code)

	})
	t.Run("old user  fail on token check", func(t *testing.T) {
		//setup
		mockStore := &mocks.Store{}
		mockSession := &mocks.SessionManager{}
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		mockSession.On("ReadCookie", rec, req).Return(testToken, nil)
		mockStore.On("Exists", testToken).Return(0, error400)
		handler := &Handler{mockStore, mockSession}
		handler.Hello(rec, req)
		mockSession.AssertExpectations(t)
		mockStore.AssertExpectations(t)
		//assertions
		assert.Equal(t, http.StatusBadRequest, rec.Code)

	})
	t.Run("old user  fail on gettin usrCountVal", func(t *testing.T) {
		//setup
		mockStore := &mocks.Store{}
		mockSession := &mocks.SessionManager{}
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		mockSession.On("ReadCookie", rec, req).Return(testToken, nil)
		mockStore.On("Exists", testToken).Return(1, nil)
		mockStore.On("Get", usrCountKey).Return(0, error500)
		handler := &Handler{mockStore, mockSession}
		handler.Hello(rec, req)
		mockSession.AssertExpectations(t)
		mockStore.AssertExpectations(t)
		//assertions
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

}
