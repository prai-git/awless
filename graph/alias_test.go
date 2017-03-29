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
	"testing"
)

func TestResourceNameToId(t *testing.T) {
	g := NewGraph()
	g.Unmarshal([]byte(`/node<inst_1>  "rdf:type"@[] /node<cloud-owl:Instance>
  /node<inst_1>  "cloud:id"@[] "inst_1"^^type:text
  /node<inst_1>  "cloud:name"@[] "redis"^^type:text
	/node<inst_2>  "rdf:type"@[] /node<cloud-owl:Instance>
  /node<inst_2>  "cloud:id"@[] "inst_2"^^type:text
  /node<inst_2>  "cloud:name"@[] "redis2"^^type:text
	/node<inst_3>  "rdf:type"@[] /node<cloud-owl:Instance>
  /node<inst_3>  "cloud:id"@[] "inst_3"^^type:text
  /node<inst_3>  "cloud:name"@[] "mongo"^^type:text
  /node<inst_3>  "cloud:creationDate"@[] "2017-01-10T16:47:18Z"^^type:text
	/node<subnet_1>  "rdf:type"@[] /node<cloud-owl:Subnet>
  /node<subnet_1>  "cloud:id"@[] "subnet_1"^^type:text
  /node<subnet_1>  "cloud:name"@[] "mongo"^^type:text`))

	tcases := []struct {
		name         string
		resourceType string
		expectID     string
		ok           bool
	}{
		{name: "redis", resourceType: "instance", expectID: "inst_1", ok: true},
		{name: "redis2", resourceType: "instance", expectID: "inst_2", ok: true},
		{name: "mongo", resourceType: "instance", expectID: "inst_3", ok: true},
		{name: "mongo", resourceType: "subnet", expectID: "subnet_1", ok: true},
		{name: "nothere", expectID: "", ok: false},
	}
	for i, tcase := range tcases {
		a := Alias(tcase.name)
		id, ok := a.ResolveToId(g, tcase.resourceType)
		if got, want := ok, tcase.ok; got != want {
			t.Fatalf("%d: got %t, want %t", i, got, want)
		}
		if ok {
			if got, want := id, tcase.expectID; got != want {
				t.Fatalf("%d: got %s, want %s", i, got, want)
			}
		}
	}
}
