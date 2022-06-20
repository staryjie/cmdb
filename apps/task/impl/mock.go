package impl

import (
	"context"

	"github.com/staryjie/cmdb/apps/secret"
)

type secretMock struct {
	secret.UnimplementedServiceServer
}

func (m *secretMock) CreateSecret(context.Context, *secret.CreateSecretRequest) (*secret.Secret, error) {
	return nil, nil
}

func (m *secretMock) QuerySecret(context.Context, *secret.QuerySecretRequest) (*secret.SecretSet, error) {
	return nil, nil
}

func (m *secretMock) DescribeSecret(context.Context, *secret.DescribeSecretRequest) (
	*secret.Secret, error) {
	sct := secret.NewDefaultSecret()
	sct.Data.ApiKey = ""
	sct.Data.ApiSecret = ""
	return sct, nil
}

func (m *secretMock) DeleteSecret(context.Context, *secret.DeleteSecretRequest) (*secret.Secret, error) {
	return nil, nil
}
