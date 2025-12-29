package http

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CryptoGu1/books-rest-clean-arch/internal/domain"
	"github.com/CryptoGu1/books-rest-clean-arch/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_signup(t *testing.T) {
	type mockBehavior func(s *mocks.AuthService, user domain.SingUpInput)

	testTable := []struct {
		name                string
		inputBody           string
		inputUser           domain.SingUpInput
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "ok",
			inputBody: `{"name": "Test", "email": "test@test.kz", "password": "test1234"}`,
			inputUser: domain.SingUpInput{
				Name:     "Test",
				Email:    "test@test.kz",
				Password: "test1234",
			},
			mockBehavior: func(s *mocks.AuthService, user domain.SingUpInput) {
				s.On("SignUp", mock.Anything, user).Return(1, nil)
			},
			expectedStatusCode:  201,
			expectedRequestBody: `{"id":1}`,
		},
		{
			name:      "Empty Field",
			inputBody: `{"email": "test@test.kz", "password": "test1234"}`,

			mockBehavior:        func(s *mocks.AuthService, user domain.SingUpInput) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"error":"invalid input body"}`,
		},
		{
			name:      "service fail",
			inputBody: `{"name": "Test", "email": "test@test.kz", "password": "test1234"}`,
			inputUser: domain.SingUpInput{
				Name:     "Test",
				Email:    "test@test.kz",
				Password: "test1234",
			},
			mockBehavior: func(s *mocks.AuthService, user domain.SingUpInput) {
				s.On("SignUp", mock.Anything, user).Return(1, errors.New("service error"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"handler": "sign-up","problem": "service error",}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			authService := mocks.NewAuthService(t)
			testCase.mockBehavior(authService, testCase.inputUser)

			var bookService BookService
			handler := NewHandler(bookService, authService, []byte("secret"))
			e := echo.New()

			req := httptest.NewRequest(http.MethodPost, "/auth/sign-up", bytes.NewBufferString(testCase.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)

			err := handler.signUp(c)

			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedStatusCode, rec.Code)

		})
	}

}
