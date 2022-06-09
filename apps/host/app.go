package host

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	request "github.com/infraboard/mcube/http/request"
	pb_request "github.com/infraboard/mcube/pb/request"
	resource "github.com/staryjie/cmdb/apps/resource"
	"github.com/staryjie/cmdb/utils"
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

func (s *HostSet) ResourceIds() (ids []string) {
	for i := range s.Items {
		ids = append(ids, s.Items[i].Base.Id)
	}
	return
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
