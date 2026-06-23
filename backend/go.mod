module go-kratos-mcp-demo

go 1.26.1

replace (
	github.com/armon/go-metrics => github.com/hashicorp/go-metrics v0.4.1
	github.com/bufbuild/protovalidate-go => buf.build/go/protovalidate v0.10.1
)

require (
	connectrpc.com/connect v1.20.0
	github.com/go-kratos/kratos/v2 v2.9.2
	github.com/google/uuid v1.6.0
	github.com/google/wire v0.7.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.29.0
	github.com/leafspace/protoc-gen-go-mcp v0.0.0-20260623072059-452d056ae213
	github.com/mark3labs/mcp-go v0.54.1
	github.com/stretchr/testify v1.11.1
	github.com/tx7do/kratos-swagger-ui v0.0.0-20250528131001-09c0dbdb208d
	github.com/tx7do/kratos-transport/transport/mcp v1.3.3
	golang.org/x/net v0.54.0
	google.golang.org/genproto/googleapis/api v0.0.0-20260427160629-7cedc36a6bc4
	google.golang.org/grpc v1.81.0
	google.golang.org/protobuf v1.36.11
)

require (
	buf.build/gen/go/redpandadata/common/protocolbuffers/go v1.34.2-20240917150400-3f349e63f44a.2 // indirect
	dario.cat/mergo v1.0.2 // indirect
	github.com/cenkalti/backoff/v5 v5.0.3 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/go-kratos/aegis v0.2.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-playground/form/v4 v4.3.0 // indirect
	github.com/google/jsonschema-go v0.4.2 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/redpanda-data/common-go/api v0.0.0-20250801174835-9eea07f1ea06 // indirect
	github.com/santhosh-tekuri/jsonschema/v6 v6.0.2 // indirect
	github.com/spf13/cast v1.10.0 // indirect
	github.com/swaggest/swgui v1.8.4 // indirect
	github.com/tx7do/kratos-transport/broker v1.3.3 // indirect
	github.com/tx7do/kratos-transport/tracing v1.1.2 // indirect
	github.com/tx7do/kratos-transport/transport v1.3.4 // indirect
	github.com/tx7do/kratos-transport/transport/keepalive v1.3.4 // indirect
	github.com/vearutop/statigz v1.5.0 // indirect
	github.com/yosida95/uritemplate/v3 v3.0.2 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel v1.43.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.43.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.43.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.43.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.43.0 // indirect
	go.opentelemetry.io/otel/metric v1.43.0 // indirect
	go.opentelemetry.io/otel/sdk v1.43.0 // indirect
	go.opentelemetry.io/otel/trace v1.43.0 // indirect
	go.opentelemetry.io/proto/otlp v1.10.0 // indirect
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/sys v0.44.0 // indirect
	golang.org/x/text v0.37.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260427160629-7cedc36a6bc4 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
