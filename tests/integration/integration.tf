resource "aws_cloudwatch_log_group" "integration_test" {
  name = "integration_test"
}

module "integration_test" {
  source            = "../../"
  accounts_info     = "self"
  project_name      = "grace"
  appenv            = "integration-test"
  master_account_id = "123456789012"
  master_role_name  = "role"
  tenant_role_name  = "tenant-role"
  source_file       = var.source_file
}
