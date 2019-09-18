package accounts

import (
	"os"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	awstest "github.com/gruntwork-io/terratest/modules/aws"
)

const defaultRegion = "us-east-1"

var rID = regexp.MustCompile(`^\d{12}$`)

func TestAccounts(t *testing.T) {
	masterRoleName := os.Getenv("TF_VAR_master_role_name")
	masterAccountID := os.Getenv("TF_VAR_master_account_id")
	arn := "arn:aws:iam::" + masterAccountID + ":role/" + masterRoleName
	var (
		sess *session.Session
		err  error
	)
	if len(masterRoleName) > 0 && rID.MatchString(masterAccountID) {
		sess, err = awstest.NewAuthenticatedSessionFromRole(defaultRegion, arn)
	} else {
		t.Skip("Skipping because master account role not set")
	}
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	options := Options{
		MasterAccountID: masterAccountID,
		MasterRoleName:  masterRoleName,
	}

	accounts, err := Accounts(sess, options)
	if err != nil {
		t.Fatalf("Accounts() failed: %v", err)
	}
	if len(accounts) < 1 {
		t.Fatal("expected at least one account")
	}
	if !rID.MatchString(*accounts[0].Id) {
		t.Fatalf("expected first account ID to be 12 digit number.  Got: %v", *accounts[0].Id)
	}
}

func TestAccountsInvalid(t *testing.T) {
	var (
		sess *session.Session
		err  error
	)
	//Use default session instead of orgRole
	sess, err = awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	options := Options{
		AccountsInfo: "invalid",
	}
	_, err = Accounts(sess, options)
	if err == nil {
		t.Fatalf("expected failure for invalid accounts_info")
	} else if err.Error() != "invalid accounts_info" {
		t.Fatalf("expected 'invalid accounts_info' error.  Got: %v", err)
	}
}

func TestAccountsSelf(t *testing.T) {
	appenv := os.Getenv("appenv")
	if appenv == "" {
		t.Skip("skipping if appenv not set")
	}
	var (
		sess *session.Session
		err  error
	)
	//Use default session instead of orgRole
	sess, err = awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	options := Options{
		AccountsInfo:  "self",
		MgmtAccountID: awstest.GetAccountId(t),
	}
	accounts, err := Accounts(sess, options)
	if err != nil {
		t.Fatalf("Accounts() failed: %v", err)
	}
	if len(accounts) != 1 {
		t.Fatalf("Accounts(\"self\") failed: expected one account. Got: %v", len(accounts))
	}
	if !rID.MatchString(*accounts[0].Id) {
		t.Fatalf("Accounts(\"self\") failed: expected account ID to be 12 digit number.  Got: %v", *accounts[0].Id)
	}
	if *accounts[0].Id != awstest.GetAccountId(t) {
		t.Fatalf("Accounts(\"self\") failed: expected account ID to be %v.  Got: %v", awstest.GetAccountId(t), *accounts[0].Id)
	}
	accountName := "grace-" + appenv + "-management"
	if *accounts[0].Name != accountName {
		t.Fatalf("Accounts(\"self\") failed: expected %v.  Got: %v", accountName, *accounts[0].Name)
	}
}

func TestAccountsS3(t *testing.T) {
	t.Skip("skipping test of accounts_info is s3 URI")
	appenv := os.Getenv("appenv")
	var (
		sess *session.Session
		err  error
	)
	//Use default session instead of orgRole
	sess, err = awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	uri := "s3://grace-" + appenv + "-inventory/accounts_" + appenv + ".json"
	options := Options{
		AccountsInfo:  uri,
		MgmtAccountID: awstest.GetAccountId(t),
	}
	_, err = Accounts(sess, options)
	if err != nil {
		t.Fatalf("Accounts(\"s3://\") failed: %v", err)
	}
}

func TestAccountsList(t *testing.T) {
	accountsInfo := os.Getenv("TF_VAR_accounts_info")

	if !rIDList.MatchString(accountsInfo) {
		t.Skip("skipping Accounts() with list since accounts_info is not a list of account IDs.")
	}
	//accountsInfo := "357295571838,650758800860,408627306697"
	var (
		sess *session.Session
		err  error
	)
	//Use default session instead of orgRole
	sess, err = awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	options := Options{
		AccountsInfo: accountsInfo,
	}
	accounts, err := Accounts(sess, options)
	if err != nil {
		t.Fatalf("Accounts() failed: %v", err)
	}
	t.Logf("Accounts returned: %v", accounts)
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
}
