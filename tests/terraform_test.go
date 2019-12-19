package testing

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestNow(t *testing.T) {
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
}
