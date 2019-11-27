package inv

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
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
	"github.com/aws/aws-sdk-go/service/glacier"
	"github.com/aws/aws-sdk-go/service/glacier/glacieriface"
	"github.com/aws/aws-sdk-go/service/organizations"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"gotest.tools/assert"
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

func TestNew(t *testing.T) {
	// expected := Inv{}
	_, err := New()
	assert.Error(t, err, `env: required environment variable "s3_bucket" is not set`)
	// assert.DeepEqual(t, expected, actual)
}

func TestGetCurrentIdentity(t *testing.T) {
	expected := sts.GetCallerIdentityOutput{
		Account: aws.String("a"),
		Arn:     aws.String("b"),
		UserId:  aws.String("c"),
	}
	svc := stsSvc{Client: &mockStsClient{Response: expected}}
	actual, _ := svc.getCurrentIdentity()

	if !reflect.DeepEqual(&expected, actual) {
		t.Errorf("failed to get caller identity, expected: %#v, got: %#v", expected, actual)
	}
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
	if err != nil {
		t.Errorf("failed to walk accounts: %v", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("failed to match accounts, expected: %#v, got: %#v", expected, actual)
	}
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
	if err != nil {
		t.Errorf("failed to walk accounts: %v", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("failed to match accounts, expected: %#v, got: %#v", expected, actual)
	}
}

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
	if err != nil {
		t.Errorf("failed to queryVaults: %v", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("failed to match objects, expected: %s, got: %s", expected, actual)
	}
}
