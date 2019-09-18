resource "aws_s3_bucket" "bucket" {
  bucket        = "${local.app_name}"
  acl           = "private"
  force_destroy = true

  versioning {
    enabled = true
  }

  logging {
    target_bucket = "${local.logging_bucket}"
    target_prefix = "${local.app_name}-logs/"
  }

  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        kms_master_key_id = "${aws_kms_key.kms_key.arn}"
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
    Name = "GRACE Inventory Report"
  }
}
