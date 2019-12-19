package testing

import (
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestNow(t *testing.T) {
	localstack := `/opt/code/localstack/bin/localstack`
	cmd := exec.Command(localstack, "start", "--host")
	go func(t *testing.T) {
		err := cmd.Run()
		if err != nil {
			t.Fatalf("failed to execute localstack: %v", err)
		}
	}(t)
	opts := &terraform.Options{

		// The path to where our Terraform code is located
		TerraformDir: "scenarios/one",

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"appenv":            "development",
			"tenant_role_name":  "0",
			"master_role_name":  "1",
			"master_account_id": "2",
		},

		// Disable colors in Terraform commands so its easier to parse stdout/stderr
		NoColor: true,
	}
	defer terraform.Destroy(t, opts)
	t.Logf("output: %s\n", terraform.InitAndApply(t, opts))

	pattern := `/tmp/localstack/data/*.json`
	matches, err := filepath.Glob(pattern)
	if err != nil {
		t.Fatalf("failed to glob files %s: %v", pattern, err)
	}
	for _, m := range matches {
		t.Logf("found file: %s\n", m)
	}
}
