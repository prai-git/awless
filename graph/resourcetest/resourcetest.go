package resourcetest

import (
	"github.com/wallix/awless/cloud/properties"
	"github.com/wallix/awless/graph"
)

type rBuilder struct {
	id, typ string
	props   map[string]interface{}
}

func New(id, typ string) *rBuilder {
	return &rBuilder{id: id, typ: typ, props: make(map[string]interface{})}
}

func Instance(id string) *rBuilder {
	return New(id, "instance").Prop(properties.ID, id)
}

func Subnet(id string) *rBuilder {
	return New(id, "subnet").Prop(properties.ID, id)
}

func VPC(id string) *rBuilder {
	return New(id, "vpc").Prop(properties.ID, id)
}

func SecGroup(id string) *rBuilder {
	return New(id, "securitygroup").Prop(properties.ID, id)
}

func Keypair(id string) *rBuilder {
	return New(id, "keypair").Prop(properties.ID, id)
}

func InternetGw(id string) *rBuilder {
	return New(id, "internetgateway").Prop(properties.ID, id)
}

func RouteTable(id string) *rBuilder {
	return New(id, "routetable").Prop(properties.ID, id)
}

func LoadBalancer(id string) *rBuilder {
	return New(id, "loadbalancer").Prop(properties.ID, id)
}

func TargetGroup(id string) *rBuilder {
	return New(id, "targetgroup").Prop(properties.ID, id)
}

func Policy(id string) *rBuilder {
	return New(id, "policy").Prop(properties.ID, id)
}

func Group(id string) *rBuilder {
	return New(id, "group").Prop(properties.ID, id)
}

func Role(id string) *rBuilder {
	return New(id, "role").Prop(properties.ID, id)
}

func User(id string) *rBuilder {
	return New(id, "user").Prop(properties.ID, id)
}

func Listener(id string) *rBuilder {
	return New(id, "listener").Prop(properties.ID, id)
}

func Bucket(id string) *rBuilder {
	return New(id, "bucket").Prop(properties.ID, id)
}

func Zone(id string) *rBuilder {
	return New(id, "zone").Prop(properties.ID, id)
}

func Record(id string) *rBuilder {
	return New(id, "record").Prop(properties.ID, id)
}

func (b *rBuilder) Prop(key string, value interface{}) *rBuilder {
	b.props[key] = value
	return b
}

func (b *rBuilder) Build() *graph.Resource {
	res := graph.InitResource(b.id, b.typ)
	res.Properties = b.props
	return res
}
