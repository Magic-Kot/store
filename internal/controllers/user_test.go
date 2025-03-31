package controllers

import (
	"bytes"
	"context"
	"net/http/httptest"
	"testing"

	"github.com/Magic-Kot/store/internal/models"
	"github.com/Magic-Kot/store/internal/services/user"
	mock_user "github.com/Magic-Kot/store/internal/services/user/mocks"
	"github.com/Magic-Kot/store/pkg/logging"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/magiconair/properties/assert"
	"go.uber.org/mock/gomock"
)

func TestHandler_signUp(t *testing.T) {
	// Init Test Table
	var ctx = context.Background()
	type mockBehavior func(r *mock_user.MockUserRepository, user models.UserLogin)

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            models.UserLogin
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"login": "username", "password": "qwerty"}`,
			inputUser: models.UserLogin{
				Username: "username",
				Password: "qwerty",
			},
			mockBehavior: func(r *mock_user.MockUserRepository, user models.UserLogin) {
				r.EXPECT().CreateUser(ctx, user.Username, user.Password).Return(1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"successfully created user, id: 1"}`,
		},
		{
			name:                 "Wrong Input",
			inputBody:            `{"login": "username"}`,
			inputUser:            models.UserLogin{},
			mockBehavior:         func(r *mock_user.MockUserRepository, user models.UserLogin) {},
			expectedStatusCode:   400,
			expectedResponseBody: `"invalid request"`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_user.NewMockUserRepository(c)
			test.mockBehavior(repo, test.inputUser)

			auth := mock_user.NewMockAuthRepository(c)

			// create logger
			logCfg := logging.LoggerDeps{
				LogLevel: "debug",
			}
			logger, _ := logging.NewLogger(&logCfg)

			services := user.UserService{UserRepository: repo, AuthRepository: auth}
			handler := ApiController{UserService: services, logger: logger, validator: validator.New()}

			// Init Endpoint
			r := echo.New()
			r.POST("/sign-up", handler.SignUp)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-up",
				bytes.NewBufferString(test.inputBody))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}
