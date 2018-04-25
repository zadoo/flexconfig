# Go Flexible Configuration

Go package providing a uniform interface to configuration, independent of how
or where that configuration is defined.

https://github.com/zadoo/flexconfig

## Installation

Install:

```shell
go get gopkg.in/yaml.v2 gopkg.in/ini.v1 github.com/coreos/etcd/clientv3
go get github.com/zadoo/flexconfig
```

Import:

```go
import "github.com/zadoo/flexconfig"
```

Test:

```shell
ETCDCTL_ENDPOINTS=0.0.0.0:2379 go test
```

Change 0.0.0.0 to the IP address of your ectd server, or expect 5 test failures
if you do not intend to use an etcd configuration store.

## Quickstart

```go
package main

import (
	"fmt"
	"log"

	"github.com/zadoo/flexconfig"
)

func main() {
	// Connect to the etcd running on the local machine and use a prefix
	// of "/example/" for all property keys accessed by this program.
	cfgstore, err := flexconfig.NewFlexConfigStore(
		flexconfig.FlexConfigStoreEtcd,
		[]string{"127.0.0.1:2379"},
		"/example")
	if err != nil {
		log.Print("Failed connecting to config store: ", err)
	}

	// Create Config. Search for config files for the application "example".
	// Convert all environment variables beginning with "EXAMPLE_" to
	// properties. Use a configuration store. Also, all command line
	// variables of the form --<name>=<value> set properties (e.g.
	// --example.log.filepath=/var/log/myapp.log).
	cfg, err := flexconfig.NewFlexibleConfiguration(
		flexconfig.ConfigurationParameters{
			ApplicationName:             "example",
			EnvironmentVariablePrefixes: []string{"EXAMPLE_"},
			ConfigurationStore:          cfgstore,
		})
	if err != nil {
		log.Fatal("Failed creating configuration: ", err)
	}

	logPath := cfg.Get("log.path")
	fmt.Printf("%s = %s\n", "log.path", logPath)

	fooPluginName := cfg.Get("example.plugins.foo.name")
	fmt.Printf("%s = %s\n", "example.plugins.foo.name", fooPluginName)
}
```

The program can be invoked as follows (assuming it is called main.go):

```shell
EXAMPLE_PLUGINS_FOO_NAME=univeral.xyz go run main.go --log.path=/var/xyz/www.log
```

and it will print:

```shell
log.path = /var/xyz/www.log
example.plugins.foo.name = universal.xyz
```

Alternatively, create a directory named ".example" (this is the application
name passed to NewFlexibleConfiguration), and add a file having a suffix of
".conf" (e.g. config.conf) to the directory with the following contents:

```
log:
  path: /var/xyz/www.log
example:
  plugins:
    foo:
      name: universal.xyz
```

and invoke it with no arguments or environment variables:

```shell
go run main.go
```

Experiment with a combination of file configuration values, environment
variables and command line arguments to see the priority each source has
in defining the configuration. Also, omit the application name, or array
of environment variable prefixes to control which sources contribute to the
final configuration.

If no instance of etcd is available, remove the code to connect to the
FlexConfigStore and delete the line setting ConfigurationStore in the request
to create a new Config.

## Copyright

Copyright (C) 2018 The flexconfig authors.

flexconfig package released under Apache 2.0 License.
See [LICENSE](https://github.com/zadoo/flexconfig/blob/master/LICENSE) for details.
