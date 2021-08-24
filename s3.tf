locals {
  useAccessLogging = length(var.access_logging_bucket) > 0 ? [1] : []
}

#tfsec:ignore:AWS002
resource "aws_s3_bucket" "bucket" {
  bucket        = local.app_name
  acl           = "private"
  force_destroy = true

  versioning {
    enabled = true
  }

  dynamic "logging" {
    for_each = local.useAccessLogging
    content {
      target_bucket = var.access_logging_bucket
      target_prefix = "${local.app_name}-logs/"
    }
  }

  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        kms_master_key_id = aws_kms_key.kms_key.arn
        sse_algorithm     = "aws:kms"
      }
    }
  }

  lifecycle_rule {
    id      = "delete"
    enabled = true

    expiration {
      days = 7
    }
  }

  tags = {
    Name = "${upper(var.project_name)} Inventory Report"
  }
}

resource "aws_s3_bucket_public_access_block" "bucket" {
  bucket = aws_s3_bucket.bucket.id

  block_public_acls       = true
  ignore_public_acls      = true
  block_public_policy     = true
  restrict_public_buckets = true
}
