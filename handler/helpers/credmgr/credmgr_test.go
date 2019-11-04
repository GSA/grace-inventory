package credmgr

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/organizations"
)

var defaultRegion = "us-east-1"

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

// New(sess *session.Session, mgmtAccount string, accounts []*organizations.Account) *CredMgr
func TestNew(t *testing.T) {
	sess, err := mockNewSession(&aws.Config{Region: aws.String(defaultRegion)})
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	t.Run("nil pointer test", func(t *testing.T) {
		c := New(nil, "", "", nil)
		if c == nil {
			t.Fatal("CredMgr was nil")
		}
	})
	t.Run("accounts test", func(t *testing.T) {
		currUser := "testUser"
		currAcct := "testAccount"
		c := New(sess, "", "", []*organizations.Account{{Id: &currAcct, Name: &currUser}})
		if len(c.creds) == 0 {
			t.Fatal("accounts expected 1, got 0")
		}
		if c.creds[currUser] == nil {
			t.Fatalf("account '%s' does not exist", currUser)
		}
	})
}

// func (mgr *CredMgr) Cred(account string) (*credentials.Credentials, error)
func TestCred(t *testing.T) {
	sess, err := mockNewSession(&aws.Config{Region: aws.String(defaultRegion)})
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	currUser := "testUser"
	currAcct := "testAccount"
	c := New(sess, "", "", []*organizations.Account{{Id: &currAcct, Name: &currUser}})
	_, err = c.Cred(currUser)
	if err != nil {
		t.Fatalf("failed to get cred for user %s", currUser)
	}
	_, err = c.Cred("invalid_user")
	if err == nil {
		t.Fatalf("Cred should fail if user isn't in map")
	}
}
