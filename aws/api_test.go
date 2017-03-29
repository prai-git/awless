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

package aws

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/s3"
	p "github.com/wallix/awless/cloud/properties"
	"github.com/wallix/awless/graph"
)

func TestBuildAccessRdfGraph(t *testing.T) {
	managedPolicies := []*iam.ManagedPolicyDetail{
		{PolicyId: awssdk.String("managed_policy_1"), PolicyName: awssdk.String("nmanaged_policy_1")},
		{PolicyId: awssdk.String("managed_policy_2"), PolicyName: awssdk.String("nmanaged_policy_2")},
		{PolicyId: awssdk.String("managed_policy_3"), PolicyName: awssdk.String("nmanaged_policy_3")},
	}

	groups := []*iam.GroupDetail{
		{GroupId: awssdk.String("group_1"), GroupName: awssdk.String("ngroup_1"), GroupPolicyList: []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_1")}}, AttachedManagedPolicies: []*iam.AttachedPolicy{{PolicyName: awssdk.String("nmanaged_policy_1")}}},
		{GroupId: awssdk.String("group_2"), GroupName: awssdk.String("ngroup_2"), GroupPolicyList: []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_1")}}, AttachedManagedPolicies: []*iam.AttachedPolicy{{PolicyName: awssdk.String("nmanaged_policy_2")}}},
		{GroupId: awssdk.String("group_3"), GroupName: awssdk.String("ngroup_3"), GroupPolicyList: []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_2")}}, AttachedManagedPolicies: []*iam.AttachedPolicy{{PolicyName: awssdk.String("nmanaged_policy_3")}}},
		{GroupId: awssdk.String("group_4"), GroupName: awssdk.String("ngroup_4"), GroupPolicyList: []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_4")}}},
	}

	roles := []*iam.RoleDetail{
		{RoleId: awssdk.String("role_1"), RolePolicyList: []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_1")}}, AttachedManagedPolicies: []*iam.AttachedPolicy{{PolicyName: awssdk.String("nmanaged_policy_1")}}},
		{RoleId: awssdk.String("role_2"), RolePolicyList: []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_1")}}},
		{RoleId: awssdk.String("role_3"), RolePolicyList: []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_2")}}, AttachedManagedPolicies: []*iam.AttachedPolicy{{PolicyName: awssdk.String("nmanaged_policy_2")}}},
		{RoleId: awssdk.String("role_4"), RolePolicyList: []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_4")}}},
	}

	usersDetails := []*iam.UserDetail{
		{
			UserId:                  awssdk.String("usr_1"),
			GroupList:               []*string{awssdk.String("ngroup_1"), awssdk.String("ngroup_2")},
			AttachedManagedPolicies: []*iam.AttachedPolicy{{PolicyName: awssdk.String("nmanaged_policy_1")}},
			UserPolicyList:          []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_1")}, {PolicyName: awssdk.String("npolicy_2")}},
		},
		{
			UserId:         awssdk.String("usr_2"),
			GroupList:      []*string{awssdk.String("ngroup_1")},
			UserPolicyList: []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_1")}},
		},
		{
			UserId:                  awssdk.String("usr_3"),
			GroupList:               []*string{awssdk.String("ngroup_1"), awssdk.String("ngroup_4")},
			AttachedManagedPolicies: []*iam.AttachedPolicy{{PolicyName: awssdk.String("nmanaged_policy_1")}, {PolicyName: awssdk.String("nmanaged_policy_2")}},
			UserPolicyList:          []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_1")}, {PolicyName: awssdk.String("npolicy_4")}},
		},
		{
			UserId:         awssdk.String("usr_4"),
			GroupList:      []*string{awssdk.String("ngroup_2")},
			UserPolicyList: []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_2")}},
		},
		{
			UserId:         awssdk.String("usr_5"),
			GroupList:      []*string{awssdk.String("ngroup_2")},
			UserPolicyList: []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_2")}},
		},
		{
			UserId:                  awssdk.String("usr_6"),
			GroupList:               []*string{awssdk.String("ngroup_2")},
			AttachedManagedPolicies: []*iam.AttachedPolicy{{PolicyName: awssdk.String("nmanaged_policy_3")}},
			UserPolicyList:          []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_2")}},
		},
		{
			UserId:         awssdk.String("usr_7"),
			GroupList:      []*string{awssdk.String("ngroup_2"), awssdk.String("ngroup_4")},
			UserPolicyList: []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_2")}, {PolicyName: awssdk.String("npolicy_4")}},
		},
		{
			UserId:         awssdk.String("usr_8"),
			GroupList:      []*string{awssdk.String("ngroup_4")},
			UserPolicyList: []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_4")}},
		},
		{
			UserId:         awssdk.String("usr_9"),
			GroupList:      []*string{awssdk.String("ngroup_4")},
			UserPolicyList: []*iam.PolicyDetail{{PolicyName: awssdk.String("npolicy_4")}},
		},
		{
			UserId: awssdk.String("usr_10"), //users not in any groups
		},
		{
			UserId: awssdk.String("usr_11"),
		},
	}

	users := []*iam.User{
		{
			UserId:           awssdk.String("usr_1"),
			PasswordLastUsed: awssdk.Time(time.Unix(1486139077, 0).UTC()),
		},
		{
			UserId: awssdk.String("usr_2"),
		},
		{
			UserId: awssdk.String("usr_3"),
		},
		{
			UserId: awssdk.String("usr_4"),
		},
		{
			UserId: awssdk.String("usr_5"),
		},
		{
			UserId: awssdk.String("usr_6"),
		},
		{
			UserId: awssdk.String("usr_7"),
		},
		{
			UserId: awssdk.String("usr_8"),
		},
		{
			UserId: awssdk.String("usr_9"),
		},
		{
			UserId: awssdk.String("usr_10"), //users not in any groups
		},
		{
			UserId: awssdk.String("usr_11"),
		},
	}

	mock := &mockIam{groups: groups, usersDetails: usersDetails, roles: roles, managedPolicies: managedPolicies, users: users}
	access := Access{IAMAPI: mock, region: "eu-west-1"}

	g, err := access.FetchResources()
	if err != nil {
		t.Fatal(err)
	}

	resources, err := g.GetAllResources("policy", "group", "role", "user")
	if err != nil {
		t.Fatal(err)
	}

	// Sort slice properties in resources
	for _, res := range resources {
		if p, ok := res.Properties[p.InlinePolicies].([]string); ok {
			sort.Strings(p)
		}
	}

	expected := map[string]*graph.Resource{
		"managed_policy_1": polRes("managed_policy_1").prop(p.Name, "nmanaged_policy_1").build(),
		"managed_policy_2": polRes("managed_policy_2").prop(p.Name, "nmanaged_policy_2").build(),
		"managed_policy_3": polRes("managed_policy_3").prop(p.Name, "nmanaged_policy_3").build(),
		"group_1":          grpRes("group_1").prop(p.Name, "ngroup_1").prop(p.InlinePolicies, []string{"npolicy_1"}).build(),
		"group_2":          grpRes("group_2").prop(p.Name, "ngroup_2").prop(p.InlinePolicies, []string{"npolicy_1"}).build(),
		"group_3":          grpRes("group_3").prop(p.Name, "ngroup_3").prop(p.InlinePolicies, []string{"npolicy_2"}).build(),
		"group_4":          grpRes("group_4").prop(p.Name, "ngroup_4").prop(p.InlinePolicies, []string{"npolicy_4"}).build(),
		"role_1":           rolRes("role_1").prop(p.InlinePolicies, []string{"npolicy_1"}).build(),
		"role_2":           rolRes("role_2").prop(p.InlinePolicies, []string{"npolicy_1"}).build(),
		"role_3":           rolRes("role_3").prop(p.InlinePolicies, []string{"npolicy_2"}).build(),
		"role_4":           rolRes("role_4").prop(p.InlinePolicies, []string{"npolicy_4"}).build(),
		"usr_1":            usrRes("usr_1").prop(p.InlinePolicies, []string{"npolicy_1", "npolicy_2"}).prop(p.PasswordLastUsed, time.Unix(1486139077, 0).UTC()).build(),
		"usr_2":            usrRes("usr_2").prop(p.InlinePolicies, []string{"npolicy_1"}).build(),
		"usr_3":            usrRes("usr_3").prop(p.InlinePolicies, []string{"npolicy_1", "npolicy_4"}).build(),
		"usr_4":            usrRes("usr_4").prop(p.InlinePolicies, []string{"npolicy_2"}).build(),
		"usr_5":            usrRes("usr_5").prop(p.InlinePolicies, []string{"npolicy_2"}).build(),
		"usr_6":            usrRes("usr_6").prop(p.InlinePolicies, []string{"npolicy_2"}).build(),
		"usr_7":            usrRes("usr_7").prop(p.InlinePolicies, []string{"npolicy_2", "npolicy_4"}).build(),
		"usr_8":            usrRes("usr_8").prop(p.InlinePolicies, []string{"npolicy_4"}).build(),
		"usr_9":            usrRes("usr_9").prop(p.InlinePolicies, []string{"npolicy_4"}).build(),
		"usr_10":           usrRes("usr_10").build(),
		"usr_11":           usrRes("usr_11").build(),
	}

	expectedChildren := map[string][]string{}

	expectedAppliedOn := map[string][]string{
		"group_1":          {"usr_1", "usr_2", "usr_3"},
		"group_2":          {"usr_1", "usr_4", "usr_5", "usr_6", "usr_7"},
		"group_4":          {"usr_3", "usr_7", "usr_8", "usr_9"},
		"managed_policy_1": {"group_1", "role_1", "usr_1", "usr_3"},
		"managed_policy_2": {"group_2", "role_3", "usr_3"},
		"managed_policy_3": {"group_3", "usr_6"},
	}

	compareResources(t, g, resources, expected, expectedChildren, expectedAppliedOn)
}

func TestBuildInfraRdfGraph(t *testing.T) {
	instances := []*ec2.Instance{
		{InstanceId: awssdk.String("inst_1"), SubnetId: awssdk.String("sub_1"), VpcId: awssdk.String("vpc_1"), Tags: []*ec2.Tag{{Key: awssdk.String("Name"), Value: awssdk.String("instance1-name")}}},
		{InstanceId: awssdk.String("inst_2"), SubnetId: awssdk.String("sub_2"), VpcId: awssdk.String("vpc_1"), SecurityGroups: []*ec2.GroupIdentifier{{GroupId: awssdk.String("secgroup_1")}}},
		{InstanceId: awssdk.String("inst_3"), SubnetId: awssdk.String("sub_3"), VpcId: awssdk.String("vpc_2")},
		{InstanceId: awssdk.String("inst_4"), SubnetId: awssdk.String("sub_3"), VpcId: awssdk.String("vpc_2"), SecurityGroups: []*ec2.GroupIdentifier{{GroupId: awssdk.String("secgroup_1")}, {GroupId: awssdk.String("secgroup_2")}}, KeyName: awssdk.String("my_key_pair")},
		{InstanceId: awssdk.String("inst_5"), SubnetId: nil, VpcId: nil}, // terminated instance (no vpc, subnet ids)
	}

	vpcs := []*ec2.Vpc{
		{VpcId: awssdk.String("vpc_1")},
		{VpcId: awssdk.String("vpc_2")},
	}

	securityGroups := []*ec2.SecurityGroup{
		{GroupId: awssdk.String("secgroup_1"), GroupName: awssdk.String("my_secgroup"), VpcId: awssdk.String("vpc_1")},
		{GroupId: awssdk.String("secgroup_2"), VpcId: awssdk.String("vpc_1")},
	}

	subnets := []*ec2.Subnet{
		{SubnetId: awssdk.String("sub_1"), VpcId: awssdk.String("vpc_1")},
		{SubnetId: awssdk.String("sub_2"), VpcId: awssdk.String("vpc_1")},
		{SubnetId: awssdk.String("sub_3"), VpcId: awssdk.String("vpc_2")},
		{SubnetId: awssdk.String("sub_4"), VpcId: nil}, // edge case subnet with no vpc id
	}

	keypairs := []*ec2.KeyPairInfo{
		{KeyName: awssdk.String("my_key_pair")},
	}

	igws := []*ec2.InternetGateway{
		{InternetGatewayId: awssdk.String("igw_1"), Attachments: []*ec2.InternetGatewayAttachment{{VpcId: awssdk.String("vpc_2")}}},
	}

	routeTables := []*ec2.RouteTable{
		{RouteTableId: awssdk.String("rt_1"), VpcId: awssdk.String("vpc_1"), Associations: []*ec2.RouteTableAssociation{{RouteTableId: awssdk.String("rt_1"), SubnetId: awssdk.String("sub_1")}}},
	}

	//ELB
	lbPages := [][]*elbv2.LoadBalancer{
		{
			{LoadBalancerArn: awssdk.String("lb_1"), LoadBalancerName: awssdk.String("my_loadbalancer"), VpcId: awssdk.String("vpc_1")},
			{LoadBalancerArn: awssdk.String("lb_2"), VpcId: awssdk.String("vpc_2")},
		},
		{
			{LoadBalancerArn: awssdk.String("lb_3"), VpcId: awssdk.String("vpc_1"), SecurityGroups: []*string{awssdk.String("secgroup_1"), awssdk.String("secgroup_2")}},
		},
	}
	targetGroups := []*elbv2.TargetGroup{
		{TargetGroupArn: awssdk.String("tg_1"), VpcId: awssdk.String("vpc_1"), LoadBalancerArns: []*string{awssdk.String("lb_1"), awssdk.String("lb_3")}},
		{TargetGroupArn: awssdk.String("tg_2"), VpcId: awssdk.String("vpc_2"), LoadBalancerArns: []*string{awssdk.String("lb_2")}},
	}
	listeners := map[string][]*elbv2.Listener{
		"lb_1": {{ListenerArn: awssdk.String("list_1"), LoadBalancerArn: awssdk.String("lb_1")}, {ListenerArn: awssdk.String("list_1.2"), LoadBalancerArn: awssdk.String("lb_1")}},
		"lb_2": {{ListenerArn: awssdk.String("list_2"), LoadBalancerArn: awssdk.String("lb_2")}},
		"lb_3": {{ListenerArn: awssdk.String("list_3"), LoadBalancerArn: awssdk.String("lb_3")}},
	}
	targetHealths := map[string][]*elbv2.TargetHealthDescription{
		"tg_1": {{HealthCheckPort: awssdk.String("80"), Target: &elbv2.TargetDescription{Id: awssdk.String("inst_1"), Port: awssdk.Int64(443)}}},
		"tg_2": {{Target: &elbv2.TargetDescription{Id: awssdk.String("inst_2"), Port: awssdk.Int64(80)}}, {Target: &elbv2.TargetDescription{Id: awssdk.String("inst_3"), Port: awssdk.Int64(80)}}},
	}

	mock := &mockEc2{vpcs: vpcs, securityGroups: securityGroups, subnets: subnets, instances: instances, keyPairs: keypairs, internetGateways: igws, routeTables: routeTables}
	mockLb := &mockELB{loadBalancerPages: lbPages, targetGroups: targetGroups, listeners: listeners, targetHealths: targetHealths}
	infra := Infra{EC2API: mock, ELBV2API: mockLb, RDSAPI: &mockRDS{}, region: "eu-west-1"}
	InfraService = &infra

	g, err := infra.FetchResources()
	if err != nil {
		t.Fatal(err)
	}
	resources, err := g.GetAllResources("region", "instance", "vpc", "securitygroup", "subnet", "keypair", "internetgateway", "routetable", "loadbalancer", "targetgroup", "listener")
	if err != nil {
		t.Fatal(err)
	}

	// Sort slice properties in resources
	for _, res := range resources {
		if p, ok := res.Properties[p.SecurityGroups].([]string); ok {
			sort.Strings(p)
		}
		if p, ok := res.Properties[p.Vpcs].([]string); ok {
			sort.Strings(p)
		}
	}

	expected := map[string]*graph.Resource{
		"eu-west-1":   testRes("eu-west-1", "region").build(),
		"inst_1":      instRes("inst_1").prop(p.Subnet, "sub_1").prop(p.Vpc, "vpc_1").prop(p.Name, "instance1-name").build(),
		"inst_2":      instRes("inst_2").prop(p.Subnet, "sub_2").prop(p.Vpc, "vpc_1").prop(p.SecurityGroups, []string{"secgroup_1"}).build(),
		"inst_3":      instRes("inst_3").prop(p.Subnet, "sub_3").prop(p.Vpc, "vpc_2").build(),
		"inst_4":      instRes("inst_4").prop(p.Subnet, "sub_3").prop(p.Vpc, "vpc_2").prop(p.SecurityGroups, []string{"secgroup_1", "secgroup_2"}).prop(p.SSHKey, "my_key_pair").build(),
		"inst_5":      instRes("inst_5").build(),
		"vpc_1":       vpcRes("vpc_1").build(),
		"vpc_2":       vpcRes("vpc_2").build(),
		"secgroup_1":  sGrpRes("secgroup_1").prop(p.Name, "my_secgroup").prop(p.Vpc, "vpc_1").build(),
		"secgroup_2":  sGrpRes("secgroup_2").prop(p.Vpc, "vpc_1").build(),
		"sub_1":       subRes("sub_1").prop(p.Vpc, "vpc_1").build(),
		"sub_2":       subRes("sub_2").prop(p.Vpc, "vpc_1").build(),
		"sub_3":       subRes("sub_3").prop(p.Vpc, "vpc_2").build(),
		"sub_4":       subRes("sub_4").build(),
		"my_key_pair": keyRes("my_key_pair").build(),
		"igw_1":       igwRes("igw_1").prop(p.Vpcs, []string{"vpc_2"}).build(),
		"rt_1":        rtRes("rt_1").prop(p.Vpc, "vpc_1").prop(p.Main, false).build(),
		"lb_1":        lbRes("lb_1").prop(p.Name, "my_loadbalancer").prop(p.Vpc, "vpc_1").build(),
		"lb_2":        lbRes("lb_2").prop(p.Vpc, "vpc_2").build(),
		"lb_3":        lbRes("lb_3").prop(p.Vpc, "vpc_1").build(),
		"tg_1":        tgRes("tg_1").prop(p.Vpc, "vpc_1").build(),
		"tg_2":        tgRes("tg_2").prop(p.Vpc, "vpc_2").build(),
		"list_1":      listRes("list_1").prop(p.LoadBalancer, "lb_1").build(),
		"list_1.2":    listRes("list_1.2").prop(p.LoadBalancer, "lb_1").build(),
		"list_2":      listRes("list_2").prop(p.LoadBalancer, "lb_2").build(),
		"list_3":      listRes("list_3").prop(p.LoadBalancer, "lb_3").build(),
	}

	expectedChildren := map[string][]string{
		"eu-west-1": {"igw_1", "my_key_pair", "vpc_1", "vpc_2"},
		"lb_1":      {"list_1", "list_1.2"},
		"lb_2":      {"list_2"},
		"lb_3":      {"list_3"},
		"sub_1":     {"inst_1"},
		"sub_2":     {"inst_2"},
		"sub_3":     {"inst_3", "inst_4"},
		"vpc_1":     {"lb_1", "lb_3", "rt_1", "secgroup_1", "secgroup_2", "sub_1", "sub_2", "tg_1"},
		"vpc_2":     {"lb_2", "sub_3", "tg_2"},
	}

	expectedAppliedOn := map[string][]string{
		"igw_1":       {"vpc_2"},
		"lb_1":        {"tg_1"},
		"lb_2":        {"tg_2"},
		"lb_3":        {"tg_1"},
		"my_key_pair": {"inst_4"},
		"rt_1":        {"sub_1"},
		"secgroup_1":  {"inst_2", "inst_4", "lb_3"},
		"secgroup_2":  {"inst_4", "lb_3"},
		"tg_1":        {"inst_1"},
		"tg_2":        {"inst_2", "inst_3"},
	}

	compareResources(t, g, resources, expected, expectedChildren, expectedAppliedOn)
}

func TestBuildStorageRdfGraph(t *testing.T) {
	buckets := map[string][]*s3.Bucket{
		"us-west-1": {
			{Name: awssdk.String("bucket_us_1")},
			{Name: awssdk.String("bucket_us_2")},
			{Name: awssdk.String("bucket_us_3")},
		},
		"eu-west-1": {
			{Name: awssdk.String("bucket_eu_1")},
			{Name: awssdk.String("bucket_eu_2")},
		},
	}
	objects := map[string][]*s3.Object{
		"bucket_us_1": {
			{Key: awssdk.String("obj_1")},
			{Key: awssdk.String("obj_2")},
		},
		"bucket_us_2": {},
		"bucket_us_3": {
			{Key: awssdk.String("obj_3")},
		},
		"bucket_eu_1": {
			{Key: awssdk.String("obj_4")},
		},
		"bucket_eu_2": {
			{Key: awssdk.String("obj_5")},
			{Key: awssdk.String("obj_6")},
		},
	}
	bucketsACL := map[string][]*s3.Grant{
		"bucket_us_1": {
			{Permission: awssdk.String("Read"), Grantee: &s3.Grantee{ID: awssdk.String("usr_1")}},
		},
		"bucket_us_3": {
			{Permission: awssdk.String("Write"), Grantee: &s3.Grantee{ID: awssdk.String("usr_2")}},
		},
		"bucket_eu_1": {
			{Permission: awssdk.String("Write"), Grantee: &s3.Grantee{ID: awssdk.String("usr_2")}},
		},
		"bucket_eu_2": {
			{Permission: awssdk.String("Write"), Grantee: &s3.Grantee{ID: awssdk.String("usr_1")}},
		},
	}

	mocks3 := &mockS3{bucketsPerRegion: buckets, objectsPerBucket: objects, bucketsACL: bucketsACL}
	StorageService = mocks3
	storage := Storage{S3API: mocks3, region: "eu-west-1"}

	g, err := storage.FetchResources()
	if err != nil {
		t.Fatal(err)
	}
	resources, err := g.GetAllResources("region", "bucket")
	if err != nil {
		t.Fatal(err)
	}

	expected := map[string]*graph.Resource{
		"eu-west-1":   testRes("eu-west-1", "region").build(),
		"bucket_eu_1": buckRes("bucket_eu_1").prop(p.Grants, []*graph.Grant{{GranteeID: "usr_2", Permission: "Write"}}).build(),
		"bucket_eu_2": buckRes("bucket_eu_2").prop(p.Grants, []*graph.Grant{{GranteeID: "usr_1", Permission: "Write"}}).build(),
	}
	expectedChildren := map[string][]string{
		"eu-west-1":   {"bucket_eu_1", "bucket_eu_2"},
		"bucket_eu_1": {"obj_4"},
		"bucket_eu_2": {"obj_5", "obj_6"},
	}
	expectedAppliedOn := map[string][]string{}

	compareResources(t, g, resources, expected, expectedChildren, expectedAppliedOn)
}

func TestBuildDnsRdfGraph(t *testing.T) {
	zonePages := [][]*route53.HostedZone{
		{
			{Id: awssdk.String("/hostedzone/12345"), Name: awssdk.String("my.first.domain")},
			{Id: awssdk.String("/hostedzone/23456"), Name: awssdk.String("my.second.domain")},
		},
		{{Id: awssdk.String("/hostedzone/34567"), Name: awssdk.String("my.third.domain")}},
	}
	recordPages := map[string][][]*route53.ResourceRecordSet{
		"/hostedzone/12345": {
			{
				{Type: awssdk.String("A"), TTL: awssdk.Int64(10), Name: awssdk.String("subdomain1.my.first.domain"), ResourceRecords: []*route53.ResourceRecord{{Value: awssdk.String("1.2.3.4")}, {Value: awssdk.String("2.3.4.5")}}},
				{Type: awssdk.String("A"), TTL: awssdk.Int64(10), Name: awssdk.String("subdomain2.my.first.domain"), ResourceRecords: []*route53.ResourceRecord{{Value: awssdk.String("3.4.5.6")}}},
			},
			{
				{Type: awssdk.String("CNAME"), TTL: awssdk.Int64(60), Name: awssdk.String("subdomain3.my.first.domain"), ResourceRecords: []*route53.ResourceRecord{{Value: awssdk.String("4.5.6.7")}}},
			},
		},
		"/hostedzone/23456": {
			{
				{Type: awssdk.String("A"), TTL: awssdk.Int64(30), Name: awssdk.String("subdomain1.my.second.domain"), ResourceRecords: []*route53.ResourceRecord{{Value: awssdk.String("5.6.7.8")}}},
				{Type: awssdk.String("CNAME"), TTL: awssdk.Int64(10), Name: awssdk.String("subdomain3.my.second.domain"), ResourceRecords: []*route53.ResourceRecord{{Value: awssdk.String("6.7.8.9")}}},
			},
		},
	}
	mockRoute53 := &mockRoute53{zonePages: zonePages, recordPages: recordPages}

	dns := Dns{Route53API: mockRoute53, region: "eu-west-1"}

	g, err := dns.FetchResources()
	if err != nil {
		t.Fatal(err)
	}

	resources, err := g.GetAllResources("zone", "record")
	if err != nil {
		t.Fatal(err)
	}
	// Sort slice properties in resources
	for _, res := range resources {
		if p, ok := res.Properties[p.Records].([]string); ok {
			sort.Strings(p)
		}
	}

	expected := map[string]*graph.Resource{
		"/hostedzone/12345": zoneRes("/hostedzone/12345").prop(p.Name, "my.first.domain").build(),
		"/hostedzone/23456": zoneRes("/hostedzone/23456").prop(p.Name, "my.second.domain").build(),
		"/hostedzone/34567": zoneRes("/hostedzone/34567").prop(p.Name, "my.third.domain").build(),
		"awls-91fa0a45":     recRes("awls-91fa0a45").prop(p.Name, "subdomain1.my.first.domain").prop(p.Type, "A").prop(p.TTL, 10).prop(p.Records, []string{"1.2.3.4", "2.3.4.5"}).build(),
		"awls-920c0a46":     recRes("awls-920c0a46").prop(p.Name, "subdomain2.my.first.domain").prop(p.Type, "A").prop(p.TTL, 10).prop(p.Records, []string{"3.4.5.6"}).build(),
		"awls-be1e0b6a":     recRes("awls-be1e0b6a").prop(p.Name, "subdomain3.my.first.domain").prop(p.Type, "CNAME").prop(p.TTL, 60).prop(p.Records, []string{"4.5.6.7"}).build(),
		"awls-9c420a99":     recRes("awls-9c420a99").prop(p.Name, "subdomain1.my.second.domain").prop(p.Type, "A").prop(p.TTL, 30).prop(p.Records, []string{"5.6.7.8"}).build(),
		"awls-c9b80bbe":     recRes("awls-c9b80bbe").prop(p.Name, "subdomain3.my.second.domain").prop(p.Type, "CNAME").prop(p.TTL, 10).prop(p.Records, []string{"6.7.8.9"}).build(),
	}
	expectedChildren := map[string][]string{
		"/hostedzone/12345": {"awls-91fa0a45", "awls-920c0a46", "awls-be1e0b6a"},
		"/hostedzone/23456": {"awls-9c420a99", "awls-c9b80bbe"},
	}
	expectedAppliedOn := map[string][]string{}

	compareResources(t, g, resources, expected, expectedChildren, expectedAppliedOn)
}

func TestBuildEmptyRdfGraphWhenNoData(t *testing.T) {
	expect := `/node<eu-west-1>	"rdf:type"@[]	/node<cloud-owl:Region>`
	access := Access{IAMAPI: &mockIam{}, region: "eu-west-1"}

	g, err := access.FetchResources()
	if err != nil {
		t.Fatal(err)
	}

	result := g.MustMarshal()
	if result != expect {
		t.Fatalf("got [%s]\nwant [%s]", result, expect)
	}

	infra := Infra{EC2API: &mockEc2{}, ELBV2API: &mockELB{}, RDSAPI: &mockRDS{}, region: "eu-west-1"}

	g, err = infra.FetchResources()
	if err != nil {
		t.Fatal(err)
	}

	result = g.MustMarshal()
	if result != expect {
		t.Fatalf("got [%s]\nwant [%s]", result, expect)
	}
}

func diffText(actual, expected string) error {
	actuals := strings.Split(actual, "\n")
	expecteds := strings.Split(expected, "\n")

	if len(actuals) != len(expecteds) {
		return fmt.Errorf("text diff: not same number of lines:\ngot \n%s\n\nwant\n%s\n", actual, expected)
	}

	for i := 0; i < len(actuals); i++ {
		if actuals[i] != expecteds[i] {
			return fmt.Errorf("text diff: diff at line %d\ngot:\n%q\nwant:\n%q", i+1, actuals[i], expecteds[i])
		}
	}

	return nil
}

func mustGetChildrenId(g *graph.Graph, res *graph.Resource) []string {
	var collect []string

	err := g.Accept(&graph.ChildrenVisitor{From: res, IncludeFrom: false, Each: func(res *graph.Resource, depth int) error {
		if depth == 1 {
			collect = append(collect, res.Id())
		}
		return nil
	}})
	if err != nil {
		panic(err)
	}
	return collect
}

func mustGetAppliedOnId(g *graph.Graph, res *graph.Resource) []string {
	resources, err := g.ListResourcesAppliedOn(res)
	if err != nil {
		panic(err)
	}
	var ids []string
	for _, r := range resources {
		ids = append(ids, r.Id())
	}
	return ids
}

func compareResources(t *testing.T, g *graph.Graph, resources []*graph.Resource, expected map[string]*graph.Resource, expectedChildren, expectedAppliedOn map[string][]string) {
	if got, want := len(resources), len(expected); got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	for _, got := range resources {
		want := expected[got.Id()]
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got \n%#v\nwant\n%#v", got, want)
		}
		children := mustGetChildrenId(g, got)
		sort.Strings(children)
		if g, w := children, expectedChildren[got.Id()]; !reflect.DeepEqual(g, w) {
			t.Fatalf("children: got %v, want %v", g, w)
		}
		appliedOn := mustGetAppliedOnId(g, got)
		sort.Strings(appliedOn)
		if g, w := appliedOn, expectedAppliedOn[got.Id()]; !reflect.DeepEqual(g, w) {
			t.Fatalf("appliedOn: got %v, want %v", g, w)
		}
	}
}
