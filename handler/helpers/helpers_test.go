package helpers

import (
	"errors"
	"reflect"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/aws/aws-sdk-go/service/glacier"
	"github.com/aws/aws-sdk-go/service/glacier/glacieriface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"

	awstest "github.com/gruntwork-io/terratest/modules/aws"
)

const defaultRegion = "us-east-1"

type mockS3Client struct {
	s3iface.S3API
}

func (m mockS3Client) ListBuckets(in *s3.ListBucketsInput) (*s3.ListBucketsOutput, error) {
	return &s3.ListBucketsOutput{Buckets: []*s3.Bucket{{}}}, nil
}

type mockCFClient struct {
	cloudformationiface.CloudFormationAPI
}

func (m mockIamClient) DescribeStacksPages(in *cloudformation.DescribeStacksInput, fn func(*cloudformation.DescribeStacksOutput, bool) bool) error {
	fn(&cloudformation.DescribeStacksOutput{Stacks: []*cloudformation.Stack{{}}}, true)
	return nil
}

// func Buckets(sess *session.Session, cred *credentials.Credentials) ([]*s3.Bucket, error)
func TestBuckets(t *testing.T) {
	svc := mockS3Client{}
	expected := []*s3.Bucket{{}}
	got, err := Buckets(svc)
	if err != nil {
		t.Fatalf("Buckets() failed: %v", err)
	}
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("Buckets() failed.\nExpected %#v (%T)\nGot: %#v (%T)\n", expected, expected, got, got)
	}
}

// func Stacks(sess *session.Session, cred *credentials.Credentials) ([]*cloudformation.Stack, error)
func TestStacks(t *testing.T) {
	expected := []*cloudformation.Stack{{}}
	svc := mockCFClient{}
	got, err := Stacks(svc)
	if err != nil {
		t.Fatalf("Stacks() failed: %v", err)
	}
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("Stacks() failed.\nExpected %#v (%T)\nGot: %#v (%T)\n", expected, expected, got, got)
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

type mockGlacierClient struct {
	glacieriface.GlacierAPI
	pages int
}

func (m *mockGlacierClient) ListVaultsPages(in *glacier.ListVaultsInput, fn func(*glacier.ListVaultsOutput, bool) bool) error {
	for i := 0; i < m.pages; i++ {
		if !fn(m.listVaultsPagesR(in, i)) {
			return nil
		}
	}
	return errors.New("function is continuing to request pages after all pages have been consumed")
}

func (m *mockGlacierClient) listVaultsPagesR(in *glacier.ListVaultsInput, index int) (out *glacier.ListVaultsOutput, lastPage bool) {
	var (
		limit = 10 // default items per page
		items []*glacier.DescribeVaultOutput
	)

	if in.Limit != nil {
		l, err := strconv.Atoi(aws.StringValue(in.Limit))
		if err == nil {
			limit = l
		}
	}

	for i := 0; i < limit; i++ {
		items = append(items, &glacier.DescribeVaultOutput{
			NumberOfArchives: aws.Int64(int64(index)), // store page index
			SizeInBytes:      aws.Int64(int64(i)),     // store element index
		})
	}

	out = &glacier.ListVaultsOutput{
		VaultList: items,
	}
	if index+1 == m.pages {
		lastPage = true
	}
	return
}

func TestVaultsErr(t *testing.T) {
	svc := &GlacierSvc{Client: &mockGlacierClient{pages: 0}}
	_, err := svc.Vaults()
	if err == nil {
		t.Error("err value was nil when failure was expected")
	}
}

func TestVaultsPagination(t *testing.T) {
	tt := []struct {
		Name                  string
		Pages                 int
		ExpectedLength        int
		ExpectedLastPageIndex int
		ExpectedLastElemIndex int
	}{
		{Name: "validate single page request", Pages: 1, ExpectedLength: 10, ExpectedLastPageIndex: 0, ExpectedLastElemIndex: 9},
		{Name: "validate two page request", Pages: 2, ExpectedLength: 20, ExpectedLastPageIndex: 1, ExpectedLastElemIndex: 9},
		{Name: "validate three page request", Pages: 3, ExpectedLength: 30, ExpectedLastPageIndex: 2, ExpectedLastElemIndex: 9},
	}

	for _, st := range tt {
		tc := st
		t.Run(tc.Name, func(t *testing.T) {
			svc := &GlacierSvc{Client: &mockGlacierClient{pages: tc.Pages}}
			items, err := svc.Vaults()
			if err != nil {
				t.Fatalf("Vaults() failed: %v", err)
			}

			length := len(items)
			if length != tc.ExpectedLength {
				t.Errorf("items length invalid, expected: %d, got: %d", tc.ExpectedLength, length)
			}

			lastItem := items[length-1]
			lastPageIndex := int(aws.Int64Value(lastItem.NumberOfArchives))
			if lastPageIndex != tc.ExpectedLastPageIndex {
				t.Errorf("last page index invalid, expected: %d, got: %d", tc.ExpectedLastPageIndex, lastPageIndex)
			}

			lastElemIndex := int(aws.Int64Value(lastItem.SizeInBytes))
			if lastElemIndex != tc.ExpectedLastElemIndex {
				t.Errorf("last element index invalid, expected: %d, got: %d", tc.ExpectedLastElemIndex, lastElemIndex)
			}
		})
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
