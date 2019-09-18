terraform {
  backend "s3" {
    region = "us-east-1"
  }
}

provider "aws" {}

data "aws_caller_identity" "current" {}

locals {
  app_name       = "grace-${var.appenv}-inventory"
  account_id     = "${data.aws_caller_identity.current.account_id}"
  logging_bucket = "${"${var.appenv}" == "integration-testing" ? "grace-development-access-logs" : "grace-${var.appenv}-access-logs"}"
}
