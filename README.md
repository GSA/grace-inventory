# <a name="top">GRACE Inventory Lambda Function</a> [![GoDoc](https://godoc.org/github.com/GSA/grace-inventory?status.svg)](https://godoc.org/github.com/GSA/grace-circleci-builder) [![Go Report Card](https://goreportcard.com/badge/gojp/goreportcard)](https://goreportcard.com/report/github.com/GSA/grace-inventory)

**Lint Checks/Unit Tests:** [![CircleCI](https://circleci.com/gh/GSA/grace-inventory.svg?style=shield)](https://circleci.com/gh/GSA/grace-inventory)

**Integration Tests:** [![CircleCI](https://circleci.com/gh/GSA/grace-inventory-tests.svg?style=shield&circle-token=f86712ce5167665fe0d4a23d4af4fe7e9a20f7de)](https://circleci.com/gh/GSA/grace-inventory-tests)

Lambda function to create an inventory report of AWS services as an Excel
spreadsheet in an S3 bucket. Includes Terraform code to deploy the Lambda
function and create S3 bucket and necessary IAM roles/permissions. The Lambda
function can inventory all accounts in an AWS organization, specified
organizational units, a specified list of AWS accounts or simply services in the
account the Lambda function is installed in.

## Table of Contents

- [Security Compliance](#security-compliance)
- [Inventoried Services](#inventoried-services)
- [Repository contents](#repository-contents)
- [Usage](#usage)
    - [Download](#download)
    - [Build/Compile Locally](#buildcompile-locally)
        - [Prerequisites](#prerequisites)
        - [Build](#build)
    - [Example Usage](#example-usage)
    - [Intermittent Error](#intermittent-error)
- [Terraform Module Inputs](#terraform-module-inputs)
- [Terraform Module Outputs](#terraform-module-outputs)
- [Environment Variables](#environment-variables)
    - [CircleCI Environment Variables](#circleci-environment-variables)
    - [Lambda Function Environment Variables](#lambda-function-environment-variables)
- [Public domain](#public-domain)

[top](#top)

## Security Compliance

**Component ATO status:** draft

**Relevant controls:**

Control    | CSP/AWS | HOST/OS | App/DB | How is it implemented?
---------- | ------- | ------- | ------ | ----------------------
[CM-8](https://nvd.nist.gov/800-53/Rev4/control/CM-8) | ╳ | | | Employs an automated Lambda function triggered by a scheduled CloudWatch event (every 24 hours, by default).  Inventories supported services in specified AWS accounts and stores results in an Excel Spreadsheet on an S3 bucket.
[CM-8(2)](https://nvd.nist.gov/800-53/Rev4/control/CM-8) | ╳ | | | Automated by scheduled CloudWatch event (every 24 hours, by default).  Can be triggered more often to maintain an up-to-date, complete, accurate and readily available inventory of the AWS cloud service components.

[top](#top)

## Supported Services

    - Organization Accounts
    - IAM Roles
    - IAM Groups
    - IAM Policies
    - IAM Users
    - S3 Buckets
    - Glacier Vaults
    - EC2 Instances
    - Amazon Machine Images (AMI)
    - EBS Volumes
    - Snapshots
    - VPCs
    - Subnets
    - Security Groups
    - IP Addresses
    - Key Pairs
    - Elastic Load Balancers (elbv2)
    - CloudFormation Stacks
    - CloudWatch Alarms
    - Config Service Rules
    - KMS Keys
    - RDS Instances and Snapshots
    - Secrets Manager Secrets
    - SNS Subscriptions and Topics
    - SSM Parameter Stores

[top](#top)

## Repository contents

- **./**: Terraform module to deploy and configure Lambda function, S3 Bucket and IAM roles and policies
- **handler**: Go code for Lambda function
- **examples**: Examples of how to use the terraform module
- **tests**: Root module for testing deployment of Lambda function

[top](#top)

## Usage

To use the Terraform module to deploy the lambda function, you will need to either
download the binary release from Github or compile the handler locally.

[top](#top)

### Download (Recommended)

```bash
mkdir -p release
cd release
curl -L https://github.com/GSA/grace-inventory/releases/download/v0.1.3/grace-inventory-lambda.zip -o grace-inventory-lambda.zip
cd ..
```

[top](#top)

### Build/Compile Locally (Not Recommended)

#### Prerequisites

- Install the following prerequisites:
    1. [Go](https://golang.org/)
    1. [GolangCI-Lint](https://github.com/golangci/golangci-lint)
    1. [gosec](https://github.com/securego/gosec)
    1. [make](https://www.gnu.org/software/make/)

[top](#top)

#### Build

- After installing all required prerequisites: compile the lambda function and
put it in a zip compressed archive in `./release/grace-inventory-lambda.zip` by
entering the following at a command prompt:

```bash
make build_handler
```

#### Alternative Build (Not Recommended)

If your IAM permissions prevent the tests from succeeding, you can build manually:

```bash
mkdir -p release
cd handler
GOOS=linux GOARCH=amd64 go build -o ../release/grace-inventory-lambda -v
cd ../release
zip -j grace-inventory-lambda.zip grace-inventory-lambda
rm grace-inventory-lambda
cd ..
```

[top](#top)

### Example Usage

To inventory a single AWS account to which the Lambda function is deployed,
include the following in your root terraform module:

```
module "example_self" {
  source       = "github.com/GSA/grace-inventory?ref=v0.1.3"
  source_file  = "../../release/grace-inventory-lambda.zip"
  appenv       = "environment"
  project_name = "your-project"
}
```

Ensure the `source_file` parameter is the path to the zip archive containing
the compiled Lambda function handler downloaded or compiled earlier.


See the [examples](examples) directory for more examples.

**Note:** The S3 bucket created to store the inventory spreadsheets has logging
enabled and requires a pre-existing bucket with a name in the form of:
`${var.project_name}-${var.appenv}-access-logs`. The LogDelivery group must have
WRITE and READ_ACP permissions on the bucket (`acl = "log-delivery-write"`).
If your logging bucket has a different name or does not exist, you will have to
create one or fork this repository and edit the logging configuration of the S3
bucket on
[Line 10 of `iam.tf`](https://github.com/GSA/grace-inventory/blob/9b46d0bfbf40d6b9a5237afb9a45621a2f1a85d9/s3.tf#L10)

**Note:** The `DEFAULT_REGION` for the lambda function to write to the S3 Bucket
will be the first region listed in the `regions` attribute.  By default, this is
`us-east-1`.  If you want to place the S3 bucket in a different region, then you
will need to set the `regions` attribute with your desired region first in the
comma delimited list.

### Intermittent Error

The KMS key policy depends on the IAM role, however, even though Terraform creates
the IAM role first, there is sometimes a delay in the configuration reaching
eventual consistency within AWS. This can result in the following error:

```
Error: MalformedPolicyDocumentException: Policy contains a statement with one or more invalid principals.
	status code: 400, request id: 2425f0db-3033-448c-8e20-347eec8cac03

  on .terraform/modules/example_self/kms.tf line 1, in resource "aws_kms_key" "kms_key":
   1: resource "aws_kms_key" "kms_key" {
```

Re-applying the terraform (`terraform apply`) will usually resolve the problem.
Terraform considers this a retriable error, so you can also increase the
`max_retries` in the aws provider:

```
provider "aws" {
  max_retries = 5
}
```

[top](#top)

## Terraform Module Inputs

| Name | Description | Type | Default | Required |
|------|-------------|:----:|:-----:|:-----:|
| source\_file | \(optional\) full or relative path to zipped binary of lambda handler | string | `"../release/grace-inventory-lambda.zip"` | no |
| appenv | \(optional\) The environment in which the script is running \(development \| test \| production\) | string | `"development"` | no |
| project_name | \(required\) project name \(e.g. grace, fcs, fas, etc.\). Used as prefix for AWS S3 bucket name | string | `"grace"` | yes |
| access\_logging\_bucket | \(optional\) the S3 bucket that will receiving on-access logs for the inventory bucket | string | `""` | no |
| accounts\_info | \(optional\) Determines which accounts to parse.  Can be "self", comma delimited list of Account IDs or an S3 URI containing JSON output of `aws organizations list-accounts`.  If empty, tries to query accounts with `organizations:ListAccounts` | string | `"self"` | no |
| master\_account\_id | \(optional\) Account ID of AWS Master Payer Account | string | `""` | no |
| master\_role\_name | \(optional\) Role assumed by lambda function to query organizations in Master Payer account | string | `""` | no |
| organizational\_units | \(optional\) comma delimited list of organizational units to query for accounts. If set it will only query accounts in those organizational units | string | `""` | no |
| regions | \(optional\) Comma delimited list of AWS regions to inventory.  **Note:** The first region listed will be used by the lambda function as the `DEFAULT_REGION`. | string | `"us-east-1,us-east-2,us-west-1,us-west-2"` | no |
| schedule\_expression | \(optional\) Cloudwatch schedule expression for when to run inventory | string | `"cron(5 3 ? * MON-FRI *)"` | no |
| tenant\_role\_name | \(optional\) Role assumed by lambda function to query tenant accounts | string | `"OrganizationAccountAccessRole"` | no |
| lambda_memory | \(optional\) The number of megabytes of RAM for the lambda | number | 2048 | no |
| sheets | \(optional\) A comma delimited list of sheets | string | `""` | no |

[top](#top)

## Terraform Module Outputs

| Name | Description |
|------|-------------|
| lambda\_function\_arn | The Amazon Resource Name \(ARN\) identifying the Lambda Function |
| lambda\_function\_kms\_key\_arn | The ARN for the KMS encryption key |
| lambda\_function\_last\_modified | The date this resource was last modified |
| s3\_bucket\_id | The name of the S3 bucket where inventory reports are saved |

[top](#top)

## Environment Variables

### Lambda Function Environment Variables

| Name                 | Description |
| -------------------- | ------------|
| s3_bucket            | (required) S3 Bucket to store inventory reports |
| kms_key_id           | (required) ID of KMS key for encrypting/decrypting S3 bucket objects |
| regions              | (required) comma delimited list of regions to be inventoried |
| accounts_info        | (optional) If `accounts_info` is empty or not set, the function will try to query accounts via the Organizations API.  If set to "self", then it will only inventory its own account.  If set to an S3 URI for a file containing the json output of the `aws organizations list-accounts` command, it will query all accounts listed.  If set to a comma separated list of account IDs, it will query those accounts. |
| master_account       | (optional) Account ID of master payer account |
| organizational_units | (optional) comma delimited list of organizational units to query for accounts. If set it will only query accounts in those organizational units |
| tenant_role_name            | (optional) Role name used to inventory tenant accounts |
| master_role_name            | (optional) Role name to assume in master payer account for querying organizations |
| sheets | (optional) A comma delimited list of sheets that should be generated (see [sheets](#sheets))

[top](#top)

## Sheets

| Name | Permission | Description |
| ---- | ---------- | ----------- |
| Roles | iam:ListRoles | queries IAM Roles |
| Groups | iam:ListGroups | queries IAM Groups |
| Policies | iam:ListPolicies | queries IAM Policies |
| Users | iam:ListUsers | queries IAM Users |
| Buckets | s3:ListBuckets | queries S3 Buckets |
| Instances | ec2:DescribeInstances | queries EC2 Instances |
| Images | ec2:DescribeImages | queries EC2 Images |
| Volumes | ec2:DescribeVolumes | queries EC2 Volumes |
| Snapshots | ec2:DescribeSnapshots | queries EC2 Snapshots |
| VPCs | ec2:DescribeVpcs | queries EC2 VPCs |
| VpcPeers | ec2:DescribeVpcPeeringConnectionsPages | queries EC2 Vpc Peers |
| Subnets | ec2:DescribeSubnets | queries EC2 Subnets |
| SecurityGroups | ec2:DescribeSecurityGroups | queries EC2 Security Groups |
| Addresses | ec2:DescribeAddresses | queries EC2 Addresses |
| KeyPairs | ec2:DescribeKeyPairs | queries EC2 Key Pairs |
| Stacks | cloudformation:DescribeStacks | queries Cloud Formation Stacks |
| Alarms | cloudwatch:DescribeAlarms | queries CloudWatch Alarms |
| ConfigRules | config:DescribeConfig | queries AWS Config rules |
| LoadBlancers | elasticloadbalancing:DescribeLoadBalancers | queries Elastic Load Balancers |
| Vaults | glacier:ListVaults | queries Glacier Vaults |
| Keys | kms:ListKeys | queries KMS Keys |
| DBInstances | rds:DescribeDBInstances | queries RDS Database Instances |
| DBSnapshots | rds:DescribeDBSnapshots | queries RDS Database Snapshots |
| Secrets | secretsmanager:ListSecrets | queries Secrets Manager secrets |
| Subscriptions | sns:ListSubscriptions | queries Simple Notification Service Subscriptions |
| Topics | sns:ListTopics | queries Simple Notification Service Topics |
| Parameters | ssm:DescribeParameters | queries AWS Systems Manager Parameters |

## Public domain

This project is in the worldwide [public domain](LICENSE.md). As stated in [CONTRIBUTING](CONTRIBUTING.md):

> This project is in the public domain within the United States, and copyright and related rights in the work worldwide are waived through the [CC0 1.0 Universal public domain dedication](https://creativecommons.org/publicdomain/zero/1.0/).
>
> All contributions to this project will be released under the CC0 dedication. By submitting a pull request, you are agreeing to comply with this waiver of copyright interest.

[top](#top)
