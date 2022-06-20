package impl

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/staryjie/cmdb/apps/resource"
	"github.com/staryjie/cmdb/apps/secret"
	"github.com/staryjie/cmdb/apps/task"
	"github.com/staryjie/cmdb/conf"
)

// 创建Task
func (s *service) CreateTask(ctx context.Context, req *task.CreateTaskRequst) (
	*task.Task, error) {
	// 创建task实例
	t, err := task.CreateTask(req)
	if err != nil {
		return nil, err
	}

	// 1.查询secret
	sct, err := s.secret.DescribeSecret(ctx, secret.NewDescribeSecretRequest(req.SecretId))
	if err != nil {
		return nil, err
	}
	t.SecretDescription = sct.Data.Description

	// 解密API secret
	if err := sct.Data.DecryptAPISecret(conf.C().App.EncryptKey); err != nil {
		return nil, err
	}

	// 启动task，修改task的状态
	t.Run()

	var taskCancel context.CancelFunc
	switch req.Type {
	// 资源同步
	case task.Type_RESOURCE_SYNC:
		// 根据secret所属的厂商，初始化对应厂商的operator
		switch sct.Data.Vendor {
		// 腾讯云
		case resource.Vendor_TENCENT:
			// 操作哪种资源
			switch req.ResourceType {
			// 同步主机数据
			case resource.Type_HOST:
				taskExecCtx, cancel := context.WithTimeout(
					context.Background(),
					time.Duration(req.Timeout)*time.Second,
				)
				taskCancel = cancel

				go s.syncHost(taskExecCtx, newSyncHostRequest(sct, t))
			case resource.Type_RDS:
			case resource.Type_BILL:
			}
		// 阿里云
		case resource.Vendor_ALIYUN:
		// 华为云
		case resource.Vendor_HUAWEI:
		// 亚马逊
		case resource.Vendor_AMAZON:
		// vsphere虚拟机
		case resource.Vendor_VSPHERE:
		default:
			return nil, fmt.Errorf("Unknow resource type: %s", sct.Data.Vendor)
		}

		// 2. 利用secret的信息, 初始化一个operater
		// 使用operator进行资源的操作, 比如同步

		// 调用host service 把数据入库
	// 资源释放
	case task.Type_RESOURCE_RELEASE:
	default:
		return nil, fmt.Errorf("unknow task type: %s", req.Type)
	}

	// 将数据保存到数据库
	if err := s.insert(ctx, t); err != nil {
		if taskCancel != nil {
			taskCancel()
		}
		return nil, err
	}

	return t, nil
}

func (s *service) QueryTask(ctx context.Context, req *task.QueryTaskRequest) (*task.TaskSet, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryTask not implemented")
}

func (s *service) DescribeTask(ctx context.Context, req *task.DescribeTaskRequest) (*task.Task, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DescribeTask not implemented")
}
