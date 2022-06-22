package cvm

import (
	"context"

	"github.com/infraboard/mcube/flowcontrol/tokenbucket"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	"github.com/staryjie/cmdb/apps/host"
)

type pagger struct {
	req        *cvm.DescribeInstancesRequest
	op         *CVMOperator
	hasNext    bool
	pageSize   int64
	pageNumber int64
	log        logger.Logger
	tb         *tokenbucket.Bucket
}

func NewPagger(rate float64, op *CVMOperator) *pagger {
	p := &pagger{
		op:         op,
		hasNext:    true,
		pageNumber: 1,
		pageSize:   20,
		log:        zap.L().Named("CVM"),
		tb:         tokenbucket.NewBucketWithRate(rate, 3),
	}

	p.req = cvm.NewDescribeInstancesRequest()
	p.req.Limit = &p.pageSize
	p.req.Offset = p.offset()

	return p
}

// 需要在请求数据的时候动态计算（当前请求页的数据是否满页）
func (p *pagger) Next() bool {
	return p.hasNext
}

// 修改Req 执行真正的下一页的offset
func (p *pagger) nextReq() *cvm.DescribeInstancesRequest {
	// 等待分配令牌
	p.tb.Wait(1)
	// 修改req指向真正的下一页的offset
	p.req.Offset = p.offset()
	p.req.Limit = &p.pageSize

	return p.req
}

// 设置速率限制
func (p *pagger) SetRate(r float64) {
	p.tb.SetRate(r)
}

// 设置pageSize
func (p *pagger) SetPageSize(ps int64) {
	p.pageSize = ps
}

// 根据分页参数计算offset
func (p *pagger) offset() *int64 {
	offset := (p.pageNumber - 1) * p.pageSize
	return &offset
}

func (p *pagger) Scan(ctx context.Context, set *host.HostSet) error {
	p.log.Debugf("Query Page: %d", p.pageNumber)
	hs, err := p.op.Query(ctx, p.nextReq())
	if err != nil {
		return err
	}

	// 查询的数据赋值给set
	// *set = *hs.Clone()
	for i := range hs.Items {
		set.Add(set.Items[i])
	}

	// 获取当前页没有数据，则没有下一页
	// 还可以根据当前页数据是否小于pageSize，如果小于pageSize，则没有下一页
	if hs.Length() < p.pageSize {
		p.hasNext = false
	}

	// 修改指针指向下一页
	p.pageNumber++

	return nil
}
