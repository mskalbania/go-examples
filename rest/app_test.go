package rest

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http/httptest"
	"testing"
)

func TestHealthExposed(t *testing.T) {
	//given
	gin.SetMode(gin.TestMode)
	healthMock, authMock, userMock := setupMocks()
	router := setupRouter(authMock, healthMock, userMock)

	//when
	rq := httptest.NewRequest("GET", "/health", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, rq)

	//then
	healthMock.AssertCalled(t, "Health", mock.Anything)
	authMock.assertNotCalled(t)
}

func TestUserAPIExposed(t *testing.T) {
	//given
	gin.SetMode(gin.TestMode)
	healthMock, authMock, userMock := setupMocks()
	router := setupRouter(authMock, healthMock, userMock)

	tests := []struct {
		method                string
		path                  string
		expectedHandlerCalled func()
	}{
		{"GET", "/api/v1/users", func() {
			userMock.AssertCalled(t, "GetUsers", mock.Anything)
		}},
		{"POST", "/api/v1/users", func() {
			userMock.AssertCalled(t, "CreateUser", mock.Anything)
		}},
		{"DELETE", "/api/v1/users/abc", func() {
			userMock.AssertCalled(t, "DeleteUser", mock.Anything)
		}},
		{"PUT", "/api/v1/users/abc", func() {
			userMock.AssertCalled(t, "UpdateUser", mock.Anything)
		}},
		{"GET", "/api/v1/users/abc", func() {
			userMock.AssertCalled(t, "GetUserById", mock.Anything)
		}},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("'%s%s'", test.method, test.path), func(t *testing.T) {
			rq := httptest.NewRequest(test.method, test.path, nil)
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, rq)

			test.expectedHandlerCalled()
			authMock.assertCalled(t)
		})
	}
}

func setupMocks() (*HealthMock, *AuthenticationMock, *UserMock) {
	healthMock := new(HealthMock)
	authMock := new(AuthenticationMock)
	userMock := new(UserMock)

	healthMock.On("Health", mock.Anything).Return()
	authMock.On("RequireAPIToken").Return()
	userMock.On("GetUsers", mock.Anything).Return()
	userMock.On("GetUserById", mock.Anything).Return()
	userMock.On("CreateUser", mock.Anything).Return()
	userMock.On("DeleteUser", mock.Anything).Return()
	userMock.On("UpdateUser", mock.Anything).Return()

	return healthMock, authMock, userMock

}

type HealthMock struct {
	mock.Mock
}

func (h *HealthMock) Health(ctx *gin.Context) {
	_ = h.Called(ctx)
}

type AuthenticationMock struct {
	mock.Mock
	called bool
}

func (a *AuthenticationMock) assertCalled(t *testing.T) {
	require.True(t, a.called, "authentication not called")
	a.called = false
}

func (a *AuthenticationMock) assertNotCalled(t *testing.T) {
	require.False(t, a.called, "authentication called")
}

func (a *AuthenticationMock) RequireAPIToken() gin.HandlerFunc {
	_ = a.Called()
	return func(context *gin.Context) {
		a.called = true
	}
}

type UserMock struct {
	mock.Mock
}

func (u *UserMock) GetUsers(context *gin.Context) {
	_ = u.Called(context)
}

func (u *UserMock) GetUserById(context *gin.Context) {
	_ = u.Called(context)
}

func (u *UserMock) CreateUser(context *gin.Context) {
	_ = u.Called(context)
}

func (u *UserMock) DeleteUser(context *gin.Context) {
	_ = u.Called(context)
}

func (u *UserMock) UpdateUser(context *gin.Context) {
	_ = u.Called(context)
}
