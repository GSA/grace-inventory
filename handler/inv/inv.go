package inv

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"runtime"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/configservice"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/glacier"
	"github.com/aws/aws-sdk-go/service/glacier/glacieriface"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/organizations"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/caarlos0/env"

	"github.com/GSA/grace-inventory/handler/helpers"
	"github.com/GSA/grace-inventory/handler/helpers/accounts"
	"github.com/GSA/grace-inventory/handler/helpers/credmgr"
	"github.com/GSA/grace-inventory/handler/helpers/sessionmgr"
	"github.com/GSA/grace-inventory/handler/spreadsheet"
)

// config ... struct for holding environment variables
type config struct {
	BucketID        string   `env:"s3_bucket,required"`
	KmsKeyID        string   `env:"kms_key_id,required"`
	Regions         []string `env:"regions,required" envSeparator:","`
	AccountsInfo    string   `env:"accounts_info" envDefault:"self"`
	MasterAccountID string   `env:"master_account_id" envDefault:""`
	OrgUnits        []string `env:"organizational_units" envSeparator:","`
	MasterRoleName  string   `env:"master_role_name" envDefault:""`
	TenantRoleName  string   `env:"tenant_role_name" envDefault:""`
}

type queryFunc func() ([]*spreadsheet.Payload, error)

type queryError struct {
	M string
	E error
}

func (q queryError) Error() string {
	return q.M
}
func newQueryErrorf(err error, format string, params ...interface{}) queryError {
	return queryError{E: err, M: fmt.Sprintf(format, params...)}
}

type done struct {
	Name string
}

func getCallerFunc() string {
	pc := make([]uintptr, 1)
	if runtime.Callers(3, pc) == 0 {
		return "nil"
	}
	frame, _ := runtime.CallersFrames([]uintptr{pc[0]}).Next()
	return frame.Function
}

var knownErrors = map[string]interface{}{
	"AccessDenied":          nil,
	"AccessDeniedException": nil,
	"AuthorizationError":    nil,
	"UnauthorizedOperation": nil,
}

func isKnownError(err error) bool {
	if queryErr, ok := err.(queryError); ok {
		if awsErr, ok := queryErr.E.(awserr.Error); ok {
			_, ok := knownErrors[awsErr.Code()]
			return ok
		}
	}
	return false
}

func logDuration() func() {
	caller := getCallerFunc()
	start := time.Now()
	log.Printf("calling %s\n", caller)
	return func() {
		log.Printf("%s took %s\n", caller, time.Since(start))
	}
}

// Inv ... is used to manage the spreadsheet and sessions required to generate the AWS report
type Inv struct {
	spreadsheet     *spreadsheet.Spreadsheet
	mgmtAccount     string
	bucketID        string
	kmsKeyID        string
	defaultRegion   string
	regions         []string
	accountsInfo    string
	masterAccountID string
	orgUnits        []string
	masterRoleName  string
	tenantRoleName  string
	sessionMgr      *sessionmgr.SessionMgr
	credMgr         *credmgr.CredMgr
	accounts        []*organizations.Account
	out             chan interface{}
	errc            chan error
	queries         map[string]queryFunc
	running         []string
}

// New ... returns an *Inv, after storing all known queryFunc and creating the *SessionMgr
func New() (*Inv, error) {
	cfg := config{}
	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}
	defaultRegion := cfg.Regions[0]
	inv := &Inv{
		bucketID:        cfg.BucketID,
		kmsKeyID:        cfg.KmsKeyID,
		defaultRegion:   defaultRegion,
		regions:         cfg.Regions,
		accountsInfo:    cfg.AccountsInfo,
		masterAccountID: cfg.MasterAccountID,
		orgUnits:        cfg.OrgUnits,
		masterRoleName:  cfg.MasterRoleName,
		tenantRoleName:  cfg.TenantRoleName,
		out:             make(chan interface{}),
		errc:            make(chan error),
	}
	//store available queries for referencing
	inv.queries = map[string]queryFunc{
		SheetRoles:          inv.queryRoles,
		SheetGroups:         inv.queryGroups,
		SheetPolicies:       inv.queryPolicies,
		SheetUsers:          inv.queryUsers,
		SheetBuckets:        inv.queryBuckets,
		SheetInstances:      inv.queryInstances,
		SheetImages:         inv.queryImages,
		SheetVolumes:        inv.queryVolumes,
		SheetSnapshots:      inv.querySnapshots,
		SheetVpcs:           inv.queryVpcs,
		SheetSubnets:        inv.querySubnets,
		SheetSecurityGroups: inv.querySecurityGroups,
		SheetAddresses:      inv.queryAddresses,
		SheetKeyPairs:       inv.queryKeyPairs,
		SheetStacks:         inv.queryStacks,
		SheetAlarms:         inv.queryAlarms,
		SheetConfigRules:    inv.queryConfigRules,
		SheetLoadBalancers:  inv.queryLoadBalancers,
		SheetVaults:         inv.queryVaults,
		SheetKeys:           inv.queryKeys,
		SheetDBInstances:    inv.queryDBInstances,
		SheetDBSnapshots:    inv.queryDBSnapshots,
		SheetSecrets:        inv.querySecrets,
		SheetSubscriptions:  inv.querySubscriptions,
		SheetTopics:         inv.queryTopics,
		SheetParameters:     inv.queryParameters,
	}

	sess, err := session.NewSession(&aws.Config{Region: &defaultRegion})
	if err != nil {
		return nil, err
	}
	svc := &stsSvc{Client: sts.New(sess)}
	identity, err := svc.getCurrentIdentity()
	if err != nil {
		return nil, err
	}
	// Set mgmtAccount to the current account
	inv.mgmtAccount = *identity.Account
	inv.sessionMgr = sessionmgr.New(defaultRegion, cfg.Regions)
	err = inv.sessionMgr.Init()
	if err != nil {
		return nil, err
	}
	return inv, nil
}

// Run ... starts the report process, the corresponding queryFunc for each sheet in the spreadsheet
// will be ran and the results added to that sheet. Run is a blocking function and will hold the cursor
// until all queries have been ran and the spreadsheet has been saved to the bucket
func (inv *Inv) Run(s *spreadsheet.Spreadsheet) error {
	inv.spreadsheet = s
	inv.query(map[string]queryFunc{SheetAccounts: inv.queryAccounts})

	err := inv.aggregate()
	if err != nil {
		return err
	}
	return inv.save()
}

// query ... enumerates over the funcs map provided, spawning each func in a new go routine
// then sending the results over the out channel to be collected by 'aggregate'. As each func
// is called appends the name of the sheet to 'running' which is used to determine whether
// all sheets have been completed successfully
func (inv *Inv) query(funcs map[string]queryFunc) {
	for name, fn := range funcs {
		inv.running = append(inv.running, name)
		go func(fn queryFunc, name string, out chan interface{}, errc chan error) {
			payloads, err := fn()
			if err != nil {
				errc <- err
			}
			for _, p := range payloads {
				out <- p
			}
			out <- &done{name}
		}(fn, name, inv.out, inv.errc)
	}
}

// runAllQueries ... executes remaining queries, excluding Accounts
func (inv *Inv) runAllQueries() {
	queries := make(map[string]queryFunc)

	for _, v := range inv.spreadsheet.Sheets {
		if fn, ok := inv.queries[v.Name]; ok {
			queries[v.Name] = fn
		}
	}

	inv.query(queries)
}

// nolint: gocyclo
// aggregate ... waits for results to be sent on the 'out' channel, then calls 'UpdateSheet'
// passing the corresponding sheet name for the 'spreadsheet.Payload.Items' type. As sheets
// are completed, removes the sheet name from 'running' to prevent infinitely looping
func (inv *Inv) aggregate() error {
	// while there are incomplete sheets, loop and wait for completion
	for len(inv.running) > 0 {
		select {
		case obj := <-inv.out:
			switch val := obj.(type) {
			case *spreadsheet.Payload:
				sheet, err := typeToSheet(val.Items)
				if err != nil {
					return err
				}
				if sheet == "" {
					// if the sheet name is empty, the payload is empty
					// stop processing further, we'll wait for the next one
					break
				}
				if sheet == SheetAccounts {
					// Use accounts to facilitate the creation of the credMgr
					sess, err := inv.sessionMgr.Default()
					if err != nil {
						return err
					}
					inv.credMgr = credmgr.New(sess, inv.mgmtAccount, inv.tenantRoleName, inv.accounts)

					inv.runAllQueries()
				}
				inv.spreadsheet.UpdateSheet(sheet, val)
			case *done:
				// Once a sheet is complete, remove it from the slice
				for i, v := range inv.running {
					if val.Name == v {
						log.Printf("Sheet %q has completed\n", val.Name)
						inv.running = append(inv.running[:i], inv.running[i+1:]...)
					}
				}
			}
		// if any errors occur, return and break the loop
		case err := <-inv.errc:
			return err
		}
	}
	return nil
}

type stsSvc struct {
	Client stsiface.STSAPI
}

// getCurrentIdentity ... returns the response from GetCallerIdentity
func (svc *stsSvc) getCurrentIdentity() (*sts.GetCallerIdentityOutput, error) {
	return svc.Client.GetCallerIdentity(&sts.GetCallerIdentityInput{})
}

// nolint: gocyclo
// typeToSheet ... converts a slice type to a sheet name
func typeToSheet(items interface{}) (string, error) {
	var sheet string

	s := reflect.ValueOf(items)
	if s.Kind() != reflect.Slice {
		return "", errors.New("items is not a sheet")
	}

	if s.Len() == 0 {
		//Empty slice - this isn't an error, but we don't need to do anything
		return "", nil
	}
	switch val := s.Index(0).Interface().(type) {
	case *organizations.Account:
		sheet = SheetAccounts
	case *iam.Role:
		sheet = SheetRoles
	case *iam.Group:
		sheet = SheetGroups
	case *iam.Policy:
		sheet = SheetPolicies
	case *iam.User:
		sheet = SheetUsers
	case *s3.Bucket:
		sheet = SheetBuckets
	case *ec2.Instance:
		sheet = SheetInstances
	case *ec2.Image:
		sheet = SheetImages
	case *ec2.Volume:
		sheet = SheetVolumes
	case *ec2.Snapshot:
		sheet = SheetSnapshots
	case *ec2.Vpc:
		sheet = SheetVpcs
	case *ec2.Subnet:
		sheet = SheetSubnets
	case *ec2.SecurityGroup:
		sheet = SheetSecurityGroups
	case *ec2.Address:
		sheet = SheetAddresses
	case *ec2.KeyPairInfo:
		sheet = SheetKeyPairs
	case *cloudformation.Stack:
		sheet = SheetStacks
	case *cloudwatch.MetricAlarm:
		sheet = SheetAlarms
	case *configservice.ConfigRule:
		sheet = SheetConfigRules
	case *glacier.DescribeVaultOutput:
		sheet = SheetVaults
	case *helpers.KmsKey:
		sheet = SheetKeys
	case *rds.DBInstance:
		sheet = SheetDBInstances
	case *rds.DBSnapshot:
		sheet = SheetDBSnapshots
	case *secretsmanager.SecretListEntry:
		sheet = SheetSecrets
	case *sns.Subscription:
		sheet = SheetSubscriptions
	case *helpers.SnsTopic:
		sheet = SheetTopics
	case *ssm.ParameterMetadata:
		sheet = SheetParameters
	default:
		log.Printf("Unknown sheet type: %T", val)
		return "", errors.New("unknown type")
	}
	return sheet, nil
}

// walkAccounts ... loops over all organization accounts, skipping suspended accounts, and calling 'fn'
// passing the *credential.Credential for each account, using the default session, collecting all returned payloads
func (inv *Inv) walkAccounts(fn func(string, *credentials.Credentials, *session.Session) (*spreadsheet.Payload, error)) ([]*spreadsheet.Payload, error) {
	var payloads []*spreadsheet.Payload
	for _, a := range inv.accounts {
		if aws.StringValue(a.Status) == "SUSPENDED" {
			continue
		}
		cred, err := inv.credMgr.Cred(aws.StringValue(a.Name))
		if err != nil {
			return nil, err
		}
		sess, err := inv.sessionMgr.Default()
		if err != nil {
			return nil, err
		}
		payload, err := fn(aws.StringValue(a.Name), cred, sess)
		if err != nil {
			if isKnownError(err) {
				log.Printf("walkAccounts got an error when called by %s -> %v\n", getCallerFunc(), err)
				continue
			}
			return nil, err
		}
		payloads = append(payloads, payload)
	}
	return payloads, nil
}

// walkSessions ... loops over all organization accounts, skipping suspended accounts,
// then looping over all sessions in the SessionMgr calling 'fn', collecting all returned payloads
func (inv *Inv) walkSessions(fn func(string, *credentials.Credentials, *session.Session) (*spreadsheet.Payload, error)) ([]*spreadsheet.Payload, error) {
	var payloads []*spreadsheet.Payload
	for _, a := range inv.accounts {
		if aws.StringValue(a.Status) == "SUSPENDED" {
			continue
		}
		cred, err := inv.credMgr.Cred(aws.StringValue(a.Name))
		if err != nil {
			return nil, err
		}
		for _, s := range inv.sessionMgr.All() {
			payload, err := fn(aws.StringValue(a.Name), cred, s)
			if err != nil {
				if isKnownError(err) {
					log.Printf("walkSessions got an error when called by %s -> %v\n", getCallerFunc(), err)
					continue
				}
				return nil, err
			}
			payloads = append(payloads, payload)
		}
	}
	return payloads, nil
}

// save - saves the report to S3 with the filename provided to New
func (inv *Inv) save() error {
	sess, err := inv.sessionMgr.Default()
	if err != nil {
		return err
	}
	reader, err := inv.spreadsheet.Bytes()
	if err != nil {
		return err
	}

	svc := s3.New(sess)
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(inv.bucketID),
		ContentType:          aws.String("application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"),
		Key:                  aws.String(inv.spreadsheet.Name),
		Body:                 reader,
		SSEKMSKeyId:          aws.String(inv.kmsKeyID),
		ServerSideEncryption: aws.String("aws:kms"),
	})
	if err != nil {
		return fmt.Errorf("failed to upload report to bucket: %v", err)
	}
	return nil
}

// queryAccounts ... Queries organization accounts, pushes them onto a slice of interface,
// then returns a slice of *spreadsheet.Payload
func (inv *Inv) queryAccounts() ([]*spreadsheet.Payload, error) {
	defer logDuration()()
	sess, err := inv.sessionMgr.Default()
	if err != nil {
		return nil, newQueryErrorf(err, "failed to get default session from sessionMgr: %v", err)
	}
	options := accounts.Options{
		AccountsInfo:    inv.accountsInfo,
		MgmtAccountID:   inv.mgmtAccount,
		MasterAccountID: inv.masterAccountID,
		MasterRoleName:  inv.masterRoleName,
		TenantRoleName:  inv.tenantRoleName,
		OrgUnits:        inv.orgUnits,
	}
	accounts, err := accounts.Accounts(sess, options)
	if err != nil {
		return nil, newQueryErrorf(err, "failed to get Accounts: %v", err)
	}
	var items []interface{}
	for i, a := range accounts {
		// Use Account ID if name/alias is not set
		if aws.StringValue(a.Name) == "" {
			accounts[i].Name = a.Id
		}
		items = append(items, a)
	}
	inv.accounts = accounts
	return []*spreadsheet.Payload{
		{Static: nil, Items: items},
	}, nil
}

// queryRoles ... queries IAM Roles for all organization accounts
// pushes them onto a slice of interface, then returns a slice of *spreadsheet.Payload
func (inv *Inv) queryRoles() ([]*spreadsheet.Payload, error) {
	defer logDuration()()
	return inv.walkAccounts(func(account string, cred *credentials.Credentials, sess *session.Session) (*spreadsheet.Payload, error) {
		svc := helpers.IamSvc{
			Client: iam.New(sess, &aws.Config{Credentials: cred}),
		}
		roles, err := svc.Roles()
		if err != nil {
			return nil, newQueryErrorf(err, "failed to get Roles for account: %s -> %v", account, err)
		}
		var items []interface{}
		for _, r := range roles {
			items = append(items, r)
		}
		return &spreadsheet.Payload{Static: []string{account}, Items: items}, nil
	})
}

// queryGroups ... queries IAM Groups for all organization accounts
// pushes them onto a slice of interface, then returns a slice of *spreadsheet.Payload
func (inv *Inv) queryGroups() ([]*spreadsheet.Payload, error) {
	defer logDuration()()
	return inv.walkAccounts(func(account string, cred *credentials.Credentials, sess *session.Session) (*spreadsheet.Payload, error) {
		svc := helpers.IamSvc{
			Client: iam.New(sess, &aws.Config{Credentials: cred}),
		}
		groups, err := svc.Groups()
		if err != nil {
			return nil, newQueryErrorf(err, "failed to get Groups for account: %s -> %v", account, err)
		}
		var items []interface{}
		for _, g := range groups {
			items = append(items, g)
		}
		return &spreadsheet.Payload{Static: []string{account}, Items: items}, nil
	})
}

// queryPolicies ... queries IAM Groups for all organization accounts
// pushes them onto a slice of interface, then returns a slice of *spreadsheet.Payload
func (inv *Inv) queryPolicies() ([]*spreadsheet.Payload, error) {
	defer logDuration()()
	return inv.walkAccounts(func(account string, cred *credentials.Credentials, sess *session.Session) (*spreadsheet.Payload, error) {
		svc := helpers.IamSvc{
			Client: iam.New(sess, &aws.Config{Credentials: cred}),
		}
		policies, err := svc.Policies()
		if err != nil {
			return nil, newQueryErrorf(err, "failed to get Policies for account: %s -> %v", account, err)
		}
		var items []interface{}
		for _, p := range policies {
			items = append(items, p)
		}
		return &spreadsheet.Payload{Static: []string{account}, Items: items}, nil
	})
}

// queryUsers ... queries IAM users for all organization accounts
// pushes them onto a slice of interface, then returns a slice of *spreadsheet.Payload
func (inv *Inv) queryUsers() ([]*spreadsheet.Payload, error) {
	defer logDuration()()
	return inv.walkAccounts(func(account string, cred *credentials.Credentials, sess *session.Session) (*spreadsheet.Payload, error) {
		svc := helpers.IamSvc{
			Client: iam.New(sess, &aws.Config{Credentials: cred}),
		}
		users, err := svc.Users()
		if err != nil {
			return nil, newQueryErrorf(err, "failed to get Users for account: %s -> %v", account, err)
		}
		var items []interface{}
		for _, u := range users {
			items = append(items, u)
		}
		return &spreadsheet.Payload{Static: []string{account}, Items: items}, nil
	})
}

// queryBuckets ... queries S3 buckets for all organization accounts
// pushes them onto a slice of interface, then returns a slice of *spreadsheet.Payload
func (inv *Inv) queryBuckets() ([]*spreadsheet.Payload, error) {
	defer logDuration()()
	return inv.walkAccounts(func(account string, cred *credentials.Credentials, sess *session.Session) (*spreadsheet.Payload, error) {
		svc := s3.New(sess, &aws.Config{Credentials: cred})
		buckets, err := helpers.Buckets(svc)
		if err != nil {
			return nil, newQueryErrorf(err, "failed to get Buckets for account: %s -> %v", account, err)
		}
		var items []interface{}
		for _, b := range buckets {
			items = append(items, b)
		}
		return &spreadsheet.Payload{Static: []string{account}, Items: items}, nil
	})
}

// queryInstances ... queries EC2 instances for all organization accounts and
// all sessions/regions in SessionMgr, pushes them onto a slice of interface
// then returns a slice of *spreadsheet.Payload
func (inv *Inv) queryInstances() ([]*spreadsheet.Payload, error) {
	defer logDuration()()
	return inv.walkAccounts(func(account string, cred *credentials.Credentials, sess *session.Session) (*spreadsheet.Payload, error) {
		svc := helpers.Ec2Svc{
			Client: ec2.New(sess, &aws.Config{Credentials: cred}),
		}
		instances, err := svc.Instances()
		if err != nil {
			return nil, newQueryErrorf(err, "failed to get Instances for account: %s, region: %s -> %v", account, *sess.Config.Region, err)
		}
		var items []interface{}
		for _, i := range instances {
			items = append(items, i)
		}
		return &spreadsheet.Payload{Static: []string{account, *sess.Config.Region}, Items: items}, nil
	})
}

// queryImages ... queries Amazon machine images (AMI) for all organization
// accounts and all sessions/regions in SessionMgr, pushes them onto a slice of
// interface then returns a slice of *spreadsheet.Payload
func (inv *Inv) queryImages() ([]*spreadsheet.Payload, error) {
	defer logDuration()()
	return inv.walkSessions(func(account string, cred *credentials.Credentials, sess *session.Session) (*spreadsheet.Payload, error) {
		svc := helpers.Ec2Svc{
			Client: ec2.New(sess, &aws.Config{Credentials: cred}),
		}
		images, err := svc.Images()
		if err != nil {
			return nil, newQueryErrorf(err, "failed to get Images for account: %s, region: %s -> %v", account, *sess.Config.Region, err)
		}
		var items []interface{}
		for _, i := range images {
			items = append(items, i)
		}
		return &spreadsheet.Payload{Static: []string{account, *sess.Config.Region}, Items: items}, nil
	})
}

// queryVolumes ... queries Elastic Block Storage (EBS) volumes for all
// organization accounts and all sessions/regions in SessionMgr, pushes them
// onto a slice of interface then returns a slice of *spreadsheet.Payload
func (inv *Inv) queryVolumes() ([]*spreadsheet.Payload, error) {
	defer logDuration()()
	return inv.walkSessions(func(account string, cred *credentials.Credentials, sess *session.Session) (*spreadsheet.Payload, error) {
		svc := helpers.Ec2Svc{
			Client: ec2.New(sess, &aws.Config{Credentials: cred}),
		}
		volumes, err := svc.Volumes()
		if err != nil {
			return nil, newQueryErrorf(err, "failed to get Volumes for account: %s, region: %s -> %v", account, *sess.Config.Region, err)
		}
		var items []interface{}
		for _, v := range volumes {
			items = append(items, v)
		}
		return &spreadsheet.Payload{Static: []string{account, *sess.Config.Region}, Items: items}, nil
	})
}

// querySnapshots ... queries EBS snapshots for all organization accounts and
// all sessions/regions in SessionMgr, pushes them onto a slice of interface
// then returns a slice of *spreadsheet.Payload
func (inv *Inv) querySnapshots() ([]*spreadsheet.Payload, error) {
	defer logDuration()()
	return inv.walkSessions(func(account string, cred *credentials.Credentials, sess *session.Session) (*spreadsheet.Payload, error) {
		svc := helpers.Ec2Svc{
			Client: ec2.New(sess, &aws.Config{Credentials: cred}),
		}
		snapshots, err := svc.Snapshots()
		if err != nil {
			return nil, newQueryErrorf(err, "failed to get Snapshots for account: %s, region: %s -> %v", account, *sess.Config.Region, err)
		}
		var items []interface{}
		for _, s := range snapshots {
			items = append(items, s)
		}
		return &spreadsheet.Payload{Static: []string{account, *sess.Config.Region}, Items: items}, nil
	})
}

// queryVpcs ... queries VPCs for all organization accounts and
// all sessions/regions in SessionMgr, pushes them onto a slice of interface
// then returns a slice of *spreadsheet.Payload
func (inv *Inv) queryVpcs() ([]*spreadsheet.Payload, error) {
	defer logDuration()()
	return inv.walkSessions(func(account string, cred *credentials.Credentials, sess *session.Session) (*spreadsheet.Payload, error) {
		svc := helpers.Ec2Svc{
			Client: ec2.New(sess, &aws.Config{Credentials: cred}),
		}
		vpcs, err := svc.Vpcs()
		if err != nil {
			return nil, newQueryErrorf(err, "failed to get VPCs for account: %s, region: %s -> %v", account, *sess.Config.Region, err)
		}
		var items []interface{}
		for _, v := range vpcs {
			items = append(items, v)
		}
		return &spreadsheet.Payload{Static: []string{account, *sess.Config.Region}, Items: items}, nil
	})
}

// querySubnets ... queries subnets for all organization accounts and
// all sessions/regions in SessionMgr, pushes them onto a slice of interface
// then returns a slice of *spreadsheet.Payload
func (inv *Inv) querySubnets() ([]*spreadsheet.Payload, error) {
	defer logDuration()()
	return inv.walkSessions(func(account string, cred *credentials.Credentials, sess *session.Session) (*spreadsheet.Payload, error) {
		svc := helpers.Ec2Svc{
			Client: ec2.New(sess, &aws.Config{Credentials: cred}),
		}
		subnets, err := svc.Subnets()
		if err != nil {
			return nil, newQueryErrorf(err, "failed to get Subnets for account: %s, region: %s -> %v", account, *sess.Config.Region, err)
		}
		var items []interface{}
		for _, s := range subnets {
			items = append(items, s)
		}
		return &spreadsheet.Payload{Static: []string{account, *sess.Config.Region}, Items: items}, nil
	})
}

// querySecurityGroups ... queries security groups for all organization accounts and
// all sessions/regions in SessionMgr, pushes them onto a slice of interface
// then returns a slice of *spreadsheet.Payload
func (inv *Inv) querySecurityGroups() ([]*spreadsheet.Payload, error) {
	defer logDuration()()
	return inv.walkSessions(func(account string, cred *credentials.Credentials, sess *session.Session) (*spreadsheet.Payload, error) {
		svc := helpers.Ec2Svc{
			Client: ec2.New(sess, &aws.Config{Credentials: cred}),
		}
		groups, err := svc.SecurityGroups()
		if err != nil {
			return nil, newQueryErrorf(err, "failed to get Security Groups for account: %s, region: %s -> %v", account, *sess.Config.Region, err)
		}
		var items []interface{}
		for _, g := range groups {
			items = append(items, g)
		}
		return &spreadsheet.Payload{Static: []string{account, *sess.Config.Region}, Items: items}, nil
	})
}

// queryAddresses ... queries EC2 DescribeAddresses for all organization
// accounts and all sessions/regions in SessionMgr, pushes them onto a slice of
// interface then returns a slice of *spreadsheet.Payload
func (inv *Inv) queryAddresses() ([]*spreadsheet.Payload, error) {
	defer logDuration()()
	return inv.walkSessions(func(account string, cred *credentials.Credentials, sess *session.Session) (*spreadsheet.Payload, error) {
		svc := helpers.Ec2Svc{
			Client: ec2.New(sess, &aws.Config{Credentials: cred}),
		}
		addresses, err := svc.Addresses()
		if err != nil {
			return nil, newQueryErrorf(err, "failed to get EC2 Addresses for account: %s, region: %s -> %v", account, *sess.Config.Region, err)
		}
		var items []interface{}
		for _, g := range addresses {
			items = append(items, g)
		}
		return &spreadsheet.Payload{Static: []string{account, *sess.Config.Region}, Items: items}, nil
	})
}

// queryKeyPairs ... queries EC2 DescribeKeyPairs for all organization
// accounts and all sessions/regions in SessionMgr, pushes them onto a slice of
// interface then returns a slice of *spreadsheet.Payload
func (inv *Inv) queryKeyPairs() ([]*spreadsheet.Payload, error) {
	defer logDuration()()
	return inv.walkSessions(func(account string, cred *credentials.Credentials, sess *session.Session) (*spreadsheet.Payload, error) {
		svc := helpers.Ec2Svc{
			Client: ec2.New(sess, &aws.Config{Credentials: cred}),
		}
		keyPairs, err := svc.KeyPairs()
		if err != nil {
			return nil, newQueryErrorf(err, "failed to get EC2 KeyPairs for account: %s, region: %s -> %v", account, *sess.Config.Region, err)
		}
		var items []interface{}
		for _, g := range keyPairs {
			items = append(items, g)
		}
		return &spreadsheet.Payload{Static: []string{account, *sess.Config.Region}, Items: items}, nil
	})
}

// queryStacks ... queries CloudFormation Stacks for all organization accounts and
// all sessions/regions in SessionMgr, pushes them onto a slice of interface
// then returns a slice of *spreadsheet.Payload
func (inv *Inv) queryStacks() ([]*spreadsheet.Payload, error) {
	defer logDuration()()
	return inv.walkSessions(func(account string, cred *credentials.Credentials, sess *session.Session) (*spreadsheet.Payload, error) {
		svc := cloudformation.New(sess, &aws.Config{Credentials: cred})
		stacks, err := helpers.Stacks(svc)
		if err != nil {
			return nil, newQueryErrorf(err, "failed to get CloudFormation Stacks for account: %s, region: %s -> %v", account, *sess.Config.Region, err)
		}
		var items []interface{}
		for _, g := range stacks {
			items = append(items, g)
		}
		return &spreadsheet.Payload{Static: []string{account, *sess.Config.Region}, Items: items}, nil
	})
}

// queryAlarms ... queries CloudWatch Alarms for all organization accounts and
// all sessions/regions in SessionMgr, pushes them onto a slice of interface
// then returns a slice of *spreadsheet.Payload
func (inv *Inv) queryAlarms() ([]*spreadsheet.Payload, error) {
	defer logDuration()()
	return inv.walkSessions(func(account string, cred *credentials.Credentials, sess *session.Session) (*spreadsheet.Payload, error) {
		alarms, err := helpers.Alarms(sess, cred)
		if err != nil {
			return nil, newQueryErrorf(err, "failed to get CloudWatch Alarms for account: %s, region: %s -> %v", account, *sess.Config.Region, err)
		}
		var items []interface{}
		for _, g := range alarms {
			items = append(items, g)
		}
		return &spreadsheet.Payload{Static: []string{account, *sess.Config.Region}, Items: items}, nil
	})
}

// queryConfigRules ... queries Config Rules for all organization accounts and
// all sessions/regions in SessionMgr, pushes them onto a slice of interface
// then returns a slice of *spreadsheet.Payload
func (inv *Inv) queryConfigRules() ([]*spreadsheet.Payload, error) {
	defer logDuration()()
	return inv.walkSessions(func(account string, cred *credentials.Credentials, sess *session.Session) (*spreadsheet.Payload, error) {
		rules, err := helpers.ConfigRules(sess, cred)
		if err != nil {
			return nil, newQueryErrorf(err, "failed to get Config Rules for account: %s, region: %s -> %v", account, *sess.Config.Region, err)
		}
		var items []interface{}
		for _, g := range rules {
			items = append(items, g)
		}
		return &spreadsheet.Payload{Static: []string{account, *sess.Config.Region}, Items: items}, nil
	})
}

// queryLoadBalancers ... queries ELBv2 Load Balancers for all organization accounts and
// all sessions/regions in SessionMgr, pushes them onto a slice of interface
// then returns a slice of *spreadsheet.Payload
func (inv *Inv) queryLoadBalancers() ([]*spreadsheet.Payload, error) {
	defer logDuration()()
	return inv.walkSessions(func(account string, cred *credentials.Credentials, sess *session.Session) (*spreadsheet.Payload, error) {
		loadBalancers, err := helpers.LoadBalancers(sess, cred)
		if err != nil {
			return nil, newQueryErrorf(err, "failed to get ELBv2 Load Balancers for account: %s, region: %s -> %v", account, *sess.Config.Region, err)
		}
		var items []interface{}
		for _, g := range loadBalancers {
			items = append(items, g)
		}
		return &spreadsheet.Payload{Static: []string{account, *sess.Config.Region}, Items: items}, nil
	})
}

var glacierCreator = glacierClientCreator

func glacierClientCreator(p client.ConfigProvider, cfgs ...*aws.Config) glacieriface.GlacierAPI {
	return glacier.New(p, cfgs...)
}

// queryVaults ... queries Glacier Vaults for all organization accounts and
// all sessions/regions in SessionMgr, pushes them onto a slice of interface
// then returns a slice of *spreadsheet.Payload
func (inv *Inv) queryVaults() ([]*spreadsheet.Payload, error) {
	defer logDuration()()
	return inv.walkSessions(func(account string, cred *credentials.Credentials, sess *session.Session) (*spreadsheet.Payload, error) {
		svc := &helpers.GlacierSvc{Client: glacierCreator(sess, &aws.Config{Credentials: cred})}
		vaults, err := svc.Vaults()
		if err != nil {
			return nil, newQueryErrorf(err, "failed to get Glacier Vaults for account: %s, region: %s -> %v", account, *sess.Config.Region, err)
		}
		var items []interface{}
		for _, g := range vaults {
			items = append(items, g)
		}
		return &spreadsheet.Payload{Static: []string{account, aws.StringValue(sess.Config.Region)}, Items: items}, nil
	})
}

// queryKeys ... queries KMS Keys for all organization accounts and
// all sessions/regions in SessionMgr, pushes them onto a slice of interface
// then returns a slice of *spreadsheet.Payload
func (inv *Inv) queryKeys() ([]*spreadsheet.Payload, error) {
	defer logDuration()()
	return inv.walkSessions(func(account string, cred *credentials.Credentials, sess *session.Session) (*spreadsheet.Payload, error) {
		keys, err := helpers.Keys(sess, cred)
		if err != nil {
			return nil, newQueryErrorf(err, "failed to get KMS Keys for account: %s, region: %s -> %v", account, *sess.Config.Region, err)
		}
		var items []interface{}
		for _, g := range keys {
			items = append(items, g)
		}
		return &spreadsheet.Payload{Static: []string{account, *sess.Config.Region}, Items: items}, nil
	})
}

// queryDBInstances ... queries RDS DBInstances for all organization accounts and
// all sessions/regions in SessionMgr, pushes them onto a slice of interface
// then returns a slice of *spreadsheet.Payload
func (inv *Inv) queryDBInstances() ([]*spreadsheet.Payload, error) {
	defer logDuration()()
	return inv.walkSessions(func(account string, cred *credentials.Credentials, sess *session.Session) (*spreadsheet.Payload, error) {
		instances, err := helpers.DBInstances(sess, cred)
		if err != nil {
			return nil, newQueryErrorf(err, "failed to get RDS DBInstances for account: %s, region: %s -> %v", account, *sess.Config.Region, err)
		}
		var items []interface{}
		for _, g := range instances {
			items = append(items, g)
		}
		return &spreadsheet.Payload{Static: []string{account, *sess.Config.Region}, Items: items}, nil
	})
}

// queryDBSnapshots ... queries RDS DBSnapshots for all organization accounts and
// all sessions/regions in SessionMgr, pushes them onto a slice of interface
// then returns a slice of *spreadsheet.Payload
func (inv *Inv) queryDBSnapshots() ([]*spreadsheet.Payload, error) {
	defer logDuration()()
	return inv.walkSessions(func(account string, cred *credentials.Credentials, sess *session.Session) (*spreadsheet.Payload, error) {
		snapshots, err := helpers.DBSnapshots(sess, cred)
		if err != nil {
			return nil, newQueryErrorf(err, "failed to get RDS DBSnapshots for account: %s, region: %s -> %v", account, *sess.Config.Region, err)
		}
		var items []interface{}
		for _, g := range snapshots {
			items = append(items, g)
		}
		return &spreadsheet.Payload{Static: []string{account, *sess.Config.Region}, Items: items}, nil
	})
}

// querySecrets ... queries SecretsManager Secrets for all organization accounts and
// all sessions/regions in SessionMgr, pushes them onto a slice of interface
// then returns a slice of *spreadsheet.Payload
func (inv *Inv) querySecrets() ([]*spreadsheet.Payload, error) {
	defer logDuration()()
	return inv.walkSessions(func(account string, cred *credentials.Credentials, sess *session.Session) (*spreadsheet.Payload, error) {
		secrets, err := helpers.Secrets(sess, cred)
		if err != nil {
			return nil, newQueryErrorf(err, "failed to get Secrets for account: %s, region: %s -> %v", account, *sess.Config.Region, err)
		}
		var items []interface{}
		for _, g := range secrets {
			items = append(items, g)
		}
		return &spreadsheet.Payload{Static: []string{account, *sess.Config.Region}, Items: items}, nil
	})
}

// querySubscriptions ... queries SNS Subscriptions for all organization accounts and
// all sessions/regions in SessionMgr, pushes them onto a slice of interface
// then returns a slice of *spreadsheet.Payload
func (inv *Inv) querySubscriptions() ([]*spreadsheet.Payload, error) {
	defer logDuration()()
	return inv.walkSessions(func(account string, cred *credentials.Credentials, sess *session.Session) (*spreadsheet.Payload, error) {
		secrets, err := helpers.Subscriptions(sess, cred)
		if err != nil {
			return nil, newQueryErrorf(err, "failed to get SNS Subscriptions for account: %s, region: %s -> %v", account, *sess.Config.Region, err)
		}
		var items []interface{}
		for _, g := range secrets {
			items = append(items, g)
		}
		return &spreadsheet.Payload{Static: []string{account, *sess.Config.Region}, Items: items}, nil
	})
}

// queryTopics ... queries SNS Topics for all organization accounts and
// all sessions/regions in SessionMgr, pushes them onto a slice of interface
// then returns a slice of *spreadsheet.Payload
func (inv *Inv) queryTopics() ([]*spreadsheet.Payload, error) {
	defer logDuration()()
	return inv.walkSessions(func(account string, cred *credentials.Credentials, sess *session.Session) (*spreadsheet.Payload, error) {
		secrets, err := helpers.Topics(sess, cred)
		if err != nil {
			return nil, newQueryErrorf(err, "failed to get SNS Topics for account: %s, region: %s -> %v", account, *sess.Config.Region, err)
		}
		var items []interface{}
		for _, g := range secrets {
			items = append(items, g)
		}
		return &spreadsheet.Payload{Static: []string{account, *sess.Config.Region}, Items: items}, nil
	})
}

// queryParameters ... queries SSM Parameter stores for all organization accounts and
// all sessions/regions in SessionMgr, pushes them onto a slice of interface
// then returns a slice of *spreadsheet.Payload
func (inv *Inv) queryParameters() ([]*spreadsheet.Payload, error) {
	defer logDuration()()
	return inv.walkSessions(func(account string, cred *credentials.Credentials, sess *session.Session) (*spreadsheet.Payload, error) {
		secrets, err := helpers.Parameters(sess, cred)
		if err != nil {
			return nil, newQueryErrorf(err, "failed to get SSM Parameters for account: %s, region: %s -> %v", account, *sess.Config.Region, err)
		}
		var items []interface{}
		for _, g := range secrets {
			items = append(items, g)
		}
		return &spreadsheet.Payload{Static: []string{account, *sess.Config.Region}, Items: items}, nil
	})
}
