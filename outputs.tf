output "lambda_function_arn" {
  value       = aws_lambda_function.lambda_function.arn
  description = "The Amazon Resource Name (ARN) identifying the Lambda Function"
}

output "lambda_function_last_modified" {
  value       = aws_lambda_function.lambda_function.last_modified
  description = "The date this resource was last modified"
}

output "lambda_function_kms_key_arn" {
  value       = aws_lambda_function.lambda_function.kms_key_arn
  description = "The ARN for the KMS encryption key"
}

output "s3_bucket_id" {
  value       = aws_s3_bucket.bucket.id
  description = "The name of the S3 bucket where inventry reports are saved"
}

