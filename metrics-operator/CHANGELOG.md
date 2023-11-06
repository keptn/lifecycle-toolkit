# Changelog

## [0.8.3](https://github.com/keptn/lifecycle-toolkit/compare/metrics-operator-v0.8.2...metrics-operator-v0.8.3) (2023-10-30)


### Features

* add test and lint cmd to makefiles ([#2176](https://github.com/keptn/lifecycle-toolkit/issues/2176)) ([c55e0a9](https://github.com/keptn/lifecycle-toolkit/commit/c55e0a9f368c82ad3032eb676edd59e68b29fad6))
* aggregation functions support for metrics controller ([#1802](https://github.com/keptn/lifecycle-toolkit/issues/1802)) ([678c4c9](https://github.com/keptn/lifecycle-toolkit/commit/678c4c9efaa53eb2bb64ef31c98a08d92eccd810))
* create new Keptn umbrella Helm chart ([#2214](https://github.com/keptn/lifecycle-toolkit/issues/2214)) ([41bd47b](https://github.com/keptn/lifecycle-toolkit/commit/41bd47b7748c4d645243a4dae165651bbfd3533f))
* generalize helm chart ([#2282](https://github.com/keptn/lifecycle-toolkit/issues/2282)) ([81334eb](https://github.com/keptn/lifecycle-toolkit/commit/81334ebec4d8afda27902b6e854c4c637a3daa87))
* **metrics-operator:** add basicauth to prometheus provider ([#2154](https://github.com/keptn/lifecycle-toolkit/issues/2154)) ([bab605e](https://github.com/keptn/lifecycle-toolkit/commit/bab605e39f40df79d615532fc7592f0bd809993d))
* **metrics-operator:** add helm chart for metrics operator ([#2189](https://github.com/keptn/lifecycle-toolkit/issues/2189)) ([a5ae3de](https://github.com/keptn/lifecycle-toolkit/commit/a5ae3ded2229444c1a8e6c3d3ebc5abbcb7187e3))
* **metrics-operator:** add query to the analysis result ([#2188](https://github.com/keptn/lifecycle-toolkit/issues/2188)) ([233aac4](https://github.com/keptn/lifecycle-toolkit/commit/233aac4e91c44f663db08e4827fa0aa693556ed7))
* **metrics-operator:** add support for user-friendly duration string for specifying time frame ([#2147](https://github.com/keptn/lifecycle-toolkit/issues/2147)) ([34e5384](https://github.com/keptn/lifecycle-toolkit/commit/34e5384bb434836658a7bf375c4fc8de765023b6))
* **metrics-operator:** expose analysis results as Prometheus Metric ([#2137](https://github.com/keptn/lifecycle-toolkit/issues/2137)) ([47b756c](https://github.com/keptn/lifecycle-toolkit/commit/47b756c7dc146709e9a1378e89592b9a2cdbbae5))
* **metrics-operator:** implement interface for analysis value retrieval in DQL provider ([#2194](https://github.com/keptn/lifecycle-toolkit/issues/2194)) ([3d7f737](https://github.com/keptn/lifecycle-toolkit/commit/3d7f737d4a891263ba1122d7fa3cf578299e252a))
* **metrics-operator:** remove omitempty tags to get complete representation of AnalysisResult ([#2078](https://github.com/keptn/lifecycle-toolkit/issues/2078)) ([a08b9ca](https://github.com/keptn/lifecycle-toolkit/commit/a08b9cae35a4e1ac224d6cb76b64003363c6e915))
* move helm docs into values files ([#2281](https://github.com/keptn/lifecycle-toolkit/issues/2281)) ([bd1a37b](https://github.com/keptn/lifecycle-toolkit/commit/bd1a37b324e25d07e88e7c4d1ad8150a7b3d4dac))
* update `KeptnMetric` to store multiple metrics in status ([#1900](https://github.com/keptn/lifecycle-toolkit/issues/1900)) ([2252b2d](https://github.com/keptn/lifecycle-toolkit/commit/2252b2daa5b26e7335a72ac4cd42086de50c0279))


### Bug Fixes

* change klt to keptn for annotations and certs ([#2229](https://github.com/keptn/lifecycle-toolkit/issues/2229)) ([608a75e](https://github.com/keptn/lifecycle-toolkit/commit/608a75ebb73006b82b370b40e86b83ee874764e8))
* helm charts image registry, image pull policy and install action ([#2361](https://github.com/keptn/lifecycle-toolkit/issues/2361)) ([76ed884](https://github.com/keptn/lifecycle-toolkit/commit/76ed884498971c87c48cdab6fea822dfcf3e6e2f))
* **lifecycle-operator:** remove hardcoded keptn namespace ([#2141](https://github.com/keptn/lifecycle-toolkit/issues/2141)) ([f10b447](https://github.com/keptn/lifecycle-toolkit/commit/f10b4470bdc4346e6ccd17fecc92c8bd5675c7e5))
* **metrics-operator:** convert SLI names to valid K8s resource names ([#2125](https://github.com/keptn/lifecycle-toolkit/issues/2125)) ([6da3276](https://github.com/keptn/lifecycle-toolkit/commit/6da3276e1fecac3bf004afb3860f2e6983b9dec0))
* **metrics-operator:** fix log message for AnalysisDefinition lookup ([#2092](https://github.com/keptn/lifecycle-toolkit/issues/2092)) ([598fed3](https://github.com/keptn/lifecycle-toolkit/commit/598fed3bad2c147c791d4d8f43bcbc33a53f448d))
* **metrics-operator:** fix panic due to write attempt on closed channel ([#2119](https://github.com/keptn/lifecycle-toolkit/issues/2119)) ([33eb9d7](https://github.com/keptn/lifecycle-toolkit/commit/33eb9d7da65dc012f1da5fdc27b1c33f88be210f))
* **metrics-operator:** flush status when analysis is finished ([#2122](https://github.com/keptn/lifecycle-toolkit/issues/2122)) ([276b609](https://github.com/keptn/lifecycle-toolkit/commit/276b6094af7af4646d2fb9cba884e2c60eec4e97))
* **metrics-operator:** introduce `.status.state` in Analysis ([#2061](https://github.com/keptn/lifecycle-toolkit/issues/2061)) ([b08b4d8](https://github.com/keptn/lifecycle-toolkit/commit/b08b4d8adca2cac13466bd3227fe23249fd5d12c))
* **metrics-operator:** use context with timeout for fetching analysis values ([#2213](https://github.com/keptn/lifecycle-toolkit/issues/2213)) ([6945069](https://github.com/keptn/lifecycle-toolkit/commit/6945069de9f7d34822af2f80f3b654579a565c02))
* update kustomization.yaml to avoid usage of deprecated patches/configs ([#2004](https://github.com/keptn/lifecycle-toolkit/issues/2004)) ([8d70fac](https://github.com/keptn/lifecycle-toolkit/commit/8d70fac1f9469107257976659fb8b7b414d0455b))
* update outdated CRDs in helm chart templates ([#2123](https://github.com/keptn/lifecycle-toolkit/issues/2123)) ([34c9d11](https://github.com/keptn/lifecycle-toolkit/commit/34c9d11a1dd34b181d2d1a9e5c61fd75638aaebf))


### Other

* adapt Makefile command to run unit tests ([#2072](https://github.com/keptn/lifecycle-toolkit/issues/2072)) ([2db2569](https://github.com/keptn/lifecycle-toolkit/commit/2db25691748beedbb02ed92806d327067c422285))
* **metrics-operator:** improve logging ([#2269](https://github.com/keptn/lifecycle-toolkit/issues/2269)) ([2e35273](https://github.com/keptn/lifecycle-toolkit/commit/2e35273b2d03397c114e46b6c57b83ff208fbe6a))
* **metrics-operator:** inject ProviderFactory into KeptnMetric reconciler ([#2062](https://github.com/keptn/lifecycle-toolkit/issues/2062)) ([579dc10](https://github.com/keptn/lifecycle-toolkit/commit/579dc100822dc1496214e1a01f9f2126a541af8c))
* **metrics-operator:** refactor fetching resouce namespaces during analysis ([#2105](https://github.com/keptn/lifecycle-toolkit/issues/2105)) ([38c8332](https://github.com/keptn/lifecycle-toolkit/commit/38c8332b3f6d59170cf2de65ab1461bac9f6f742))
* regenerate CRDs ([#2074](https://github.com/keptn/lifecycle-toolkit/issues/2074)) ([63f5dc1](https://github.com/keptn/lifecycle-toolkit/commit/63f5dc1bc3dfd696de3730ed3949c0f99abdecc0))
* update k8s version ([#1701](https://github.com/keptn/lifecycle-toolkit/issues/1701)) ([010d7cd](https://github.com/keptn/lifecycle-toolkit/commit/010d7cd48c2e26993e25de607f30b40513c9cd61))
* update release please config to work with umbrella chart ([#2357](https://github.com/keptn/lifecycle-toolkit/issues/2357)) ([6ff3a5f](https://github.com/keptn/lifecycle-toolkit/commit/6ff3a5f64e394504fd5e7b67f0ac0a608428c1be))


### Docs

* add first iteration of analysis documentation ([#2167](https://github.com/keptn/lifecycle-toolkit/issues/2167)) ([366ee1f](https://github.com/keptn/lifecycle-toolkit/commit/366ee1f77e466b5939e32603e374292001758cd5))
* **metrics-operator:** usage of SLI and SLO converters ([#2013](https://github.com/keptn/lifecycle-toolkit/issues/2013)) ([57bc225](https://github.com/keptn/lifecycle-toolkit/commit/57bc225f8f3990f7bc9aeab077f3bd6ea511db22))


### Dependency Updates

* **metrics-operator:** replace grpc version with v1.58.3 ([#2353](https://github.com/keptn/lifecycle-toolkit/issues/2353)) ([51269d4](https://github.com/keptn/lifecycle-toolkit/commit/51269d4b7de60a3d87392c1eaaa9042bcda75c0b))
* replace otel libraries with newer versions ([#2312](https://github.com/keptn/lifecycle-toolkit/issues/2312)) ([adda244](https://github.com/keptn/lifecycle-toolkit/commit/adda244abee4efbc2d263763064769e16d1ac421))
* update dependency kubernetes-sigs/kustomize to v5.2.1 ([#2308](https://github.com/keptn/lifecycle-toolkit/issues/2308)) ([6653a47](https://github.com/keptn/lifecycle-toolkit/commit/6653a47d4156c0e60aa471f11a643a2664669023))
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
* update golang.org/x/exp digest to 3e424a5 ([#2243](https://github.com/keptn/lifecycle-toolkit/issues/2243)) ([61ca8b7](https://github.com/keptn/lifecycle-toolkit/commit/61ca8b7167e39a8f69f283c13bf026716017a0ee))
* update golang.org/x/exp digest to 7918f67 ([#2246](https://github.com/keptn/lifecycle-toolkit/issues/2246)) ([a05a915](https://github.com/keptn/lifecycle-toolkit/commit/a05a91504f165fe326aac5c36862dd0c84ae18fe))
* update golang.org/x/exp digest to 9212866 ([#2039](https://github.com/keptn/lifecycle-toolkit/issues/2039)) ([2fba7c9](https://github.com/keptn/lifecycle-toolkit/commit/2fba7c94ff0fa6aa14facc7dc22b1a7558a88a18))
* update golang.org/x/exp digest to 9212866 ([#2133](https://github.com/keptn/lifecycle-toolkit/issues/2133)) ([3390f17](https://github.com/keptn/lifecycle-toolkit/commit/3390f17ed901979cced46c3d05183cc21716099a))
* update kubernetes packages (patch) ([#2102](https://github.com/keptn/lifecycle-toolkit/issues/2102)) ([b2853f9](https://github.com/keptn/lifecycle-toolkit/commit/b2853f9ecdfb4b7b81d0b88cf782b82c9958c5cb))
* update module github.com/datadog/datadog-api-client-go/v2 to v2.17.0 ([#2107](https://github.com/keptn/lifecycle-toolkit/issues/2107)) ([a048036](https://github.com/keptn/lifecycle-toolkit/commit/a04803625526b48b1492f344a5641a1603ae3d4d))
* update module github.com/datadog/datadog-api-client-go/v2 to v2.17.0 ([#2134](https://github.com/keptn/lifecycle-toolkit/issues/2134)) ([f2f3a0c](https://github.com/keptn/lifecycle-toolkit/commit/f2f3a0c1e09c7d44c3e0be2c6fb6e4907aa9583f))
* update module github.com/datadog/datadog-api-client-go/v2 to v2.18.0 ([#2301](https://github.com/keptn/lifecycle-toolkit/issues/2301)) ([d3ecf2e](https://github.com/keptn/lifecycle-toolkit/commit/d3ecf2eb09921b9d6ec9461767ea15a04b68a8c5))
* update module github.com/go-logr/logr to v1.3.0 ([#2346](https://github.com/keptn/lifecycle-toolkit/issues/2346)) ([bc06204](https://github.com/keptn/lifecycle-toolkit/commit/bc06204b97c765d0f5664fd66f441af86f21e191))
* update module github.com/open-feature/go-sdk to v1.8.0 ([#2208](https://github.com/keptn/lifecycle-toolkit/issues/2208)) ([1494568](https://github.com/keptn/lifecycle-toolkit/commit/14945687c5cc0fc864cd2b1ac67c423fdde77312))
* update module github.com/prometheus/client_golang to v1.17.0 ([#2207](https://github.com/keptn/lifecycle-toolkit/issues/2207)) ([de8b958](https://github.com/keptn/lifecycle-toolkit/commit/de8b9587cb95b1ee9c2be7a66320d284491102d9))
* update module github.com/prometheus/client_model to v0.5.0 ([#2247](https://github.com/keptn/lifecycle-toolkit/issues/2247)) ([cb7a4b3](https://github.com/keptn/lifecycle-toolkit/commit/cb7a4b3f20d1137062bf960b741f23ccf3a07915))
* update module github.com/prometheus/common to v0.45.0 ([#2304](https://github.com/keptn/lifecycle-toolkit/issues/2304)) ([5705847](https://github.com/keptn/lifecycle-toolkit/commit/5705847860120adc4bcfd9677ccc19a9788ad3e4))
* update module golang.org/x/net to v0.15.0 ([#2065](https://github.com/keptn/lifecycle-toolkit/issues/2065)) ([50ce9c0](https://github.com/keptn/lifecycle-toolkit/commit/50ce9c09914f505ffaf33eee41564afa65661215))
* update module golang.org/x/net to v0.15.0 ([#2135](https://github.com/keptn/lifecycle-toolkit/issues/2135)) ([214af30](https://github.com/keptn/lifecycle-toolkit/commit/214af301ae7180fb94050ca50a9e6d090301fa98))
* update module golang.org/x/net to v0.16.0 ([#2249](https://github.com/keptn/lifecycle-toolkit/issues/2249)) ([e89ea71](https://github.com/keptn/lifecycle-toolkit/commit/e89ea71bc1a2d69828179c64ffe3c34ce359dd94))
* update module golang.org/x/net to v0.17.0 ([#2267](https://github.com/keptn/lifecycle-toolkit/issues/2267)) ([8443874](https://github.com/keptn/lifecycle-toolkit/commit/8443874254cda9e5f4c662cab1a3e5e3b3277435))
* update module k8s.io/apimachinery to v0.28.3 ([#2298](https://github.com/keptn/lifecycle-toolkit/issues/2298)) ([f2f8dfe](https://github.com/keptn/lifecycle-toolkit/commit/f2f8dfec6e47517f2c476d6425c22db875f9bd3c))
* update module sigs.k8s.io/controller-runtime to v0.16.3 ([#2306](https://github.com/keptn/lifecycle-toolkit/issues/2306)) ([3d634a7](https://github.com/keptn/lifecycle-toolkit/commit/3d634a79996be6cb50805c745c51309c2f091a61))
* update module sigs.k8s.io/yaml to v1.4.0 ([#2347](https://github.com/keptn/lifecycle-toolkit/issues/2347)) ([a8d9170](https://github.com/keptn/lifecycle-toolkit/commit/a8d9170a8b58cab0833cc388489c500160793ecc))

## [0.8.2](https://github.com/keptn/lifecycle-toolkit/compare/metrics-operator-v0.8.1...metrics-operator-v0.8.2) (2023-09-06)


### Features

* add `aggregation` field in `KeptnMetric` ([#1780](https://github.com/keptn/lifecycle-toolkit/issues/1780)) ([c0b66ea](https://github.com/keptn/lifecycle-toolkit/commit/c0b66eae296e0502608dd66c5fe7eb8f890497e6))
* add `interval` field for `kubectl get KeptnMetric` ([#1689](https://github.com/keptn/lifecycle-toolkit/issues/1689)) ([1599ee9](https://github.com/keptn/lifecycle-toolkit/commit/1599ee939c94c6dbdf6207796b7b4e5211c2cbf7))
* add `step` field in `KeptnMetric` ([#1755](https://github.com/keptn/lifecycle-toolkit/issues/1755)) ([03ca7dd](https://github.com/keptn/lifecycle-toolkit/commit/03ca7ddde4ce787d0bfddaba2bb3f7b422ff5d6a))
* metrics-operator monorepo setup ([#1791](https://github.com/keptn/lifecycle-toolkit/issues/1791)) ([51445eb](https://github.com/keptn/lifecycle-toolkit/commit/51445ebd24b0914d34b0339ab05ec939440aa4a3))
* **metrics-operator:** adapt to changes in DQL API ([#1948](https://github.com/keptn/lifecycle-toolkit/issues/1948)) ([88d693a](https://github.com/keptn/lifecycle-toolkit/commit/88d693af82d41656a991546b4cef555c5171ffdf))
* **metrics-operator:** add analysis controller ([#1875](https://github.com/keptn/lifecycle-toolkit/issues/1875)) ([017e08b](https://github.com/keptn/lifecycle-toolkit/commit/017e08b0a65679ca417e363f2223b7f4fef3bc55))
* **metrics-operator:** add Analysis CRD ([#1839](https://github.com/keptn/lifecycle-toolkit/issues/1839)) ([9521a16](https://github.com/keptn/lifecycle-toolkit/commit/9521a16ce4946d3169993780f2d2a4f3a75d0445))
* **metrics-operator:** add AnalysisDefinition CRD ([#1823](https://github.com/keptn/lifecycle-toolkit/issues/1823)) ([adf4621](https://github.com/keptn/lifecycle-toolkit/commit/adf4621c2e8147bc0e4ee7f1719859007290c978))
* **metrics-operator:** add AnalysisValueTemplate CRD  ([#1822](https://github.com/keptn/lifecycle-toolkit/issues/1822)) ([f25b24d](https://github.com/keptn/lifecycle-toolkit/commit/f25b24dfef07a600c0fbcd4bdb540efe58cff387))
* **metrics-operator:** add new provider interface ([#1943](https://github.com/keptn/lifecycle-toolkit/issues/1943)) ([66320f8](https://github.com/keptn/lifecycle-toolkit/commit/66320f893c1bfc88b2e5a03d77caecaa0b42681f))
* **metrics-operator:** convert corner cases in SLO convertor ([#1999](https://github.com/keptn/lifecycle-toolkit/issues/1999)) ([95e0953](https://github.com/keptn/lifecycle-toolkit/commit/95e0953c2c95ac22fe62d592c4ff6dd186e6a260))
* **metrics-operator:** introduce range operators in AnalysisDefinition ([#1976](https://github.com/keptn/lifecycle-toolkit/issues/1976)) ([7fb8952](https://github.com/keptn/lifecycle-toolkit/commit/7fb8952c514909ce2c0202e01f1cf501de2c8d55))
* **metrics-operator:** introduce scoring logic for Analysis evaluations ([#1872](https://github.com/keptn/lifecycle-toolkit/issues/1872)) ([b6f2172](https://github.com/keptn/lifecycle-toolkit/commit/b6f2172637481a395fcd48d428ab39780f5c0fa8))
* **metrics-operator:** introduce SLI -&gt; AnalysisValueTemplate converter ([#1939](https://github.com/keptn/lifecycle-toolkit/issues/1939)) ([6f2d261](https://github.com/keptn/lifecycle-toolkit/commit/6f2d2614c2fcd6b83edfa77db01255a666e0b07b))
* **metrics-operator:** introduce SLO -&gt; AnalysisDefinition converter ([#1955](https://github.com/keptn/lifecycle-toolkit/issues/1955)) ([9c9929c](https://github.com/keptn/lifecycle-toolkit/commit/9c9929c6f0cf1d51baedc6ddb40f0b0f25ddc228))
* **metrics-operator:** support combination of OR criteria in SLO converter ([#2023](https://github.com/keptn/lifecycle-toolkit/issues/2023)) ([aa430e7](https://github.com/keptn/lifecycle-toolkit/commit/aa430e7f106949e1ee3108bc08c2e604d3e25a9b))
* **metrics-operator:** update datadog api to support `range.step` ([#1842](https://github.com/keptn/lifecycle-toolkit/issues/1842)) ([1d957b7](https://github.com/keptn/lifecycle-toolkit/commit/1d957b724c5db7679bf32a33441f96108537d8e3))
* **metrics-operator:** update dql provider to include range ([#1919](https://github.com/keptn/lifecycle-toolkit/issues/1919)) ([39db23e](https://github.com/keptn/lifecycle-toolkit/commit/39db23e60c6d6f0a8f298bae045abf374877d7f1))
* **metrics-operator:** update dynatrace api to support `range.step` ([#1812](https://github.com/keptn/lifecycle-toolkit/issues/1812)) ([4407fc4](https://github.com/keptn/lifecycle-toolkit/commit/4407fc4f3878ba5897991734df7e38fc273531a9))
* monorepo setup for lifecycle-operator, scheduler and runtimes ([#1857](https://github.com/keptn/lifecycle-toolkit/issues/1857)) ([84e243a](https://github.com/keptn/lifecycle-toolkit/commit/84e243a213ffba86eddd51ccc4bf4dbd61140069))
* update Datadog API to query metrics for range ([#1615](https://github.com/keptn/lifecycle-toolkit/issues/1615)) ([3b370ab](https://github.com/keptn/lifecycle-toolkit/commit/3b370abd2c887ea939f8a22c4c36babecb114265))
* update Dynatrace provider to query metrics over a range ([#1658](https://github.com/keptn/lifecycle-toolkit/issues/1658)) ([0f0cddb](https://github.com/keptn/lifecycle-toolkit/commit/0f0cddb2b6cd97aadff4da266c8d15f8ecab2881))
* update prometheus api to support `range.step` ([#1801](https://github.com/keptn/lifecycle-toolkit/issues/1801)) ([e64fcd6](https://github.com/keptn/lifecycle-toolkit/commit/e64fcd6d63cb87d76e83fc7e718f387953a0ff02))


### Bug Fixes

* **metrics-operator:** fix url encoding in DT metrics queries ([#1893](https://github.com/keptn/lifecycle-toolkit/issues/1893)) ([5437df9](https://github.com/keptn/lifecycle-toolkit/commit/5437df9419a71b4f9b2ae21534839704c5a05762))
* **metrics-operator:** flaky test in SLI converter ([#1954](https://github.com/keptn/lifecycle-toolkit/issues/1954)) ([cadb170](https://github.com/keptn/lifecycle-toolkit/commit/cadb170f4a7f26443a77634f9d7fb936610b9677))
* **metrics-operator:** flaky test in SLI converter ([#1961](https://github.com/keptn/lifecycle-toolkit/issues/1961)) ([d02a8ef](https://github.com/keptn/lifecycle-toolkit/commit/d02a8efa37f8cb5319f4f0b521b1b6c2696c5461))
* **metrics-operator:** make Fail target in AnalysisDefinition optional ([#1903](https://github.com/keptn/lifecycle-toolkit/issues/1903)) ([df874e2](https://github.com/keptn/lifecycle-toolkit/commit/df874e25dcfde7ae848dbb65c9fbe846d3ce1f61))


### Other

* add status field docs to all CRDs ([#1807](https://github.com/keptn/lifecycle-toolkit/issues/1807)) ([650ecba](https://github.com/keptn/lifecycle-toolkit/commit/650ecba95624ed3dc2bd61bf1f86578f450223a5))
* remove debug log containing secret ([#1967](https://github.com/keptn/lifecycle-toolkit/issues/1967)) ([75baefd](https://github.com/keptn/lifecycle-toolkit/commit/75baefd7b45de3dbe3e8ba0fc473ee25351a31e7))
* rename operator folder to lifecycle-operator ([#1819](https://github.com/keptn/lifecycle-toolkit/issues/1819)) ([97a2d25](https://github.com/keptn/lifecycle-toolkit/commit/97a2d25919c0a02165dd0dc6c7c82d57ad200139))


### Docs

* document `timeframe` feature for `KeptnMetric` ([#1703](https://github.com/keptn/lifecycle-toolkit/issues/1703)) ([077f0d5](https://github.com/keptn/lifecycle-toolkit/commit/077f0d5d0a49bc5b1f0e800274343660b8218c65))


### Performance

* **metrics-operator:** improve performance of storing analysis results ([#1905](https://github.com/keptn/lifecycle-toolkit/issues/1905)) ([efe3380](https://github.com/keptn/lifecycle-toolkit/commit/efe33803c4878599db958dcb76c7217b9e82c77f))


### Dependency Updates

* update dependency kubernetes-sigs/controller-tools to v0.12.1 ([#1765](https://github.com/keptn/lifecycle-toolkit/issues/1765)) ([ba79a32](https://github.com/keptn/lifecycle-toolkit/commit/ba79a32ef6acc9de8fb5d618b9ede7d6f96ce15e))
* update dependency kubernetes-sigs/controller-tools to v0.13.0 ([#1930](https://github.com/keptn/lifecycle-toolkit/issues/1930)) ([8b34b63](https://github.com/keptn/lifecycle-toolkit/commit/8b34b63404d0339633ef41ff1cf2005deae8d2b7))
* update dependency kubernetes-sigs/kustomize to v5.1.1 ([#1853](https://github.com/keptn/lifecycle-toolkit/issues/1853)) ([354ab3f](https://github.com/keptn/lifecycle-toolkit/commit/354ab3f980c2569e17a0354ece417df40317d120))
* update github.com/keptn/lifecycle-toolkit/klt-cert-manager digest to 0b618c4 ([#1654](https://github.com/keptn/lifecycle-toolkit/issues/1654)) ([c749313](https://github.com/keptn/lifecycle-toolkit/commit/c749313bfad7bd98b8d0ae7cc6dd2ea56f23e041))
* update github.com/keptn/lifecycle-toolkit/klt-cert-manager digest to 440c308 ([#2017](https://github.com/keptn/lifecycle-toolkit/issues/2017)) ([c365734](https://github.com/keptn/lifecycle-toolkit/commit/c365734fa7e3e40b2ae4c97c61628892d040dacc))
* update github.com/keptn/lifecycle-toolkit/klt-cert-manager digest to 88a54f9 ([#1794](https://github.com/keptn/lifecycle-toolkit/issues/1794)) ([fc976eb](https://github.com/keptn/lifecycle-toolkit/commit/fc976eb07ed9a5e49ed7d4ab1dbf187cee583e64))
* update github.com/keptn/lifecycle-toolkit/klt-cert-manager digest to 8dbec2d ([#1995](https://github.com/keptn/lifecycle-toolkit/issues/1995)) ([2f51445](https://github.com/keptn/lifecycle-toolkit/commit/2f5144540c4b3876e800bff29c30bfded334be40))
* update github.com/keptn/lifecycle-toolkit/klt-cert-manager digest to bb133cf ([#1963](https://github.com/keptn/lifecycle-toolkit/issues/1963)) ([c7697bf](https://github.com/keptn/lifecycle-toolkit/commit/c7697bf54d5fe18b7c62c5b11801c6c83079b0a3))
* update github.com/keptn/lifecycle-toolkit/klt-cert-manager digest to cba2de5 ([#1762](https://github.com/keptn/lifecycle-toolkit/issues/1762)) ([b77bcea](https://github.com/keptn/lifecycle-toolkit/commit/b77bceae39d6e4372b879afa326e7658d4ccdd89))
* update golang.org/x/exp digest to d852ddb ([#2011](https://github.com/keptn/lifecycle-toolkit/issues/2011)) ([12ee7b6](https://github.com/keptn/lifecycle-toolkit/commit/12ee7b6f37c2422a37da1e53e4b9e250db6d8ca3))
* update kubernetes packages (patch) ([#1786](https://github.com/keptn/lifecycle-toolkit/issues/1786)) ([cba2de5](https://github.com/keptn/lifecycle-toolkit/commit/cba2de5a5cd04c094131552aaf92c2b85ac23d21))
* update kubernetes packages to v0.26.8 (patch) ([#1945](https://github.com/keptn/lifecycle-toolkit/issues/1945)) ([6ce03d6](https://github.com/keptn/lifecycle-toolkit/commit/6ce03d600cbb3d3d3988573c616ec7f3830ba324))
* update module github.com/datadog/datadog-api-client-go/v2 to v2.15.0 ([#1803](https://github.com/keptn/lifecycle-toolkit/issues/1803)) ([ff62c60](https://github.com/keptn/lifecycle-toolkit/commit/ff62c6096f793307b54c0462dddf3473f2c2017d))
* update module github.com/datadog/datadog-api-client-go/v2 to v2.16.0 ([#1957](https://github.com/keptn/lifecycle-toolkit/issues/1957)) ([00f3cd3](https://github.com/keptn/lifecycle-toolkit/commit/00f3cd3f1d9f27046f31cd36924c0bed855faed2))
* update module github.com/open-feature/go-sdk to v1.7.0 ([#1854](https://github.com/keptn/lifecycle-toolkit/issues/1854)) ([8b90008](https://github.com/keptn/lifecycle-toolkit/commit/8b90008fc081e0da8029ffb1d28869a01efd15e6))
* update module golang.org/x/net to v0.12.0 ([#1662](https://github.com/keptn/lifecycle-toolkit/issues/1662)) ([49318bf](https://github.com/keptn/lifecycle-toolkit/commit/49318bfc40497a120304de9d831dfe033259220f))
* update module golang.org/x/net to v0.14.0 ([#1855](https://github.com/keptn/lifecycle-toolkit/issues/1855)) ([3186188](https://github.com/keptn/lifecycle-toolkit/commit/31861889bf7b227f489b941ac4a52db86551fcc2))
* update module golang.org/x/net to v0.14.0 ([#2020](https://github.com/keptn/lifecycle-toolkit/issues/2020)) ([14573cd](https://github.com/keptn/lifecycle-toolkit/commit/14573cd1b3a0fec12ceb7bd3f23c3fa8432d8528))
