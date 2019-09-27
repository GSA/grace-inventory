variable "appenv" {
  type        = string
  description = "(optional) The environment in which the script is running (development | test | production)"
  default     = "integration-testing"
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

variable "master_account_id" {
  type        = string
  description = "(optional) Account ID of AWS Master Payer Account"
  default     = ""
}

