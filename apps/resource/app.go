package resource

import "strings"

const (
	AppName = "resource"
)

// 判断查询请求是否携带Tag标签
func (r *SearchRequest) HasTag() bool {
	return len(r.Tags) > 0
}

// Tag的比较操作符, 内比promethues 的Tag比较操作, 官网也能找到该4种操作符号
type Operator string

const (
	// SQL 比较操作  =
	Operator_EQUAL = "="
	// SQL 比较操作  !=
	Operator_NOT_EQUAL = "!="
	// SQL 比较操作  LIKE
	Operator_LIKE_EQUAL = "=~"
	// SQL 比较操作  NOT LIKE
	Operator_NOT_LIKE_EQUAL = "!~"
)

// 多个值比较的关系说明:
//    app=~app1,app2  你不能说 app1和app2是 AND关系, 一定是OR关系    是一种白名单策略(包含策略)
//    app!~app3,app4  tag_key=app tag_value NOT LIKE (app3,app4), 是一种黑名单策略(排除策略)
func (s *TagSelector) RelationShip() string {
	switch s.Operator {
	case Operator_EQUAL, Operator_LIKE_EQUAL:
		return " OR "
	case Operator_NOT_EQUAL, Operator_NOT_LIKE_EQUAL:
		return " AND "
	default:
		return " OR "
	}
}

func NewResourceSet() *ResourceSet {
	return &ResourceSet{
		Items: []*Resource{},
	}
}

func (s *ResourceSet) Add(item *Resource) {
	s.Items = append(s.Items, item)
}

func NewDefaultResource() *Resource {
	return &Resource{
		Base:        &Base{},
		Information: &Information{},
	}
}

func (i *Information) LoadPrivateIPString(s string) {
	if s != "" {
		i.PrivateIp = strings.Split(s, ",")
	}
}

func (i *Information) LoadPublicIPString(s string) {
	if s != "" {
		i.PublicIp = strings.Split(s, ",")
	}
}
