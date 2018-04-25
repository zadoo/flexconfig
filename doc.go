/*
Package flexconfig provides a uniform interface to retrieve configuration
properties, independent of how or where the properties are specified. Files,
environment variables, command line arguments, and a configuration store can
be used alone or in combination with each other to define the configuration for
the application.

An application can retrieve and set property values using a key with a
canonical form. The canonical form is all lowercase with hierarchical fields
separated by dots. For example,
    animal.bear.polar.habitat

Initialization of a configuration is done by an application near the
beginning of its main function by calling NewFlexibleConfiguration. The
application can specify several parameters about how and where configuration
properties will be obtained from. Configuration sources include (in priority order, lowest to highest):
    - directories on the local file system
    - environment variables
    - command line arguments
    - configuration store

The static configuration sources (not including configuration store) are read
from lowest priority to highest priority, creating configuration properties
from data found in the sources. If the same canonical property key is set in
a data source read later, its value will override a value from a data source
read earlier.

Configuration directory names are derived from the ApplicationName in the
ConfigurationParameters. A non-empty ApplicationName indicates that file-based
properties will be searched for. Multiple directories are searched as described
under ConfigurationParameters. The subset of files read in these directories
is controlled by AcceptedFileSuffixes, where the suffix ".conf" is used if
none are specified. The contents of the files may have formats that include
JSON, YAML, and INI.

Hierarchical properties (multiple fields separated by dots) are defined by
parsing JSON and YAML files. Arrays defined in these files result in
property names that include fields consisting of digits. For example, the
YAML file:

    myapp:
      plugins:
        - name: foo
          loglevel: debug
        - name: bar
          server:
            address: 192.168.1.1

will result in the following properties in the configuration:

    myapp.plugins.0.name
    myapp.plugins.0.loglevel
    myapp.plugins.1.name
    myapp.plugins.1.server.address

Environment variables will be searched if EnvironmentVariablePrefixes
includes non-empty members. The members specify prefixes for environment
variable names. For instance, specifying a prefix of "SUN" will match both

    SUNSHINE_DAILY
    SUN_MICROSYSTEMS

whereas, specifying "SUN_" will only match

    SUN_MICROSYSTEMS

Environment variable names are converted into the canonical form before
storing in the configuration.

Command line arguments are checked for property definitions without the
application needing to manage arguments beyond calling NewFlexibleConfiguration.
Any argument beginning with a double dash (--) and being all lowercase is used
to create configuration properties. Arguments are checked up to the point where
all arguments have been read, or an argument consisting entirely of "--" is
found. The argument name matches the canonical form for names while the
value can be anything. For example, given the following commandline:
    hello -v --smiley.face=true -n --happy.day -- --dont.check.sig=true "Hi!"
will result in the following properties being defined in the configuration:
    smiley.face: true
Note, the '=' separating name and value is required with no intermediate spaces.
Single or double quotes can enclose the value of a property.

Accessing the configuration to obtain property values will consult the
configuration store first, if it has been configured. If this results in an
error or an empty value, the in-memory configuration read from files, env vars,
and the command line, is consulted.
*/
package flexconfig
