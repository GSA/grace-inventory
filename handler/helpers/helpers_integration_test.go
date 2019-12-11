// +build integration

package helpers

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/configservice"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/glacier"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/ssm"
	awstest "github.com/gruntwork-io/terratest/modules/aws"
)

const defaultRegion = "us-east-1"

// func Roles(sess *session.Session, cred *credentials.Credentials) ([]*iam.Role, error)
func TestIntegrationRoles(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	svc := IamSvc{Client: iam.New(sess)}
	_, err = svc.Roles()
	if err != nil {
		t.Fatalf("Roles() failed: %v", err)
	}
}

// func Groups(sess *session.Session, cred *credentials.Credentials) ([]*iam.Group, error)
func TestIntegrationGroups(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	svc := IamSvc{Client: iam.New(sess)}
	_, err = svc.Groups()
	if err != nil {
		t.Fatalf("Groups() failed: %v", err)
	}
}

// func Policies(sess *session.Session, cred *credentials.Credentials) ([]*iam.Policy, error)
func TestIntegrationPolicies(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	svc := IamSvc{Client: iam.New(sess)}
	_, err = svc.Policies()
	if err != nil {
		t.Fatalf("Policies() failed: %v", err)
	}
}

// func Users(sess *session.Session, cred *credentials.Credentials) ([]*iam.User, error)
func TestIntegrationUsers(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	svc := IamSvc{Client: iam.New(sess)}
	_, err = svc.Users()
	if err != nil {
		t.Fatalf("Users() failed: %v", err)
	}
}

// func Buckets(sess *session.Session, cred *credentials.Credentials) ([]*s3.Bucket, error)
func TestIntegrationBuckets(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	svc := s3.New(sess)
	_, err = Buckets(svc)
	if err != nil {
		t.Fatalf("Buckets() failed: %v", err)
	}
}

// func Instances(sess *session.Session, cred *credentials.Credentials) ([]*ec2.Instance, error)
func TestIntegrationInstances(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	svc := Ec2Svc{Client: ec2.New(sess)}
	_, err = svc.Instances()
	if err != nil {
		t.Fatalf("Instances() failed: %v", err)
	}
}

// func Images(sess *session.Session, cred *credentials.Credentials) ([]*ec2.Image, error)
func TestIntegrationImages(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	svc := Ec2Svc{Client: ec2.New(sess)}
	_, err = svc.Images()
	if err != nil {
		t.Fatalf("Images() failed: %v", err)
	}
}

// func Volumes(sess *session.Session, cred *credentials.Credentials) ([]*ec2.Volume, error)
func TestIntegrationVolumes(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	svc := Ec2Svc{Client: ec2.New(sess)}
	_, err = svc.Volumes()
	if err != nil {
		t.Fatalf("Volumes() failed: %v", err)
	}
}

// func Snapshots(sess *session.Session, cred *credentials.Credentials) ([]*ec2.Snapshot, error)
func TestIntegrationSnapshots(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	svc := Ec2Svc{Client: ec2.New(sess)}
	_, err = svc.Snapshots()
	if err != nil {
		t.Fatalf("Snapshots() failed: %v", err)
	}
}

// func Vpcs(sess *session.Session, cred *credentials.Credentials) ([]*ec2.Vpc, error)
func TestIntegrationVpcs(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	svc := Ec2Svc{Client: ec2.New(sess)}
	_, err = svc.Vpcs()
	if err != nil {
		t.Fatalf("Vpcs() failed: %v", err)
	}
}

// func Subnets(sess *session.Session, cred *credentials.Credentials) ([]*ec2.Subnet, error)
func TestIntegrationSubnets(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	svc := Ec2Svc{Client: ec2.New(sess)}
	_, err = svc.Subnets()
	if err != nil {
		t.Fatalf("Subnets() failed: %v", err)
	}
}

// func SecurityGroups(sess *session.Session, cred *credentials.Credentials) ([]*ec2.SecurityGroup, error)
func TestIntegrationSecurityGroups(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	svc := Ec2Svc{Client: ec2.New(sess)}
	_, err = svc.SecurityGroups()
	if err != nil {
		t.Fatalf("SecurityGroups() failed: %v", err)
	}
}

// func Addresses(sess *session.Session, cred *credentials.Credentials) ([]*ec2.Address, error)
func TestIntegrationAddresses(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	svc := Ec2Svc{Client: ec2.New(sess)}
	_, err = svc.Addresses()
	if err != nil {
		t.Fatalf("Addresses() failed: %v", err)
	}
}

// func KeyPairs(sess *session.Session, cred *credentials.Credentials) ([]*ec2.KeyPairInfo, error)
func TestIntegrationKeyPairs(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	svc := Ec2Svc{Client: ec2.New(sess)}
	_, err = svc.KeyPairs()
	if err != nil {
		t.Fatalf("KeyPairs() failed: %v", err)
	}
}

// func Stacks(sess *session.Session, cred *credentials.Credentials) ([]*cloudformation.Stack, error)
func TestIntegrationStacks(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	_, err = Stacks(cloudformation.New(sess))
	if err != nil {
		t.Fatalf("Stacks() failed: %v", err)
	}
}

// func Alarms(sess *session.Session, cred *credentials.Credentials) ([]*cloudwatch.MetricAlarm, error)
func TestIntegrationAlarms(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	_, err = Alarms(cloudwatch.New(sess))
	if err != nil {
		t.Fatalf("Alarms() failed: %v", err)
	}
}

// func ConfigRules(sess *session.Session, cred *credentials.Credentials) ([]*configservice.ConfigRule, error)
func TestIntegrationConfigRules(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	_, err = ConfigRules(configservice.New(sess))
	if err != nil {
		t.Fatalf("ConfigRules() failed: %v", err)
	}
}

// func LoadBalancers(sess *session.Session, cred *credentials.Credentials) ([]*elbv2.LoadBalancer, error)
func TestIntegrationLoadBalancers(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	_, err = LoadBalancers(elbv2.New(sess))
	if err != nil {
		t.Fatalf("LoadBalancers() failed: %v", err)
	}
}

// func Vaults(sess *session.Session, cred *credentials.Credentials) ([]*glacier.DescribeVaultOutput, error)
func TestVaults(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	svc := &GlacierSvc{Client: glacier.New(sess)}
	_, err = svc.Vaults()
	if err != nil {
		t.Fatalf("Vaults() failed: %v", err)
	}
}

// func Keys(cfg client.ConfigProvider, cred *credentials.Credentials) ([]*KmsKey, error) {
func TestIntegrationKeys(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	_, err = Keys(kms.New(sess))
	if err != nil {
		t.Fatalf("Keys() failed: %v", err)
	}
}

// func DBInstances(cfg client.ConfigProvider, cred *credentials.Credentials) ([]*rds.DBInstance, error) {
func TestIntegrationDBInstances(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	svc := &RDSSvc{Client: rds.New(sess)}
	_, err = svc.DBInstances()
	if err != nil {
		t.Fatalf("DBInstances() failed: %v", err)
	}
}

// func DBSnapshots(cfg client.ConfigProvider, cred *credentials.Credentials) ([]*rds.DBInstance, error) {
func TestIntegrationDBSnapshots(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	svc := &RDSSvc{Client: rds.New(sess)}
	_, err = svc.DBSnapshots()
	if err != nil {
		t.Fatalf("DBSnapshots() failed: %v", err)
	}
}

// func Secrets(cfg client.ConfigProvider, cred *credentials.Credentials) ([]*secretsmanager.SecretListEntry, error) {
func TestIntegrationSecrets(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	svc := &SecretsManagerSvc{Client: secretsmanager.New(sess)}
	_, err = svc.Secrets()
	if err != nil {
		t.Fatalf("Secrets() failed: %v", err)
	}
}

// func Subscriptions(cfg client.ConfigProvider, cred *credentials.Credentials) ([]*sns.Subscription, error) {
func TestIntegrationSubscriptions(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	_, err = Subscriptions(sns.New(sess))
	if err != nil {
		t.Fatalf("Subscriptions() failed: %v", err)
	}
}

// func Topics(cfg client.ConfigProvider, cred *credentials.Credentials) ([]*SnsTopic, error) {
func TestIntegrationTopics(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	_, err = Topics(sns.New(sess))
	if err != nil {
		t.Fatalf("Topics() failed: %v", err)
	}
}

//func Parameters(cfg client.ConfigProvider, cred *credentials.Credentials) ([]*ssm.ParameterMetadata, error) {
func TestIntegrationParameters(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	_, err = Parameters(ssm.New(sess))
	if err != nil {
		t.Fatalf("Parameters() failed: %v", err)
	}
}
