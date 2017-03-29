package graph

import (
	"testing"

	"github.com/wallix/awless/cloud/properties"
)

func TestFilterGraph(t *testing.T) {
	g := NewGraph()
	g.Unmarshal([]byte(
		`/node<inst_1>  "rdf:type"@[] /node<cloud-owl:Instance>
  /node<inst_1>  "cloud:id"@[] "inst_1"^^type:text
  /node<inst_2>  "rdf:type"@[] /node<cloud-owl:Instance>
  /node<inst_2>  "cloud:id"@[] "inst_2"^^type:text
  /node<inst_2>  "cloud:name"@[] "redis"^^type:text
  /node<sub_1>  "rdf:type"@[] /node<cloud-owl:Subnet>
  /node<sub_1>  "cloud:id"@[] "sub_1"^^type:text`))
	filtered, err := g.Filter("subnet")
	if err != nil {
		t.Fatal(err)
	}
	subnets, _ := filtered.GetAllResources("subnet")
	if got, want := len(subnets), 1; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	instances, _ := filtered.GetAllResources("instance")
	if got, want := len(instances), 0; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	filterFn := func(r *Resource) bool {
		if r.Properties[properties.ID] == "inst_1" {
			return true
		}
		return false
	}
	filtered, _ = g.Filter("instance", filterFn)
	instances, _ = filtered.GetAllResources("instance")
	if got, want := len(instances), 1; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := instances[0].Properties[properties.ID], "inst_1"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	subnets, _ = filtered.GetAllResources("subnet")
	if got, want := len(subnets), 0; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	filterOne := func(r *Resource) bool {
		if r.Properties[properties.ID] == "inst_2" {
			return true
		}
		return false
	}
	filterTwo := func(r *Resource) bool {
		if r.Properties[properties.Name] == "redis" {
			return true
		}
		return false
	}
	filtered, _ = g.Filter("instance", filterOne, filterTwo)
	instances, _ = filtered.GetAllResources("instance")
	if got, want := len(instances), 1; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := instances[0].Id(), "inst_2"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := instances[0].Properties["Name"], "redis"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	subnets, _ = filtered.GetAllResources("subnet")
	if got, want := len(subnets), 0; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	filtered, _ = g.Filter("instance",
		BuildPropertyFilterFunc(properties.ID, "inst"),
		BuildPropertyFilterFunc(properties.Name, "Redis"),
	)
	instances, _ = filtered.GetAllResources("instance")
	if got, want := len(instances), 1; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := instances[0].Id(), "inst_2"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := instances[0].Properties["Name"], "redis"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	subnets, _ = filtered.GetAllResources("subnet")
	if got, want := len(subnets), 0; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
}
