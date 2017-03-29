package aws

import (
	"github.com/wallix/awless/cloud/properties"
	"github.com/wallix/awless/graph"
)

type rBuilder struct {
	id, typ string
	props   map[string]interface{}
}

func testRes(id, typ string) *rBuilder {
	return &rBuilder{id: id, typ: typ, props: make(map[string]interface{})}
}

func instRes(id string) *rBuilder {
	return testRes(id, "instance").prop(properties.ID, id)
}

func subRes(id string) *rBuilder {
	return testRes(id, "subnet").prop(properties.ID, id)
}

func vpcRes(id string) *rBuilder {
	return testRes(id, "vpc").prop(properties.ID, id)
}

func sGrpRes(id string) *rBuilder {
	return testRes(id, "securitygroup").prop(properties.ID, id)
}

func keyRes(id string) *rBuilder {
	return testRes(id, "keypair").prop(properties.ID, id)
}

func igwRes(id string) *rBuilder {
	return testRes(id, "internetgateway").prop(properties.ID, id)
}

func rtRes(id string) *rBuilder {
	return testRes(id, "routetable").prop(properties.ID, id)
}

func lbRes(id string) *rBuilder {
	return testRes(id, "loadbalancer").prop(properties.ID, id)
}

func tgRes(id string) *rBuilder {
	return testRes(id, "targetgroup").prop(properties.ID, id)
}

func polRes(id string) *rBuilder {
	return testRes(id, "policy").prop(properties.ID, id)
}

func grpRes(id string) *rBuilder {
	return testRes(id, "group").prop(properties.ID, id)
}

func rolRes(id string) *rBuilder {
	return testRes(id, "role").prop(properties.ID, id)
}

func usrRes(id string) *rBuilder {
	return testRes(id, "user").prop(properties.ID, id)
}

func listRes(id string) *rBuilder {
	return testRes(id, "listener").prop(properties.ID, id)
}

func buckRes(id string) *rBuilder {
	return testRes(id, "bucket").prop(properties.ID, id)
}

func zoneRes(id string) *rBuilder {
	return testRes(id, "zone").prop(properties.ID, id)
}

func recRes(id string) *rBuilder {
	return testRes(id, "record").prop(properties.ID, id)
}

func (b *rBuilder) prop(key string, value interface{}) *rBuilder {
	b.props[key] = value
	return b
}

func (b *rBuilder) build() *graph.Resource {
	res := graph.InitResource(b.id, b.typ)
	res.Properties = b.props
	return res
}
