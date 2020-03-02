package integration

import (
	"log"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/GSA/grace-tftest/aws/cloudwatchevents"
	"github.com/GSA/grace-tftest/aws/iam"
	"github.com/GSA/grace-tftest/aws/lambda"
	"github.com/GSA/grace-tftest/aws/s3"

	"github.com/GSA/grace-tftest/aws/kms"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
)

var appName = "grace-integration-test-inventory"

func TestAll(t *testing.T) {
	url := "http://localhost:" + os.Getenv("MOTO_PORT")
	sess, err := session.NewSession(&aws.Config{
		Endpoint:         aws.String(url),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	})
	if err != nil {
		t.Fatalf("failed to connect to moto: %s -> %v", url, err)
	}

	log.Println("Starting integration test: default")
	keyID, keyArn := testKmsKey(t, sess)
	roleArn := testRole(t, sess, keyArn)
	ruleArn := testCWERule(t, sess)
	regions := ""
	lambdaArn := testLambda(t, sess, roleArn, ruleArn, keyArn, keyID, regions)
	testCWETarget(t, sess, lambdaArn)
	testBucket(t, sess, keyArn)
}

func testLambda(t *testing.T, cfg client.ConfigProvider, roleArn string,
	ruleArn string, keyArn string, keyID string, regions string) string {
	log.Println("entered testLambda")
	defer log.Println("exited testLambda")

	handlerName := "grace-inventory-lambda"
	svc := lambda.New(cfg, appName)

	svc.Config.
		Role(roleArn).
		Handler(handlerName).
		KeyArn(keyArn).
		Runtime("go1.x").
		Timeout(900).
		//Env("s3_bucket", bucketName).
		//Env("kms_key_id", keyID).
		//Env("regions", regions). Moto failed to return this value
		Assert(t, nil)

	lambdaArn := aws.StringValue(svc.Config.Get(t).FunctionArn)

	svc.Policy.Statement(t, nil).
		Sid("AllowExecutionFromCloudWatch").
		Action("lambda:InvokeFunction").
		Effect("Allow").
		Principal("Service", "events.amazonaws.com").
		Resource(lambdaArn+":$LATEST").
		Condition("ArnLike", "AWS:SourceArn", ruleArn).
		Assert(t)

	return lambdaArn
}

func testKmsKey(t *testing.T, cfg client.ConfigProvider) (string, string) {
	log.Println("entered testKmsKey")
	defer log.Println("exited testKmsKey")

	svc := kms.
		New(cfg).
		Alias.
		Name("alias/" + appName).
		Assert(t)

	stmt := svc.
		Policy(t).
		Statement(t, nil)

	stmt.
		Sid("Enable IAM User Permissions").
		Action("kms:*").
		Effect("Allow").
		Resource("*").
		Assert(t)

	stmt.
		Sid("Allow use of the key").
		Effect("Allow").
		Resource("*").
		Action(
			"kms:Encrypt",
			"kms:Decrypt",
			"kms:ReEncrypt*",
			"kms:GenerateDataKey*",
			"kms:DescribeKey",
		).
		Assert(t)

	return aws.StringValue(svc.Selected().TargetKeyId), aws.StringValue(svc.Key(t).Arn)
}

// nolint: funlen
func testRole(t *testing.T, cfg client.ConfigProvider, keyArn string) string {
	log.Println("entered testRole")
	defer log.Println("exited testRole")

	role := iam.
		New(cfg).
		Role.
		Name(appName).
		Assert(t)

	roleArn := aws.StringValue(role.Selected().Arn)
	role.
		Attached().
		Name(appName).
		Assert(t)

	stmt := iam.
		New(cfg).
		Policy.
		Name(appName).
		Assert(t).
		Statement(t, nil)

	stmt.
		Effect("Allow").
		Resource("*").
		Action(
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
			"logs:PutLogEvents",
		).
		Assert(t)

	// Broken with moto and fake roles
	//stmt.Action("sts:AssumeRole").Effect("Allow").
	//	Resource(roleArn, "role", "tenant-role").Assert(t)

	stmt.
		Effect("Allow").
		Action("s3:GetObject", "s3:PutObject").
		Resource("arn:aws:s3:::" + appName + "/*").
		Assert(t)

	stmt.
		Effect("Allow").
		Action("kms:Encrypt").
		Resource(keyArn).
		Assert(t)

	return roleArn
}

func testCWERule(t *testing.T, cfg client.ConfigProvider) string {
	log.Println("entered testCWERule")
	defer log.Println("exited testCWERule")

	rule := cloudwatchevents.
		New(cfg).
		Rule.
		Name(appName).
		State("enabled").
		SchedExpr("cron(5 3 ? * MON-FRI *)").
		Assert(t).
		Selected()

	if rule == nil {
		return ""
	}
	return aws.StringValue(rule.Arn)
}

func testCWETarget(t *testing.T, cfg client.ConfigProvider, targetArn string) {
	log.Println("entered testCWETarget")
	defer log.Println("exited testCWETarget")

	svc := cloudwatchevents.
		New(cfg).
		Rule.
		Name(appName).
		Assert(t)

	svc.
		Target().
		Arn(targetArn).
		Assert(t)
}

func testBucket(t *testing.T, cfg client.ConfigProvider, keyArn string) {
	log.Println("entered testBucket")
	defer log.Println("exited testBucket")

	svc := s3.New(cfg).
		Bucket.
		Name(appName).
		Assert(t)

	svc.
		Lifecycle().
		Status("enabled").
		Method("delete").
		ExpDays(7).
		Assert(t)

	// NO MOTO SUPPORT
	// svc.Encryption().IsSSE().ID(keyArn).Alg("aws:kms").Assert(t)
}
