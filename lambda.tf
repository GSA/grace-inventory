resource "aws_lambda_function" "lambda_function" {
  filename         = var.source_file
  function_name    = local.app_name
  description      = "Creates report of AWS Services in Organization accounts and saves to Excel spreadsheet in S3 bucket"
  role             = aws_iam_role.iam_role.arn
  handler          = "grace-inventory-lambda"
  source_code_hash = filebase64sha256(var.source_file)
  kms_key_arn      = aws_kms_key.kms_key.arn
  runtime          = "go1.x"
  timeout          = 900

  environment {
    variables = {
      accounts_info     = var.accounts_info
      kms_key_id        = aws_kms_key.kms_key.key_id
      master_role_name  = var.master_role_name
      master_account_id = var.master_account_id
      // organizational_units = "${organizational_units}"
      regions          = var.regions
      s3_bucket        = aws_s3_bucket.bucket.bucket
      tenant_role_name = var.tenant_role_name
    }
  }
}

resource "aws_lambda_permission" "lambda_permission" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.lambda_function.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.cwe_rule.arn
}

