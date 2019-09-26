data "aws_caller_identity" "current" {
}

locals {
  app_name       = "${var.project_name}-${var.appenv}-inventory"
  account_id     = data.aws_caller_identity.current.account_id
  logging_bucket = var.appenv == "integration-testing" ? "grace-development-access-logs" : "${var.project_name}-${var.appenv}-access-logs"
}

