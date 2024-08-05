package test

import (
	"strings"
	"testing"
)

func TestTable(t *testing.T) {
	testData := []struct {
		inputString    string
		expectedString string
	}{
		{"ab  ", "ab"},
		{"  ab", "ab"},
	}
	for _, testCase := range testData {
		t.Logf("Given input string '%s'", testCase.inputString)
		{
			t.Logf("When trim space function called")
			{
				result := strings.TrimSpace(testCase.inputString)
				if result == testCase.expectedString {
					t.Logf("Then result matches expected string")
				} else {
					t.Errorf("Then result doesn't match, expected - '%s', actual - '%s'", testCase.expectedString, result)
				}
			}
		}
	}
}
