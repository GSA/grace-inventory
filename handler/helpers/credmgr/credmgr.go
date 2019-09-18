package credmgr

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/organizations"
)

// CredMgr ... stores credentials for each *organizations.Account provided to New
type CredMgr struct {
	creds map[string]*credentials.Credentials
}

// New ... returns a *CredMgr after creating new *credential.Credential for all *organization.Account provided
func New(cfg *session.Session, mgmtAccount string, tenantRoleName string, accounts []*organizations.Account) *CredMgr {
	c := &CredMgr{creds: make(map[string]*credentials.Credentials)}
	// prevent nil pointer crash, if session is nil
	if cfg == nil {
		return c
	}
	for _, a := range accounts {
		if aws.StringValue(a.Id) == mgmtAccount {
			c.creds[aws.StringValue(a.Name)] = cfg.Config.Credentials
		} else {
			arn := "arn:aws:iam::" + aws.StringValue(a.Id) + ":role/" + tenantRoleName
			c.creds[aws.StringValue(a.Name)] = stscreds.NewCredentials(cfg, arn)
		}
	}
	return c
}

// Cred ... returns the *credential.Credential for the account name provided, if found
func (mgr *CredMgr) Cred(account string) (*credentials.Credentials, error) {
	if val, ok := mgr.creds[account]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("could not find a credential for %s", account)
}
