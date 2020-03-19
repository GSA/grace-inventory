variable "schedule_expression" {
  type        = string
  description = "(optional) Cloudwatch schedule expression for when to run inventory"
  default     = "cron(5 3 ? * MON-FRI *)"
}

variable "tenant_role_name" {
  type        = string
  description = "(optional) Role assumed by lambda function to query tenant accounts"
  default     = "OrganizationAccountAccessRole"
}

variable "master_role_name" {
  type        = string
  description = "(optional) Role assumed by lambda function to query organizations in Master Payer account"
  default     = ""
}

variable "regions" {
  type        = string
  description = "(optional) Comma delimited list of AWS regions to inventory"
  default     = "us-east-1,us-east-2,us-west-1,us-west-2"
}

variable "appenv" {
  type        = string
  description = "(optional) The environment in which the script is running (development | test | production)"
  default     = "development"
}

variable "master_account_id" {
  type        = string
  description = "(optional) Account ID of AWS Master Payer Account"
  default     = ""
}

variable "accounts_info" {
  type        = string
  description = "(optional) Determines which accounts to parse.  Can be \"self\", comma delimited list of Account IDs or an S3 URI containing JSON output of `aws organizations list-accounts`.  If empty, tries to query accounts with `organizations:ListAccounts`"
  default     = "self"
}

variable "organizational_units" {
  type        = string
  description = "(optional) comma delimited list of organizational units to query for accounts. If set it will only query accounts in those organizational units"
  default     = ""
}

variable "source_file" {
  type        = string
  description = "(optional) full or relative path to zipped binary of lambda handler"
  default     = "../release/grace-inventory-lambda.zip"
}

variable "project_name" {
  type        = string
  description = "(required) project name (e.g. grace, fcs, fas, etc.). Used as prefix for AWS S3 bucket name"
  default     = "grace"
}

variable "access_logging_bucket" {
  type        = string
  description = "(optional) the S3 bucket that will receiving on-access logs for the inventory bucket"
  default     = ""
}

