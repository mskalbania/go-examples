package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type AuthenticationSuite struct {
	suite.Suite
	authentication Authentication
	ctx            *gin.Context
	recorder       *httptest.ResponseRecorder
}

func TestAuthenticationSuite(t *testing.T) {
	suite.Run(t, new(AuthenticationSuite))
}

func (s *AuthenticationSuite) BeforeTest(suiteName, testName string) {
	gin.SetMode(gin.TestMode)
	s.recorder = httptest.NewRecorder()
	s.ctx, _ = gin.CreateTestContext(s.recorder)
	s.authentication = NewAuthentication()
}

func (s *AuthenticationSuite) TestAuthenticationSuccessful() {
	//given
	rq := httptest.NewRequest("GET", "/", nil)
	rq.Header.Add(apiKeyHeader, "token")
	s.ctx.Request = rq

	//when
	s.authentication.RequireAPIToken()(s.ctx)

	//
	require.Equal(s.T(), http.StatusOK, s.recorder.Code)
}

func (s *AuthenticationSuite) TestAuthenticationMissingToken() {
	//given
	rq := httptest.NewRequest("GET", "/", nil)
	s.ctx.Request = rq

	//when
	s.authentication.RequireAPIToken()(s.ctx)

	//
	require.Equal(s.T(), http.StatusUnauthorized, s.recorder.Code)
	require.Contains(s.T(), s.recorder.Body.String(), "missing api key")
}

func (s *AuthenticationSuite) TestAuthenticationInvalidToken() {
	//given
	rq := httptest.NewRequest("GET", "/", nil)
	rq.Header.Add(apiKeyHeader, "invalid")
	s.ctx.Request = rq

	//when
	s.authentication.RequireAPIToken()(s.ctx)

	//
	require.Equal(s.T(), http.StatusUnauthorized, s.recorder.Code)
	require.Contains(s.T(), s.recorder.Body.String(), "invalid api key")
}
