package cvm

import (
	"context"

	"github.com/staryjie/cmdb/apps/host"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

// 腾讯云CVM查看实例列表API文档: https://cloud.tencent.com/document/api/213/15728
func (o *CVMOperator) Query(ctx context.Context, req *cvm.DescribeInstancesRequest) (
	*host.HostSet, error) {
	resp, err := o.client.DescribeInstancesWithContext(ctx, req)
	if err != nil {
		return nil, err
	}

	// b, _ := json.Marshal(r)
	// return string(b)

	// --> function --> 面过程 --> 代码的阅读这需要按流程来阅读代码
	// 维护代码的 需要了解很多细节
	// 面向--> 封装

	// 1. 打印日志
	// 2. debug
	o.log.Debug(resp.ToJsonString())

	set := o.transferSet(resp.Response.InstanceSet)

	return set, nil
}
