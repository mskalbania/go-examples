package test

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestTable(t *testing.T) {
	tests := []struct {
		inputString    string
		expectedString string
	}{
		{"ab  ", "ab"},
		{"  ab", "ab"},
	}
	for _, test := range tests {
		//this runs subtests
		t.Run(fmt.Sprintf("input-'%s'", test.inputString), func(t *testing.T) {
			require.Equal(t, test.expectedString, strings.TrimSpace(test.inputString))
		})
	}
}
