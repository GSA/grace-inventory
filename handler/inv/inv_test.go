package inv

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/GSA/grace-inventory/handler/helpers/credmgr"
	"github.com/GSA/grace-inventory/handler/helpers/sessionmgr"
	"github.com/GSA/grace-inventory/handler/spreadsheet"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/awstesting/mock"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/glacier"
	"github.com/aws/aws-sdk-go/service/glacier/glacieriface"
	"github.com/aws/aws-sdk-go/service/organizations"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/google/go-cmp/cmp"
	"gotest.tools/v3/assert"
)

///////////////////
// Mock Services //
///////////////////

type mockStsClient struct {
	stsiface.STSAPI
	Response sts.GetCallerIdentityOutput
}

func (m *mockStsClient) GetCallerIdentity(*sts.GetCallerIdentityInput) (*sts.GetCallerIdentityOutput, error) {
	return &m.Response, nil
}

func mockNewSession(cfgs ...*aws.Config) (*session.Session, error) {
	// server is the mock server that simply writes a 200 status back to the client
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	for _, cfg := range cfgs {
		cfg.DisableSSL = aws.Bool(true)
		cfg.Endpoint = aws.String(server.URL)
	}
	return session.NewSession(cfgs...)
}

////////////////////////////
// inv Package Unit Tests //
////////////////////////////

func TestIsKnownError(t *testing.T) {
	err := fmt.Errorf("%s", "test")
	tt := map[string]struct {
		in       error
		expected bool
	}{
		"nil input": {},
		"non-query error": {
			in:       err,
			expected: false,
		},
		"non-awserr query error": {
			in:       newQueryErrorf(err, "%s", "test"),
			expected: false,
		},
		"unknown awserr": {
			in:       newQueryErrorf(awserr.New("test", "test", err), "%s", "test"),
			expected: false,
		},
		"known awserr AccessDenied": {
			in:       newQueryErrorf(awserr.New("AccessDenied", "test", err), "%s", "test"),
			expected: true,
		},
	}
	for name, tc := range tt {
		tc := tc
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, isKnownError(tc.in))
		})
	}
}

// TestNew ... not very useful without a way of mocking/stubbing the API calls
// to get the current identity and sessions. This unit test will fail if you
// have valid credentials configured in your environment
func TestNew(t *testing.T) {
	tt := map[string]struct {
		env         map[string]string
		expected    Inv
		expectedErr string
	}{
		"environment variables not set": {
			expectedErr: `required environment variable "s3_bucket" is not set`,
		},
		"no credentials": {
			env: map[string]string{
				"s3_bucket":  "test",
				"kms_key_id": "test",
				"regions":    "us-east-1",
			},
			expectedErr: "NoCredentialProviders: no valid providers in chain. Deprecated.\n\tFor verbose messaging see aws.Config.CredentialsChainVerboseErrors",
		},
		"happy path": {
			env: map[string]string{
				"s3_bucket":  "test",
				"kms_key_id": "test",
				"regions":    "us-east-1",
			},
		},
	}
	// Remove test cases that would fail during integration testing because environment variables are set
	if os.Getenv("s3_bucket") != "" {
		t.Log("deleting test case for environment variables not set")
		delete(tt, "environment variables not set")
	}
	if os.Getenv("AWS_ACCESS_KEY_ID") != "" {
		t.Log("deleting test case for no credentials")
		delete(tt, "no credentials")
	} else {
		t.Log("deleting test case for happy path")
		delete(tt, "happy path")
	}
	for name, tc := range tt {
		tc := tc
		t.Run(name, func(t *testing.T) {
			for k, v := range tc.env {
				err := os.Setenv(k, v)
				if err != nil {
					t.Fatalf("error setting environment variable: %v", err)
				}
			}

			actual, err := New()

			envBucket := os.Getenv("s3_bucket")
			envKmsKey := os.Getenv("kms_key_id")
			// maps are dynamically randomized in memory
			// we must cleanup the ENV before running the
			// next test
			for k := range tc.env {
				enverr := os.Unsetenv(k)
				if enverr != nil {
					t.Fatalf("error removing environment variable: %v", enverr)
				}
			}

			if tc.expectedErr == "" {
				assert.NilError(t, err)
				assert.Equal(t, actual.bucketID, envBucket)
				assert.Equal(t, actual.kmsKeyID, envKmsKey)
			} else {
				assert.ErrorContains(t, err, tc.expectedErr)
			}
		})
	}
}

func TestGetCurrentIdentity(t *testing.T) {
	expected := sts.GetCallerIdentityOutput{
		Account: aws.String("a"),
		Arn:     aws.String("b"),
		UserId:  aws.String("c"),
	}
	svc := stsSvc{Client: &mockStsClient{Response: expected}}
	actual, _ := svc.getCurrentIdentity()

	assert.DeepEqual(t, actual, &expected, cmp.AllowUnexported(sts.GetCallerIdentityOutput{}))
}

func TestWalkAccounts(t *testing.T) {
	sessMgr := sessionmgr.New("us-east-1", []string{"us-east-1", "us-west-1"})
	sessMgr.Sessioner(mockNewSession)
	err := sessMgr.Init()
	if err != nil {
		t.Fatalf("failed to instantiate session manager: %v", err)
	}

	expected := []*organizations.Account{
		{Id: aws.String("a"), Name: aws.String("a")},
		{Id: aws.String("b"), Name: aws.String("b")},
		{Id: aws.String("c"), Name: aws.String("c")},
	}

	inv := &Inv{
		accounts:   expected,
		credMgr:    credmgr.New(mock.Session, "", "", expected),
		sessionMgr: sessMgr,
	}

	var actual []*organizations.Account
	_, err = inv.walkAccounts(func(name string, credentials *credentials.Credentials, sess *session.Session) (*spreadsheet.Payload, error) {
		actual = append(actual, &organizations.Account{
			Id:   aws.String(name),
			Name: aws.String(name),
		})
		return nil, nil
	})
	assert.NilError(t, err)
	assert.DeepEqual(t, actual, expected, cmp.AllowUnexported(organizations.Account{}))
}

func TestWalkSessions(t *testing.T) {
	regions := []string{"us-east-1", "us-west-1"}
	sessMgr := sessionmgr.New("us-east-1", regions)
	sessMgr.Sessioner(mockNewSession)
	err := sessMgr.Init()
	if err != nil {
		t.Fatalf("failed to instantiate session manager: %v", err)
	}

	accounts := []*organizations.Account{
		{Id: aws.String("a"), Name: aws.String("a")},
		{Id: aws.String("b"), Name: aws.String("b")},
		{Id: aws.String("c"), Name: aws.String("c")},
	}

	inv := &Inv{
		accounts:   accounts,
		credMgr:    credmgr.New(mock.Session, "", "", accounts),
		sessionMgr: sessMgr,
	}

	// walkSessions calls iterates over each region
	// for each account
	var expected []*organizations.Account
	for i := 0; i < len(accounts); i++ {
		for j := 0; j < len(regions); j++ {
			expected = append(expected, accounts[i])
		}
	}

	var actual []*organizations.Account
	_, err = inv.walkSessions(func(name string, credentials *credentials.Credentials, sess *session.Session) (*spreadsheet.Payload, error) {
		t.Logf("region: %s\n", aws.StringValue(sess.Config.Region))
		actual = append(actual, &organizations.Account{
			Id:   aws.String(name),
			Name: aws.String(name),
		})
		return nil, nil
	})
	assert.NilError(t, err)
	assert.DeepEqual(t, actual, expected, cmp.AllowUnexported(organizations.Account{}))
}

/////////////////////////////////////
// Mocks for querying EC2 Services //
/////////////////////////////////////

type mockEc2Client struct {
	ec2iface.EC2API
}

func (m mockEc2Client) DescribeInstancesPages(in *ec2.DescribeInstancesInput, fn func(*ec2.DescribeInstancesOutput, bool) bool) error {
	fn(&ec2.DescribeInstancesOutput{
		Reservations: []*ec2.Reservation{
			{
				Instances: []*ec2.Instance{
					{InstanceId: aws.String("i-1234567890abcdef0")},
				},
			}}}, true)
	return nil
}

func (m mockEc2Client) DescribeImages(in *ec2.DescribeImagesInput) (*ec2.DescribeImagesOutput, error) {
	return &ec2.DescribeImagesOutput{
		Images: []*ec2.Image{
			{ImageId: aws.String("ami-5731123e")},
		},
	}, nil
}

func (m mockEc2Client) DescribeVolumesPages(in *ec2.DescribeVolumesInput, fn func(*ec2.DescribeVolumesOutput, bool) bool) error {
	fn(&ec2.DescribeVolumesOutput{
		Volumes: []*ec2.Volume{
			{VolumeId: aws.String("vol-049df61146c4d7901")},
		}}, true)
	return nil
}

func (m mockEc2Client) DescribeSnapshotsPages(in *ec2.DescribeSnapshotsInput, fn func(*ec2.DescribeSnapshotsOutput, bool) bool) error {
	fn(&ec2.DescribeSnapshotsOutput{
		Snapshots: []*ec2.Snapshot{
			{SnapshotId: aws.String("snap-1234567890abcdef0")},
		}}, true)
	return nil
}

func (m mockEc2Client) DescribeVpcsPages(in *ec2.DescribeVpcsInput, fn func(*ec2.DescribeVpcsOutput, bool) bool) error {
	fn(&ec2.DescribeVpcsOutput{
		Vpcs: []*ec2.Vpc{
			{VpcId: aws.String("vpc-0e9801d129EXAMPLE")},
		}}, true)
	return nil
}

func (m mockEc2Client) DescribeSubnetsPages(in *ec2.DescribeSubnetsInput, fn func(*ec2.DescribeSubnetsOutput, bool) bool) error {
	fn(&ec2.DescribeSubnetsOutput{
		Subnets: []*ec2.Subnet{
			{SubnetId: aws.String("subnet-0bb1c79de3EXAMPLE")},
		}}, true)
	return nil
}

func (m mockEc2Client) DescribeSecurityGroupsPages(in *ec2.DescribeSecurityGroupsInput, fn func(*ec2.DescribeSecurityGroupsOutput, bool) bool) error {
	fn(&ec2.DescribeSecurityGroupsOutput{
		SecurityGroups: []*ec2.SecurityGroup{
			{GroupId: aws.String("sg-903004f8")},
		}}, true)
	return nil
}

func (m mockEc2Client) DescribeAddresses(in *ec2.DescribeAddressesInput) (*ec2.DescribeAddressesOutput, error) {
	return &ec2.DescribeAddressesOutput{
		Addresses: []*ec2.Address{
			{
				InstanceId:     aws.String("i-1234567890abcdef0"),
				PublicIp:       aws.String("127.0.0.1"),
				PublicIpv4Pool: aws.String("amazon"),
				Domain:         aws.String("standard"),
			},
		},
	}, nil
}

func (m mockEc2Client) DescribeKeyPairs(in *ec2.DescribeKeyPairsInput) (*ec2.DescribeKeyPairsOutput, error) {
	return &ec2.DescribeKeyPairsOutput{
		KeyPairs: []*ec2.KeyPairInfo{
			{KeyName: aws.String("test")},
		},
	}, nil
}

func mockEc2Creator(client.ConfigProvider, ...*aws.Config) ec2iface.EC2API {
	return mockEc2Client{}
}

func mockInv(t *testing.T) *Inv {
	regions := []string{"us-east-1", "us-west-1"}
	sessMgr := sessionmgr.New("us-east-1", regions)
	sessMgr.Sessioner(mockNewSession)
	err := sessMgr.Init()
	if err != nil {
		t.Fatalf("failed to instantiate session manager: %v", err)
	}

	accounts := []*organizations.Account{
		{Id: aws.String("a"), Name: aws.String("a")},
		{Id: aws.String("b"), Name: aws.String("b")},
		{Id: aws.String("c"), Name: aws.String("c")},
	}

	inv := &Inv{
		accounts:   accounts,
		credMgr:    credmgr.New(mock.Session, "", "", accounts),
		sessionMgr: sessMgr,
	}
	return inv
}

func TestQueryInstances(t *testing.T) {
	inv := mockInv(t)
	ec2Creator = mockEc2Creator
	_, err := inv.queryInstances()
	assert.NilError(t, err)
}

func TestQueryImages(t *testing.T) {
	inv := mockInv(t)
	ec2Creator = mockEc2Creator
	_, err := inv.queryImages()
	assert.NilError(t, err)
}

func TestQueryVolumes(t *testing.T) {
	inv := mockInv(t)
	ec2Creator = mockEc2Creator
	_, err := inv.queryVolumes()
	assert.NilError(t, err)
}

func TestQuerySnapshots(t *testing.T) {
	inv := mockInv(t)
	ec2Creator = mockEc2Creator
	_, err := inv.querySnapshots()
	assert.NilError(t, err)
}

func TestQueryVpcs(t *testing.T) {
	inv := mockInv(t)
	ec2Creator = mockEc2Creator
	_, err := inv.queryVpcs()
	assert.NilError(t, err)
}

func TestQuerySubnets(t *testing.T) {
	inv := mockInv(t)
	ec2Creator = mockEc2Creator
	_, err := inv.querySubnets()
	assert.NilError(t, err)
}

func TestQuerySecurityGroups(t *testing.T) {
	inv := mockInv(t)
	ec2Creator = mockEc2Creator
	_, err := inv.querySecurityGroups()
	assert.NilError(t, err)
}

func TestQueryAddresses(t *testing.T) {
	inv := mockInv(t)
	ec2Creator = mockEc2Creator
	_, err := inv.queryAddresses()
	assert.NilError(t, err)
}

func TestQueryKeyPairs(t *testing.T) {
	inv := mockInv(t)
	ec2Creator = mockEc2Creator
	_, err := inv.queryKeyPairs()
	assert.NilError(t, err)
}

///////////////////////////////////
// Mocks for testing queryVaults //
///////////////////////////////////

type mockGlacierClient struct {
	glacieriface.GlacierAPI
}

func (m mockGlacierClient) ListVaultsPages(in *glacier.ListVaultsInput, fn func(*glacier.ListVaultsOutput, bool) bool) error {
	fn(&glacier.ListVaultsOutput{
		VaultList: []*glacier.DescribeVaultOutput{
			{VaultARN: aws.String("a"), VaultName: aws.String("a")},
			{VaultARN: aws.String("b"), VaultName: aws.String("b")},
			{VaultARN: aws.String("c"), VaultName: aws.String("c")},
			{VaultARN: aws.String("d"), VaultName: aws.String("d")},
		},
	}, true)
	return nil
}

func mockGlacierCreator(client.ConfigProvider, ...*aws.Config) glacieriface.GlacierAPI {
	return mockGlacierClient{}
}

func genericListVaultsOutput() []interface{} {
	in := []*glacier.DescribeVaultOutput{
		{VaultARN: aws.String("a"), VaultName: aws.String("a")},
		{VaultARN: aws.String("b"), VaultName: aws.String("b")},
		{VaultARN: aws.String("c"), VaultName: aws.String("c")},
		{VaultARN: aws.String("d"), VaultName: aws.String("d")},
	}
	var out []interface{}
	for _, i := range in {
		out = append(out, i)
	}
	return out
}

func TestQueryVaults(t *testing.T) {
	regions := []string{"us-east-1", "us-west-1"}
	sessMgr := sessionmgr.New("us-east-1", regions)
	sessMgr.Sessioner(mockNewSession)
	err := sessMgr.Init()
	if err != nil {
		t.Fatalf("failed to instantiate session manager: %v", err)
	}

	accounts := []*organizations.Account{
		{Id: aws.String("a"), Name: aws.String("a")},
		{Id: aws.String("b"), Name: aws.String("b")},
		{Id: aws.String("c"), Name: aws.String("c")},
	}

	inv := &Inv{
		accounts:   accounts,
		credMgr:    credmgr.New(mock.Session, "", "", accounts),
		sessionMgr: sessMgr,
	}

	var expected []*spreadsheet.Payload
	for i := 0; i < len(accounts); i++ {
		for j := 0; j < len(regions); j++ {
			expected = append(expected, &spreadsheet.Payload{
				Static: []string{
					aws.StringValue(accounts[i].Name),
					regions[j],
				},
				Items: genericListVaultsOutput(),
			})
		}
	}

	glacierCreator = mockGlacierCreator
	actual, err := inv.queryVaults()
	assert.NilError(t, err)
	assert.DeepEqual(t, actual, expected, cmp.AllowUnexported(spreadsheet.Payload{}, glacier.DescribeVaultOutput{}))
}
