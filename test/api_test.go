package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func init() {
	Routes()
}

func TestApi(t *testing.T) {
	t.Log("Given request to service")
	{
		request, _ := http.NewRequest("GET", "/users", nil)
		t.Log("When service called")
		{
			recorder := httptest.NewRecorder()
			http.DefaultServeMux.ServeHTTP(recorder, request)

			if recorder.Code == 200 {
				t.Log("Then code matches expected - 200")
			} else {
				t.Errorf("Then code doesn't match, actual - %v", recorder.Code)
			}
			var users []struct {
				Id   string `json:"id"`
				Name string `json:"name"`
			}
			json.NewDecoder(recorder.Body).Decode(&users)
			if len(users) != 2 {
				t.Errorf("Then expected 2 users, actual - %v", len(users))
			} else {
				t.Logf("Then expected 2 users")
			}
			if users[0].Id == "id1" && users[0].Name == "john" {
				t.Logf("Then user1 matches expected")
			} else {
				t.Errorf("Then user1 doesn't match, actual - %v", users[0])
			}
			if users[1].Id == "id2" && users[1].Name == "doe" {
				t.Logf("Then user2 matches expected")
			} else {
				t.Errorf("Then user2 doesn't match, actual - %v", users[1])
			}
		}
	}
}
