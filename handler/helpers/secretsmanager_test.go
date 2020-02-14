package helpers

import (
	"errors"
	"reflect"
	"testing"

	"github.com/GSA/grace-inventory-lambda/handler/inv"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
)

type mockedSecretsManager struct {
	secretsmanageriface.SecretsManagerAPI
	mocked

	describeSecretsPages []*secretsmanager.ListSecretsOutput
	describeSecretsErr   error
}

func (m mockedSecretsManager) ListSecretsPages(inp *secretsmanager.ListSecretsInput, f func(*secretsmanager.ListSecretsOutput, bool) bool) error {
	m.Called("DescribeSecretsPages")
	for i, p := range m.describeSecretsPages {
		f(p, (i == (len(m.describeSecretsPages) - 1)))
	}
	return m.describeSecretsErr
}

func TestSecretsManagerSvc_Secrets(t *testing.T) {
	page1 := []*secretsmanager.SecretListEntry{
		{Name: aws.String("qwe")},
	}
	page2 := []*secretsmanager.SecretListEntry{
		{Name: aws.String("asd")},
		{Name: aws.String("zxc")},
	}
	allPages := append(page1, page2...)
	tests := []struct {
		name  string
		err   error
		pages []*secretsmanager.ListSecretsOutput

		want      []*secretsmanager.SecretListEntry
		wantErr   bool
		wantCalls []string
	}{
		{
			name: "error",
			err:  errors.New("tst error"),

			wantErr:   true,
			wantCalls: []string{"DescribeSecretsPages"},
		},
		{
			name: "ok onepage",
			err:  nil,
			pages: []*secretsmanager.ListSecretsOutput{
				{SecretList: page1},
			},

			want:      page1,
			wantCalls: []string{"DescribeSecretsPages"},
		},
		{
			name: "ok multipage",
			err:  nil,
			pages: []*secretsmanager.ListSecretsOutput{
				{SecretList: page1},
				{SecretList: page2},
			},

			want:      allPages,
			wantCalls: []string{"DescribeSecretsPages"},
		},
	}
	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			m := mockedSecretsManager{
				mocked: mocked{mockCalls: &mockCalls{}},

				describeSecretsPages: tc.pages,
				describeSecretsErr:   tc.err,
			}
			svc := SecretsManagerSvc{
				Client: &m,
			}
			got, err := svc.Secrets()
			if (err != nil) != tc.wantErr {
				t.Errorf("EC2Svc.Secrets() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			_, err = inv.TypeToSheet(got)
			if err != nil {
				t.Fatalf("inv.TypeToSheet failed: %v", err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("EC2Svc.Secrets() = %v, want %v", got, tc.want)
			}
			if !reflect.DeepEqual(m.CallsList(), tc.wantCalls) {
				t.Errorf("Call list mismatch: %v, expected %v", m.CallsList(), tc.wantCalls)
			}
		})
	}
}
