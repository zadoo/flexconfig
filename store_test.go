package flexconfig

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func Test_flexConfigStoreNew(t *testing.T) {
	endpointstr := os.Getenv(etcdEndpointEnvironmentVariable)
	if len(endpointstr) == 0 {
		endpointstr = defaultEtcdEndpoint
	}

	endpoints := strings.Split(endpointstr, ",")

	fcs, err := NewFlexConfigStore(FlexConfigStoreEtcd,
		endpoints, etcdTestPrefix)
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

	val, err := fcs.Get(prop)
	if err != nil {
		t.Errorf("Error getting property: %v", err)
	}

	if val != propval {
		t.Errorf("Did not get back same value as set: %v", val)
	}

	err = fcs.Delete(prop)
	if err != nil {
		t.Errorf("Error deleting property: %v", err)
	}
}

func Test_badStoreType(t *testing.T) {
	endpointstr := os.Getenv(etcdEndpointEnvironmentVariable)
	if len(endpointstr) == 0 {
		endpointstr = defaultEtcdEndpoint
	}

	endpoints := strings.Split(endpointstr, ",")

	_, err := NewFlexConfigStore(999,
		endpoints, etcdTestPrefix)
	if err == nil {
		t.Errorf("Unexpected success creating store")
	}
}

func Test_typeString(t *testing.T) {
	str := fmt.Sprintf("%s", FlexConfigStoreEtcd)
	if str != "etcd" {
		t.Errorf("Unexpected string for known type: %s", str)
	}

	str = fmt.Sprintf("%s", FlexConfigStoreUnknown)
	if str != "unknown" {
		t.Errorf("Unexpected string for known type: %s", str)
	}

	str = fmt.Sprintf("%s", FlexConfigStoreType(42))
	if str != "unknown" {
		t.Errorf("Unexpected string for known type: %s", str)
	}
}
