# Changelog

## [0.8.3](https://github.com/keptn/lifecycle-toolkit/compare/lifecycle-operator-v0.8.2...lifecycle-operator-v0.8.3) (2023-10-31)


### Features

* adapt code to use KeptnWorkloadVersion instead of KeptnWorkloadInstance ([#2255](https://github.com/keptn/lifecycle-toolkit/issues/2255)) ([c06fae1](https://github.com/keptn/lifecycle-toolkit/commit/c06fae13daa2aa98a3daf71abafe0e8ce4e5f4a3))
* add test and lint cmd to makefiles ([#2176](https://github.com/keptn/lifecycle-toolkit/issues/2176)) ([c55e0a9](https://github.com/keptn/lifecycle-toolkit/commit/c55e0a9f368c82ad3032eb676edd59e68b29fad6))
* create new Keptn umbrella Helm chart ([#2214](https://github.com/keptn/lifecycle-toolkit/issues/2214)) ([41bd47b](https://github.com/keptn/lifecycle-toolkit/commit/41bd47b7748c4d645243a4dae165651bbfd3533f))
* generalize helm chart ([#2282](https://github.com/keptn/lifecycle-toolkit/issues/2282)) ([81334eb](https://github.com/keptn/lifecycle-toolkit/commit/81334ebec4d8afda27902b6e854c4c637a3daa87))
* **lifecycle-operator:** add helm chart for lifecycle operator ([#2200](https://github.com/keptn/lifecycle-toolkit/issues/2200)) ([9f0853f](https://github.com/keptn/lifecycle-toolkit/commit/9f0853fca2b92c9636e76dc77666148d86078af7))
* **lifecycle-operator:** automatically decide for scheduler installation based on k8s version ([#2212](https://github.com/keptn/lifecycle-toolkit/issues/2212)) ([25976ea](https://github.com/keptn/lifecycle-toolkit/commit/25976ead3fb1d95634ee3a00a7d37b3e98b2ec06))
* **lifecycle-operator:** introduce functions for SchedulingGates functionality ([#2140](https://github.com/keptn/lifecycle-toolkit/issues/2140)) ([b40503e](https://github.com/keptn/lifecycle-toolkit/commit/b40503ef6c867033994383767ad5149eb08ab8da))
* **lifecycle-operator:** introduce metric showing readiness of operator ([#2152](https://github.com/keptn/lifecycle-toolkit/issues/2152)) ([c0e3f48](https://github.com/keptn/lifecycle-toolkit/commit/c0e3f48dd0e34084c7d2d8e469e73c3f2865ea48))
* **lifecycle-operator:** introduce option to enable lifecycle orchestration only for specific namespaces ([#2244](https://github.com/keptn/lifecycle-toolkit/issues/2244)) ([12caf03](https://github.com/keptn/lifecycle-toolkit/commit/12caf031d336c7a34e495b36daccb5ec3524ae49))
* **lifecycle-operator:** introduce v1alpha4 API version for KeptnWorkloadInstance ([#2250](https://github.com/keptn/lifecycle-toolkit/issues/2250)) ([d95dc10](https://github.com/keptn/lifecycle-toolkit/commit/d95dc1037ce22296aff65d6ad6fa420e96172d5d))
* **metrics-operator:** add support for user-friendly duration string for specifying time frame ([#2147](https://github.com/keptn/lifecycle-toolkit/issues/2147)) ([34e5384](https://github.com/keptn/lifecycle-toolkit/commit/34e5384bb434836658a7bf375c4fc8de765023b6))
* move helm docs into values files ([#2281](https://github.com/keptn/lifecycle-toolkit/issues/2281)) ([bd1a37b](https://github.com/keptn/lifecycle-toolkit/commit/bd1a37b324e25d07e88e7c4d1ad8150a7b3d4dac))
* support scheduling gates in integration tests ([#2149](https://github.com/keptn/lifecycle-toolkit/issues/2149)) ([3ff67d5](https://github.com/keptn/lifecycle-toolkit/commit/3ff67d5632f287613f337c7418aa5e28e616d536))


### Bug Fixes

* change klt to keptn for annotations and certs ([#2229](https://github.com/keptn/lifecycle-toolkit/issues/2229)) ([608a75e](https://github.com/keptn/lifecycle-toolkit/commit/608a75ebb73006b82b370b40e86b83ee874764e8))
* helm charts image registry, image pull policy and install action ([#2361](https://github.com/keptn/lifecycle-toolkit/issues/2361)) ([76ed884](https://github.com/keptn/lifecycle-toolkit/commit/76ed884498971c87c48cdab6fea822dfcf3e6e2f))
* **lifecycle-operator:** make sure the CloudEvents endpoint from the KeptnConfig is applied ([#2289](https://github.com/keptn/lifecycle-toolkit/issues/2289)) ([b5d9fc0](https://github.com/keptn/lifecycle-toolkit/commit/b5d9fc0b182ff3d1a777dabec74314df3157edbb))
* **lifecycle-operator:** remove hardcoded keptn namespace ([#2141](https://github.com/keptn/lifecycle-toolkit/issues/2141)) ([f10b447](https://github.com/keptn/lifecycle-toolkit/commit/f10b4470bdc4346e6ccd17fecc92c8bd5675c7e5))
* update kustomization.yaml to avoid usage of deprecated patches/configs ([#2004](https://github.com/keptn/lifecycle-toolkit/issues/2004)) ([8d70fac](https://github.com/keptn/lifecycle-toolkit/commit/8d70fac1f9469107257976659fb8b7b414d0455b))
* update outdated CRDs in helm chart templates ([#2123](https://github.com/keptn/lifecycle-toolkit/issues/2123)) ([34c9d11](https://github.com/keptn/lifecycle-toolkit/commit/34c9d11a1dd34b181d2d1a9e5c61fd75638aaebf))


### Other

* adapt Makefile command to run unit tests ([#2072](https://github.com/keptn/lifecycle-toolkit/issues/2072)) ([2db2569](https://github.com/keptn/lifecycle-toolkit/commit/2db25691748beedbb02ed92806d327067c422285))
* **lifecycle-operator:** improve logging ([#2253](https://github.com/keptn/lifecycle-toolkit/issues/2253)) ([8dd3394](https://github.com/keptn/lifecycle-toolkit/commit/8dd3394087cf0d445ec0b3bad0a54242ad9f4f26))
* **lifecycle-operator:** refactor pod mutating webhook ([#2233](https://github.com/keptn/lifecycle-toolkit/issues/2233)) ([c2cc89a](https://github.com/keptn/lifecycle-toolkit/commit/c2cc89a3ad3ac0fef3410adb1c0b24aa10e8dc66))
* **lifecycle-operator:** remove direct dependency on jsonpatch ([#2187](https://github.com/keptn/lifecycle-toolkit/issues/2187)) ([d7fce2a](https://github.com/keptn/lifecycle-toolkit/commit/d7fce2a320bd34cd41d564e0d528675e5a1cd93e))
* **lifecycle-operator:** remove spans created by webhook ([#2331](https://github.com/keptn/lifecycle-toolkit/issues/2331)) ([9f21fb6](https://github.com/keptn/lifecycle-toolkit/commit/9f21fb62284e806f6356341315873f98e0c4fd29))
* **lifecycle-operator:** remove spans for reconciliation loops, adjust log levels ([#2310](https://github.com/keptn/lifecycle-toolkit/issues/2310)) ([d73008c](https://github.com/keptn/lifecycle-toolkit/commit/d73008ccaa5e028f1551392b9c68a4ea0315350e))
* regenerate CRDs ([#2074](https://github.com/keptn/lifecycle-toolkit/issues/2074)) ([63f5dc1](https://github.com/keptn/lifecycle-toolkit/commit/63f5dc1bc3dfd696de3730ed3949c0f99abdecc0))
* update k8s version ([#1701](https://github.com/keptn/lifecycle-toolkit/issues/1701)) ([010d7cd](https://github.com/keptn/lifecycle-toolkit/commit/010d7cd48c2e26993e25de607f30b40513c9cd61))
* update release please config to work with umbrella chart ([#2357](https://github.com/keptn/lifecycle-toolkit/issues/2357)) ([6ff3a5f](https://github.com/keptn/lifecycle-toolkit/commit/6ff3a5f64e394504fd5e7b67f0ac0a608428c1be))


### Docs

* adapt KeptnTask example to changes in API  ([#2124](https://github.com/keptn/lifecycle-toolkit/issues/2124)) ([bcc64e8](https://github.com/keptn/lifecycle-toolkit/commit/bcc64e814d7735bc330d2d0b3b52eccf7a51dbbe))


### Dependency Updates

* update dependency kubernetes-sigs/kustomize to v5.2.1 ([#2308](https://github.com/keptn/lifecycle-toolkit/issues/2308)) ([6653a47](https://github.com/keptn/lifecycle-toolkit/commit/6653a47d4156c0e60aa471f11a643a2664669023))
* update ghcr.io/keptn/deno-runtime docker tag to v1.0.2 ([#2367](https://github.com/keptn/lifecycle-toolkit/issues/2367)) ([6c17203](https://github.com/keptn/lifecycle-toolkit/commit/6c1720356fab6b4a9d1c0dae30e76e6d5c135c70))
* update ghcr.io/keptn/python-runtime docker tag to v1.0.1 ([#2368](https://github.com/keptn/lifecycle-toolkit/issues/2368)) ([134191a](https://github.com/keptn/lifecycle-toolkit/commit/134191a523c6d278771ad1f3421e4ae68dad4de9))
* update ghcr.io/keptn/scheduler docker tag to v0.8.3 ([#2374](https://github.com/keptn/lifecycle-toolkit/issues/2374)) ([16a4a14](https://github.com/keptn/lifecycle-toolkit/commit/16a4a147905fe19b319010e880730ee46e6c5965))
* update github.com/keptn/lifecycle-toolkit/klt-cert-manager digest to 010d7cd ([#2106](https://github.com/keptn/lifecycle-toolkit/issues/2106)) ([2ef614a](https://github.com/keptn/lifecycle-toolkit/commit/2ef614ad08dbeed1301889ed560375a2bb9e737c))
* update github.com/keptn/lifecycle-toolkit/klt-cert-manager digest to 066be3e ([#2274](https://github.com/keptn/lifecycle-toolkit/issues/2274)) ([c6d9c52](https://github.com/keptn/lifecycle-toolkit/commit/c6d9c524afa61e7c87553b89ebab1a2b8daa1438))
* update github.com/keptn/lifecycle-toolkit/klt-cert-manager digest to 099a457 ([#2169](https://github.com/keptn/lifecycle-toolkit/issues/2169)) ([643ae4e](https://github.com/keptn/lifecycle-toolkit/commit/643ae4e69ac527a342eed20c0e078c6b0e8cdd97))
* update github.com/keptn/lifecycle-toolkit/klt-cert-manager digest to 3077e31 ([#2313](https://github.com/keptn/lifecycle-toolkit/issues/2313)) ([cf52945](https://github.com/keptn/lifecycle-toolkit/commit/cf529455f2e99119e48ba433d28e8aecc31cad00))
* update github.com/keptn/lifecycle-toolkit/klt-cert-manager digest to 4342d33 ([#2177](https://github.com/keptn/lifecycle-toolkit/issues/2177)) ([2b5267c](https://github.com/keptn/lifecycle-toolkit/commit/2b5267c0a88b7f68167314d624c54453a326b5ce))
* update github.com/keptn/lifecycle-toolkit/klt-cert-manager digest to 469578e ([#2038](https://github.com/keptn/lifecycle-toolkit/issues/2038)) ([d240e56](https://github.com/keptn/lifecycle-toolkit/commit/d240e56fbc0b42caad04d8393ec59e55f1013efa))
* update github.com/keptn/lifecycle-toolkit/klt-cert-manager digest to 5efa650 ([#2155](https://github.com/keptn/lifecycle-toolkit/issues/2155)) ([fa8c891](https://github.com/keptn/lifecycle-toolkit/commit/fa8c8912825ad0bbc3f75b7a037e856bac6dad93))
* update github.com/keptn/lifecycle-toolkit/klt-cert-manager digest to 608a75e ([#2231](https://github.com/keptn/lifecycle-toolkit/issues/2231)) ([26ff714](https://github.com/keptn/lifecycle-toolkit/commit/26ff714800bb605bfe58b61da432237132edf072))
* update github.com/keptn/lifecycle-toolkit/klt-cert-manager digest to 6566e7d ([#2143](https://github.com/keptn/lifecycle-toolkit/issues/2143)) ([9e7fe83](https://github.com/keptn/lifecycle-toolkit/commit/9e7fe8353dd8c84fae96081c2bb7522ab7ff7f5a))
* update github.com/keptn/lifecycle-toolkit/klt-cert-manager digest to 8dd3394 ([#2271](https://github.com/keptn/lifecycle-toolkit/issues/2271)) ([b29fc99](https://github.com/keptn/lifecycle-toolkit/commit/b29fc999aef2c214b7b45a1161b226e85e3eaffe))
* update github.com/keptn/lifecycle-toolkit/klt-cert-manager digest to a15b038 ([#2205](https://github.com/keptn/lifecycle-toolkit/issues/2205)) ([1592926](https://github.com/keptn/lifecycle-toolkit/commit/1592926d6c70062cd632969f991531aa9b00f0de))
* update github.com/keptn/lifecycle-toolkit/klt-cert-manager digest to a656512 ([#2230](https://github.com/keptn/lifecycle-toolkit/issues/2230)) ([f11fdb9](https://github.com/keptn/lifecycle-toolkit/commit/f11fdb959b4e207d3704361870d515f61ad92360))
* update github.com/keptn/lifecycle-toolkit/klt-cert-manager digest to b2853f9 ([#2094](https://github.com/keptn/lifecycle-toolkit/issues/2094)) ([b9019cd](https://github.com/keptn/lifecycle-toolkit/commit/b9019cd96a161c4e0c4dd08e3ddbabd152ea921c))
* update github.com/keptn/lifecycle-toolkit/klt-cert-manager digest to c1166ff ([#2242](https://github.com/keptn/lifecycle-toolkit/issues/2242)) ([aa53137](https://github.com/keptn/lifecycle-toolkit/commit/aa531375032468a6e0d1b3a9f6eb3e6e9b0c998b))
* update github.com/keptn/lifecycle-toolkit/klt-cert-manager digest to f2f3a0c ([#2132](https://github.com/keptn/lifecycle-toolkit/issues/2132)) ([2039d36](https://github.com/keptn/lifecycle-toolkit/commit/2039d36f427e22bfe692fade207319747ee15083))
* update github.com/keptn/lifecycle-toolkit/klt-cert-manager digest to f2f8dfe ([#2297](https://github.com/keptn/lifecycle-toolkit/issues/2297)) ([e13b9be](https://github.com/keptn/lifecycle-toolkit/commit/e13b9be3217fdd5bf3a646dd3d6ba49438cbb9e6))
* update github.com/keptn/lifecycle-toolkit/klt-cert-manager digest to f3bbb96 ([#2342](https://github.com/keptn/lifecycle-toolkit/issues/2342)) ([89ddb2f](https://github.com/keptn/lifecycle-toolkit/commit/89ddb2f427561bbc41ea2e4b762ac3a14aab3bc5))
* update github.com/keptn/lifecycle-toolkit/klt-cert-manager digest to fda2315 ([#2300](https://github.com/keptn/lifecycle-toolkit/issues/2300)) ([bffbaf2](https://github.com/keptn/lifecycle-toolkit/commit/bffbaf2edf4710e085db1ee956b3ccc1b6599275))
* update kubernetes packages (patch) ([#2102](https://github.com/keptn/lifecycle-toolkit/issues/2102)) ([b2853f9](https://github.com/keptn/lifecycle-toolkit/commit/b2853f9ecdfb4b7b81d0b88cf782b82c9958c5cb))
* update module github.com/argoproj/argo-rollouts to v1.6.0 ([#2064](https://github.com/keptn/lifecycle-toolkit/issues/2064)) ([d5c428a](https://github.com/keptn/lifecycle-toolkit/commit/d5c428a7e31c00362b8280da5acd91ade89c1fa8))
* update module github.com/go-logr/logr to v1.3.0 ([#2346](https://github.com/keptn/lifecycle-toolkit/issues/2346)) ([bc06204](https://github.com/keptn/lifecycle-toolkit/commit/bc06204b97c765d0f5664fd66f441af86f21e191))
* update module github.com/onsi/ginkgo/v2 to v2.12.1 ([#2156](https://github.com/keptn/lifecycle-toolkit/issues/2156)) ([dbf2867](https://github.com/keptn/lifecycle-toolkit/commit/dbf2867133067b162e82b71b6547c3dfac95d0af))
* update module github.com/onsi/ginkgo/v2 to v2.13.0 ([#2272](https://github.com/keptn/lifecycle-toolkit/issues/2272)) ([0df464d](https://github.com/keptn/lifecycle-toolkit/commit/0df464dd8e4fc7729deeb5bae4938b236902d661))
* update module github.com/onsi/gomega to v1.28.0 ([#2209](https://github.com/keptn/lifecycle-toolkit/issues/2209)) ([c0726d0](https://github.com/keptn/lifecycle-toolkit/commit/c0726d0b0e9d9732123aaf8b1ad012bc24672b84))
* update module github.com/onsi/gomega to v1.28.1 ([#2343](https://github.com/keptn/lifecycle-toolkit/issues/2343)) ([64b1508](https://github.com/keptn/lifecycle-toolkit/commit/64b1508f0e383aa7fbc406e17e2cc66546601e53))
* update module github.com/prometheus/client_golang to v1.17.0 ([#2207](https://github.com/keptn/lifecycle-toolkit/issues/2207)) ([de8b958](https://github.com/keptn/lifecycle-toolkit/commit/de8b9587cb95b1ee9c2be7a66320d284491102d9))
* update module golang.org/x/net to v0.15.0 ([#2065](https://github.com/keptn/lifecycle-toolkit/issues/2065)) ([50ce9c0](https://github.com/keptn/lifecycle-toolkit/commit/50ce9c09914f505ffaf33eee41564afa65661215))
* update module golang.org/x/net to v0.16.0 ([#2249](https://github.com/keptn/lifecycle-toolkit/issues/2249)) ([e89ea71](https://github.com/keptn/lifecycle-toolkit/commit/e89ea71bc1a2d69828179c64ffe3c34ce359dd94))
* update module golang.org/x/net to v0.17.0 ([#2267](https://github.com/keptn/lifecycle-toolkit/issues/2267)) ([8443874](https://github.com/keptn/lifecycle-toolkit/commit/8443874254cda9e5f4c662cab1a3e5e3b3277435))
* update module google.golang.org/grpc to v1.58.0 ([#2066](https://github.com/keptn/lifecycle-toolkit/issues/2066)) ([6fae5a7](https://github.com/keptn/lifecycle-toolkit/commit/6fae5a7ebf356625b4754b7890f7c71dbb4ac0a6))
* update module google.golang.org/grpc to v1.58.1 ([#2115](https://github.com/keptn/lifecycle-toolkit/issues/2115)) ([d08df40](https://github.com/keptn/lifecycle-toolkit/commit/d08df40188bc633037c49a1468a70eefc960a4a1))
* update module google.golang.org/grpc to v1.58.2 ([#2163](https://github.com/keptn/lifecycle-toolkit/issues/2163)) ([5efa650](https://github.com/keptn/lifecycle-toolkit/commit/5efa6502403daa37bdfc51fa8600da6b1f845ac2))
* update module google.golang.org/grpc to v1.58.3 ([#2275](https://github.com/keptn/lifecycle-toolkit/issues/2275)) ([66e86c0](https://github.com/keptn/lifecycle-toolkit/commit/66e86c03272d75207bd3b42014d88b1b912b9198))
* update module google.golang.org/grpc to v1.59.0 ([#2302](https://github.com/keptn/lifecycle-toolkit/issues/2302)) ([fda2315](https://github.com/keptn/lifecycle-toolkit/commit/fda231552475eaf0f60457ad42a26c4ed3473008))
* update module k8s.io/apimachinery to v0.28.3 ([#2298](https://github.com/keptn/lifecycle-toolkit/issues/2298)) ([f2f8dfe](https://github.com/keptn/lifecycle-toolkit/commit/f2f8dfec6e47517f2c476d6425c22db875f9bd3c))
* update module sigs.k8s.io/controller-runtime to v0.16.3 ([#2306](https://github.com/keptn/lifecycle-toolkit/issues/2306)) ([3d634a7](https://github.com/keptn/lifecycle-toolkit/commit/3d634a79996be6cb50805c745c51309c2f091a61))
* update opentelemetry-go monorepo (minor) ([#2108](https://github.com/keptn/lifecycle-toolkit/issues/2108)) ([4e5d29e](https://github.com/keptn/lifecycle-toolkit/commit/4e5d29e681f78590b4406ba7b74cc46ca6107e4b))
* update opentelemetry-go monorepo (minor) ([#2210](https://github.com/keptn/lifecycle-toolkit/issues/2210)) ([d577311](https://github.com/keptn/lifecycle-toolkit/commit/d5773111c327f5d30ec24437d16cf5d4454dd69e))

## [0.8.2](https://github.com/keptn/lifecycle-toolkit/compare/lifecycle-operator-v0.8.1...lifecycle-operator-v0.8.2) (2023-09-06)


### Features

* add cloud events support ([#1843](https://github.com/keptn/lifecycle-toolkit/issues/1843)) ([5b47120](https://github.com/keptn/lifecycle-toolkit/commit/5b471203e412a919903876212ac45c04f180e482))
* **lifecycle-operator:** clean up KeptnTask API by removing duplicated attributes ([#1965](https://github.com/keptn/lifecycle-toolkit/issues/1965)) ([257b220](https://github.com/keptn/lifecycle-toolkit/commit/257b220a6171ccc82d1b471002b6cf773ec9bd09))
* **metrics-operator:** add analysis controller ([#1875](https://github.com/keptn/lifecycle-toolkit/issues/1875)) ([017e08b](https://github.com/keptn/lifecycle-toolkit/commit/017e08b0a65679ca417e363f2223b7f4fef3bc55))
* **metrics-operator:** add Analysis CRD ([#1839](https://github.com/keptn/lifecycle-toolkit/issues/1839)) ([9521a16](https://github.com/keptn/lifecycle-toolkit/commit/9521a16ce4946d3169993780f2d2a4f3a75d0445))
* monorepo setup for lifecycle-operator, scheduler and runtimes ([#1857](https://github.com/keptn/lifecycle-toolkit/issues/1857)) ([84e243a](https://github.com/keptn/lifecycle-toolkit/commit/84e243a213ffba86eddd51ccc4bf4dbd61140069))


### Bug Fixes

* **lifecycle-operator:** avoid setting the overall state of an App or WorkloadInstance between state transitions ([#1871](https://github.com/keptn/lifecycle-toolkit/issues/1871)) ([ee0b085](https://github.com/keptn/lifecycle-toolkit/commit/ee0b085b05b2b9781457eba34d5d1050b3c7a604))


### Other

* **main:** release lifecycle-operator-and-scheduler libraries ([#1979](https://github.com/keptn/lifecycle-toolkit/issues/1979)) ([12d0f40](https://github.com/keptn/lifecycle-toolkit/commit/12d0f40725e466825c4a0d483fa344e5888b03ae))
* rename operator folder to lifecycle-operator ([#1819](https://github.com/keptn/lifecycle-toolkit/issues/1819)) ([97a2d25](https://github.com/keptn/lifecycle-toolkit/commit/97a2d25919c0a02165dd0dc6c7c82d57ad200139))


### Docs

* fix typos and grammar issues ([#1925](https://github.com/keptn/lifecycle-toolkit/issues/1925)) ([5570d55](https://github.com/keptn/lifecycle-toolkit/commit/5570d555bfc4bbdcbfc66b2725d5352090e5b937))
* implement KLT -&gt; Keptn name change ([#2001](https://github.com/keptn/lifecycle-toolkit/issues/2001)) ([440c308](https://github.com/keptn/lifecycle-toolkit/commit/440c3082e5400f89d791724651984ba2bc0a4724))


### Dependency Updates

* update dependency kubernetes-sigs/controller-tools to v0.13.0 ([#1930](https://github.com/keptn/lifecycle-toolkit/issues/1930)) ([8b34b63](https://github.com/keptn/lifecycle-toolkit/commit/8b34b63404d0339633ef41ff1cf2005deae8d2b7))
* update dependency kubernetes-sigs/kustomize to v5.1.1 ([#1853](https://github.com/keptn/lifecycle-toolkit/issues/1853)) ([354ab3f](https://github.com/keptn/lifecycle-toolkit/commit/354ab3f980c2569e17a0354ece417df40317d120))
* update github.com/keptn/lifecycle-toolkit/klt-cert-manager digest to 440c308 ([#2017](https://github.com/keptn/lifecycle-toolkit/issues/2017)) ([c365734](https://github.com/keptn/lifecycle-toolkit/commit/c365734fa7e3e40b2ae4c97c61628892d040dacc))
* update github.com/keptn/lifecycle-toolkit/klt-cert-manager digest to 88a54f9 ([#1794](https://github.com/keptn/lifecycle-toolkit/issues/1794)) ([fc976eb](https://github.com/keptn/lifecycle-toolkit/commit/fc976eb07ed9a5e49ed7d4ab1dbf187cee583e64))
* update github.com/keptn/lifecycle-toolkit/klt-cert-manager digest to 8dbec2d ([#1995](https://github.com/keptn/lifecycle-toolkit/issues/1995)) ([2f51445](https://github.com/keptn/lifecycle-toolkit/commit/2f5144540c4b3876e800bff29c30bfded334be40))
* update github.com/keptn/lifecycle-toolkit/klt-cert-manager digest to bb133cf ([#1963](https://github.com/keptn/lifecycle-toolkit/issues/1963)) ([c7697bf](https://github.com/keptn/lifecycle-toolkit/commit/c7697bf54d5fe18b7c62c5b11801c6c83079b0a3))
* update kubernetes packages to v0.26.8 (patch) ([#1945](https://github.com/keptn/lifecycle-toolkit/issues/1945)) ([6ce03d6](https://github.com/keptn/lifecycle-toolkit/commit/6ce03d600cbb3d3d3988573c616ec7f3830ba324))
* update module github.com/onsi/ginkgo/v2 to v2.12.0 ([#2019](https://github.com/keptn/lifecycle-toolkit/issues/2019)) ([41e878f](https://github.com/keptn/lifecycle-toolkit/commit/41e878ff8bbb438efa4b221470a571687dd392e9))
* update module github.com/onsi/gomega to v1.27.10 ([#1796](https://github.com/keptn/lifecycle-toolkit/issues/1796)) ([8f14bff](https://github.com/keptn/lifecycle-toolkit/commit/8f14bffe27485a36e0b05b770a01e357402d92f7))
* update module golang.org/x/net to v0.14.0 ([#1855](https://github.com/keptn/lifecycle-toolkit/issues/1855)) ([3186188](https://github.com/keptn/lifecycle-toolkit/commit/31861889bf7b227f489b941ac4a52db86551fcc2))
* update module google.golang.org/grpc to v1.57.0 ([#1861](https://github.com/keptn/lifecycle-toolkit/issues/1861)) ([fdcbdf5](https://github.com/keptn/lifecycle-toolkit/commit/fdcbdf50365dfd69d16c679c6814e89570a8a0e2))
* update opentelemetry-go monorepo (minor) ([#1931](https://github.com/keptn/lifecycle-toolkit/issues/1931)) ([a0a1a7e](https://github.com/keptn/lifecycle-toolkit/commit/a0a1a7e97906ab56ed85da7ab9b6d1e13c902397))


### Refactoring

* **lifecycle-operator:** eventing and telemetry ([#1844](https://github.com/keptn/lifecycle-toolkit/issues/1844)) ([0130576](https://github.com/keptn/lifecycle-toolkit/commit/0130576a17a78453019c150af849c06553d799a6))
* **lifecycle-operator:** refactor event emitter unit tests ([#1867](https://github.com/keptn/lifecycle-toolkit/issues/1867)) ([2558f74](https://github.com/keptn/lifecycle-toolkit/commit/2558f742031e4d38a8006ce9894f41bddac7cd3f))
