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
	"errors"
)

// FlexConfigStoreType is an enumerated type defining the type of a
// FlexConfigStore.
type FlexConfigStoreType int

const (
	// FlexConfigStoreUnknown is a value of FlexConfigStoreType
	// indicating the type of configuration store is unknown.
	FlexConfigStoreUnknown FlexConfigStoreType = iota

	// FlexConfigStoreEtcd is a value of FlexConfigStoreType
	// indicating the type of configuration store is etcd.
	FlexConfigStoreEtcd
)

var (
	// ErrStoreUnsupportedType indicates that the value of
	// FlexConfigStoreType is not supported.
	ErrStoreUnsupportedType = errors.New("Unsupported FlexConfigStoreType")

	// ErrStoreKeyRequired indicates that the requested function requires a
	// property key as a parameter.
	ErrStoreKeyRequired = errors.New("Key is required")
)

// FlexConfigStore describes the interface to a flexible configuration store.
// A Config optionally uses a FlexConfigStore to provide dynamic configuration
// properties.
type FlexConfigStore interface {
	// Get returns a single configuration property value given the key
	// for the property.
	Get(key string) (string, error)

	// GetAll returns all properties in the FlexConfigStore having the
	// prefix used by the store.
	GetAll() ([]KeyValue, error)

	// Set creates or modifies the property indicated by key with the
	// specified value.
	Set(key, val string) error

	// Delete removes the specified property.
	Delete(key string) error

	// GetPrefix returns the "namespace" prefix specified when the
	// FlexConfigStore was created by calling NewFlexConfigStore.
	GetPrefix() string
}

// KeyValue describes a property key and value.
type KeyValue struct {
	Key   string
	Value string
}

// NewFlexConfigStore creates a FlexConfigStore instance for the specified
// store type. Supported store types include: etcd. The store type may
// require passing zero or more endpoints for the store in order to instantiate
// it, as well as a prefix used to differentiate from other uses of the store.
func NewFlexConfigStore(
	storeType FlexConfigStoreType,
	endpoints []string,
	prefix string) (FlexConfigStore, error) {
	switch storeType {
	case FlexConfigStoreEtcd:
		return newEtcdFlexConfigStore(endpoints, prefix)
	default:
		return nil, ErrStoreUnsupportedType
	}
}

// String returns the string representation of the FlexConfigStoreType.
func (fcst FlexConfigStoreType) String() string {
	switch fcst {
	case FlexConfigStoreEtcd:
		return "etcd"
	default:
		return "unknown"
	}
}
