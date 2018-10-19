package flexconfig

/*
Copyright 2018 The flexconfig Authors

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

import (
	"os"
	"strings"
	"testing"
)

const (
	etcdEndpointEnvironmentVariable = "ETCDCTL_ENDPOINTS"
	defaultEtcdEndpoint             = "http://127.0.0.1:2379"
	etcdTestPrefix                  = "/test"
)

func getEndpointList() []string {
	endpointstr := os.Getenv(etcdEndpointEnvironmentVariable)
	if len(endpointstr) == 0 {
		endpointstr = defaultEtcdEndpoint
	}

	endpoints := strings.Split(endpointstr, ",")

	return endpoints
}

func Test_etcd(t *testing.T) {
	fcs, err := newEtcdFlexConfigStore(getEndpointList(), etcdTestPrefix)
	if err != nil {
		t.Errorf("Error creating store, have you defined ETCDCTL_ENDPOINTS?: %v", err)
		return
	}

	prefix := fcs.GetPrefix()
	if prefix != etcdTestPrefix {
		t.Errorf("Prefix does not match: %s", prefix)
	}

	prop := "moo"
	propval := "bar"

	err = fcs.Set(prop, propval)
	if err != nil {
		t.Errorf("Error setting property: %v", err)
	}

	val, err := fcs.Get(prop)
	if err != nil {
		t.Errorf("Error getting property: %v", err)
	}

	if val != propval {
		t.Errorf("Did not get back same value as set: %v", val)
	}
}

func Test_etcd_getAll(t *testing.T) {
	fcs, err := newEtcdFlexConfigStore(getEndpointList(), etcdTestPrefix)
	if err != nil {
		t.Errorf("Error creating store: %v", err)
		return
	}

	prop := "moo"
	propval := "bar"

	err = fcs.Set(prop, propval)
	if err != nil {
		t.Errorf("Error setting property: %v", err)
	}

	prop = "spiffy.lake"
	propval = "123"

	err = fcs.Set(prop, propval)
	if err != nil {
		t.Errorf("Error setting property: %v", err)
	}

	list, err := fcs.GetAll()
	if err != nil {
		t.Errorf("Error getting all: %v", err)
	}

	if len(list) != 2 {
		t.Errorf("Unexpected number of GetAll result: %d", len(list))
	}

	if list[0].Key != "moo" && list[0].Key != prop {
		t.Errorf("Unexpected key: %s", list[0].Key)
	}

	if list[1].Key != "moo" && list[1].Key != prop {
		t.Errorf("Unexpected key: %s", list[1].Key)
	}

	err = fcs.Delete("moo")
	if err != nil {
		t.Errorf("Error deleting property: %v", err)
	}

	val, err := fcs.Get("moo")
	if err != nil {
		t.Errorf("Error getting property: %v", err)
	}

	if val != "" {
		t.Errorf("Expecting empty value but found: %v", val)
	}

	err = fcs.Delete(prop)
	if err != nil {
		t.Errorf("Error deleting property: %v", err)
	}
}

func Test_etcd_prefix(t *testing.T) {
	endpointstr := os.Getenv(etcdEndpointEnvironmentVariable)
	if len(endpointstr) == 0 {
		endpointstr = defaultEtcdEndpoint
	}

	endpoints := strings.Split(endpointstr, ",")

	fcs, err := newEtcdFlexConfigStore(endpoints, "test/")
	if err != nil {
		t.Errorf("Error creating store: %v", err)
		return
	}

	prefix := fcs.GetPrefix()
	if prefix != etcdTestPrefix {
		t.Errorf("Prefix does not match: %s", prefix)
	}
}

func Test_etcd_badNames(t *testing.T) {
	endpointstr := os.Getenv(etcdEndpointEnvironmentVariable)
	if len(endpointstr) == 0 {
		endpointstr = defaultEtcdEndpoint
	}

	endpoints := strings.Split(endpointstr, ",")

	fcs, err := newEtcdFlexConfigStore(endpoints, "")
	if err != nil {
		t.Errorf("Error creating store: %v", err)
		return
	}

	_, err = fcs.Get("")
	if err == nil {
		t.Errorf("Unexpected success getting empty property name")
	}

	if err.Error() != "Key is required" {
		t.Errorf("Unexpected error string: %s", err.Error())
	}

	_, err = fcs.GetAll()
	if err != nil {
		t.Errorf("Error calling GetAll: %v", err)
	}

	err = fcs.Set("", "value")
	if err == nil {
		t.Errorf("Unexpected success setting empty property name")
	}

	if err.Error() != "Key is required" {
		t.Errorf("Unexpected error string: %s", err.Error())
	}

	err = fcs.Delete("")
	if err == nil {
		t.Errorf("Unexpected success deleting empty property name")
	}
}
