package helpers

import (
	"testing"

	awstest "github.com/gruntwork-io/terratest/modules/aws"
)

const defaultRegion = "us-east-1"

// func Roles(sess *session.Session, cred *credentials.Credentials) ([]*iam.Role, error)
func TestRoles(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	_, err = Roles(sess, nil)
	if err != nil {
		t.Fatalf("Roles() failed: %v", err)
	}
}

// func Groups(sess *session.Session, cred *credentials.Credentials) ([]*iam.Group, error)
func TestGroups(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	_, err = Groups(sess, nil)
	if err != nil {
		t.Fatalf("Groups() failed: %v", err)
	}
}

// func Policies(sess *session.Session, cred *credentials.Credentials) ([]*iam.Policy, error)
func TestPolicies(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	_, err = Policies(sess, nil)
	if err != nil {
		t.Fatalf("Policies() failed: %v", err)
	}
}

// func Users(sess *session.Session, cred *credentials.Credentials) ([]*iam.User, error)
func TestUsers(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	_, err = Users(sess, nil)
	if err != nil {
		t.Fatalf("Users() failed: %v", err)
	}
}

// func Buckets(sess *session.Session, cred *credentials.Credentials) ([]*s3.Bucket, error)
func TestBuckets(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	_, err = Buckets(sess, nil)
	if err != nil {
		t.Fatalf("Buckets() failed: %v", err)
	}
}

// func Instances(sess *session.Session, cred *credentials.Credentials) ([]*ec2.Instance, error)
func TestInstances(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	_, err = Instances(sess, nil)
	if err != nil {
		t.Fatalf("Instances() failed: %v", err)
	}
}

// func Images(sess *session.Session, cred *credentials.Credentials) ([]*ec2.Image, error)
func TestImages(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	_, err = Images(sess, nil)
	if err != nil {
		t.Fatalf("Images() failed: %v", err)
	}
}

// func Volumes(sess *session.Session, cred *credentials.Credentials) ([]*ec2.Volume, error)
func TestVolumes(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	_, err = Volumes(sess, nil)
	if err != nil {
		t.Fatalf("Volumes() failed: %v", err)
	}
}

// func Snapshots(sess *session.Session, cred *credentials.Credentials) ([]*ec2.Snapshot, error)
func TestSnapshots(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	_, err = Snapshots(sess, nil)
	if err != nil {
		t.Fatalf("Snapshots() failed: %v", err)
	}
}

// func Vpcs(sess *session.Session, cred *credentials.Credentials) ([]*ec2.Vpc, error)
func TestVpcs(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	_, err = Vpcs(sess, nil)
	if err != nil {
		t.Fatalf("Vpcs() failed: %v", err)
	}
}

// func Subnets(sess *session.Session, cred *credentials.Credentials) ([]*ec2.Subnet, error)
func TestSubnets(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	_, err = Subnets(sess, nil)
	if err != nil {
		t.Fatalf("Subnets() failed: %v", err)
	}
}

// func SecurityGroups(sess *session.Session, cred *credentials.Credentials) ([]*ec2.SecurityGroup, error)
func TestSecurityGroups(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	_, err = SecurityGroups(sess, nil)
	if err != nil {
		t.Fatalf("SecurityGroups() failed: %v", err)
	}
}

// func Addresses(sess *session.Session, cred *credentials.Credentials) ([]*ec2.Address, error)
func TestAddresses(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	_, err = Addresses(sess, nil)
	if err != nil {
		t.Fatalf("Addresses() failed: %v", err)
	}
}

// func KeyPairs(sess *session.Session, cred *credentials.Credentials) ([]*ec2.KeyPairInfo, error)
func TestKeyPairs(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	_, err = KeyPairs(sess, nil)
	if err != nil {
		t.Fatalf("KeyPairs() failed: %v", err)
	}
}

// func Stacks(sess *session.Session, cred *credentials.Credentials) ([]*cloudformation.Stack, error)
func TestStacks(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	_, err = Stacks(sess, nil)
	if err != nil {
		t.Fatalf("Stacks() failed: %v", err)
	}
}

// func Alarms(sess *session.Session, cred *credentials.Credentials) ([]*cloudwatch.MetricAlarm, error)
func TestAlarms(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	_, err = Alarms(sess, nil)
	if err != nil {
		t.Fatalf("Alarms() failed: %v", err)
	}
}

// func ConfigRules(sess *session.Session, cred *credentials.Credentials) ([]*configservice.ConfigRule, error)
func TestConfigRules(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	_, err = ConfigRules(sess, nil)
	if err != nil {
		t.Fatalf("ConfigRules() failed: %v", err)
	}
}

// func LoadBalancers(sess *session.Session, cred *credentials.Credentials) ([]*elbv2.LoadBalancer, error)
func TestLoadBalancers(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	_, err = LoadBalancers(sess, nil)
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
	_, err = Vaults(sess, nil)
	if err != nil {
		t.Fatalf("Vaults() failed: %v", err)
	}
}

// func Keys(cfg client.ConfigProvider, cred *credentials.Credentials) ([]*KmsKey, error) {
func TestKeys(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	_, err = Keys(sess, nil)
	if err != nil {
		t.Fatalf("Keys() failed: %v", err)
	}
}

// func DBInstances(cfg client.ConfigProvider, cred *credentials.Credentials) ([]*rds.DBInstance, error) {
func TestDBInstances(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	_, err = DBInstances(sess, nil)
	if err != nil {
		t.Fatalf("DBInstances() failed: %v", err)
	}
}

// func DBSnapshots(cfg client.ConfigProvider, cred *credentials.Credentials) ([]*rds.DBInstance, error) {
func TestDBSnapshots(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	_, err = DBSnapshots(sess, nil)
	if err != nil {
		t.Fatalf("DBSnapshots() failed: %v", err)
	}
}

// func Secrets(cfg client.ConfigProvider, cred *credentials.Credentials) ([]*secretsmanager.SecretListEntry, error) {
func TestSecrets(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	_, err = Secrets(sess, nil)
	if err != nil {
		t.Fatalf("Secrets() failed: %v", err)
	}
}

// func Subscriptions(cfg client.ConfigProvider, cred *credentials.Credentials) ([]*sns.Subscription, error) {
func TestSubscriptions(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	_, err = Subscriptions(sess, nil)
	if err != nil {
		t.Fatalf("Subscriptions() failed: %v", err)
	}
}

// func Topics(cfg client.ConfigProvider, cred *credentials.Credentials) ([]*SnsTopic, error) {
func TestTopics(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	_, err = Topics(sess, nil)
	if err != nil {
		t.Fatalf("Topics() failed: %v", err)
	}
}

//func Parameters(cfg client.ConfigProvider, cred *credentials.Credentials) ([]*ssm.ParameterMetadata, error) {
func TestParameters(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	_, err = Parameters(sess, nil)
	if err != nil {
		t.Fatalf("Parameters() failed: %v", err)
	}
}
