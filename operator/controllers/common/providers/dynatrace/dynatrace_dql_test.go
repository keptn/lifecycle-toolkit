package dynatrace

import (
	"testing"
)

const dqlPayload = "{\"records\":[{\"value\":{\"count\":1,\"sum\":36.78336878333334,\"min\":36.78336878333334,\"avg\":36.78336878333334,\"max\":36.78336878333334},\"metric.key\":\"dt.containers.cpu.usage_user_milli_cores\",\"timeframe\":{\"start\":\"2023-01-31T09:11:00.000Z\",\"end\":\"2023-01-31T09:12:00.000Z\"},\"Container\":\"frontend\",\"host.name\":\"default-pool-349eb8c6-gccf\",\"k8s.namespace.name\":\"hipstershop\",\"k8s.pod.uid\":\"632df64d-474c-4410-968d-666f639ad358\"}],\"types\":[{\"mappings\":{\"value\":{\"type\":\"summary_stats\"},\"metric.key\":{\"type\":\"string\"},\"timeframe\":{\"type\":\"timeframe\"},\"Container\":{\"type\":\"string\"},\"host.name\":{\"type\":\"string\"},\"k8s.namespace.name\":{\"type\":\"string\"},\"k8s.pod.uid\":{\"type\":\"string\"}},\"indexRange\":[0,1]}]}"
const dqlPayloadTooManyItems = "{\"records\":[{\"value\":{\"count\":1,\"sum\":6.293549483333334,\"min\":6.293549483333334,\"avg\":6.293549483333334,\"max\":6.293549483333334},\"metric.key\":\"dt.containers.cpu.usage_user_milli_cores\",\"timeframe\":{\"start\":\"2023-01-31T09:07:00.000Z\",\"end\":\"2023-01-31T09:08:00.000Z\"},\"Container\":\"loginservice\",\"host.name\":\"default-pool-349eb8c6-gccf\",\"k8s.namespace.name\":\"easytrade\",\"k8s.pod.uid\":\"fc084e57-11a0-4a95-b8a0-76191c31d839\"},{\"value\":{\"count\":1,\"sum\":1.0421756,\"min\":1.0421756,\"avg\":1.0421756,\"max\":1.0421756},\"metric.key\":\"dt.containers.cpu.usage_user_milli_cores\",\"timeframe\":{\"start\":\"2023-01-31T09:07:00.000Z\",\"end\":\"2023-01-31T09:08:00.000Z\"},\"Container\":\"frontendreverseproxy\",\"host.name\":\"default-pool-349eb8c6-gccf\",\"k8s.namespace.name\":\"easytrade\",\"k8s.pod.uid\":\"41b5d6e0-98fc-4dce-a1b4-bb269a03d72b\"},{\"value\":{\"count\":1,\"sum\":6.3881383000000005,\"min\":6.3881383000000005,\"avg\":6.3881383000000005,\"max\":6.3881383000000005},\"metric.key\":\"dt.containers.cpu.usage_user_milli_cores\",\"timeframe\":{\"start\":\"2023-01-31T09:07:00.000Z\",\"end\":\"2023-01-31T09:08:00.000Z\"},\"Container\":\"shippingservice\",\"host.name\":\"default-pool-349eb8c6-gccf\",\"k8s.namespace.name\":\"hipstershop\",\"k8s.pod.uid\":\"96fcf9d7-748a-47f7-b1b3-ca6427e20edd\"}],\"types\":[{\"mappings\":{\"value\":{\"type\":\"summary_stats\"},\"metric.key\":{\"type\":\"string\"},\"timeframe\":{\"type\":\"timeframe\"},\"Container\":{\"type\":\"string\"},\"host.name\":{\"type\":\"string\"},\"k8s.namespace.name\":{\"type\":\"string\"},\"k8s.pod.uid\":{\"type\":\"string\"}},\"indexRange\":[0,3]}]}"

func TestGetDQL(t *testing.T) {

}

func TestDQL_TooManyItems(t *testing.T) {

}
