terraform {
  backend "s3" {
    region = "us-east-1"
  }
}

provider "aws" {
}

// If the Lambda function is installed in a non-master/mgmt account, it can
// list all accounts and inventory each one using the OrganizationAccessRole
// if accounts_info = "" and master_account_id and master_role_name are set
// and the roles are assumable by the Lambda function's IAM role
module "integration_test" {
  // source            = "github.com/GSA/grace-inventory/terraform?ref=latest"
  source            = "../"
  accounts_info     = "self"
  appenv            = var.appenv
  master_account_id = var.master_account_id
  master_role_name  = var.master_role_name
  tenant_role_name  = var.tenant_role_name
  source_file       = "../../release/grace-inventory-lambda.zip"
}

