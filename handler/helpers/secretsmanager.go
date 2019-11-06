package helpers

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
)

// SecretsManagerSvc ...
type SecretsManagerSvc struct {
	Client secretsmanageriface.SecretsManagerAPI
}

// NewSecretsManagerSvc ...
func NewSecretsManagerSvc(cfg client.ConfigProvider, cred *credentials.Credentials) (*SecretsManagerSvc, error) {
	if cfg == nil {
		return nil, errors.New("nil ConfigProvider")
	}
	return &SecretsManagerSvc{
		Client: secretsmanager.New(cfg, &aws.Config{Credentials: cred}),
	}, nil
}

// Secrets ... pages through ListSecretsPages to get list of Secrets
func (svc SecretsManagerSvc) Secrets() ([]*secretsmanager.SecretListEntry, error) {
	var results []*secretsmanager.SecretListEntry
	err := svc.Client.ListSecretsPages(&secretsmanager.ListSecretsInput{},
		func(page *secretsmanager.ListSecretsOutput, lastPage bool) bool {
			results = append(results, page.SecretList...)
			return !lastPage
		})
	if err != nil {
		return nil, err
	}
	return results, nil
}
