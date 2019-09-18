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
	"github.com/aws/aws-sdk-go/service/organizations"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
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

// Accounts ... performs Queries or parses accounts and returns all organization accounts
func Accounts(cfg client.ConfigProvider, opt Options) ([]*organizations.Account, error) {
	if cfg == nil {
		return nil, errors.New("nil ConfigProvider")
	}
	switch str := opt.AccountsInfo; {
	case str == "":
		return queryAccounts(cfg, opt)
	case str == "self":
		return selfAccountInfo(cfg, opt)
	case strings.HasPrefix(strings.ToLower(str), "s3://"):
		return parseAccountsFromJSON(opt.AccountsInfo, cfg)
	case rIDList.MatchString(str):
		return getAccountAliases(cfg, opt)
	default:
		return nil, errors.New("invalid accounts_info")
	}
}

// queryAccounts ... selects between ListAccounts and ListAccountsForParent
func queryAccounts(cfg client.ConfigProvider, opt Options) ([]*organizations.Account, error) {
	var svc *organizations.Organizations
	if opt.MasterAccountID != "" && opt.MasterAccountID != opt.MgmtAccountID {
		arn := "arn:aws:iam::" + opt.MasterAccountID + ":role/" + opt.MasterRoleName
		cred := stscreds.NewCredentials(cfg, arn)
		svc = organizations.New(cfg, &aws.Config{Credentials: cred})
	} else {
		svc = organizations.New(cfg)
	}
	if len(opt.OrgUnits) > 0 {
		return listAccountsForParents(svc, opt.OrgUnits)
	}
	return listAccountsForMaster(svc)
}

// listAccounts ... performs ListAccounts and returns all organization accounts
func listAccountsForMaster(svc *organizations.Organizations) ([]*organizations.Account, error) {
	input := &organizations.ListAccountsInput{}
	result, err := svc.ListAccounts(input)
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
		result, err := svc.ListAccounts(input)
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
func selfAccountInfo(cfg client.ConfigProvider, opt Options) ([]*organizations.Account, error) {
	opt.AccountsInfo = opt.MgmtAccountID
	return getAccountAliases(cfg, opt)
}

// parseAccountsFromJSON ... parses account info from json S3 object
func parseAccountsFromJSON(accountsInfo string, cfg client.ConfigProvider) ([]*organizations.Account, error) {
	if cfg == nil {
		return nil, errors.New("nil ConfigProvider")
	}

	u, err := url.Parse(accountsInfo)
	if err != nil {
		return nil, err
	}
	bucket := u.Host
	key := u.Path

	downloader := s3manager.NewDownloader(cfg)
	buff := &aws.WriteAtBuffer{}
	//  Download the item from the bucket. If an error occurs, log it and exit. Otherwise, notify the user that the download succeeded.
	_, err = downloader.Download(buff,
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

func getAccountAliases(cfg client.ConfigProvider, opt Options) ([]*organizations.Account, error) {
	accountIDs := strings.Split(opt.AccountsInfo, ",")
	var accounts []*organizations.Account
	var svc *iam.IAM
	for _, acct := range accountIDs {
		var account organizations.Account
		account.Id = aws.String(acct)
		if acct == opt.MgmtAccountID {
			svc = iam.New(cfg)
		} else {
			arn := "arn:aws:iam::" + acct + ":role/" + opt.TenantRoleName
			fmt.Printf("ARN: %v", arn)
			cred := stscreds.NewCredentials(cfg, arn)
			svc = iam.New(cfg, &aws.Config{Credentials: cred})
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
func listAccountsForParents(svc *organizations.Organizations, orgUnits []string) ([]*organizations.Account, error) {
	var accounts []*organizations.Account
	for _, ou := range orgUnits {
		input := &organizations.ListAccountsForParentInput{
			ParentId: aws.String(ou),
		}
		result, err := svc.ListAccountsForParent(input)
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
			result, err := svc.ListAccountsForParent(input)
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
