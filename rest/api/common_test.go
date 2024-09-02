package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go-examples/rest/model"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAbortWithContextError(t *testing.T) {
	//given
	recorder := httptest.NewRecorder()
	testCtx, _ := gin.CreateTestContext(recorder)
	status := http.StatusTeapot
	message := "user message"
	err := fmt.Errorf("error details")

	//when
	AbortWithContextError(testCtx, status, message, err)

	//then
	require.Equal(t, status, testCtx.Writer.Status())
	require.Equal(t, true, testCtx.IsAborted())
	require.NotNil(t, testCtx.Errors.Last())
	require.Equal(t, "user message: error details", testCtx.Errors.Last().Err.Error())

	errorBody := new(model.Error)
	err = json.Unmarshal(recorder.Body.Bytes(), errorBody)
	require.NoError(t, err)
	require.Equal(t, message, errorBody.Message)
	require.NotEmpty(t, errorBody.Timestamp)
}

func TestAbort(t *testing.T) {
	//given
	recorder := httptest.NewRecorder()
	testCtx, _ := gin.CreateTestContext(recorder)
	status := http.StatusTeapot
	message := "user message"

	//when
	Abort(testCtx, status, message)

	//then
	require.Equal(t, status, testCtx.Writer.Status())
	require.Equal(t, true, testCtx.IsAborted())
	require.Nil(t, testCtx.Errors.Last())

	errorBody := new(model.Error)
	err := json.Unmarshal(recorder.Body.Bytes(), errorBody)
	require.NoError(t, err)
	require.Equal(t, message, errorBody.Message)
	require.NotEmpty(t, errorBody.Timestamp)
}
