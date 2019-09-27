# <a name="top">GRACE Inventory Lambda Function</a>

**Lint Checks:** [![CircleCI](https://circleci.com/gh/GSA/grace-inventory.svg?style=svg)](https://circleci.com/gh/GSA/grace-inventory)

**Unit/Integration Tests:** [![CircleCI](https://circleci.com/gh/GSA/grace-inventory-tests.svg?style=svg&circle-token=f86712ce5167665fe0d4a23d4af4fe7e9a20f7de)](https://circleci.com/gh/GSA/grace-inventory-tests)

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
- [Terraform Module Inputs](#terraform-module-inputs)
- [Terraform Module Outputs](#terraform-module-outputs)
- [Non-Module Installation](#non-module-installation)
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

- **handler**: Go code for Lambda function
- **terraform**: Terraform module to deploy and configure Lambda function, S3 Bucket and IAM roles and policies
    - **examples**: Examples of how to use the terraform module
    - **tests**: Root module for testing deployment of Lambda function

[top](#top)

## Usage

To use the Terraform module to deploy the lambda function, you will need to either
download the binary release from Github or compile the handler locally.

[top](#top)

### Download

```bash
mkdir -p release
cd release
curl -L https://github.com/GSA/grace-inventory/releases/download/v0.1.1/grace-inventory-lambda.zip -o grace-inventory-lambda.zip
cd ..
```

[top](#top)

### Build/Compile Locally

#### Prerequisites

- Install the following prerequisites:
    1. [Go](https://golang.org/)
    1. [Dep](https://golang.github.io/dep/docs/installation.html)
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

#### Alternative Build

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
  source       = "github.com/GSA/grace-inventory/terraform"
  source_file  = "../../release/grace-inventory-lambda.zip"
  appenv       = "environment"
  project_name = "your-project"
}
```

Ensure the `source_file` parameter is the path to the zip archive containing
the compiled Lambda function handler downloaded or compiled earlier.


See the [examples](terraform/examples) directory for more examples.

[top](#top)

## Terraform Module Inputs

| Name | Description | Type | Default | Required |
|------|-------------|:----:|:-----:|:-----:|
| source\_file | \(optional\) full or relative path to zipped binary of lambda handler | string | `"../release/grace-inventory-lambda.zip"` | no |
| appenv | \(optional\) The environment in which the script is running \(development \| test \| production\) | string | `"development"` | no |
| project_name | \(required\) project name \(e.g. grace, fcs, fas, etc.\). Used as prefix for AWS S3 bucket name | string | `"grace"` | yes |
| accounts\_info | \(optional\) Determines which accounts to parse.  Can be "self", comma delimited list of Account IDs or an S3 URI containing JSON output of `aws organizations list-accounts`.  If empty, tries to query accounts with `organizations:ListAccounts` | string | `"self"` | no |
| master\_account\_id | \(optional\) Account ID of AWS Master Payer Account | string | `""` | no |
| master\_role\_name | \(optional\) Role assumed by lambda function to query organizations in Master Payer account | string | `""` | no |
| organizational\_units | \(optional\) comma delimited list of organizational units to query for accounts. If set it will only query accounts in those organizational units | string | `""` | no |
| regions | \(optional\) Comma delimited list of AWS regions to inventory | string | `"us-east-1,us-east-2,us-west-1,us-west-2"` | no |
| schedule\_expression | \(optional\) Cloudwatch schedule expression for when to run inventory | string | `"cron(5 3 ? * MON-FRI *)"` | no |
| tenant\_role\_name | \(optional\) Role assumed by lambda function to query tenant accounts | string | `"OrganizationAccountAccessRole"` | no |

[top](#top)

## Terraform Module Outputs

| Name | Description |
|------|-------------|
| lambda\_function\_arn | The Amazon Resource Name \(ARN\) identifying the Lambda Function |
| lambda\_function\_kms\_key\_arn | The ARN for the KMS encryption key |
| lambda\_function\_last\_modified | The date this resource was last modified |
| s3\_bucket\_id | The name of the S3 bucket where inventory reports are saved |

[top](#top)

## Non-Module Installation

It is also possible to build and apply locally without using as a Terraform
module.

1. Install system dependencies.
    1. [Go](https://golang.org/)
    1. [Dep](https://golang.github.io/dep/docs/installation.html)
    1. [GolangCI-Lint](https://github.com/golangci/golangci-lint)
    1. [gosec](https://github.com/securego/gosec)
    1. [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/installing.html)
    1. [Terraform](https://www.terraform.io/)
1. [Configure AWS](https://www.terraform.io/docs/providers/aws/#authentication) with credentials for your AWS account locally.
1. Set the the environment variables specified in the CircleCI section below.
1. Copy the `terrafrom/terraform.tfvars.example` file to `terraform/terraform.tfvars` and set the values as necessary for your environment.
1. Validate and test the code

    ```bash
    make test
    ```

1. Build and deploy

    ```bash
    make deploy
    ```

[top](#top)

## Environment Variables

### CircleCI Environment Variables

| Name                              | Description |
| --------------------------------- | ------------|
| AWS_DEFAULT_REGION                | default AWS region |
| DEVELOPMENT_AWS_ACCESS_KEY_ID     | AWS access key for deployment to development environment |
| DEVELOPMENT_AWS_SECRET_ACCESS_KEY | AWS secret key for deployment to development environment |
| DEVELOPMENT_MASTER_ACCT_ID        | Account ID of master payer account |
| TEST_AWS_ACCESS_KEY_ID            | AWS access key for deployment to test environment |
| TEST_AWS_SECRET_ACCESS_KEY        | AWS secret key for deployment to test environment |
| TEST_MASTER_ACCT_ID               | Account ID of master payer account |
| PRODUCTION_AWS_ACCESS_KEY_ID      | AWS access key for deployment to production environment |
| PRODUCTION_AWS_SECRET_ACCESS_KEY  | AWS secret key for deployment to production environment |
| PRODUCTION_MASTER_ACCT_ID               | Account ID of master payer account |
| TF_VAR_regions                    | comma delimited list of regions to be inventoried |
| TF_VAR_tenant_role_name           | Role name used to inventory tenant accounts |
| TF_VAR_master_role_name           | Role name to assume in master payer account for querying organizations |
| TF_VAR_schedule_expression        | Cloudwatch schedule expression for scheduling Lambda function |
| backend_bucket                    | S3 Bucket for saving shared Terraform state file |
| backend_key                       | S3 Bucket Key for saving shared Terraform state file |

[top](#top)

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

[top](#top)

## Public domain

This project is in the worldwide [public domain](LICENSE.md). As stated in [CONTRIBUTING](CONTRIBUTING.md):

> This project is in the public domain within the United States, and copyright and related rights in the work worldwide are waived through the [CC0 1.0 Universal public domain dedication](https://creativecommons.org/publicdomain/zero/1.0/).
>
> All contributions to this project will be released under the CC0 dedication. By submitting a pull request, you are agreeing to comply with this waiver of copyright interest.

[top](#top)
