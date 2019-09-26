resource "aws_iam_role" "iam_role" {
  name        = local.app_name
  description = "Role for GRACE Inventory Lambda function"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF

}

resource "aws_iam_policy" "iam_policy" {
  name        = local.app_name
  description = "Policy to allow creating GRACE service inventory report"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "cloudformation:DescribeStacks",
        "cloudwatch:DescribeAlarms",
        "config:DescribeConfigRules",
        "ec2:DescribeAddresses",
        "ec2:DescribeImages",
        "ec2:DescribeInstances",
        "ec2:DescribeKeyPairs",
        "ec2:DescribeSecurityGroups",
        "ec2:DescribeSnapshots",
        "ec2:DescribeSubnets",
        "ec2:DescribeVolumes",
        "ec2:DescribeVpcs",
        "elasticloadbalancing:DescribeLoadBalancers",
        "glacier:ListVaults",
        "iam:GetUser",
        "iam:ListAccountAliases",
        "iam:ListGroups",
        "iam:ListPolicies",
        "iam:ListRoles",
        "iam:ListUsers",
        "kms:ListKeys",
        "kms:DescribeKey",
        "kms:ListAliases",
        "organizations:ListAccounts",
        "organizations:ListAccountsForParent",
        "rds:DescribeDBInstances",
        "rds:DescribeDBSnapshots",
        "s3:ListBucket",
        "s3:ListAllMyBuckets",
        "s3:HeadBucket",
        "secretsmanager:ListSecrets",
        "sns:GetTopicAttributes",
        "sns:ListSubscriptions",
        "sns:ListTopics",
        "ssm:DescribeParameters",
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Effect": "Allow",
      "Resource": "*"
    },
    {
      "Action": [
        "sts:AssumeRole"
      ],
      "Effect": "Allow",
      "Resource": [
        "${aws_iam_role.iam_role.arn}",
        "arn:aws:iam::${var.master_account_id}:role/${var.master_role_name}",
        "arn:aws:iam::*:role/${var.tenant_role_name}"
      ]
    },
    {
      "Action": [
        "s3:GetObject",
        "s3:PutObject"
      ],
      "Effect": "Allow",
      "Resource": "${aws_s3_bucket.bucket.arn}/*"
    },
    {
      "Effect": "Allow",
      "Action": [
        "kms:Encrypt"
      ],
      "Resource": "${aws_kms_key.kms_key.arn}"
    }
  ]
}
EOF

}

resource "aws_iam_role_policy_attachment" "iam_role_policy_attachment" {
  role       = aws_iam_role.iam_role.name
  policy_arn = aws_iam_policy.iam_policy.arn
}

