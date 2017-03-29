package graph

import "testing"

func TestAndResolver(t *testing.T) {
	t.Parallel()
	g := NewGraph()

	g.Unmarshal([]byte(`/node<inst_1>  "rdf:type"@[] /node<cloud-owl:Instance>
  /node<inst_1>  "cloud:id"@[] "inst_1"^^type:text
  /node<inst_1>  "cloud:name"@[] "redis"^^type:text
  /node<inst_2>  "rdf:type"@[] /node<cloud-owl:Instance>
  /node<inst_2>  "cloud:id"@[] "inst_2"^^type:text
  /node<sub_1>  "rdf:type"@[] /node<cloud-owl:Subnet>
  /node<sub_1>  "cloud:name"@[] "redis"^^type:text`))

	resources, err := g.ResolveResources(&And{[]Resolver{
		&ByType{Typ: "instance"},
		&ByProperty{Name: "Name", Val: "redis"},
	},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(resources), 1; got != want {
		t.Fatalf("got %d want %d", got, want)
	}
	if got, want := resources[0].Id(), "inst_1"; got != want {
		t.Fatalf("got %s want %s", got, want)
	}

	resources, err = g.ResolveResources(&And{[]Resolver{
		&ByType{Typ: "subnet"},
		&ByProperty{Name: "ID", Val: "inst_2"},
	},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(resources), 0; got != want {
		t.Fatalf("got %d want %d", got, want)
	}
}
