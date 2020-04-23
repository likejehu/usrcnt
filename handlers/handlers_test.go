package handlers

import (
	_ "errors"
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
		mockStore.On("Set", testToken, testToken).Return(error400)
		handler := &Handler{mockStore, mockSession}
		handler.Hello(rec, req)
		print("yey")
		mockSession.AssertExpectations(t)
		mockStore.AssertExpectations(t)
		//assertions
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

	})

	t.Run("succes case", func(t *testing.T) {

	})

	t.Run("succes case", func(t *testing.T) {

	})

	t.Run("succes case", func(t *testing.T) {

	})

}
