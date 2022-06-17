package impl_test

import (
	"context"
	"os"
	"testing"

	"github.com/infraboard/mcube/app"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/staryjie/cmdb/conf"

	_ "github.com/staryjie/cmdb/apps/all"
	"github.com/staryjie/cmdb/apps/secret"
)

var (
	sct secret.ServiceServer
)

func TestCreateSecret(t *testing.T) {
	req := secret.NewCreateSecretRequest()

	req.Description = "测试用例"
	req.ApiKey = os.Getenv("TX_CLOUD_SECRET_ID")
	req.ApiSecret = os.Getenv("TX_CLOUD_SECRET_KEY")
	req.AllowRegions = []string{"*"}

	ss, err := sct.CreateSecret(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ss)
}

func TestDescribeSecret(t *testing.T) {
	ss, err := sct.DescribeSecret(context.Background(), secret.NewDescribeSecretRequest("cam48hi6i3s3el6tt270"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ss)
}

func init() {
	// 通过环境变量加载测试配置
	// if err := conf.LoadConfigFromToml("../../../etc/config.toml"); err != nil {
	if err := conf.LoadConfigFromEnv(); err != nil {
		panic(err)
	}

	// 全局日志对象的初始化
	zap.DevelopmentSetup()

	// 初始化所有实例
	if err := app.InitAllApp(); err != nil {
		panic(err)
	}

	sct = app.GetGrpcApp(secret.AppName).(secret.ServiceServer)
}
