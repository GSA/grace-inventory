package helpers

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
)

type mockIamClient struct {
	iamiface.IAMAPI
}

func (m *mockIamClient) ListRolesPages(in *iam.ListRolesInput, f func(*iam.ListRolesOutput, bool) bool) error {
	f(&iam.ListRolesOutput{Roles: []*iam.Role{}}, false)
	return nil
}

// func Roles(sess *session.Session, cred *credentials.Credentials) ([]*iam.Role, error)
func TestRoles(t *testing.T) {
	svc := IamSvc{Client: &mockIamClient{}}
	expected := []*iam.Role{nil}
	roles, err := svc.Roles()
	if err != nil {
		t.Fatalf("Roles() failed: %v", err)
	}
	if !reflect.DeepEqual(expected, roles) {
		t.Errorf("Roles() failed. Expected: %#v (%T)\nGot: %#v (%T)", expected, expected, roles, roles)
	}
}

/*
// func Groups(sess *session.Session, cred *credentials.Credentials) ([]*iam.Group, error)
func TestGroups(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	_, err = Groups(sess, nil)
	if err != nil {
		t.Fatalf("Groups() failed: %v", err)
	}
}

// func Policies(sess *session.Session, cred *credentials.Credentials) ([]*iam.Policy, error)
func TestPolicies(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	_, err = Policies(sess, nil)
	if err != nil {
		t.Fatalf("Policies() failed: %v", err)
	}
}

// func Users(sess *session.Session, cred *credentials.Credentials) ([]*iam.User, error)
func TestUsers(t *testing.T) {
	sess, err := awstest.NewAuthenticatedSession(defaultRegion)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	_, err = Users(sess, nil)
	if err != nil {
		t.Fatalf("Users() failed: %v", err)
	}
}
*/
