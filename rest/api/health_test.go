package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go-examples/rest/test"
	"net/http"
	"net/http/httptest"
	"testing"
)

type HealthSuite struct {
	suite.Suite
	dbMock    *test.DatabaseMock
	healthAPI *HealthAPI
	ctx       *gin.Context
	recorder  *httptest.ResponseRecorder
}

func TestHealthSuite(t *testing.T) {
	suite.Run(t, new(HealthSuite))
}

func (suite *HealthSuite) BeforeTest(suiteName, testName string) {
	suite.dbMock = &test.DatabaseMock{}
	suite.recorder = httptest.NewRecorder()
	suite.healthAPI = NewHealthAPI(suite.dbMock)
	suite.ctx, _ = gin.CreateTestContext(suite.recorder)
}

func (suite *HealthSuite) TestHealthSuccess() {
	//given db is reachable
	suite.dbMock.On("Ping", mock.Anything).Return(nil)

	//when health is called
	suite.healthAPI.Health(suite.ctx)

	//then status is 200
	require.Equal(suite.T(), http.StatusOK, suite.recorder.Code)
}

func (suite *HealthSuite) TestHealthFailureDbNotReachable() {
	//given db is not reachable
	suite.dbMock.On("Ping", mock.Anything).Return(fmt.Errorf("db not reachable"))

	//when health is called
	suite.healthAPI.Health(suite.ctx)

	//then status is 500
	require.Equal(suite.T(), http.StatusInternalServerError, suite.recorder.Code)
}
