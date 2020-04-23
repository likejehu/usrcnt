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

var testToken = "testToken"
var notSet = "cookie is not set"
var notSetErr = session.ErrorNotSet

func TestHello(t *testing.T) {

	t.Run("succes case", func(t *testing.T) {
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

	t.Run("succes case", func(t *testing.T) {

	})

	t.Run("succes case", func(t *testing.T) {

	})

	t.Run("succes case", func(t *testing.T) {

	})

	t.Run("succes case", func(t *testing.T) {

	})

}
