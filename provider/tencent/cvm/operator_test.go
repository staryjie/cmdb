package cvm_test

import (
	"context"
	"testing"

	"github.com/infraboard/mcube/logger/zap"
	tenctnt_cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	"github.com/staryjie/cmdb/apps/host"
	"github.com/staryjie/cmdb/provider/tencent/connectivity"
	"github.com/staryjie/cmdb/provider/tencent/cvm"
)

var (
	op *cvm.CVMOperator
)

func TestQueryCvm(t *testing.T) {
	var (
		offset int64 = 0
		limit  int64 = 20
	)
	req := tenctnt_cvm.NewDescribeInstancesRequest()
	req.Offset = &offset
	req.Limit = &limit

	set, err := op.Query(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(set)

}

func TestPaggerQuery(t *testing.T) {
	p := cvm.NewPagger(op)
	set := host.NewHostSet()
	for p.Next() {
		if err := p.Scan(context.Background(), set); err != nil {
			panic(err)
		}
		t.Log(set)
	}
}

func init() {
	// 初始化cvm客户端
	err := connectivity.LoadClientFromEnv()
	if err != nil {
		panic(err)
	}

	// 初始化全局logger
	zap.DevelopmentSetup()

	op = cvm.NewCVMOperator(connectivity.C().CvmClient())
}
