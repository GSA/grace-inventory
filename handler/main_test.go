package main

import (
	"os"

	"testing"

	"gotest.tools/v3/assert"
)

func TestGetSheets(t *testing.T) {
	tt := map[string]struct {
		env      string
		expected []string
	}{
		"none": {"", defaultSheets},
		"one":  {"Buckets", []string{"Accounts", "Buckets"}},
		"two":  {"Buckets,Groups", []string{"Accounts", "Buckets", "Groups"}},
		"mix":  {"Buckets,Accounts,Groups", []string{"Accounts", "Buckets", "Groups"}},
	}

	// restore the env value of sheets, if set
	sheets := os.Getenv("sheets")
	defer func() {
		if len(sheets) > 0 {
			os.Setenv("sheets", sheets)
		}
	}()

	for name, tc := range tt {
		tc := tc
		t.Run(name, func(t *testing.T) {
			os.Setenv("sheets", tc.env)
			actual := getSheets()
			assert.DeepEqual(t, tc.expected, actual)
		})
	}
}
