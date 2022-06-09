package connectivity_test

import (
	"testing"

	"github.com/staryjie/cmdb/provider/tencent/connectivity"
)

func TestTencentCloudClient(t *testing.T) {
	conn := connectivity.C()
	if err := conn.Check(); err != nil {
		t.Fatal(err)
	}
	t.Log(conn.AccountID())
}

func init() {
	// 初始化客户端对象
	err := connectivity.LoadClientFromEnv()
	if err != nil {
		panic(err)
	}
}
