package connectivity

import (
	"fmt"

	"github.com/caarlos0/env/v6"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	sts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sts/v20180813" // 管理用户凭证

	"github.com/staryjie/cmdb/utils"
)

var (
	client *TencentCloudClient
)

func C() *TencentCloudClient {
	if client == nil {
		panic("please load config first")
	}
	return client
}

// TencentCloudClient client for all TencentCloud service
type TencentCloudClient struct {
	Region    string `env:"TX_CLOUD_REGION"`
	SecretID  string `env:"TX_CLOUD_SECRET_ID"`
	SecretKey string `env:"TX_CLOUD_SECRET_KEY"`

	accountId string
	cvmConn   *cvm.Client
}

func LoadClientFromEnv() error {
	client = &TencentCloudClient{}
	if err := env.Parse(client); err != nil {
		return err
	}

	return nil
}

// NewTencentCloudClient client
func NewTencentCloudClient(secretID, secretKey, region string) *TencentCloudClient {
	return &TencentCloudClient{
		Region:    region,
		SecretID:  secretID,
		SecretKey: secretKey,
	}
}

// UseCvmClient cvm
func (me *TencentCloudClient) CvmClient() *cvm.Client {
	if me.cvmConn != nil {
		return me.cvmConn
	}

	credential := common.NewCredential(
		me.SecretID,
		me.SecretKey,
	)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = "POST"
	cpf.HttpProfile.ReqTimeout = 300
	cpf.Language = "en-US"

	cvmConn, _ := cvm.NewClient(credential, me.Region, cpf)
	me.cvmConn = cvmConn
	return me.cvmConn
}

// 获取客户端账号ID
func (me *TencentCloudClient) Check() error {
	credential := common.NewCredential(
		me.SecretID,
		me.SecretKey,
	)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = "POST"
	cpf.HttpProfile.ReqTimeout = 300
	cpf.Language = "en-US"

	stsConn, _ := sts.NewClient(credential, me.Region, cpf)

	req := sts.NewGetCallerIdentityRequest()

	resp, err := stsConn.GetCallerIdentity(req)
	if err != nil {
		return fmt.Errorf("unable to initialize the STS client: %#v", err)
	}

	// me.accountId = *resp.Response.AccountId  // 可能是空指针，会报错
	me.accountId = utils.PtrStrV(resp.Response.AccountId)
	return nil
}

// 将私有变量暴露，在外部通过该函数访问
func (me *TencentCloudClient) AccountID() string {
	return me.accountId
}
