package test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const JSON = `{"status": "ok"}`

func TestMockServer(t *testing.T) { //exported function with keyword Test and ptr to T is signature of tests in go
	t.Log("Given server responds")
	{
		server := httptest.NewServer(http.HandlerFunc(
			func(writer http.ResponseWriter, request *http.Request) {
				if request.Method != "POST" {
					t.Errorf("Then unexpected method - %v", request.Method)
				} else {
					t.Logf("Then method is POST")
				}
				var b bytes.Buffer
				io.Copy(&b, request.Body)
				if b.String() != "{}" {
					t.Errorf("Then unexpected body - %v", b.String())
				} else {
					t.Logf("Then body is correct")
				}
				writer.WriteHeader(http.StatusOK)
				writer.Header().Add("Content-Type", "application/json")
				writer.Write([]byte(JSON))
			},
		))
		defer server.Close()

		t.Log("When call is made")
		{
			rs, err := http.Post(server.URL, "application/json", strings.NewReader("{}"))
			if err != nil {
				//reports that test failed AND also interrupts the further execution
				t.Fatalf("Then unable to make a call, err - %v", err)
			}
			defer rs.Body.Close()
			{
				if rs.StatusCode != 200 {
					//reports that test failed at this point but doesn't interrupt test execution
					t.Errorf("Then code not expected - %v", rs.StatusCode)
				} else {
					t.Logf("Then status code is 200")
				}
			}
		}
	}
}
