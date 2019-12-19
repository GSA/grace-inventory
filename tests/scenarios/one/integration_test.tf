terraform {
  backend "local" {
    path = "terraform.tfstate"
  }
}

provider "aws" {
  endpoints {
    lambda           = "http://localhost:4574"
    s3               = "http://localhost:4572"
    sns              = "http://localhost:4575"
    sqs              = "http://localhost:4576"
    ses              = "http://localhost:4579"
    cloudwatch       = "http://localhost:4582"
    cloudwatchlogs   = "http://localhost:4586"
    cloudwatchevents = "http://localhost:4587"
    sts              = "http://localhost:4592"
    iam              = "http://localhost:4593"
    kms              = "http://localhost:4599"
  }
}

// If the Lambda function is installed in a non-master/mgmt account, it can
// list all accounts and inventory each one using the OrganizationAccessRole
// if accounts_info = "" and master_account_id and master_role_name are set
// and the roles are assumable by the Lambda function's IAM role
module "integration_test" {
  // source            = "github.com/GSA/grace-inventory?ref=latest"
  source            = "../../../"
  accounts_info     = "self"
  project_name      = "grace"
  appenv            = var.appenv
  master_account_id = var.master_account_id
  master_role_name  = var.master_role_name
  tenant_role_name  = var.tenant_role_name
  source_file       = "../../../release/grace-inventory-lambda.zip"
}
