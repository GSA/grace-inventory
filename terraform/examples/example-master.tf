// If the Lambda function is installed in a master payer account, it will
// list all accounts and inventory each one using the OrganizationAccessRole
// if accounts_info = ""
module "example_master" {
  source        = "github.com/GSA/grace-inventory/terraform"
  accounts_info = ""
  source_file   = "../../release/grace-inventory-lambda.zip"
  appenv        = "development"
  project_name  = "grace"
}

