// The default behavior is to inventory only the account the lambda function
// is installed in (i.e. accounts_info = "self"
module "example_self" {
  source       = "github.com/GSA/grace-inventory/terraform"
  source_file  = "../../release/grace-inventory-lambda.zip"
  appenv       = "development"
  project_name = "grace"
}
