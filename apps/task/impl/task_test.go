package impl_test

import (
	"context"
	"os"
	"testing"

	"github.com/infraboard/mcube/app"
	"github.com/infraboard/mcube/logger/zap"

	"github.com/staryjie/cmdb/apps/resource"
	"github.com/staryjie/cmdb/apps/task"
	"github.com/staryjie/cmdb/conf"

	// 注册所有对象
	_ "github.com/staryjie/cmdb/apps/all"
)

var (
	tsk task.ServiceServer
)

func TestCreateTask(t *testing.T) {
	req := task.NewCreateTaskRequst()
	req.Type = task.Type_RESOURCE_SYNC
	req.Region = os.Getenv("TX_CLOUD_REGION")
	req.ResourceType = resource.Type_HOST
	req.SecretId = "cam48hi6i3s3el6tt270"
	taskIns, err := tsk.CreateTask(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(taskIns)
}

func init() {
	// 通过环境变量加载测试配置
	if err := conf.LoadConfigFromEnv(); err != nil {
		panic(err)
	}

	// 全局日志对象初始化
	zap.DevelopmentSetup()

	// 初始化所有实例
	if err := app.InitAllApp(); err != nil {
		panic(err)
	}

	tsk = app.GetGrpcApp(task.AppName).(task.ServiceServer)
}
