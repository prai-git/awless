/*
Copyright 2017 WALLIX

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package graph

import (
	"fmt"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/wallix/awless/cloud/properties"
)

func TestAddGraphRelation(t *testing.T) {

	t.Run("Add parent", func(t *testing.T) {
		g := NewGraph()
		g.Unmarshal([]byte(`/node<inst_1>  "rdf:type"@[] /node<cloud-owl:Instance>`))

		res, err := g.GetResource("instance", "inst_1")
		if err != nil {
			t.Fatal(err)
		}
		g.AddParentRelation(InitResource("subnet", "subnet_1"), res)

		exp := `/node<inst_1>	"rdf:type"@[]	/node<cloud-owl:Instance>
/node<subnet_1>	"parent_of"@[]	/node<inst_1>`

		if got, want := g.MustMarshal(), exp; got != want {
			t.Fatalf("got\n%q\nwant\n%q\n", got, want)
		}
	})

	t.Run("Add applies on", func(t *testing.T) {
		g := NewGraph()
		g.Unmarshal([]byte(`/node<inst_1>  "rdf:type"@[] /node<cloud-owl:Instance>`))

		res, err := g.GetResource("instance", "inst_1")
		if err != nil {
			t.Fatal(err)
		}
		g.AddAppliesOnRelation(InitResource("subnet", "subnet_1"), res)

		exp := `/node<inst_1>	"rdf:type"@[]	/node<cloud-owl:Instance>
/node<subnet_1>	"applies_on"@[]	/node<inst_1>`

		if got, want := g.MustMarshal(), exp; got != want {
			t.Fatalf("got\n%q\nwant\n%q\n", got, want)
		}
	})
}

func TestGetResource(t *testing.T) {
	g := NewGraph()

	g.Unmarshal([]byte(`/node<inst_1>  "rdf:type"@[] /node<cloud-owl:Instance>
  /node<inst_1>  "cloud:id"@[] "inst_1"^^type:text
  /node<inst_1>  "cloud:name"@[] "redis"^^type:text
  /node<inst_1>  "cloud:type"@[] "t2.micro"^^type:text
  /node<inst_1>  "net:publicIP"@[] "1.2.3.4"^^type:text
  /node<inst_1>  "cloud:state"@[] "running"^^type:text`))

	res, err := g.GetResource("instance", "inst_1")
	if err != nil {
		t.Fatal(err)
	}

	expected := Properties{properties.ID: "inst_1", properties.Type: "t2.micro", properties.PublicIP: "1.2.3.4",
		properties.State: "running",
		properties.Name:  "redis",
	}

	if got, want := res.Properties, expected; !reflect.DeepEqual(got, want) {
		t.Fatalf("got \n%#v\n\nwant \n%#v\n", got, want)
	}
}

func TestFindResources(t *testing.T) {
	t.Parallel()
	g := NewGraph()

	g.Unmarshal([]byte(`/node<inst_1>  "rdf:type"@[] /node<cloud-owl:Instance>
  /node<inst_1>  "cloud:id"@[] "inst_1"^^type:text
  /node<inst_1>  "cloud:name"@[] "redis"^^type:text
  /node<inst_2>  "rdf:type"@[] /node<cloud-owl:Instance>
  /node<inst_2>  "cloud:id"@[] "inst_2"^^type:text
  /node<sub_1>  "rdf:type"@[] /node<cloud-owl:Subnet>
  /node<sub_1>  "cloud:id"@[] "sub_1"^^type:text
  /node<sub_1>  "cloud:name"@[] "redis"^^type:text`))

	t.Run("FindResource", func(t *testing.T) {
		t.Parallel()
		res, err := g.FindResource("inst_1")
		if err != nil {
			t.Fatal(err)
		}
		if got, want := res.Properties["Name"], "redis"; got != want {
			t.Fatalf("got %s want %s", got, want)
		}

		if res, err = g.FindResource("none"); err != nil {
			t.Fatal(err)
		}
		if res != nil {
			t.Fatalf("expected nil got %v", res)
		}

		if res, err = g.FindResource("sub_1"); err != nil {
			t.Fatal(err)
		}
		if got, want := res.Type(), "subnet"; got != want {
			t.Fatalf("got %s want %s", got, want)
		}
	})
	t.Run("FindResourcesByProperty", func(t *testing.T) {
		t.Parallel()
		res, err := g.FindResourcesByProperty("ID", "inst_1")
		if err != nil {
			t.Fatal(err)
		}
		expected := []*Resource{
			{id: "inst_1", kind: "instance", Properties: map[string]interface{}{"ID": interface{}("inst_1"), "Name": interface{}("redis")}, Meta: make(Properties)},
		}
		if got, want := len(res), len(expected); got != want {
			t.Fatalf("got %d want %d", got, want)
		}
		if got, want := res[0], expected[0]; !reflect.DeepEqual(got, want) {
			t.Fatalf("got %+v want %+v", got, want)
		}
		res, err = g.FindResourcesByProperty("Name", "redis")
		if err != nil {
			t.Fatal(err)
		}
		expected = []*Resource{
			{id: "inst_1", kind: "instance", Properties: map[string]interface{}{"ID": "inst_1", "Name": "redis"}, Meta: make(Properties)},
			{id: "sub_1", kind: "subnet", Properties: map[string]interface{}{"ID": "sub_1", "Name": "redis"}, Meta: make(Properties)},
		}
		if got, want := len(res), len(expected); got != want {
			t.Fatalf("got %d want %d", got, want)
		}
		if res[0].Id() == expected[0].Id() {
			if got, want := res[0], expected[0]; !reflect.DeepEqual(got, want) {
				t.Fatalf("got %+v want %+v", got, want)
			}
			if got, want := res[1], expected[1]; !reflect.DeepEqual(got, want) {
				t.Fatalf("got %+v want %+v", got, want)
			}
		} else {
			if got, want := res[0], expected[1]; !reflect.DeepEqual(got, want) {
				t.Fatalf("got %+v want %+v", got, want)
			}
			if got, want := res[1], expected[0]; !reflect.DeepEqual(got, want) {
				t.Fatalf("got %+v want %+v", got, want)
			}
		}
	})
}

func TestGetAllResources(t *testing.T) {
	g := NewGraph()

	g.Unmarshal([]byte(`/node<inst_1>  "rdf:type"@[] /node<cloud-owl:Instance>
  /node<inst_1>  "cloud:id"@[] "inst_1"^^type:text
  /node<inst_1>  "cloud:name"@[] "redis"^^type:text
  /node<inst_2>  "rdf:type"@[] /node<cloud-owl:Instance>
  /node<inst_2>  "cloud:id"@[] "inst_2"^^type:text
  /node<inst_2>  "cloud:name"@[] "redis2"^^type:text
  /node<inst_3>  "rdf:type"@[] /node<cloud-owl:Instance>
  /node<inst_3>  "cloud:id"@[] "inst_3"^^type:text
  /node<inst_3>  "cloud:name"@[] "redis3"^^type:text
  /node<inst_3>  "cloud:created"@[] "2017-01-10T16:47:18Z"^^type:text
  /node<subnet>  "rdf:type"@[] /node<cloud-owl:Subnet>
  /node<subnet>  "cloud:id"@[] "my subnet"^^type:text`))

	time, _ := time.Parse(time.RFC3339, "2017-01-10T16:47:18Z")

	expected := []*Resource{
		{kind: "instance", id: "inst_1", Properties: Properties{properties.ID: "inst_1", "Name": "redis"}},
		{kind: "instance", id: "inst_2", Properties: Properties{properties.ID: "inst_2", "Name": "redis2"}},
		{kind: "instance", id: "inst_3", Properties: Properties{properties.ID: "inst_3", "Name": "redis3", "Created": time}},
	}
	res, err := g.GetAllResources("instance")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(res), len(expected); got != want {
		t.Fatalf("got %d want %d", got, want)
	}
	for _, r := range expected {
		found := false
		for _, r2 := range res {
			if r2.kind == r.kind && r2.id == r.id && reflect.DeepEqual(r2.Properties, r.Properties) {
				found = true
			}
		}
		if !found {
			t.Fatalf("%+v not found", r)
		}
	}
}

func TestLoadIpPermissions(t *testing.T) {
	g := NewGraph()
	g.Unmarshal([]byte(`/node<sg-1234>	"rdf:type"@[]	/node<cloud-owl:Securitygroup>
/node<sg-1234>	"cloud:id"@[]	"sg-1234"^^type:text
/node<sg-1234>	"net:inboundRules"@[]	/node<d04ab55f>
/node<d04ab55f>	"net:portRange"@[]	"22:22"^^type:text
/node<d04ab55f>	"net:protocol"@[]	"tcp"^^type:text
/node<d04ab55f>	"net:cidr"@[]	"10.10.0.0/16"^^type:text
/node<sg-1234>	"net:inboundRules"@[]	/node<36d9ff45>
/node<36d9ff45>	"net:portRange"@[]	"443:443"^^type:text
/node<36d9ff45>	"net:protocol"@[]	"tcp"^^type:text
/node<36d9ff45>	"net:cidr"@[]	"0.0.0.0/0"^^type:text
/node<sg-1234>	"net:outboundRules"@[]	/node<6172bfe3>
/node<6172bfe3>	"net:portRange"@[]	":"^^type:text
/node<6172bfe3>	"net:protocol"@[]	"any"^^type:text
/node<6172bfe3>	"net:cidr"@[]	"0.0.0.0/0"^^type:text`))
	expected := []*Resource{
		{kind: "securitygroup", id: "sg-1234", Properties: Properties{
			"ID": "sg-1234",
			"InboundRules": []*FirewallRule{
				{
					PortRange: PortRange{FromPort: int64(22), ToPort: int64(22), Any: false},
					Protocol:  "tcp",
					IPRanges:  []*net.IPNet{{IP: net.IPv4(10, 10, 0, 0), Mask: net.CIDRMask(16, 32)}},
				},
				{
					PortRange: PortRange{FromPort: int64(443), ToPort: int64(443), Any: false},
					Protocol:  "tcp",
					IPRanges:  []*net.IPNet{{IP: net.IPv4(0, 0, 0, 0), Mask: net.CIDRMask(0, 32)}},
				},
			},
			"OutboundRules": []*FirewallRule{
				{
					PortRange: PortRange{Any: true},
					Protocol:  "any",
					IPRanges:  []*net.IPNet{{IP: net.IPv4(0, 0, 0, 0), Mask: net.CIDRMask(0, 32)}},
				},
			},
		},
		},
	}
	res, err := g.GetAllResources("securitygroup")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(res), len(expected); got != want {
		t.Fatalf("got %d want %d", got, want)
	}
	if got, want := res[0].id, expected[0].id; got != want {
		t.Fatalf("got %s want %s", got, want)
	}
	if got, want := res[0].kind, expected[0].kind; got != want {
		t.Fatalf("got %s want %s", got, want)
	}
	if got, want := len(res[0].Properties), len(expected[0].Properties); got != want {
		t.Fatalf("got %d want %d", got, want)
	}
	for k := range expected[0].Properties {
		if got, want := fmt.Sprintf("%T", res[0].Properties[k]), fmt.Sprintf("%T", expected[0].Properties[k]); got != want {
			t.Fatalf("got %s want %s", got, want)
		}
	}
}
