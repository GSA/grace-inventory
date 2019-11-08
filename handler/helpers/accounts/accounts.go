package accounts

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/organizations"
	"github.com/aws/aws-sdk-go/service/organizations/organizationsiface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
)

type accountsList struct {
	Accounts []struct {
		Status          string
		Name            string
		Email           string
		JoinedMethod    string
		JoinedTimestamp float64
		ID              string
		Arn             string
	}
}

// Options ... Options for Accounts() function
type Options struct {
	AccountsInfo    string
	MgmtAccountID   string
	MasterAccountID string
	MasterRoleName  string
	TenantRoleName  string
	OrgUnits        []string
}

// Global compiled regular expressions
var rIDList = regexp.MustCompile(`^\d{12}(,\d{12})*$`)

// For unit testing, set the iamSvc, organizationsSvc and downloaderSvc to
// an iface mock svc client
type Svc struct {
	cfg              client.ConfigProvider
	iamSvc           iamiface.IAMAPI
	organizationsSvc organizationsiface.OrganizationsAPI
	downloaderSvc    s3manageriface.DownloaderAPI
}

func NewAccountsSvc(cfg client.ConfigProvider) (as *Svc, err error) {
	if cfg == nil {
		return nil, errors.New("nil ConfigProvider")
	}
	as.cfg = cfg
	as.downloaderSvc = s3manager.NewDownloader(cfg)
	return as, nil
}

// Accounts ... performs Queries or parses accounts and returns all organization accounts
func (as *Svc) AccountsList(opt Options) ([]*organizations.Account, error) {
	switch str := opt.AccountsInfo; {
	case str == "":
		return as.queryAccounts(opt)
	case str == "self":
		return as.selfAccountInfo(opt)
	case strings.HasPrefix(strings.ToLower(str), "s3://"):
		return as.parseAccountsFromJSON(opt.AccountsInfo)
	case rIDList.MatchString(str):
		return as.getAccountAliases(opt)
	default:
		return nil, errors.New("invalid accounts_info")
	}
}

// queryAccounts ... selects between ListAccounts and ListAccountsForParent
func (as *Svc) queryAccounts(opt Options) ([]*organizations.Account, error) {
	if as.organizationsSvc == nil {
		if opt.MasterAccountID != "" && opt.MasterAccountID != opt.MgmtAccountID {
			arn := "arn:aws:iam::" + opt.MasterAccountID + ":role/" + opt.MasterRoleName
			cred := stscreds.NewCredentials(as.cfg, arn)
			as.organizationsSvc = organizations.New(as.cfg, &aws.Config{Credentials: cred})
		} else {
			as.organizationsSvc = organizations.New(as.cfg)
		}
	}
	if len(opt.OrgUnits) > 0 {
		return as.listAccountsForParents(opt.OrgUnits)
	}
	return as.listAccountsForMaster()
}

// listAccounts ... performs ListAccounts and returns all organization accounts
func (as *Svc) listAccountsForMaster() ([]*organizations.Account, error) {
	input := &organizations.ListAccountsInput{}
	result, err := as.organizationsSvc.ListAccounts(input)
	if err != nil {
		return nil, err
	}
	accounts := result.Accounts
	token := ""
	if result.NextToken != nil {
		token = *result.NextToken
	}
	for token != "" {
		input.NextToken = &token
		result, err := as.organizationsSvc.ListAccounts(input)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, result.Accounts...)
		token = ""
		if result.NextToken != nil {
			token = *result.NextToken
		}
	}
	return accounts, nil
}

// selfAccountInfo ... returns current account ID and alias
func (as *Svc) selfAccountInfo(opt Options) ([]*organizations.Account, error) {
	opt.AccountsInfo = opt.MgmtAccountID
	return as.getAccountAliases(opt)
}

// parseAccountsFromJSON ... parses account info from json S3 object
func (as *Svc) parseAccountsFromJSON(accountsInfo string) ([]*organizations.Account, error) {
	u, err := url.Parse(accountsInfo)
	if err != nil {
		return nil, err
	}
	bucket := u.Host
	key := u.Path

	buff := &aws.WriteAtBuffer{}
	//  Download the item from the bucket. If an error occurs, log it and exit. Otherwise, notify the user that the download succeeded.
	_, err = as.downloaderSvc.Download(buff,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})

	if err != nil {
		return nil, err
	}

	var dat accountsList

	err = json.Unmarshal(buff.Bytes(), &dat)
	if err != nil {
		return nil, err
	}
	var accounts []*organizations.Account
	accts := dat.Accounts
	for _, acct := range accts {
		var account organizations.Account
		account.Status = aws.String(acct.Status)
		account.Name = aws.String(acct.Name)
		account.Email = aws.String(acct.Email)
		account.JoinedMethod = aws.String(acct.JoinedMethod)
		sec, dec := math.Modf(acct.JoinedTimestamp)
		t := time.Unix(int64(sec), int64(dec*(1e9)))
		account.JoinedTimestamp = aws.Time(t)
		account.Id = aws.String(acct.ID)
		account.Arn = aws.String(acct.Arn)
		accounts = append(accounts, &account)
	}
	return accounts, nil
}

func (as *Svc) getAccountAliases(opt Options) ([]*organizations.Account, error) {
	accountIDs := strings.Split(opt.AccountsInfo, ",")
	var accounts []*organizations.Account
	for _, acct := range accountIDs {
		svc := as.iamSvc
		var account organizations.Account
		account.Id = aws.String(acct)
		if svc == nil {
			if acct == opt.MgmtAccountID {
				svc = iam.New(as.cfg)
			} else {
				arn := "arn:aws:iam::" + acct + ":role/" + opt.TenantRoleName
				fmt.Printf("ARN: %v", arn)
				cred := stscreds.NewCredentials(as.cfg, arn)
				svc = iam.New(as.cfg, &aws.Config{Credentials: cred})
			}
		}
		result, err := svc.ListAccountAliases(&iam.ListAccountAliasesInput{})
		if err != nil {
			log.Printf("Error getting account alias for %v: %v", acct, err)
		} else {
			account.Name = result.AccountAliases[0]
		}
		accounts = append(accounts, &account)
	}
	return accounts, nil
}

// AccountsForParents ... performs ListAccountsForParent and returns all accounts
func (as *Svc) listAccountsForParents(orgUnits []string) ([]*organizations.Account, error) {
	var accounts []*organizations.Account
	for _, ou := range orgUnits {
		input := &organizations.ListAccountsForParentInput{
			ParentId: aws.String(ou),
		}
		result, err := as.organizationsSvc.ListAccountsForParent(input)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, result.Accounts...)
		token := ""
		if result.NextToken != nil {
			token = *result.NextToken
		}
		for token != "" {
			input.NextToken = &token
			result, err := as.organizationsSvc.ListAccountsForParent(input)
			if err != nil {
				return nil, err
			}
			accounts = append(accounts, result.Accounts...)
			token = ""
			if result.NextToken != nil {
				token = *result.NextToken
			}
		}
	}
	return accounts, nil
}
