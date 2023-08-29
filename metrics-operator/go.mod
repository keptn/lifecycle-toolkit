module github.com/keptn/lifecycle-toolkit/metrics-operator

go 1.20

require (
	github.com/DataDog/datadog-api-client-go/v2 v2.15.0
	github.com/benbjohnson/clock v1.3.5
	github.com/go-logr/logr v1.2.4
	github.com/gorilla/mux v1.8.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/keptn/lifecycle-toolkit/klt-cert-manager v0.0.0-20230807110601-88a54f97c157
	github.com/open-feature/go-sdk v1.7.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.16.0
	github.com/prometheus/common v0.44.0
	github.com/stretchr/testify v1.8.4
	gopkg.in/inf.v0 v0.9.1
	k8s.io/api v0.26.7
	k8s.io/apiextensions-apiserver v0.26.7
	k8s.io/apimachinery v0.26.7
	k8s.io/apiserver v0.26.7
	k8s.io/client-go v0.26.7
	k8s.io/component-base v0.26.7
	k8s.io/klog/v2 v2.100.1
	k8s.io/metrics v0.26.7
	sigs.k8s.io/controller-runtime v0.14.6
	sigs.k8s.io/custom-metrics-apiserver v1.25.1-0.20230116101851-63817c8ac8f2
	sigs.k8s.io/yaml v1.3.0
)

require (
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	github.com/DataDog/zstd v1.5.2 // indirect
	github.com/NYTimes/gziphandler v1.1.1 // indirect
	github.com/antlr/antlr4/runtime/Go/antlr v1.4.10 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/blang/semver/v4 v4.0.0 // indirect
	github.com/cenkalti/backoff/v4 v4.2.0 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd/v22 v22.3.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/emicklei/go-restful/v3 v3.10.1 // indirect
	github.com/evanphx/json-patch v4.12.0+incompatible // indirect
	github.com/evanphx/json-patch/v5 v5.6.0 // indirect
	github.com/felixge/httpsnoop v1.0.3 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-logr/zapr v1.2.3 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.20.0 // indirect
	github.com/go-openapi/swag v0.22.3 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/cel-go v0.12.6 // indirect
	github.com/google/gnostic v0.6.9 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.7.0 // indirect
	github.com/imdario/mergo v0.3.13 // indirect
	github.com/inconshreveable/mousetrap v1.0.1 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/onsi/ginkgo/v2 v2.9.7 // indirect
	github.com/onsi/gomega v1.27.8 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.4.0 // indirect
	github.com/prometheus/procfs v0.10.1 // indirect
	github.com/spf13/afero v1.9.5 // indirect
	github.com/spf13/cobra v1.6.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stoewer/go-strcase v1.2.0 // indirect
	go.etcd.io/etcd/api/v3 v3.5.5 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.5.5 // indirect
	go.etcd.io/etcd/client/v3 v3.5.5 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.35.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.35.0 // indirect
	go.opentelemetry.io/otel v1.11.2 // indirect
	go.opentelemetry.io/otel/exporters/otlp/internal/retry v1.11.2 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.11.2 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.11.2 // indirect
	go.opentelemetry.io/otel/metric v0.34.0 // indirect
	go.opentelemetry.io/otel/sdk v1.11.2 // indirect
	go.opentelemetry.io/otel/trace v1.11.2 // indirect
	go.opentelemetry.io/proto/otlp v0.19.0 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	go.uber.org/zap v1.24.0 // indirect
	golang.org/x/crypto v0.12.0 // indirect
	golang.org/x/exp v0.0.0-20230811145659-89c5cff77bcb // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/oauth2 v0.8.0 // indirect
	golang.org/x/sync v0.2.0 // indirect
	golang.org/x/sys v0.11.0 // indirect
	golang.org/x/term v0.11.0 // indirect
	golang.org/x/text v0.12.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	gomodules.xyz/jsonpatch/v2 v2.2.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20221202195650-67e5cbc046fd // indirect
	google.golang.org/grpc v1.52.3 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/kms v0.26.7 // indirect
	k8s.io/kube-openapi v0.0.0-20221012153701-172d655c2280 // indirect
	k8s.io/utils v0.0.0-20221128185143-99ec85e7a448 // indirect
	sigs.k8s.io/apiserver-network-proxy/konnectivity-client v0.0.37 // indirect
	sigs.k8s.io/json v0.0.0-20220713155537-f223a00ba0e2 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
)

replace golang.org/x/net => golang.org/x/net v0.14.0
