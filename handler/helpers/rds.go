package helpers

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
)

// RDSSvc ... rds service interface
type RDSSvc struct {
	Client rdsiface.RDSAPI
}

// NewRDSSvc ...
func NewRDSSvc(cfg client.ConfigProvider, cred *credentials.Credentials) (*RDSSvc, error) {
	if cfg == nil {
		return nil, errors.New("nil ConfigProvider")
	}
	return &RDSSvc{
		Client: rds.New(cfg, &aws.Config{Credentials: cred}),
	}, nil
}

// DBInstances ... pages through DescribeDBInstancesPages to get list of DBInstances
func (svc RDSSvc) DBInstances() ([]*rds.DBInstance, error) {
	var results []*rds.DBInstance
	err := svc.Client.DescribeDBInstancesPages(&rds.DescribeDBInstancesInput{},
		func(page *rds.DescribeDBInstancesOutput, lastPage bool) bool {
			results = append(results, page.DBInstances...)
			return !lastPage
		})
	if err != nil {
		return nil, err
	}
	return results, nil
}

// DBSnapshots ... pages through DescribeDBSnapshotsPages to get list of DBSnapshots
func (svc RDSSvc) DBSnapshots() ([]*rds.DBSnapshot, error) {
	var results []*rds.DBSnapshot
	err := svc.Client.DescribeDBSnapshotsPages(&rds.DescribeDBSnapshotsInput{},
		func(page *rds.DescribeDBSnapshotsOutput, lastPage bool) bool {
			results = append(results, page.DBSnapshots...)
			return !lastPage
		})
	if err != nil {
		return nil, err
	}
	return results, nil
}

