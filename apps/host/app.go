package host

import (
	context "context"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/infraboard/mcube/flowcontrol/tokenbucket"
	request "github.com/infraboard/mcube/http/request"
	pb_request "github.com/infraboard/mcube/pb/request"
	resource "github.com/staryjie/cmdb/apps/resource"
	"github.com/staryjie/cmdb/utils"
	"google.golang.org/protobuf/proto"
)

const (
	AppName = "host"
)

// use a single instance of Validate, it caches struct info
var (
	validate = validator.New()
)

func (h *Host) GenHash() error {
	// hash resource
	h.Base.ResourceHash = h.Information.Hash()
	// hash describe
	h.Base.DescribeHash = utils.Hash(h.Describe)
	return nil
}

func (d *Describe) KeyPairNameToString() string {
	return strings.Join(d.KeyPairName, ",")
}

func (d *Describe) SecurityGroupsToString() string {
	return strings.Join(d.SecurityGroups, ",")
}

func (d *Describe) LoadKeyPairNameString(s string) {
	if s != "" {
		d.KeyPairName = strings.Split(s, ",")
	}
}

func (d *Describe) LoadSecurityGroupsString(s string) {
	if s != "" {
		d.SecurityGroups = strings.Split(s, ",")
	}
}

func NewDescribeHostRequestWithID(id string) *DescribeHostRequest {
	return &DescribeHostRequest{
		DescribeBy: DescribeBy_HOST_ID,
		Value:      id,
	}
}

func NewUpdateHostRequest(id string) *UpdateHostRequest {
	return &UpdateHostRequest{
		Id:             id,
		UpdateMode:     pb_request.UpdateMode_PUT,
		UpdateHostData: &UpdateHostData{},
	}
}

func (req *UpdateHostRequest) Validate() error {
	return validate.Struct(req)
}

func (h *Host) Put(req *UpdateHostData) {
	oldRH, oldDH := h.Base.ResourceHash, h.Base.DescribeHash

	h.Information = req.Information
	h.Describe = req.Describe
	h.Information.UpdateAt = time.Now().UnixMilli()
	h.GenHash()

	if h.Base.ResourceHash != oldRH {
		h.Base.ResourceHashChanged = true
	}
	if h.Base.DescribeHash != oldDH {
		h.Base.DescribeHashChanged = true
	}
}

func NewUpdateHostDataByIns(ins *Host) *UpdateHostData {
	return &UpdateHostData{
		Information: ins.Information,
		Describe:    ins.Describe,
	}
}

func NewHostSet() *HostSet {
	return &HostSet{
		Items: []*Host{},
	}
}

func NewDefaultHost() *Host {
	return &Host{
		Base: &resource.Base{
			ResourceType: resource.Type_HOST,
		},
		Information: &resource.Information{},
		Describe:    &Describe{},
	}
}

func (s *HostSet) Add(item any) {
	s.Items = append(s.Items, item.(*Host))
	return
}

func (s *HostSet) Length() int64 {
	return int64(len(s.Items))
}

func (s *HostSet) ResourceIds() (ids []string) {
	for i := range s.Items {
		ids = append(ids, s.Items[i].Base.Id)
	}
	return
}

func (s *HostSet) Clone() *HostSet {
	return proto.Clone(s).(*HostSet)
}

func (s *HostSet) UpdateTag(tags []*resource.Tag) {
	for i := range tags {
		for j := range s.Items {
			if s.Items[j].Base.Id == tags[i].ResourceId {
				s.Items[j].Information.AddTag(tags[i])
			}
		}
	}
}

func (req *DescribeHostRequest) Where() (string, interface{}) {
	switch req.DescribeBy {
	default:
		return "r.id = ?", req.Value
	}
}

func NewQueryHostRequestFromHTTP(r *http.Request) *QueryHostRequest {
	qs := r.URL.Query()
	page := request.NewPageRequestFromHTTP(r)
	kw := qs.Get("keywords")

	return &QueryHostRequest{
		Page:     page,
		Keywords: kw,
	}
}

// 分页器接口
// for p.Next() {
// if err := p.Scan(set); err != nil {
// panic(err)
// }
// }
type Pagger interface {
	Next() bool
	SetPageSize(ps int64)
	Scan(context.Context, *HostSet) error
}

type Set interface {
	// 往set中添加任何类型的元素
	Add(any)
	// 当前集合中有多少个元素
	Length() int64
}

type PaggerV2 interface {
	Next() bool
	Offset() int64
	SetPageSize(ps int64)
	PageSize() int64
	PageNumber() int64
	SetRate(r float64)
	Scan(context.Context, Set) error
}

// 利用面向对象的继承，实现一个模板，把除了Scan()之外的其他方法都实现

type BasePaggerV2 struct {
	// 令牌桶
	hasNext bool
	tb      *tokenbucket.Bucket

	// 控制分页的核心参数
	pageNumber int64
	pageSize   int64
}

func NewBasePaggerV2() *BasePaggerV2 {
	return &BasePaggerV2{
		hasNext:    true,
		tb:         tokenbucket.NewBucketWithRate(1, 1),
		pageNumber: 1,
		pageSize:   20,
	}
}

func (p *BasePaggerV2) Next() bool {
	// 等待分配令牌
	p.tb.Wait(1)

	return p.hasNext
}

func (p *BasePaggerV2) Offset() int64 {
	return (p.pageNumber - 1) * p.pageSize
}

func (p *BasePaggerV2) SetPageSize(ps int64) {
	p.pageSize = ps
}

func (p *BasePaggerV2) PageSize() int64 {
	return p.pageSize
}

func (p *BasePaggerV2) PageNumber() int64 {
	return p.pageNumber
}

func (p *BasePaggerV2) SetRate(r float64) {
	p.tb.SetRate(r)
}

func (p *BasePaggerV2) CheckHasNext(current int64) {
	// 可以根据当前一页是满页来决定是否有下一页
	if current < p.pageSize {
		p.hasNext = false
	} else {
		// 直接调整指针指向下一页
		p.pageNumber++
	}
}
