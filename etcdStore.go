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
	"strings"
	"time"

	etcd "go.etcd.io/etcd/clientv3"
	//etcd "github.com/coreos/etcd/clientv3"
	//"golang.org/x/net/context"
	"context"
)

const (
	etcdRequestTimeoutMs = 1000
)

// etcdStruct provides a handle for an instance of an etcd flexible
// configuration store.
type etcdStruct struct {
	client *etcd.Client
	prefix string
}

// newEtcdFlexConfigStore creates a new FlexConfigStore with the specified
// endpoints for etcd, and the specified prefix indicating a namespace for
// this use of etcd, so etcd can be used for other purposes without having
// property keys clash.
//
// etcd, by default, uses '/' as a hierarchical separator. The property keys
// specified in Config using '.' as a separator are translated to the etcd
// model.
//
// Example: The configuration store is created using a prefix of "example".
// Calling the Set() method with the property key "log.filepath" will result
// in setting the etcd property "/example/log/filepath". Calling the Get()
// method returns the value for the same etcd property.
func newEtcdFlexConfigStore(
	endpoints []string,
	prefix string) (FlexConfigStore, error) {
	// Create the configuration for etcd
	etcdConfig := etcd.Config{
		Endpoints: endpoints,
		DialTimeout: time.Duration(etcdRequestTimeoutMs) *
			time.Millisecond,
	}

	client, err := etcd.New(etcdConfig)
	if err != nil {
		return nil, err
	}

	// Normalize the value of prefix to start with a slash ('/') and
	// to not have a slash at the end.
	if len(prefix) > 0 {
		if !strings.HasPrefix(prefix, "/") {
			prefix = "/" + prefix
		}

		if strings.HasSuffix(prefix, "/") {
			prefix = prefix[:len(prefix)-1]
		}
	}

	fcs := new(etcdStruct)
	fcs.client = client
	fcs.prefix = prefix

	return fcs, nil
}

// Get returns a property value from etcd given its key. The received key is
// translated to use slashes instead of dots as field separators and the prefix
// specified in the call to newEtcdFlexConfigStore is prepended.
func (fcs *etcdStruct) Get(key string) (string, error) {
	if len(key) == 0 {
		return "", ErrStoreKeyRequired
	}

	resp, err := fcs.client.Get(context.Background(),
		fcs.prefix+dotsToSlash(key))
	if err != nil {
		return "", err
	}

	// Not finding a key is not an error
	if len(resp.Kvs) == 0 {
		return "", nil
	}

	return string(resp.Kvs[0].Value), nil
}

// GetAll returns all properties having the prefix specified for this instance
// of the store. Returned keys will not include the prefix.
func (fcs *etcdStruct) GetAll() ([]KeyValue, error) {
	prefix := fcs.prefix
	if len(prefix) == 0 {
		prefix = "/"
	}

	resp, err := fcs.client.Get(context.Background(),
		prefix, etcd.WithPrefix())
	if err != nil {
		return nil, err
	}

	// Iterate through the response creating a KeyValue structure for
	// every property found. The returned keys will have the prefix
	// stripped if a prefix was set when the store was created.
	var result []KeyValue
	for _, n := range resp.Kvs {
		key := string(n.Key)
		if strings.HasPrefix(key, prefix) {
			key = key[len(prefix):]
		}

		if strings.HasPrefix(key, "/") {
			key = key[1:]
		}

		result = append(result,
			KeyValue{Key: slashToDots(key), Value: string(n.Value)})
	}

	return result, nil
}

// Set creates or modifies the property indicated by key with the specified
// value. The received key is translated to use slashes instead of dots as
// field separators and the prefix specified in the call to
// newEtcdFlexConfigStore is prepended.
func (fcs *etcdStruct) Set(key, val string) error {
	if len(key) == 0 {
		return ErrStoreKeyRequired
	}

	// counting on dotsToSlash to add initial '/' if necessary
	key = dotsToSlash(key)

	_, err := fcs.client.Put(context.Background(), fcs.prefix+key, val)
	if err != nil {
		return err
	}

	return nil
}

// Delete removes a property from the store. The received key is translated to
// use slashes instead of dots as field separators and the prefix specified in
// the call to newEtcdFlexConfigStore is prepended.
func (fcs *etcdStruct) Delete(key string) error {
	if len(key) == 0 {
		return ErrStoreKeyRequired
	}

	// counting on dotsToSlash to add initial '/' if necessary
	key = dotsToSlash(key)

	_, err := fcs.client.Delete(context.Background(),
		fcs.prefix+key, etcd.WithPrefix())
	if err != nil {
		return err
	}

	return nil
}

// GetPrefix returns the "namespace" prefix specified when the FlexConfigStore
// was created by calling NewFlexConfigStore.
func (fcs *etcdStruct) GetPrefix() string {
	return fcs.prefix
}
