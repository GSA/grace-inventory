package accounts

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/organizations"
	"github.com/aws/aws-sdk-go/service/organizations/organizationsiface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"gotest.tools/assert"
)

var rID = regexp.MustCompile(`^\d{12}$`)

// newAPIStub ... Creates a new httptest.Server to respond to AWS API calls
func newAPIStub() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
}

// newStubSession ... creates a *Session with a stub endpoint
func newStubSession(t *testing.T) *session.Session {
	stub := newAPIStub()
	creds := credentials.NewStaticCredentials("test", "test", "test")
	sess, err := session.NewSession(&aws.Config{
		Endpoint:    aws.String(stub.URL),
		Region:      aws.String("us-east-1"),
		Credentials: creds,
	})
	if err != nil {
		t.Fatal(err)
	}
	return sess
}

///////////////////
// Mock Services //
///////////////////

// mockIamSvc ... creates a mock of the AWS Identity and Access Management (IAM) service
type mockIamSvc struct {
	iamiface.IAMAPI
	Resp iam.ListAccountAliasesOutput
}

func (m mockIamSvc) ListAccountAliases(in *iam.ListAccountAliasesInput) (*iam.ListAccountAliasesOutput, error) {
	return &m.Resp, nil
}

// mockOrgSvc ... creates a mock of the AWS Organizations service
type mockOrgSvc struct {
	organizationsiface.OrganizationsAPI
	Resp organizations.ListAccountsOutput
}

func (m mockOrgSvc) ListAccounts(in *organizations.ListAccountsInput) (*organizations.ListAccountsOutput, error) {
	return &m.Resp, nil
}

// mockDownloaderSvc ... creates a mock of the AWS S3 Manager Downloader API
type mockDownloaderSvc struct {
	s3manageriface.DownloaderAPI
}

func (m mockDownloaderSvc) Download(w io.WriterAt, in *s3.GetObjectInput, fn ...func(*s3manager.Downloader)) (int64, error) {
	s := `{
	  "Accounts": [
	    {"Id": "111111111111"},
	    {"Id": "222222222222"},
	    {"Id": "333333333333"}
	  ]
	}`
	buf := bytes.NewBufferString(s)
	n, err := w.WriteAt(buf.Bytes(), 0)
	/*
		//Alternatively: Read from file
		dat, err := ioutil.ReadFile(aws.StringValue(in.Bucket) + aws.StringValue(in.Key))
		if err != nil {
			return 0, err
		}
		n, err := w.WriteAt(dat, 0)
	*/
	return int64(n), err
}

// mockStsSvc ... create a mock of the AWS Security Token Service (STS)
type mockStsSvc struct {
	stsiface.STSAPI
	TestInput func(*sts.AssumeRoleInput)
}

func (s mockStsSvc) AssumeRole(input *sts.AssumeRoleInput) (*sts.AssumeRoleOutput, error) {
	if s.TestInput != nil {
		s.TestInput(input)
	}
	expiry := time.Now().Add(60 * time.Minute)
	return &sts.AssumeRoleOutput{
		Credentials: &sts.Credentials{
			// Just reflect the role arn to the provider.
			AccessKeyId:     input.RoleArn,
			SecretAccessKey: aws.String("assumedSecretAccessKey"),
			SessionToken:    aws.String("assumedSessionToken"),
			Expiration:      &expiry,
		},
	}, nil
}

var mockSvc = Svc{
	iamSvc:           mockIamSvc{},
	organizationsSvc: mockOrgSvc{},
	downloaderSvc:    mockDownloaderSvc{},
	stsSvc:           mockStsSvc{},
}

/////////////////////////////////
// Accounts Package Unit Tests //
/////////////////////////////////

func TestNewAccountsSvc(t *testing.T) {
	// test case table
	tt := map[string]struct {
		sess        *session.Session
		expectedErr string
	}{"nil client.ConfigProvider": {
		expectedErr: "nil ConfigProvider",
	}, "happy path": {
		sess: newStubSession(t),
	}}
	// loop through test cases
	for name, tc := range tt {
		tc := tc
		t.Run(name, func(t *testing.T) {
			_, err := NewAccountsSvc(tc.sess)
			if tc.expectedErr == "" {
				assert.NilError(t, err)
			} else {
				assert.Error(t, err, tc.expectedErr)
			}
		})
	}
}

//nolint: gocyclo
func TestAccountsList(t *testing.T) {
	t.Run("queryAccounts", func(t *testing.T) {
		options := Options{
			MasterAccountID: "123456789012",
			MasterRoleName:  "masterRoleName",
		}
		mockSvc.organizationsSvc = mockOrgSvc{
			Resp: organizations.ListAccountsOutput{
				Accounts: []*organizations.Account{{Id: aws.String("123456789012")}},
			},
		}
		accounts, err := mockSvc.AccountsList(options)
		if err != nil {
			t.Fatalf("Accounts() failed: %v", err)
		}
		if len(accounts) < 1 {
			t.Fatal("expected at least one account")
		}
		if !rID.MatchString(*accounts[0].Id) {
			t.Fatalf("expected first account ID to be 12 digit number.  Got: %v", *accounts[0].Id)
		}
	})

	t.Run("invalid AccountsInfo", func(t *testing.T) {
		options := Options{
			AccountsInfo: "invalid",
		}
		_, err := mockSvc.AccountsList(options)
		if err == nil {
			t.Fatalf("expected failure for invalid accounts_info")
		} else if err.Error() != "invalid accounts_info" {
			t.Fatalf("expected 'invalid accounts_info' error.  Got: %v", err)
		}
	})

	t.Run("self", func(t *testing.T) {
		accountID := "123456789012"
		accountName := "test"
		options := Options{
			AccountsInfo:  "self",
			MgmtAccountID: "123456789012",
		}
		mockSvc.organizationsSvc = mockOrgSvc{
			Resp: organizations.ListAccountsOutput{
				Accounts: []*organizations.Account{{
					Id:   aws.String(accountID),
					Name: aws.String(accountName)}},
			},
		}
		mockSvc.iamSvc = mockIamSvc{
			Resp: iam.ListAccountAliasesOutput{
				AccountAliases: []*string{aws.String("test")},
			},
		}
		accounts, err := mockSvc.AccountsList(options)
		if err != nil {
			t.Fatalf("Accounts() failed: %v", err)
		}
		if len(accounts) != 1 {
			t.Fatalf("Accounts(\"self\") failed: expected one account. Got: %v", len(accounts))
		}
		if !rID.MatchString(*accounts[0].Id) {
			t.Fatalf("Accounts(\"self\") failed: expected account ID to be 12 digit number.  Got: %v", *accounts[0].Id)
		}
		if *accounts[0].Id != accountID {
			t.Fatalf("Accounts(\"self\") failed: expected account ID to be %v.  Got: %v", accountID, *accounts[0].Id)
		}
		if *accounts[0].Name != accountName {
			t.Fatalf("Accounts(\"self\") failed: expected %v.  Got: %v", accountName, *accounts[0].Name)
		}
	})

	t.Run("s3 accounts list", func(t *testing.T) {
		uri := "s3://test_data/accounts.json"
		options := Options{
			AccountsInfo:  uri,
			MgmtAccountID: "123456789012",
		}
		_, err := mockSvc.AccountsList(options)
		if err != nil {
			t.Fatalf("Accounts(\"s3://\") failed: %v", err)
		}
	})

	t.Run("accountsInfo accounts list", func(t *testing.T) {
		accountsInfo := "111111111111,222222222222,333333333333"
		options := Options{
			AccountsInfo: accountsInfo,
		}
		accounts, err := mockSvc.AccountsList(options)
		if err != nil {
			t.Fatalf("Accounts() failed: %v", err)
		}
		if len(accounts) < 1 {
			t.Fatalf("Accounts(\"%v\") failed: expected at least one account", accountsInfo)
		}
		if len(accounts) != 3 {
			t.Fatalf("Accounts(\"%v\") failed: expected 3 accounts", accountsInfo)
		}
		if !rID.MatchString(*accounts[0].Id) {
			t.Fatalf("Accounts(\"%v\") failed: expected first account ID to be 12 digit number.  Got: %v", accountsInfo, *accounts[0].Id)
		}
		//Need better way to check name/alias result
		//accountName := "grace-" + appenv + "-management"
		//if *accounts[0].Name != accountName {
		//	t.Fatalf("Accounts(\"%v\") failed: expected %v.  Got: %v", accountsInfo, accountName, *accounts[0].Name)
		//}
	})
}

func TestQueryAccounts(t *testing.T) {
	sess := newStubSession(t)
	// test case table
	tt := map[string]struct {
		opt         Options
		expectedErr string
		expected    []*organizations.Account
	}{
		"stub AccountsSvc nil Options": {},
		"stub AccountsSvc master account set": {
			opt: Options{
				MasterAccountID: "test_master",
				MgmtAccountID:   "test_mgmt",
			},
		},
		"stub AccountsSvc OrgUnits set": {
			opt: Options{
				OrgUnits: []string{"test_ou"},
			},
		},
	}
	// loop through test cases
	for name, tc := range tt {
		tc := tc
		t.Run(name, func(t *testing.T) {
			svc, err := NewAccountsSvc(sess)
			if err != nil {
				t.Fatal(err)
			}
			svc.stsSvc = mockStsSvc{}
			actual, err := svc.queryAccounts(tc.opt)
			if tc.expectedErr == "" {
				assert.NilError(t, err)
			} else {
				assert.Error(t, err, tc.expectedErr)
			}
			assert.DeepEqual(t, tc.expected, actual)
		})
	}
}
