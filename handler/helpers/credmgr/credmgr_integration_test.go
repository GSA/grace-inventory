// +build integration

package credmgr

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/organizations"
	awstest "github.com/gruntwork-io/terratest/modules/aws"
)

// New(sess *session.Session, mgmtAccount string, accounts []*organizations.Account) *CredMgr
func TestIntegrationNew(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
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
		currUser := awstest.GetIamCurrentUserName(t)
		currAcct := awstest.GetAccountId(t)
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
func TestIntegrationCred(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	currUser := awstest.GetIamCurrentUserName(t)
	currAcct := awstest.GetAccountId(t)
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
