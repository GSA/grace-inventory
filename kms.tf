resource "aws_kms_key" "kms_key" {
  description             = "Key for GRACE service inventory reporting S3 bucket"
  deletion_window_in_days = 7
  enable_key_rotation     = "true"
  depends_on              = [aws_iam_role.iam_role]

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "Enable IAM User Permissions",
      "Effect": "Allow",
      "Principal": {
        "AWS": "arn:aws:iam::${local.account_id}:root"
      },
      "Action": "kms:*",
      "Resource": "*"
    },
    {
      "Sid": "Allow use of the key",
      "Effect": "Allow",
      "Principal": {
        "AWS": [
          "${aws_iam_role.iam_role.arn}"
        ]
      },
      "Action": [
        "kms:Encrypt",
        "kms:Decrypt",
        "kms:ReEncrypt*",
        "kms:GenerateDataKey*",
        "kms:DescribeKey"
      ],
      "Resource": "*"
    }
  ]
}
EOF

}

resource "aws_kms_alias" "kms_alias" {
  name          = "alias/${local.app_name}"
  target_key_id = aws_kms_key.kms_key.key_id
}

