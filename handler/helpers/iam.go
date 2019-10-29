package helpers

import (
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
)

// IamSvc ... uses an SDK service iface to access SDK service client
type IamSvc struct {
	Client iamiface.IAMAPI
}

// Roles ... pages through ListRolesPages and returns all IAM roles
func (svc *IamSvc) Roles() ([]*iam.Role, error) {
	var results []*iam.Role
	err := svc.Client.ListRolesPages(&iam.ListRolesInput{},
		func(page *iam.ListRolesOutput, lastPage bool) bool {
			results = append(results, page.Roles...)
			return !lastPage
		})
	if err != nil {
		return nil, err
	}
	return results, nil
}

// Groups ... pages through ListGroupsPages and returns all IAM groups
func (svc *IamSvc) Groups() ([]*iam.Group, error) {
	var results []*iam.Group
	err := svc.Client.ListGroupsPages(&iam.ListGroupsInput{},
		func(page *iam.ListGroupsOutput, lastPage bool) bool {
			results = append(results, page.Groups...)
			return !lastPage
		})
	if err != nil {
		return nil, err
	}
	return results, nil
}

// Policies ... pages through ListPoliciesPages and returns all IAM policies
func (svc *IamSvc) Policies() ([]*iam.Policy, error) {
	var results []*iam.Policy
	err := svc.Client.ListPoliciesPages(&iam.ListPoliciesInput{},
		func(page *iam.ListPoliciesOutput, lastPage bool) bool {
			results = append(results, page.Policies...)
			return !lastPage
		})
	if err != nil {
		return nil, err
	}
	return results, nil
}

// Users ... pages through ListUsersPages and returns all IAM users
func (svc *IamSvc) Users() ([]*iam.User, error) {
	var results []*iam.User
	err := svc.Client.ListUsersPages(&iam.ListUsersInput{},
		func(page *iam.ListUsersOutput, lastPage bool) bool {
			results = append(results, page.Users...)
			return !lastPage
		})
	if err != nil {
		return nil, err
	}
	return results, nil
}
