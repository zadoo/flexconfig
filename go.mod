module github.com/zadoo/flexconfig

go 1.13

require (
	github.com/coreos/etcd v3.3.18+incompatible // indirect
	github.com/coreos/go-systemd v0.0.0-00010101000000-000000000000 // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/google/uuid v1.1.1 // indirect
	go.etcd.io/etcd v3.3.18+incompatible
	go.uber.org/zap v1.13.0 // indirect
	google.golang.org/grpc v1.26.0 // indirect
	gopkg.in/ini.v1 v1.51.1
	gopkg.in/yaml.v2 v2.2.7
)

replace github.com/coreos/go-systemd => ./go-systemd
