package helpers

import (
	"errors"
	"reflect"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
	"github.com/aws/aws-sdk-go/service/configservice"
	"github.com/aws/aws-sdk-go/service/configservice/configserviceiface"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
	"github.com/aws/aws-sdk-go/service/glacier"
	"github.com/aws/aws-sdk-go/service/glacier/glacieriface"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/kms/kmsiface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"

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

func (m mockCFClient) DescribeStacksPages(in *cloudformation.DescribeStacksInput, fn func(*cloudformation.DescribeStacksOutput, bool) bool) error {
	fn(&cloudformation.DescribeStacksOutput{Stacks: []*cloudformation.Stack{{}}}, true)
	return nil
}

type mockCWClient struct {
	cloudwatchiface.CloudWatchAPI
}

func (m mockCWClient) DescribeAlarmsPages(in *cloudwatch.DescribeAlarmsInput, fn func(*cloudwatch.DescribeAlarmsOutput, bool) bool) error {
	fn(&cloudwatch.DescribeAlarmsOutput{MetricAlarms: []*cloudwatch.MetricAlarm{{}}}, true)
	return nil
}

type mockCSClient struct {
	configserviceiface.ConfigServiceAPI
}

func (m mockCSClient) DescribeConfigRules(in *configservice.DescribeConfigRulesInput) (*configservice.DescribeConfigRulesOutput, error) {
	return &configservice.DescribeConfigRulesOutput{ConfigRules: []*configservice.ConfigRule{{}}}, nil
}

type mockElbClient struct {
	elbv2iface.ELBV2API
}

func (m mockElbClient) DescribeLoadBalancersPages(in *elbv2.DescribeLoadBalancersInput, fn func(*elbv2.DescribeLoadBalancersOutput, bool) bool) error {
	fn(&elbv2.DescribeLoadBalancersOutput{LoadBalancers: []*elbv2.LoadBalancer{{}}}, true)
	return nil
}

type mockKmsClient struct {
	kmsiface.KMSAPI
}

func (m mockKmsClient) ListKeysPages(in *kms.ListKeysInput, fn func(*kms.ListKeysOutput, bool) bool) error {
	fn(&kms.ListKeysOutput{Keys: []*kms.KeyListEntry{{}}}, true)
	return nil
}

func (m mockKmsClient) DescribeKey(in *kms.DescribeKeyInput) (*kms.DescribeKeyOutput, error) {
	return &kms.DescribeKeyOutput{KeyMetadata: &kms.KeyMetadata{}}, nil
}

func (m mockKmsClient) ListAliasesPages(in *kms.ListAliasesInput, fn func(*kms.ListAliasesOutput, bool) bool) error {
	fn(&kms.ListAliasesOutput{Aliases: []*kms.AliasListEntry{{}}}, true)
	return nil
}

type mockSsmClient struct {
	ssmiface.SSMAPI
}

func (m mockSsmClient) DescribeParametersPages(in *ssm.DescribeParametersInput, fn func(*ssm.DescribeParametersOutput, bool) bool) error {
	fn(&ssm.DescribeParametersOutput{Parameters: []*ssm.ParameterMetadata{{}}}, true)
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

// func Stacks(svc *cloudformation.CloudFormationAPI) ([]*cloudformation.Stack, error)
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

// func Alarms(svc *cloudwatch.CloudWatchAPI) ([]*cloudwatch.MetricAlarm, error)
func TestAlarms(t *testing.T) {
	expected := []*cloudwatch.MetricAlarm{{}}
	svc := mockCWClient{}
	got, err := Alarms(svc)
	if err != nil {
		t.Fatalf("Alarms() failed: %v", err)
	}
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("Alarms() failed.\nExpected %#v (%T)\nGot: %#v (%T)\n", expected, expected, got, got)
	}
}

// func ConfigRules(svc *configservice.ConfigServiceAPI) ([]*configservice.ConfigRule, error)
func TestConfigRules(t *testing.T) {
	expected := []*configservice.ConfigRule{{}}
	svc := mockCSClient{}
	got, err := ConfigRules(svc)
	if err != nil {
		t.Fatalf("ConfigRules() failed: %v", err)
	}
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("ConfigRules() failed.\nExpected %#v (%T)\nGot: %#v (%T)\n", expected, expected, got, got)
	}
}

// func LoadBalancers(svc *elbv2.ELBV2API) ([]*elbv2.LoadBalancer, error)
func TestLoadBalancers(t *testing.T) {
	expected := []*elbv2.LoadBalancer{{}}
	svc := mockElbClient{}
	got, err := LoadBalancers(svc)
	if err != nil {
		t.Fatalf("LoadBalancers() failed: %v", err)
	}
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("LoadBalancers() failed.\nExpected %#v (%T)\nGot: %#v (%T)\n", expected, expected, got, got)
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

// func Keys(svc *kms.ELBV2API) ([]*kms.Key, error)
func TestKeys(t *testing.T) {
	expected := []*KmsKey{{}}
	svc := mockKmsClient{}
	got, err := Keys(svc)
	if err != nil {
		t.Fatalf("Keys() failed: %v", err)
	}
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("Keys() failed.\nExpected %#v (%T)\nGot: %#v (%T)\n", expected, expected, got, got)
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

// func Parameters(svc *ssm.SSMAPI) ([]*ssm.Parameter, error)
func TestParameters(t *testing.T) {
	expected := []*ssm.ParameterMetadata{{}}
	svc := mockSsmClient{}
	got, err := Parameters(svc)
	if err != nil {
		t.Fatalf("Parameters() failed: %v", err)
	}
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("Parameters() failed.\nExpected %#v (%T)\nGot: %#v (%T)\n", expected, expected, got, got)
	}
}
