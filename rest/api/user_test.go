package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go-examples/rest/model"
	"go-examples/rest/repository"
	"go-examples/rest/test"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var testUserId = "some-id"
var testUserEmail = "email@example.com"

type UserSuite struct {
	suite.Suite
	repositoryMock *test.UserRepositoryMock
	userAPI        UserAPI
	ctx            *gin.Context
	recorder       *httptest.ResponseRecorder
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserSuite))
}

func (suite *UserSuite) BeforeTest(suiteName, testName string) {
	gin.SetMode(gin.TestMode)
	suite.repositoryMock = new(test.UserRepositoryMock)
	suite.userAPI = NewUserAPI(suite.repositoryMock)
	suite.recorder = httptest.NewRecorder()
	suite.ctx, _ = gin.CreateTestContext(suite.recorder)
}

func (suite *UserSuite) TestGetUsersSuccess() {
	//given
	suite.repositoryMock.On("GetAllUsers").Return([]*model.User{{testUserId, testUserEmail}}, nil)

	//when
	suite.userAPI.GetUsers(suite.ctx)

	//then
	require.Equal(suite.T(), http.StatusOK, suite.recorder.Code)
	expectedJson := fmt.Sprintf(`[{"id": "%s", "email": "%s"}]`, testUserId, testUserEmail)
	require.JSONEq(suite.T(), expectedJson, suite.recorder.Body.String())
}

func (suite *UserSuite) TestGetUsersError() {
	//given
	suite.repositoryMock.On("GetAllUsers").Return([]*model.User{}, fmt.Errorf("db error"))

	//when
	suite.userAPI.GetUsers(suite.ctx)

	//then
	require.Equal(suite.T(), http.StatusInternalServerError, suite.recorder.Code)
	require.Contains(suite.T(), suite.recorder.Body.String(), "error getting users")
}

func (suite *UserSuite) TestGetUserByIdSuccess() {
	//given
	suite.repositoryMock.On("GetUserById", testUserId).Return(&model.User{ID: testUserId, Email: testUserEmail}, nil)
	suite.ctx.Params = append(suite.ctx.Params, gin.Param{Key: "id", Value: testUserId})

	//when
	suite.userAPI.GetUserById(suite.ctx)

	//then
	require.Equal(suite.T(), http.StatusOK, suite.recorder.Code)
	expectedJson := fmt.Sprintf(`{"id": "%s", "email": "%s"}`, testUserId, testUserEmail)
	require.JSONEq(suite.T(), expectedJson, suite.recorder.Body.String())
}

func (suite *UserSuite) TestGetUserByIdNoId() {
	//given
	suite.ctx.Params = []gin.Param{}

	//when
	suite.userAPI.GetUserById(suite.ctx)

	//then
	require.Equal(suite.T(), http.StatusBadRequest, suite.recorder.Code)
	require.Contains(suite.T(), suite.recorder.Body.String(), "id is required")
}

func (suite *UserSuite) TestGetUserByIdUserNotFound() {
	//give
	suite.repositoryMock.On("GetUserById", testUserId).Return(new(model.User), repository.ErrUserNotFound)
	suite.ctx.Params = append(suite.ctx.Params, gin.Param{Key: "id", Value: testUserId})

	//when
	suite.userAPI.GetUserById(suite.ctx)

	//then
	require.Equal(suite.T(), http.StatusNotFound, suite.recorder.Code)
	require.Contains(suite.T(), suite.recorder.Body.String(), "user not found")
}

func (suite *UserSuite) TestGetUserByIdRepositoryError() {
	//given
	suite.repositoryMock.On("GetUserById", testUserId).Return(new(model.User), fmt.Errorf("db error"))
	suite.ctx.Params = append(suite.ctx.Params, gin.Param{Key: "id", Value: testUserId})

	//when
	suite.userAPI.GetUserById(suite.ctx)

	//then
	require.Equal(suite.T(), http.StatusInternalServerError, suite.recorder.Code)
	require.Contains(suite.T(), suite.recorder.Body.String(), "error getting user")
}

func (suite *UserSuite) TestCreateUserSuccess() {
	//given
	suite.repositoryMock.On("Save", &model.PostUser{Email: testUserEmail}).Return(&model.User{ID: testUserId, Email: testUserEmail}, nil)
	suite.ctx.Request, _ = http.NewRequest(http.MethodPost, "/users", strings.NewReader(fmt.Sprintf(`{"email": "%s"}`, testUserEmail)))

	//when
	suite.userAPI.CreateUser(suite.ctx)

	//then
	require.Equal(suite.T(), http.StatusCreated, suite.recorder.Code)
	expectedJson := fmt.Sprintf(`{"id": "%s", "email": "%s"}`, testUserId, testUserEmail)
	require.JSONEq(suite.T(), expectedJson, suite.recorder.Body.String())
}

func (suite *UserSuite) TestCreateUserInvalidRequest() {
	testData := []struct {
		request string
	}{
		{``},
		{`{}`},
		{`{"email": ""}`},
	}
	for _, testCase := range testData {
		//given
		suite.ctx.Request, _ = http.NewRequest(http.MethodPost, "/users", strings.NewReader(testCase.request))

		//when
		suite.userAPI.CreateUser(suite.ctx)

		//then
		require.Equal(suite.T(), http.StatusBadRequest, suite.recorder.Code)
		require.Contains(suite.T(), suite.recorder.Body.String(), "invalid request")
	}
}

func (suite *UserSuite) TestCreateUserRepositoryError() {
	//given
	suite.repositoryMock.On("Save", &model.PostUser{Email: testUserEmail}).Return(new(model.User), fmt.Errorf("db error"))
	suite.ctx.Request, _ = http.NewRequest(http.MethodPost, "/users", strings.NewReader(fmt.Sprintf(`{"email": "%s"}`, testUserEmail)))

	//when
	suite.userAPI.CreateUser(suite.ctx)

	//then
	require.Equal(suite.T(), http.StatusInternalServerError, suite.recorder.Code)
	require.Contains(suite.T(), suite.recorder.Body.String(), "error saving user")
}

func (suite *UserSuite) TestDeleteUserSuccess() {
	//given
	suite.repositoryMock.On("Delete", testUserId).Return(nil)
	suite.ctx.Params = append(suite.ctx.Params, gin.Param{Key: "id", Value: testUserId})

	//when
	suite.userAPI.DeleteUser(suite.ctx)

	//then
	require.Equal(suite.T(), http.StatusNoContent, suite.recorder.Code)
}

func (suite *UserSuite) TestDeleteUserNoId() {
	//given
	suite.ctx.Params = []gin.Param{}

	//when
	suite.userAPI.DeleteUser(suite.ctx)

	//then
	require.Equal(suite.T(), http.StatusBadRequest, suite.recorder.Code)
	require.Contains(suite.T(), suite.recorder.Body.String(), "id is required")
}

func (suite *UserSuite) TestDeleteUserRepositoryError() {
	//given
	suite.repositoryMock.On("Delete", testUserId).Return(fmt.Errorf("db error"))
	suite.ctx.Params = append(suite.ctx.Params, gin.Param{Key: "id", Value: testUserId})

	//when
	suite.userAPI.DeleteUser(suite.ctx)

	//then
	require.Equal(suite.T(), http.StatusInternalServerError, suite.recorder.Code)
	require.Contains(suite.T(), suite.recorder.Body.String(), "error deleting user")
}

func (suite *UserSuite) TestUpdateUserSuccess() {
	//given
	suite.ctx.Params = append(suite.ctx.Params, gin.Param{Key: "id", Value: testUserId})
	suite.ctx.Request, _ = http.NewRequest(http.MethodPut, fmt.Sprintf("/users/%s", testUserId), strings.NewReader(fmt.Sprintf(`{"email": "%s"}`, testUserEmail)))
	suite.repositoryMock.On("Exists", testUserId).Return(true, nil)
	suite.repositoryMock.On("Update", testUserId, &model.PostUser{Email: testUserEmail}).Return(&model.User{ID: testUserId, Email: testUserEmail}, nil)

	//when
	suite.userAPI.UpdateUser(suite.ctx)

	//then
	require.Equal(suite.T(), http.StatusOK, suite.recorder.Code)
	expectedJson := fmt.Sprintf(`{"id": "%s", "email": "%s"}`, testUserId, testUserEmail)
	require.JSONEq(suite.T(), expectedJson, suite.recorder.Body.String())
}

func (suite *UserSuite) TestUpdateUserExistsError() {
	//given
	suite.ctx.Params = append(suite.ctx.Params, gin.Param{Key: "id", Value: testUserId})
	suite.repositoryMock.On("Exists", testUserId).Return(false, fmt.Errorf("db error"))

	//when
	suite.userAPI.UpdateUser(suite.ctx)

	//then
	require.Equal(suite.T(), http.StatusInternalServerError, suite.recorder.Code)
	require.Contains(suite.T(), suite.recorder.Body.String(), "error updating user")
}

func (suite *UserSuite) TestUpdateUserDoesNotExists() {
	//given
	suite.ctx.Params = append(suite.ctx.Params, gin.Param{Key: "id", Value: testUserId})
	suite.repositoryMock.On("Exists", testUserId).Return(false, nil)

	//when
	suite.userAPI.UpdateUser(suite.ctx)

	//then
	require.Equal(suite.T(), http.StatusNotFound, suite.recorder.Code)
	require.Contains(suite.T(), suite.recorder.Body.String(), "user not found")
}

func (suite *UserSuite) TestUpdateUserInvalidRequest() {
	testData := []struct {
		request string
	}{
		{``},
		{`{}`},
		{`{"email": ""}`},
	}
	for _, testCase := range testData {
		//given
		suite.ctx.Request, _ = http.NewRequest(http.MethodPut, fmt.Sprintf("/users/%s", testUserId), strings.NewReader(testCase.request))
		suite.ctx.Params = append(suite.ctx.Params, gin.Param{Key: "id", Value: testUserId})
		suite.repositoryMock.On("Exists", testUserId).Return(true, nil)

		//when
		suite.userAPI.UpdateUser(suite.ctx)

		//then
		require.Equal(suite.T(), http.StatusBadRequest, suite.recorder.Code)
		require.Contains(suite.T(), suite.recorder.Body.String(), "invalid request")
	}
}

func (suite *UserSuite) TestUpdateUserUpdateRepositoryError() {
	//given
	suite.ctx.Params = append(suite.ctx.Params, gin.Param{Key: "id", Value: testUserId})
	suite.ctx.Request, _ = http.NewRequest(http.MethodPut, fmt.Sprintf("/users/%s", testUserId), strings.NewReader(fmt.Sprintf(`{"email": "%s"}`, testUserEmail)))
	suite.repositoryMock.On("Exists", testUserId).Return(true, nil)
	suite.repositoryMock.On("Update", testUserId, &model.PostUser{Email: testUserEmail}).Return(new(model.User), fmt.Errorf("db error"))

	//when
	suite.userAPI.UpdateUser(suite.ctx)

	//then
	require.Equal(suite.T(), http.StatusInternalServerError, suite.recorder.Code)
	require.Contains(suite.T(), suite.recorder.Body.String(), "error updating user")
}
