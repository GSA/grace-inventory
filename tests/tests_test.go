package tests

import (
	"testing"

	"github.com/GSA/grace-tftest/tester"
)

func TestNow(t *testing.T) {
	zipFile := "../release/grace-inventory-lambda.zip"
	err := tester.Run("integration", map[string]string{
		"TF_VAR_source_file": zipFile,
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
}
