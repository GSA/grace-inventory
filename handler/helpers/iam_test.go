package helpers

import (
	"reflect"
	"testing"

	"github.com/GSA/grace-inventory-lambda/handler/inv"

	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
)

type mockIamClient struct {
	iamiface.IAMAPI
}

func (m *mockIamClient) ListRolesPages(in *iam.ListRolesInput, fn func(*iam.ListRolesOutput, bool) bool) error {
	fn(&iam.ListRolesOutput{Roles: []*iam.Role{{}}}, true)
	return nil
}

func (m *mockIamClient) ListGroupsPages(in *iam.ListGroupsInput, fn func(*iam.ListGroupsOutput, bool) bool) error {
	fn(&iam.ListGroupsOutput{Groups: []*iam.Group{{}}}, true)
	return nil
}

func (m *mockIamClient) ListPoliciesPages(in *iam.ListPoliciesInput, fn func(*iam.ListPoliciesOutput, bool) bool) error {
	fn(&iam.ListPoliciesOutput{Policies: []*iam.Policy{{}}}, true)
	return nil
}

func (m *mockIamClient) ListUsersPages(in *iam.ListUsersInput, fn func(*iam.ListUsersOutput, bool) bool) error {
	fn(&iam.ListUsersOutput{Users: []*iam.User{{}}}, true)
	return nil
}

// func Roles() ([]*iam.Role, error)
func TestRoles(t *testing.T) {
	svc := IamSvc{Client: &mockIamClient{}}
	expected := []*iam.Role{{}}
	roles, err := svc.Roles()
	if err != nil {
		t.Fatalf("Roles() failed: %v", err)
	}
	if !reflect.DeepEqual(expected, roles) {
		t.Errorf("Roles() failed. Expected: %#v (%T)\nGot: %#v (%T)", expected, expected, roles, roles)
	}
	_, err = inv.TypeToSheet(expected)
	if err != nil {
		t.Fatalf("inv.TypeToSheet failed: %v", err)
	}
}

// func Groups() ([]*iam.Group, error)
func TestGroups(t *testing.T) {
	svc := IamSvc{Client: &mockIamClient{}}
	expected := []*iam.Group{{}}
	groups, err := svc.Groups()
	if err != nil {
		t.Fatalf("Groups() failed: %v", err)
	}
	if !reflect.DeepEqual(expected, groups) {
		t.Errorf("Groups() failed. Expected: %#v (%T)\nGot: %#v (%T)", expected, expected, groups, groups)
	}
	_, err = inv.TypeToSheet(expected)
	if err != nil {
		t.Fatalf("inv.TypeToSheet failed: %v", err)
	}
}

// func Policies() ([]*iam.Policy, error)
func TestPolicies(t *testing.T) {
	svc := IamSvc{Client: &mockIamClient{}}
	expected := []*iam.Policy{{}}
	policies, err := svc.Policies()
	if err != nil {
		t.Fatalf("Policies() failed: %v", err)
	}
	if !reflect.DeepEqual(expected, policies) {
		t.Errorf("Policies() failed. Expected: %#v (%T)\nGot: %#v (%T)", expected, expected, policies, policies)
	}
	_, err = inv.TypeToSheet(expected)
	if err != nil {
		t.Fatalf("inv.TypeToSheet failed: %v", err)
	}
}

// func Users() ([]*iam.User, error)
func TestUsers(t *testing.T) {
	svc := IamSvc{Client: &mockIamClient{}}
	expected := []*iam.User{{}}
	users, err := svc.Users()
	if err != nil {
		t.Fatalf("Users() failed: %v", err)
	}
	if !reflect.DeepEqual(expected, users) {
		t.Errorf("Users() failed. Expected: %#v (%T)\nGot: %#v (%T)", expected, expected, users, users)
	}
	_, err = inv.TypeToSheet(expected)
	if err != nil {
		t.Fatalf("inv.TypeToSheet failed: %v", err)
	}
}
