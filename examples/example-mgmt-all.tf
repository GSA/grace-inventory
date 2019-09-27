// If the Lambda function is installed in a non-master/mgmt account, it can
// list all accounts and inventory each one using the OrganizationAccessRole
// if accounts_info = "" and master_account_id and master_role_name are set
// and the roles are assumable by the Lambda function's IAM role
module "example_mgmt_all" {
  source            = "github.com/GSA/grace-inventory?ref=v0.1.1"
  accounts_info     = ""
  master_account_id = "111111111111"
  master_role_name  = "AssumableRole"
  source_file       = "../../release/grace-inventory-lambda.zip"
  appenv            = "development"
  project_name      = "grace"
}
