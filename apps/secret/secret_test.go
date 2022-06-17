package secret_test

import (
	"testing"

	"github.com/staryjie/cmdb/apps/secret"
)

var (
	encryptKey = "staryjie"
)

func TestSecretEncrypt(t *testing.T) {
	sct := secret.NewDefaultSecret()
	sct.Data.ApiSecret = "123456"
	// 加密APISecret
	sct.Data.EncryptAPISecret(encryptKey)
	t.Log(sct.Data.ApiSecret)

	// 解密APISecret
	sct.Data.DecryptAPISecret(encryptKey)
	t.Log(sct.Data.ApiSecret)
}
