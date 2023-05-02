# Changelog

## [0.7.0](https://github.com/keptn/lifecycle-toolkit/compare/v0.6.0...v0.7.0) (2023-03-16)


### ⚠ BREAKING CHANGES

* The different components of KLT have been renamed and use a new container image repository. For more information, please look at https://github.com/keptn/lifecycle-toolkit/issues/960
* The handling of the CRD lifecycle and metrics has been split into two different operators

### Features

* adapt examples to use KeptnMetric and KeptnMetricsProvider ([91e57ca](https://github.com/keptn/lifecycle-toolkit/commit/91e57cadba32fce6d873bc480408f90bcb8964d9))
* add CRD docs auto generator tooling ([#884](https://github.com/keptn/lifecycle-toolkit/issues/884)) ([5f63d9a](https://github.com/keptn/lifecycle-toolkit/commit/5f63d9a28a30a7022799d6debb365baadd72dd9b))
* add load test metrics ([#831](https://github.com/keptn/lifecycle-toolkit/issues/831)) ([2fa1a02](https://github.com/keptn/lifecycle-toolkit/commit/2fa1a02df06656d510cc2ddd2c868e37eb42970f))
* add YAMLLint ([#935](https://github.com/keptn/lifecycle-toolkit/issues/935)) ([48476bd](https://github.com/keptn/lifecycle-toolkit/commit/48476bd44f694ce2132b71ec92aed8259ae7fc2b))
* added the metrics-operator ([5153a05](https://github.com/keptn/lifecycle-toolkit/commit/5153a058d6eb30b6455941ee1d76dd09f98d6689))
* added the metrics-operator ([91e57ca](https://github.com/keptn/lifecycle-toolkit/commit/91e57cadba32fce6d873bc480408f90bcb8964d9))
* **cert-manager:** support certificate injection for all matching resources based on label selector ([91e57ca](https://github.com/keptn/lifecycle-toolkit/commit/91e57cadba32fce6d873bc480408f90bcb8964d9))
* fill in chart README ([#987](https://github.com/keptn/lifecycle-toolkit/issues/987)) ([2321180](https://github.com/keptn/lifecycle-toolkit/commit/23211800f2b897d1b146ad97d33b0e5f1994ad06))
* **helm-chart:** split documentation from value files ([#876](https://github.com/keptn/lifecycle-toolkit/issues/876)) ([c366739](https://github.com/keptn/lifecycle-toolkit/commit/c36673943e7ff54e2921fe6b21ad531603f367aa))
* improve naming and use new repository ([bd49357](https://github.com/keptn/lifecycle-toolkit/commit/bd493578df8825a52ec0f027583341a80b3c90f6))
* introduce lifecycle.keptn.sh/v1alpha3 API version ([91e57ca](https://github.com/keptn/lifecycle-toolkit/commit/91e57cadba32fce6d873bc480408f90bcb8964d9))
* **metrics-operator:** added conversion webhook for KeptnMetric CRDs ([91e57ca](https://github.com/keptn/lifecycle-toolkit/commit/91e57cadba32fce6d873bc480408f90bcb8964d9))
* **metrics-operator:** allow KeptnMetrics to be placed in any namespace ([91e57ca](https://github.com/keptn/lifecycle-toolkit/commit/91e57cadba32fce6d873bc480408f90bcb8964d9))
* **metrics-operator:** implement metric functionality ([91e57ca](https://github.com/keptn/lifecycle-toolkit/commit/91e57cadba32fce6d873bc480408f90bcb8964d9))
* **metrics-operator:** introduce KeptnMetricsProvider CRD ([91e57ca](https://github.com/keptn/lifecycle-toolkit/commit/91e57cadba32fce6d873bc480408f90bcb8964d9))
* **metrics-operator:** introduce migration from KeptnEvaluationProvider to KeptnMetricsProvider ([91e57ca](https://github.com/keptn/lifecycle-toolkit/commit/91e57cadba32fce6d873bc480408f90bcb8964d9))
* **operator:** accept LogLevels for all controllers  ([#790](https://github.com/keptn/lifecycle-toolkit/issues/790)) ([d175486](https://github.com/keptn/lifecycle-toolkit/commit/d175486fc10832458ebb95b17356fee4a2ccc1d7))
* **operator:** adapt KeptnEvaluationDefinition to reflect changes in KeptnMetric and KeptnMetricsProvider ([91e57ca](https://github.com/keptn/lifecycle-toolkit/commit/91e57cadba32fce6d873bc480408f90bcb8964d9))
* remove kube-rbac-proxy ([#909](https://github.com/keptn/lifecycle-toolkit/issues/909)) ([7d2621b](https://github.com/keptn/lifecycle-toolkit/commit/7d2621b70cdfd817aa9e1a408f4ed2841aef833b))
* use helmify to release our helm chart ([91e57ca](https://github.com/keptn/lifecycle-toolkit/commit/91e57cadba32fce6d873bc480408f90bcb8964d9))


### Bug Fixes

* adapted patch for mutating webhook to correctly add release namespace to exclusions ([#1044](https://github.com/keptn/lifecycle-toolkit/issues/1044)) ([d7cfc17](https://github.com/keptn/lifecycle-toolkit/commit/d7cfc171603cc85711e6d49b6c9cd857f312fc1b))
* added metric-operator prefix to related ClusterRole and ClusterRoleBindings ([#1042](https://github.com/keptn/lifecycle-toolkit/issues/1042)) ([92d16a3](https://github.com/keptn/lifecycle-toolkit/commit/92d16a3c6be57b823be17f6ac3a37134d1840438))
* fix cosign image signing after breaking changes ([#1047](https://github.com/keptn/lifecycle-toolkit/issues/1047)) ([e5abf85](https://github.com/keptn/lifecycle-toolkit/commit/e5abf85726f6a78673ba63a564c5274926726aa7))
* fix examples ([#1053](https://github.com/keptn/lifecycle-toolkit/issues/1053)) ([6f5c105](https://github.com/keptn/lifecycle-toolkit/commit/6f5c1059d427aca97513c606b224810f9446aefc))
* fix markdown issues in main ([#963](https://github.com/keptn/lifecycle-toolkit/issues/963)) ([ef35387](https://github.com/keptn/lifecycle-toolkit/commit/ef3538703ed87895828809c0204975015fc691be))
* fix some sonarcloud settings ([#1018](https://github.com/keptn/lifecycle-toolkit/issues/1018)) ([a40a8d3](https://github.com/keptn/lifecycle-toolkit/commit/a40a8d36458b880a468b8714c9fcfbb403776704))
* helm chart generation and helm pipeline ([#975](https://github.com/keptn/lifecycle-toolkit/issues/975)) ([444ba74](https://github.com/keptn/lifecycle-toolkit/commit/444ba745f7e120b7cba95291d06485002edb5f9e))
* helm chart generation fixes ([#977](https://github.com/keptn/lifecycle-toolkit/issues/977)) ([85e9d0e](https://github.com/keptn/lifecycle-toolkit/commit/85e9d0eb3da630aa4cf636dfbcb411205de24bd8))
* htmltest error for newly created documents ([#1010](https://github.com/keptn/lifecycle-toolkit/issues/1010)) ([4bf2919](https://github.com/keptn/lifecycle-toolkit/commit/4bf2919655b05890fc8803336091eaa8752fcae7))
* include namespace creation in release manifest ([#855](https://github.com/keptn/lifecycle-toolkit/issues/855)) ([d7a2b48](https://github.com/keptn/lifecycle-toolkit/commit/d7a2b486dd90ff173edbab49ff59988d58cc53c1))
* **metrics-operator:** adapt metric types from v1alpha1 ([91e57ca](https://github.com/keptn/lifecycle-toolkit/commit/91e57cadba32fce6d873bc480408f90bcb8964d9))
* **metrics-operator:** use correct port for serving metrics api ([#954](https://github.com/keptn/lifecycle-toolkit/issues/954)) ([d29ab64](https://github.com/keptn/lifecycle-toolkit/commit/d29ab64c6d295239586537c8002040d480fe17cd))
* move conversion webhooks to hub version API (v1alpha3) ([#992](https://github.com/keptn/lifecycle-toolkit/issues/992)) ([b2bb268](https://github.com/keptn/lifecycle-toolkit/commit/b2bb2685809abe7909a31518833236db8931f4c1))
* **operator:** compute deployment interval on deployment endtime ([#842](https://github.com/keptn/lifecycle-toolkit/issues/842)) ([140b2f2](https://github.com/keptn/lifecycle-toolkit/commit/140b2f28e1effd7877401bbbb8678d76a0ccab63))
* **operator:** invalid import of metrics ([91e57ca](https://github.com/keptn/lifecycle-toolkit/commit/91e57cadba32fce6d873bc480408f90bcb8964d9))
* remove missing 404 bug ([#1006](https://github.com/keptn/lifecycle-toolkit/issues/1006)) ([e8c0f38](https://github.com/keptn/lifecycle-toolkit/commit/e8c0f389e65c74f55e482ccb96127e31d475931d))
* remove required from release docs yaml parameters ([#952](https://github.com/keptn/lifecycle-toolkit/issues/952)) ([57cc938](https://github.com/keptn/lifecycle-toolkit/commit/57cc9389955576f3c708a6ec0352abe0dde2367c))
* wrong link in the local-setup ([#916](https://github.com/keptn/lifecycle-toolkit/issues/916)) ([de89309](https://github.com/keptn/lifecycle-toolkit/commit/de89309070e26cd76c5d964d9dc3a46f95897ddb))


### Dependency Updates

* update aquasecurity/trivy-action action to v0.9.2 ([#985](https://github.com/keptn/lifecycle-toolkit/issues/985)) ([6c79514](https://github.com/keptn/lifecycle-toolkit/commit/6c795141316cd463094607e6794002fb57beb8b6))
* update busybox docker tag to v1.36.0 ([#1023](https://github.com/keptn/lifecycle-toolkit/issues/1023)) ([83c1e15](https://github.com/keptn/lifecycle-toolkit/commit/83c1e1557f937d2719ba5febeb27f1defd8fa351))
* update curlimages/curl docker tag to v7.88.1 ([#1024](https://github.com/keptn/lifecycle-toolkit/issues/1024)) ([e89264d](https://github.com/keptn/lifecycle-toolkit/commit/e89264ddd5bce4d06224ee2e762cddeb36b3e2d7))
* update dawidd6/action-download-artifact action to v2.26.0 ([#903](https://github.com/keptn/lifecycle-toolkit/issues/903)) ([8c4ba83](https://github.com/keptn/lifecycle-toolkit/commit/8c4ba83cc3a1864b70379151f90b271eb39f39dc))
* update dependency argoproj/argo-cd to v2.6.2 ([#871](https://github.com/keptn/lifecycle-toolkit/issues/871)) ([9c813ac](https://github.com/keptn/lifecycle-toolkit/commit/9c813ac8a74be7c1ebe9c5eacd973273ed9ef3c8))
* update dependency argoproj/argo-cd to v2.6.3 ([#965](https://github.com/keptn/lifecycle-toolkit/issues/965)) ([4fc984f](https://github.com/keptn/lifecycle-toolkit/commit/4fc984f495c47ec24ad45e2e6d8f411c0b7bff1c))
* update dependency golangci/golangci-lint to v1.51.2 ([#765](https://github.com/keptn/lifecycle-toolkit/issues/765)) ([7b182fa](https://github.com/keptn/lifecycle-toolkit/commit/7b182fac52faee7e2be0917c4732ccf7d26fe924))
* update golang docker tag to v1.20.1 ([#844](https://github.com/keptn/lifecycle-toolkit/issues/844)) ([489f7f9](https://github.com/keptn/lifecycle-toolkit/commit/489f7f9100d97c79b57446db6ef1df957aa6b996))
* update golang.org/x/exp digest to 5e25df0 ([#833](https://github.com/keptn/lifecycle-toolkit/issues/833)) ([17c8185](https://github.com/keptn/lifecycle-toolkit/commit/17c81853b19f2af6057013c91c2d3a1c8f611f37))
* update klakegg/hugo docker tag to v0.107.0 ([#969](https://github.com/keptn/lifecycle-toolkit/issues/969)) ([018937b](https://github.com/keptn/lifecycle-toolkit/commit/018937b306af473874c09c4afb35af34c3d66ed4))
* update kubernetes packages (patch) ([#966](https://github.com/keptn/lifecycle-toolkit/issues/966)) ([7ba66c9](https://github.com/keptn/lifecycle-toolkit/commit/7ba66c936d9d9b8271c1f0b5a7d6966cb167d1af))
* update module github.com/onsi/ginkgo/v2 to v2.8.1 ([#867](https://github.com/keptn/lifecycle-toolkit/issues/867)) ([4c36b48](https://github.com/keptn/lifecycle-toolkit/commit/4c36b483ecacfb8639d26cde4cc0cf88bbb34826))
* update module github.com/onsi/gomega to v1.27.0 ([#872](https://github.com/keptn/lifecycle-toolkit/issues/872)) ([5b68118](https://github.com/keptn/lifecycle-toolkit/commit/5b6811856f24cd35e19a0af074dd689c8d176655))
* update module github.com/onsi/gomega to v1.27.1 ([#887](https://github.com/keptn/lifecycle-toolkit/issues/887)) ([4d2d0ed](https://github.com/keptn/lifecycle-toolkit/commit/4d2d0edc26bb0df43d89900cbd8324101de729ed))
* update module github.com/open-feature/flagd to v0.3.6 ([#810](https://github.com/keptn/lifecycle-toolkit/issues/810)) ([5d431b0](https://github.com/keptn/lifecycle-toolkit/commit/5d431b0e4099a1338bd949f8d2a67acdd6fdd9cd))
* update module github.com/open-feature/flagd to v0.3.7 ([#868](https://github.com/keptn/lifecycle-toolkit/issues/868)) ([8ed6455](https://github.com/keptn/lifecycle-toolkit/commit/8ed645573c3582952dd1519cd9aaf2ff336ace90))
* update module github.com/open-feature/go-sdk to v1.3.0 ([#767](https://github.com/keptn/lifecycle-toolkit/issues/767)) ([576a353](https://github.com/keptn/lifecycle-toolkit/commit/576a353326bf8dd60f3cf04e44342b86325a7bb2))
* update module github.com/prometheus/common to v0.40.0 ([#907](https://github.com/keptn/lifecycle-toolkit/issues/907)) ([d90355d](https://github.com/keptn/lifecycle-toolkit/commit/d90355d54c5af8ca4ce3d3fc9036e966dba65314))
* update module github.com/spf13/afero to v1.9.4 ([#911](https://github.com/keptn/lifecycle-toolkit/issues/911)) ([36cfe90](https://github.com/keptn/lifecycle-toolkit/commit/36cfe909611edc168824a7f570bed73f4c019264))
* update module k8s.io/klog/v2 to v2.90.1 ([#982](https://github.com/keptn/lifecycle-toolkit/issues/982)) ([90052bc](https://github.com/keptn/lifecycle-toolkit/commit/90052bc059af2d67f1835fe5ba72b5fa3eb77941))
* update sigstore/cosign-installer action to v3 ([#973](https://github.com/keptn/lifecycle-toolkit/issues/973)) ([e92259a](https://github.com/keptn/lifecycle-toolkit/commit/e92259a26da97f5b9f3e8cdcdb8797e254430abf))


### Docs

* adapt KeptnEvaluationDefinition and introduce KeptnMetricsProvider ([#944](https://github.com/keptn/lifecycle-toolkit/issues/944)) ([d56bfa4](https://github.com/keptn/lifecycle-toolkit/commit/d56bfa4bceb8b5bb6040fa7410ddfa745440cf7f))
* adapt metrics documentation and example ([#941](https://github.com/keptn/lifecycle-toolkit/issues/941)) ([82488ec](https://github.com/keptn/lifecycle-toolkit/commit/82488ec782c56295708c6f509d9d5be3f0b33fda))
* add "Intro to KLT"; edit "Getting Started" ([#785](https://github.com/keptn/lifecycle-toolkit/issues/785)) ([27932ff](https://github.com/keptn/lifecycle-toolkit/commit/27932ff7de4418bb314065a1b62ae401b80133b1))
* add cert-manager to jaeger installation script ([#1020](https://github.com/keptn/lifecycle-toolkit/issues/1020)) ([6dc6cee](https://github.com/keptn/lifecycle-toolkit/commit/6dc6ceefe4b6aa9191229be77d2f97466d32f07a))
* add CONTRIBUTING.md file for docs ([#758](https://github.com/keptn/lifecycle-toolkit/issues/758)) ([17fd319](https://github.com/keptn/lifecycle-toolkit/commit/17fd319cbd494c8663179f625ddde05d2279c3a3))
* add docs publishing information ([#949](https://github.com/keptn/lifecycle-toolkit/issues/949)) ([4351e18](https://github.com/keptn/lifecycle-toolkit/commit/4351e18c4097370520e63e48b947200a210f5380))
* add htmltest verification for documentation ([#932](https://github.com/keptn/lifecycle-toolkit/issues/932)) ([f342ccc](https://github.com/keptn/lifecycle-toolkit/commit/f342ccc0775eb41139ba0d679526dd95127bdfe8))
* add KLT components diagram ([#1016](https://github.com/keptn/lifecycle-toolkit/issues/1016)) ([dcf49cf](https://github.com/keptn/lifecycle-toolkit/commit/dcf49cfd0f90b5a648f14aa10c7bc4820acbf1ed))
* add Netlify configuration and advanced build ([#892](https://github.com/keptn/lifecycle-toolkit/issues/892)) ([81cd1f2](https://github.com/keptn/lifecycle-toolkit/commit/81cd1f2d1fd11e451269e580e10ea57cfbadff71))
* added more detailed explanation of how to make use of secrets in KeptnTasks ([#959](https://github.com/keptn/lifecycle-toolkit/issues/959)) ([06fa5fd](https://github.com/keptn/lifecycle-toolkit/commit/06fa5fd8a5d4134ea185e8909a6c2968f200bbda))
* change Development url ([#923](https://github.com/keptn/lifecycle-toolkit/issues/923)) ([335722d](https://github.com/keptn/lifecycle-toolkit/commit/335722dabb44a7d9b5d82d8b78e4e0f022462123))
* enhance contributors guide ([#866](https://github.com/keptn/lifecycle-toolkit/issues/866)) ([60bd934](https://github.com/keptn/lifecycle-toolkit/commit/60bd934058c34cb7e654f631c5dbe63ed2439606))
* fix broken link in README.md ([#913](https://github.com/keptn/lifecycle-toolkit/issues/913)) ([09a4f94](https://github.com/keptn/lifecycle-toolkit/commit/09a4f94055ae3c75682b084cfd62f87ea90203f8))
* improve netlify build ([#920](https://github.com/keptn/lifecycle-toolkit/issues/920)) ([39a002d](https://github.com/keptn/lifecycle-toolkit/commit/39a002d343df6248fe8caea78298f180e1260a09))
* initial list of related technologies of Keptn ([#795](https://github.com/keptn/lifecycle-toolkit/issues/795)) ([d4bd002](https://github.com/keptn/lifecycle-toolkit/commit/d4bd00262bdaa86458bc4c0eac459cc5575dec35))
* migrator for KeptnEvaluationProvider -&gt; KeptnMetricsProvider ([#945](https://github.com/keptn/lifecycle-toolkit/issues/945)) ([5bac785](https://github.com/keptn/lifecycle-toolkit/commit/5bac7858e87ef0b825adad0b0ff35bf6ae75d412))
* set up directories for contribution guide ([#1004](https://github.com/keptn/lifecycle-toolkit/issues/1004)) ([a3aa4e5](https://github.com/keptn/lifecycle-toolkit/commit/a3aa4e5b2d76443559727da1752921196ccffac4))
* update README and CONTRIBUTING instructions ([#991](https://github.com/keptn/lifecycle-toolkit/issues/991)) ([e42b750](https://github.com/keptn/lifecycle-toolkit/commit/e42b750a64f3681efdfa64dd55fe3ade61f53c53))
* use helm charts instead of manifests + document KeptnConfig CRD ([#747](https://github.com/keptn/lifecycle-toolkit/issues/747)) ([338c0fa](https://github.com/keptn/lifecycle-toolkit/commit/338c0fa2042ef74cb253d49ce050c2f61ea24f95))
* website build improvements ([#806](https://github.com/keptn/lifecycle-toolkit/issues/806)) ([03ce732](https://github.com/keptn/lifecycle-toolkit/commit/03ce732d0cc72988c49b012df70c776cfdc8eb06))


### Other

* add Hugo caching ([#958](https://github.com/keptn/lifecycle-toolkit/issues/958)) ([b2f24fe](https://github.com/keptn/lifecycle-toolkit/commit/b2f24fe4448edd24a3711e522caf393464ee877d))
* added sonar-project.properties file and adapted codecov settings ([#989](https://github.com/keptn/lifecycle-toolkit/issues/989)) ([ca1c6ba](https://github.com/keptn/lifecycle-toolkit/commit/ca1c6bad8e9f6983c2a781ea761201cabeeff954))
* adjust manifest limits ([#891](https://github.com/keptn/lifecycle-toolkit/issues/891)) ([32ce1b0](https://github.com/keptn/lifecycle-toolkit/commit/32ce1b01ea71fc0d52f5848144af6675289a39f0))
* close issues and PRs if they get stale ([#1041](https://github.com/keptn/lifecycle-toolkit/issues/1041)) ([89e03c2](https://github.com/keptn/lifecycle-toolkit/commit/89e03c21476cc6cd98ca6e1c1bef95384c8495f4))
* fix golangci-lint errors ([#905](https://github.com/keptn/lifecycle-toolkit/issues/905)) ([a133fdd](https://github.com/keptn/lifecycle-toolkit/commit/a133fdd99515765642d354c3a0cea51408333d99))
* improve Makefiles usage ([#921](https://github.com/keptn/lifecycle-toolkit/issues/921)) ([2761a4d](https://github.com/keptn/lifecycle-toolkit/commit/2761a4dad36f452b2dd575ab5ec1572b68602165))
* improve markdownlint ([#946](https://github.com/keptn/lifecycle-toolkit/issues/946)) ([d5d1675](https://github.com/keptn/lifecycle-toolkit/commit/d5d1675010cf0d3b8b506ef0a24c19d284d67727))
* move to new theme repo for docs ([74903a4](https://github.com/keptn/lifecycle-toolkit/commit/74903a481b69d3eb36c67652ea48b495b4f9fb3d))
* **operator:** remove KeptnMetric and KeptnEvaluationProvider from klt operator ([91e57ca](https://github.com/keptn/lifecycle-toolkit/commit/91e57cadba32fce6d873bc480408f90bcb8964d9))
* polish examples and integration tests ([#956](https://github.com/keptn/lifecycle-toolkit/issues/956)) ([72d3c9e](https://github.com/keptn/lifecycle-toolkit/commit/72d3c9ee086c203431120f6899a274180882fac4))
* release 0.7.0 ([#843](https://github.com/keptn/lifecycle-toolkit/issues/843)) ([bade181](https://github.com/keptn/lifecycle-toolkit/commit/bade181b735c7e069c510424ad5350476e41eeba))
* remove generated fake folder from sonar checks ([#1021](https://github.com/keptn/lifecycle-toolkit/issues/1021)) ([ec4ccb9](https://github.com/keptn/lifecycle-toolkit/commit/ec4ccb976117f88fc70afbdadd0b8c93da81edff))
* remove golang exp dependency ([#919](https://github.com/keptn/lifecycle-toolkit/issues/919)) ([c5c3fdf](https://github.com/keptn/lifecycle-toolkit/commit/c5c3fdfc822f8da629c1114f78ce31861e4c286a))
* run CI also on epic branches ([#853](https://github.com/keptn/lifecycle-toolkit/issues/853)) ([a2f7cce](https://github.com/keptn/lifecycle-toolkit/commit/a2f7cce17a7622ca8d5cbd9daaacc711d96b2660))
* set new documentation approach live ([#1007](https://github.com/keptn/lifecycle-toolkit/issues/1007)) ([f3511f1](https://github.com/keptn/lifecycle-toolkit/commit/f3511f1f5efec86fb1c86a6c7e39790d662417f9))
* switch to registry.k8s.io in yaml files in prometheus example ([#870](https://github.com/keptn/lifecycle-toolkit/issues/870)) ([909a1d6](https://github.com/keptn/lifecycle-toolkit/commit/909a1d6fd8788545e6d7cbee1351c7d574e1f39c))
* upgraded metrics operator to go 1.19 ([#1017](https://github.com/keptn/lifecycle-toolkit/issues/1017)) ([c2238fa](https://github.com/keptn/lifecycle-toolkit/commit/c2238fa2765bf5295720c9777e80f16f2b3ee289))

## [0.6.0](https://github.com/keptn/lifecycle-toolkit/compare/v0.5.0...v0.6.0) (2023-02-14)


### ⚠ BREAKING CHANGES

* The dependency on cert-manager has been removed in favor of a custom implementation. With these changes, the operator will be waiting for a certificate to be ready before registering the controllers. The certificate is generated as a k8s secret in the lifecycle-toolkit namespace and then loaded into an empty dir volume. The Keptn certificate manager will make sure to renew it (every 12 hours) and will take care of its validity every time the controller manager deployment resource changes.

### Features

* add cert-manager logic ([#528](https://github.com/keptn/lifecycle-toolkit/issues/528)) ([c0ece7a](https://github.com/keptn/lifecycle-toolkit/commit/c0ece7a9eae679f7bbb13328d961dcfce72c2fc8))
* add KeptnAppCreationRequestTimeout field to KeptnConfig API ([#735](https://github.com/keptn/lifecycle-toolkit/issues/735)) ([eda3f23](https://github.com/keptn/lifecycle-toolkit/commit/eda3f230af598977ba2a0d826eef7eafeb17c822))
* add KeptnConfig API ([#651](https://github.com/keptn/lifecycle-toolkit/issues/651)) ([9784216](https://github.com/keptn/lifecycle-toolkit/commit/9784216548364d18e941dbc0fc8a261e0396722b))
* add metadata to helm chart ([#737](https://github.com/keptn/lifecycle-toolkit/issues/737)) ([b5c5801](https://github.com/keptn/lifecycle-toolkit/commit/b5c580124b748ad8bce4fd5405d72dcf249d9498))
* add prometheus metrics evaluation example ([#677](https://github.com/keptn/lifecycle-toolkit/issues/677)) ([e5f644c](https://github.com/keptn/lifecycle-toolkit/commit/e5f644c5bf37c569fc2b328a0e6681488b1af8d0))
* add validating webhook for Keptn Metrics ([#668](https://github.com/keptn/lifecycle-toolkit/issues/668)) ([a4cc579](https://github.com/keptn/lifecycle-toolkit/commit/a4cc579b91a6156604b33a86f53af287cabd2989))
* annotate K8s Events ([#589](https://github.com/keptn/lifecycle-toolkit/issues/589)) ([4ea7da9](https://github.com/keptn/lifecycle-toolkit/commit/4ea7da92576d8fc16bc73ab37b711910e57859d4))
* automatically update documentation repository ([#610](https://github.com/keptn/lifecycle-toolkit/issues/610)) ([a84d4e4](https://github.com/keptn/lifecycle-toolkit/commit/a84d4e43b20e6cdb7641468e997a3b08ffe06d77))
* configurable imagePullPolicy via Helm ([#740](https://github.com/keptn/lifecycle-toolkit/issues/740)) ([b6b4160](https://github.com/keptn/lifecycle-toolkit/commit/b6b4160e8ae46e0d16fc06c5807c03e8599489b1))
* create an helm overlay ([#697](https://github.com/keptn/lifecycle-toolkit/issues/697)) ([9668ce8](https://github.com/keptn/lifecycle-toolkit/commit/9668ce8761c5526d625a8e26f26b244c2e93cc0c))
* **operator:** add KeptnConfig API ([#694](https://github.com/keptn/lifecycle-toolkit/issues/694)) ([4971a8b](https://github.com/keptn/lifecycle-toolkit/commit/4971a8b3915e9de152965e5c8cbc81de6bf03db9))
* **operator:** add logic to keptnmetrics controller ([#647](https://github.com/keptn/lifecycle-toolkit/issues/647)) ([ed5e200](https://github.com/keptn/lifecycle-toolkit/commit/ed5e20032c4a86c36d7cce4e76d1f8d0bf7a3933))
* **operator:** added adapter for custom metrics ([#682](https://github.com/keptn/lifecycle-toolkit/issues/682)) ([64cb972](https://github.com/keptn/lifecycle-toolkit/commit/64cb972a45e9377a40daf5c29e511ca9f578d773))
* **operator:** added Dynatrace DQL provider ([#783](https://github.com/keptn/lifecycle-toolkit/issues/783)) ([d19b533](https://github.com/keptn/lifecycle-toolkit/commit/d19b533e4469b21a59fb7f022373fc28ac11deec))
* **operator:** evaluation controller uses KeptnMetric as SLI provider ([#661](https://github.com/keptn/lifecycle-toolkit/issues/661)) ([da8fcee](https://github.com/keptn/lifecycle-toolkit/commit/da8fceedfe0a82ebd2072fecda7688a47c545aa5))
* **operator:** expose KeptnMetrics as OTel metrics ([#684](https://github.com/keptn/lifecycle-toolkit/issues/684)) ([eab9397](https://github.com/keptn/lifecycle-toolkit/commit/eab93970ec13f8c6486da89bf248883972534936))
* **operator:** introduce KeptnMetrics CRD and controller ([#643](https://github.com/keptn/lifecycle-toolkit/issues/643)) ([96170df](https://github.com/keptn/lifecycle-toolkit/commit/96170df5a10090de6618e986019cbc98e319bcb1))
* wire the new cert-manager into lifecycle operator ([#529](https://github.com/keptn/lifecycle-toolkit/issues/529)) ([752ea58](https://github.com/keptn/lifecycle-toolkit/commit/752ea5870b59ceb3a339d31f34a4c252dcd204d5))


### Bug Fixes

* add cert-manager to missing pipelines + fix linter issues ([#702](https://github.com/keptn/lifecycle-toolkit/issues/702)) ([a4ab1e3](https://github.com/keptn/lifecycle-toolkit/commit/a4ab1e36a2c20f83cf65b0d6b5b0d6c97186d2fc))
* broken link to examples folder in README.md ([#671](https://github.com/keptn/lifecycle-toolkit/issues/671)) ([4ff944b](https://github.com/keptn/lifecycle-toolkit/commit/4ff944b67c1742e31a79993a4edbc74d3f9a7b8b))
* fix klt-cert-manager release-local Makefile target ([#669](https://github.com/keptn/lifecycle-toolkit/issues/669)) ([a3b0f7b](https://github.com/keptn/lifecycle-toolkit/commit/a3b0f7be40309efbc6cedab8f420fc1bf2ccf8a1))
* fixed helm chart generation to include crds directly in the template ([#801](https://github.com/keptn/lifecycle-toolkit/issues/801)) ([f46e603](https://github.com/keptn/lifecycle-toolkit/commit/f46e603782badb1dbb0761725920b245c0efb97e))
* fixed helm chart patch ([#775](https://github.com/keptn/lifecycle-toolkit/issues/775)) ([fd3e2b0](https://github.com/keptn/lifecycle-toolkit/commit/fd3e2b087a949d3d9f2bb0db8fef9ff38bc647f1))
* **operator:** adapt resource requests and limits ([#835](https://github.com/keptn/lifecycle-toolkit/issues/835)) ([8249de6](https://github.com/keptn/lifecycle-toolkit/commit/8249de661f7843d6a244c9c4d0b62a5374f8b39f))
* **operator:** disable cache for secrets ([#727](https://github.com/keptn/lifecycle-toolkit/issues/727)) ([6ddbb6d](https://github.com/keptn/lifecycle-toolkit/commit/6ddbb6d2bdea53cd9152ed76ba1314ca66ad1bbc))
* **operator:** dynamically create tracers during reconciliation ([#804](https://github.com/keptn/lifecycle-toolkit/issues/804)) ([68f188e](https://github.com/keptn/lifecycle-toolkit/commit/68f188e403d4a0ca143d38328f6064fd0d925861))
* **operator:** fix calculation of deployment interval metrics ([#822](https://github.com/keptn/lifecycle-toolkit/issues/822)) ([a798eed](https://github.com/keptn/lifecycle-toolkit/commit/a798eed0e6358faed7d40ea60920ec09858665f3))
* **operator:** prevent re-execution of workload tasks that have been cancelled in a previous KLT version ([#718](https://github.com/keptn/lifecycle-toolkit/issues/718)) ([d89e179](https://github.com/keptn/lifecycle-toolkit/commit/d89e17909d9294dc172b62c621cfb8edf2eef533))
* **operator:** refactored metric adapter for helm generation ([#725](https://github.com/keptn/lifecycle-toolkit/issues/725)) ([e271162](https://github.com/keptn/lifecycle-toolkit/commit/e271162a3120618b617a4d3e501c2ab02de071fd))
* security pipeline issues ([#700](https://github.com/keptn/lifecycle-toolkit/issues/700)) ([ef5a7c5](https://github.com/keptn/lifecycle-toolkit/commit/ef5a7c5f816dca9a6767ca166e263d36843c720d))
* updated path to observability folder ([#780](https://github.com/keptn/lifecycle-toolkit/issues/780)) ([f2f09ea](https://github.com/keptn/lifecycle-toolkit/commit/f2f09ea4fcaca32db485f60d500c0f5a8ff29a68))


### Performance

* **operator:** only check for KeptnApp pre-evaluation if KWI has not entered its first phase yet ([#701](https://github.com/keptn/lifecycle-toolkit/issues/701)) ([a9f41d7](https://github.com/keptn/lifecycle-toolkit/commit/a9f41d7c42d2bad64da325881bd6a61f37f70b6b))
* requeue `KeptnMetric` and process them only when deadline is met ([#681](https://github.com/keptn/lifecycle-toolkit/issues/681)) ([39dd3f8](https://github.com/keptn/lifecycle-toolkit/commit/39dd3f842607ed2e9c93cfef72c7cf53a6a92ad5))


### Docs

* add context, update secret wording ([#781](https://github.com/keptn/lifecycle-toolkit/issues/781)) ([29b00cb](https://github.com/keptn/lifecycle-toolkit/commit/29b00cbb161202376524fb36ec5f8db8fa616489))
* add keptn certificate manager infos ([#652](https://github.com/keptn/lifecycle-toolkit/issues/652)) ([8cfb221](https://github.com/keptn/lifecycle-toolkit/commit/8cfb221d7ecb42093ee5c2f752fc3837e9d5a318))
* added breaking change message ([#726](https://github.com/keptn/lifecycle-toolkit/issues/726)) ([ebdebad](https://github.com/keptn/lifecycle-toolkit/commit/ebdebad0976d4a7b6f7b99e94331651199b181a7))
* added documentation to enable Slack notification post deployment task ([#787](https://github.com/keptn/lifecycle-toolkit/issues/787)) ([#788](https://github.com/keptn/lifecycle-toolkit/issues/788)) ([28a7319](https://github.com/keptn/lifecycle-toolkit/commit/28a7319f0f94271bbe76038ebba7dbdc2c38ada5))
* adjustments to folder structure ([#660](https://github.com/keptn/lifecycle-toolkit/issues/660)) ([1ec07ba](https://github.com/keptn/lifecycle-toolkit/commit/1ec07ba867a20b256fe7340a1ae63a39db706972))
* change port for KTL docs ([#713](https://github.com/keptn/lifecycle-toolkit/issues/713)) ([517e148](https://github.com/keptn/lifecycle-toolkit/commit/517e148fd13c771e1b5c4f8406315828f8fc8e6b))
* contribution guide for the community ([#709](https://github.com/keptn/lifecycle-toolkit/issues/709)) ([8b37dd7](https://github.com/keptn/lifecycle-toolkit/commit/8b37dd7cd7fc4646552995da6829bf8cfccedb6e))
* describe how to use Custom Metrics API integration ([#706](https://github.com/keptn/lifecycle-toolkit/issues/706)) ([d33af19](https://github.com/keptn/lifecycle-toolkit/commit/d33af197e6b4f64724ee9c0470d6347b86be01e8))
* documentation for restartable applications feature ([#645](https://github.com/keptn/lifecycle-toolkit/issues/645)) ([672bfa8](https://github.com/keptn/lifecycle-toolkit/commit/672bfa8f18923174d6ca9a73aaf1862cdf798e18))
* fix broken edit for developer ([#756](https://github.com/keptn/lifecycle-toolkit/issues/756)) ([ce6b7f0](https://github.com/keptn/lifecycle-toolkit/commit/ce6b7f0fccbb4015fb4f7f809a473a08384d2d9e))
* fix build status readme badge ([#590](https://github.com/keptn/lifecycle-toolkit/issues/590)) ([88e7ac4](https://github.com/keptn/lifecycle-toolkit/commit/88e7ac49eb3876b352f2418c853d9f306cf59386))
* modify getting-started.md ([#768](https://github.com/keptn/lifecycle-toolkit/issues/768)) ([f2263b5](https://github.com/keptn/lifecycle-toolkit/commit/f2263b50e5172f3d052b42f4a8f1eba997b8e21e))
* set up directory for Architecture docs ([#773](https://github.com/keptn/lifecycle-toolkit/issues/773)) ([7c3696e](https://github.com/keptn/lifecycle-toolkit/commit/7c3696e98b915ad2179bdf14cd093247e07056e7))
* set up directory for CRD Reference ([#800](https://github.com/keptn/lifecycle-toolkit/issues/800)) ([4748728](https://github.com/keptn/lifecycle-toolkit/commit/47487284f5c2c17b80cdaed736c7bb78fc1e5d7e))
* update broken task file ([#757](https://github.com/keptn/lifecycle-toolkit/issues/757)) ([082f091](https://github.com/keptn/lifecycle-toolkit/commit/082f091bd8781dcf9323da89fb3731abcc2d1ca4))
* update getting started wording ([#782](https://github.com/keptn/lifecycle-toolkit/issues/782)) ([dc6f664](https://github.com/keptn/lifecycle-toolkit/commit/dc6f664c5a3132e7a01078c25c4c4bb059cde8bf))


### Dependency Updates

* update amannn/action-semantic-pull-request action to v5.1.0 ([#837](https://github.com/keptn/lifecycle-toolkit/issues/837)) ([fa9bb6a](https://github.com/keptn/lifecycle-toolkit/commit/fa9bb6a4bfe58dacb570ae20e39623d75136162d))
* update anchore/sbom-action action to v0.13.3 ([#715](https://github.com/keptn/lifecycle-toolkit/issues/715)) ([bc75f00](https://github.com/keptn/lifecycle-toolkit/commit/bc75f001fb2bd6d8142254de1dc67d4065033247))
* update aquasecurity/trivy-action action to v0.9.0 ([#763](https://github.com/keptn/lifecycle-toolkit/issues/763)) ([1a79def](https://github.com/keptn/lifecycle-toolkit/commit/1a79def6e6075cf28d39fe8fdac2f54eccd01c2f))
* update aquasecurity/trivy-action action to v0.9.1 ([#834](https://github.com/keptn/lifecycle-toolkit/issues/834)) ([3db24f1](https://github.com/keptn/lifecycle-toolkit/commit/3db24f1200e69d25afd9b80bec62f0d8d7c4b1d4))
* update dawidd6/action-download-artifact action to v2.24.4 ([#836](https://github.com/keptn/lifecycle-toolkit/issues/836)) ([2296d8f](https://github.com/keptn/lifecycle-toolkit/commit/2296d8fe318aa9e9d1aeba2ff86fe74906e504d3))
* update dawidd6/action-download-artifact action to v2.25.0 ([#838](https://github.com/keptn/lifecycle-toolkit/issues/838)) ([d70753b](https://github.com/keptn/lifecycle-toolkit/commit/d70753bc252ca2f02150aa40cbca3ed8275b2c67))
* update dependency argoproj/argo-cd to v2.5.10 ([#766](https://github.com/keptn/lifecycle-toolkit/issues/766)) ([e4046ae](https://github.com/keptn/lifecycle-toolkit/commit/e4046ae677edb1213b9f1ad6ca837501dc11b3ba))
* update dependency argoproj/argo-cd to v2.5.6 ([#624](https://github.com/keptn/lifecycle-toolkit/issues/624)) ([278c74b](https://github.com/keptn/lifecycle-toolkit/commit/278c74bce5356a7e00f93604a10b56cde79a388e))
* update dependency argoproj/argo-cd to v2.5.7 ([#649](https://github.com/keptn/lifecycle-toolkit/issues/649)) ([ca5c106](https://github.com/keptn/lifecycle-toolkit/commit/ca5c106d389eb4dc790c8beb3daaaf51cbfbdb20))
* update dependency argoproj/argo-cd to v2.5.9 ([#714](https://github.com/keptn/lifecycle-toolkit/issues/714)) ([3e79f3d](https://github.com/keptn/lifecycle-toolkit/commit/3e79f3d7bccc491d1285f2acbfb5e8c9b21d3468))
* update dependency argoproj/argo-cd to v2.6.1 ([#816](https://github.com/keptn/lifecycle-toolkit/issues/816)) ([44acfbb](https://github.com/keptn/lifecycle-toolkit/commit/44acfbb5f05e2917b22fbd6fcdf69ed8cce9ffd6))
* update dependency cert-manager/cert-manager to v1.11.0 ([#627](https://github.com/keptn/lifecycle-toolkit/issues/627)) ([8720282](https://github.com/keptn/lifecycle-toolkit/commit/8720282a9c9817976d75457de392b7bd9989ca72))
* update dependency helm/helm to v3.10.3 ([#722](https://github.com/keptn/lifecycle-toolkit/issues/722)) ([491874c](https://github.com/keptn/lifecycle-toolkit/commit/491874c9ed7490ff54b88d0b511da67b14f25b1f))
* update dependency helm/helm to v3.11.0 ([#730](https://github.com/keptn/lifecycle-toolkit/issues/730)) ([ca85d3d](https://github.com/keptn/lifecycle-toolkit/commit/ca85d3de42879d9a5e3b5341918bd90a7e5e5274))
* update dependency helm/helm to v3.11.1 ([#819](https://github.com/keptn/lifecycle-toolkit/issues/819)) ([9f6b93f](https://github.com/keptn/lifecycle-toolkit/commit/9f6b93f674d3b88d464b1357bef3790858168b78))
* update dependency kubernetes-sigs/controller-tools to v0.11.2 ([#741](https://github.com/keptn/lifecycle-toolkit/issues/741)) ([bd0d218](https://github.com/keptn/lifecycle-toolkit/commit/bd0d2183f453fab9f25c016d185d070ca168e3c6))
* update dependency kubernetes-sigs/controller-tools to v0.11.3 ([#777](https://github.com/keptn/lifecycle-toolkit/issues/777)) ([207d2ca](https://github.com/keptn/lifecycle-toolkit/commit/207d2ca1b777a0f9cf3881c98345bf8e91a1033e))
* update dependency pyyaml to v5.4.1 ([#642](https://github.com/keptn/lifecycle-toolkit/issues/642)) ([d854d7a](https://github.com/keptn/lifecycle-toolkit/commit/d854d7ae379cc0401e49c1c8cfccece39bb82619))
* update dependency pyyaml to v6 ([#648](https://github.com/keptn/lifecycle-toolkit/issues/648)) ([850cf7a](https://github.com/keptn/lifecycle-toolkit/commit/850cf7a8e32f31f02bb3fd2cdc879293b062a432))
* update docker/build-push-action action to v4 ([#736](https://github.com/keptn/lifecycle-toolkit/issues/736)) ([cad355a](https://github.com/keptn/lifecycle-toolkit/commit/cad355a2a8ee15296015db62cb86e1038c97fba3))
* update golang docker tag to v1.19.5 ([#587](https://github.com/keptn/lifecycle-toolkit/issues/587)) ([567a211](https://github.com/keptn/lifecycle-toolkit/commit/567a2115b156fe950f52b2408247d859fb168b72))
* update golang docker tag to v1.20.0 ([#742](https://github.com/keptn/lifecycle-toolkit/issues/742)) ([c9873ae](https://github.com/keptn/lifecycle-toolkit/commit/c9873ae18e5f353e5a0539eda680bac0444817a8))
* update golang.org/x/exp digest to 1de6713 ([#641](https://github.com/keptn/lifecycle-toolkit/issues/641)) ([a1417dd](https://github.com/keptn/lifecycle-toolkit/commit/a1417ddfd250c1a9528a0204ecf48444b3adfd6b))
* update golang.org/x/exp digest to 46f607a ([#760](https://github.com/keptn/lifecycle-toolkit/issues/760)) ([77196c7](https://github.com/keptn/lifecycle-toolkit/commit/77196c70fed0f1f030f03de2900554138cd1e230))
* update golang.org/x/exp digest to 54bba9f ([#753](https://github.com/keptn/lifecycle-toolkit/issues/753)) ([79b9021](https://github.com/keptn/lifecycle-toolkit/commit/79b9021b777439443630b0f682499a070aefb62c))
* update golang.org/x/exp digest to a67bb56 ([#691](https://github.com/keptn/lifecycle-toolkit/issues/691)) ([cbe2ed3](https://github.com/keptn/lifecycle-toolkit/commit/cbe2ed33ea7c14dbe073936fa653edef2739d883))
* update golang.org/x/exp digest to a684f29 ([#815](https://github.com/keptn/lifecycle-toolkit/issues/815)) ([65495a8](https://github.com/keptn/lifecycle-toolkit/commit/65495a84dac2e6db10833be12a889ed759eb3dbd))
* update golang.org/x/exp digest to a68e582 ([#653](https://github.com/keptn/lifecycle-toolkit/issues/653)) ([d6fabf5](https://github.com/keptn/lifecycle-toolkit/commit/d6fabf55d6a45d7bd5594fb59121eba2447e82d2))
* update golang.org/x/exp digest to f062dba ([#710](https://github.com/keptn/lifecycle-toolkit/issues/710)) ([9135eaf](https://github.com/keptn/lifecycle-toolkit/commit/9135eaf169eb2c710f1484390374ca8c6b3157f3))
* update kubernetes packages (patch) ([#704](https://github.com/keptn/lifecycle-toolkit/issues/704)) ([7370933](https://github.com/keptn/lifecycle-toolkit/commit/7370933f6b4d6461ee17c9d16c196c8e1bb3eb3c))
* update kubernetes packages to v0.25.6 (patch) ([#663](https://github.com/keptn/lifecycle-toolkit/issues/663)) ([2bcc1dd](https://github.com/keptn/lifecycle-toolkit/commit/2bcc1dd248aa6c173e9e0566d1b6601afac09031))
* update module github.com/benbjohnson/clock to v1.3.0 ([#705](https://github.com/keptn/lifecycle-toolkit/issues/705)) ([cd989be](https://github.com/keptn/lifecycle-toolkit/commit/cd989bebc0641683fc8e073800019bd859bc91cd))
* update module github.com/onsi/ginkgo/v2 to v2.7.0 ([#611](https://github.com/keptn/lifecycle-toolkit/issues/611)) ([9ace485](https://github.com/keptn/lifecycle-toolkit/commit/9ace4859188654ac74c2bede289aba3dfcd026a3))
* update module github.com/onsi/ginkgo/v2 to v2.7.1 ([#729](https://github.com/keptn/lifecycle-toolkit/issues/729)) ([59c853c](https://github.com/keptn/lifecycle-toolkit/commit/59c853ce54ed01076cd8b0a36513173efaf62f02))
* update module github.com/onsi/ginkgo/v2 to v2.8.0 ([#754](https://github.com/keptn/lifecycle-toolkit/issues/754)) ([5613491](https://github.com/keptn/lifecycle-toolkit/commit/5613491e4cde12d10c106aa341520181b5f182d8))
* update module github.com/onsi/gomega to v1.26.0 ([#672](https://github.com/keptn/lifecycle-toolkit/issues/672)) ([3b092bd](https://github.com/keptn/lifecycle-toolkit/commit/3b092bd1faf46a04523819f83025ca67585767d8))
* update module github.com/open-feature/flagd to v0.3.4 ([#716](https://github.com/keptn/lifecycle-toolkit/issues/716)) ([026b25d](https://github.com/keptn/lifecycle-toolkit/commit/026b25d70b83180b8cb7ce80ec40e659ac4b2f24))
* update module github.com/spf13/afero to v1.9.3 ([#654](https://github.com/keptn/lifecycle-toolkit/issues/654)) ([2831566](https://github.com/keptn/lifecycle-toolkit/commit/283156638a960725e8bcbc11ae82277f930c9398))
* update module google.golang.org/grpc to v1.52.0 ([#626](https://github.com/keptn/lifecycle-toolkit/issues/626)) ([ba65776](https://github.com/keptn/lifecycle-toolkit/commit/ba657761fcecd1e52835e5e5d27a785329b597c4))
* update module google.golang.org/grpc to v1.52.3 ([#711](https://github.com/keptn/lifecycle-toolkit/issues/711)) ([6c3009d](https://github.com/keptn/lifecycle-toolkit/commit/6c3009dfb39ade9d2cc765f0897687e98571a19e))
* update module k8s.io/component-helpers to v0.25.6 ([#676](https://github.com/keptn/lifecycle-toolkit/issues/676)) ([bda60c8](https://github.com/keptn/lifecycle-toolkit/commit/bda60c8301ce7cc247c3748fa25304c284750f75))
* update module k8s.io/klog/v2 to v2.90.0 ([#685](https://github.com/keptn/lifecycle-toolkit/issues/685)) ([98164ed](https://github.com/keptn/lifecycle-toolkit/commit/98164ed938a4f247a892e1ebf1a47bb0e5af1f67))
* update module sigs.k8s.io/controller-runtime to v0.14.2 ([#723](https://github.com/keptn/lifecycle-toolkit/issues/723)) ([31fac62](https://github.com/keptn/lifecycle-toolkit/commit/31fac62d81e67354235f55b40ecdfc921a75b265))
* update module sigs.k8s.io/controller-runtime to v0.14.3 ([#762](https://github.com/keptn/lifecycle-toolkit/issues/762)) ([37e783e](https://github.com/keptn/lifecycle-toolkit/commit/37e783ee09107a9eac0d1e0392fb4958d63355d4))
* update module sigs.k8s.io/controller-runtime to v0.14.4 ([#811](https://github.com/keptn/lifecycle-toolkit/issues/811)) ([0191385](https://github.com/keptn/lifecycle-toolkit/commit/0191385137751d8cc3b4b3a2fe35403fe2ad5771))


### Other

* add maturity status for each feature ([#825](https://github.com/keptn/lifecycle-toolkit/issues/825)) ([d5849c9](https://github.com/keptn/lifecycle-toolkit/commit/d5849c9e4f6d1f7c05c30fb2534284ff3b8cc25b))
* add StackScribe as Code Owner of docs ([#821](https://github.com/keptn/lifecycle-toolkit/issues/821)) ([c90db1e](https://github.com/keptn/lifecycle-toolkit/commit/c90db1e6d848248c05032572dfad86fbda610bf0))
* fail linter CI jobs when checks are failing ([#630](https://github.com/keptn/lifecycle-toolkit/issues/630)) ([a5e0eaa](https://github.com/keptn/lifecycle-toolkit/commit/a5e0eaaf040026fedfa08f065b97b8de315b4132))
* fix markdown linter errors ([#824](https://github.com/keptn/lifecycle-toolkit/issues/824)) ([5df8789](https://github.com/keptn/lifecycle-toolkit/commit/5df87899038bd6205b51d87f0954d2b3b616868a))
* **operator:** fix linter issues ([#579](https://github.com/keptn/lifecycle-toolkit/issues/579)) ([64603fb](https://github.com/keptn/lifecycle-toolkit/commit/64603fbe728631b7f84e79552b39ad34957e60ea))
* **scheduler:** added new scheduler tests ([#634](https://github.com/keptn/lifecycle-toolkit/issues/634)) ([2e47b92](https://github.com/keptn/lifecycle-toolkit/commit/2e47b9227b7622af703311b1be69b7e0706c4700))
* stop linter CI job on error ([#631](https://github.com/keptn/lifecycle-toolkit/issues/631)) ([c5463c6](https://github.com/keptn/lifecycle-toolkit/commit/c5463c621267859ac56d53a4ec6b0458e62db9f9))
* update codeowners to have default owners for docs folder ([#827](https://github.com/keptn/lifecycle-toolkit/issues/827)) ([82351f4](https://github.com/keptn/lifecycle-toolkit/commit/82351f4458d715607f356a04e343477b2937f564))
* update observability demo link ([#666](https://github.com/keptn/lifecycle-toolkit/issues/666)) ([53e53f1](https://github.com/keptn/lifecycle-toolkit/commit/53e53f1ccb81e397f415dcaea53a6a7c1589daf3))

## [0.5.0](https://github.com/keptn/lifecycle-toolkit/compare/v0.4.1...v0.5.0) (2023-01-10)


### ⚠ BREAKING CHANGES

* Evaluation and Task statuses in KeptnWorkloadVersion/KeptnAppVersion use the same structure
* **operator:** With API version `v1alpha2`, `KeptnEvaluationProvider` uses a Secret Selector instead of `SecretName`.
* We are introducing scarf.sh for download tracking to have valuable information that will show us adoption rates for the toolkit.

### Features

* add a new operator for self-hosted certificate management ([#523](https://github.com/keptn/lifecycle-toolkit/issues/523)) ([90bbdba](https://github.com/keptn/lifecycle-toolkit/commit/90bbdba2ab560cc2650ba45b2126ebcb0c90a1da))
* add revision field to KeptnApp ([#494](https://github.com/keptn/lifecycle-toolkit/issues/494)) ([23ddfa3](https://github.com/keptn/lifecycle-toolkit/commit/23ddfa3b6a9a445b99eea1332b776299b4e4558a))
* generate SBOMs for container images on release ([#571](https://github.com/keptn/lifecycle-toolkit/issues/571)) ([72fe001](https://github.com/keptn/lifecycle-toolkit/commit/72fe001eee0c97a8efaa4a572aec05095b14c9be))
* introduce API groups in KLT ([#547](https://github.com/keptn/lifecycle-toolkit/issues/547)) ([b482d96](https://github.com/keptn/lifecycle-toolkit/commit/b482d96d6e76564e0f76e1aad6af3ed2b50be84a))
* introduce API version v1alpha2 ([#491](https://github.com/keptn/lifecycle-toolkit/issues/491)) ([229bcc9](https://github.com/keptn/lifecycle-toolkit/commit/229bcc9778780e9a5c5233f8753692cf578b60b5))
* **operator:** add version conversion rule for KeptnEvaluationProvider ([#531](https://github.com/keptn/lifecycle-toolkit/issues/531)) ([f1e9fe1](https://github.com/keptn/lifecycle-toolkit/commit/f1e9fe15023679ba1920fceb3306afadee025582))
* **operator:** emit K8s events with detailed messages for failed evaluations ([#477](https://github.com/keptn/lifecycle-toolkit/issues/477)) ([1b3a56f](https://github.com/keptn/lifecycle-toolkit/commit/1b3a56fc41ae6322dcf7f0fdc947cf1100ed49bd))
* **operator:** fix ownership information for keptnworkloads ([#520](https://github.com/keptn/lifecycle-toolkit/issues/520)) ([1e642c7](https://github.com/keptn/lifecycle-toolkit/commit/1e642c7bc254601a14c374ae5b928e53a03a3e52))
* **operator:** introduce evaluation support for Dynatrace ([#194](https://github.com/keptn/lifecycle-toolkit/issues/194)) ([c6483cc](https://github.com/keptn/lifecycle-toolkit/commit/c6483cc8e569e3cae5284315eeb873e051aea7d4))
* **operator:** support restartability of KeptnApp ([#544](https://github.com/keptn/lifecycle-toolkit/issues/544)) ([99070c2](https://github.com/keptn/lifecycle-toolkit/commit/99070c2a5a98c25fc520a5c8c6825917fdc7726c))
* set default of 1 for KeptnApp revision field ([#513](https://github.com/keptn/lifecycle-toolkit/issues/513)) ([a5cb3f2](https://github.com/keptn/lifecycle-toolkit/commit/a5cb3f23c164f9c9a4f5b1f7a8bf5251e9b50885))
* use scarf.sh registry for all container images ([#507](https://github.com/keptn/lifecycle-toolkit/issues/507)) ([647c6a1](https://github.com/keptn/lifecycle-toolkit/commit/647c6a1f84ead8bed394e32528d5dd85854765b6))


### Bug Fixes

* **dashboards:** use fixed color mode for succeeded AppVersion/WorkloadVersion tiles ([#515](https://github.com/keptn/lifecycle-toolkit/issues/515)) ([8cdb23e](https://github.com/keptn/lifecycle-toolkit/commit/8cdb23ee61cc7ee22be2b7326bbf202ee3ddf09f))


### Dependency Updates

* add more renovate annotations for auto updates ([#519](https://github.com/keptn/lifecycle-toolkit/issues/519)) ([1555d63](https://github.com/keptn/lifecycle-toolkit/commit/1555d63cb034cf70ea4664fddd046a5320dc3fc0))
* add yaml files to renovate ([#527](https://github.com/keptn/lifecycle-toolkit/issues/527)) ([2ee11f5](https://github.com/keptn/lifecycle-toolkit/commit/2ee11f542f3f3532f65b1c864f28316744f4dc89))
* update dawidd6/action-download-artifact action to v2.24.3 ([#560](https://github.com/keptn/lifecycle-toolkit/issues/560)) ([5220e9d](https://github.com/keptn/lifecycle-toolkit/commit/5220e9d414921ef209b5d2d0438004d2db4cd152))
* update denoland/deno docker tag to v1.28.3 ([#467](https://github.com/keptn/lifecycle-toolkit/issues/467)) ([59fa6b0](https://github.com/keptn/lifecycle-toolkit/commit/59fa6b05d13a468daf004ff4b180c3c161040f96))
* update denoland/deno docker tag to v1.29.1 ([#534](https://github.com/keptn/lifecycle-toolkit/issues/534)) ([3b316f9](https://github.com/keptn/lifecycle-toolkit/commit/3b316f930a1fa96c4ff1d644ec7c7565d3e7f040))
* update dependency argoproj/argo-cd to v2.5.4 ([#492](https://github.com/keptn/lifecycle-toolkit/issues/492)) ([6f16dac](https://github.com/keptn/lifecycle-toolkit/commit/6f16dac49514fbd5f04c08e65cfe567db6b5c0f1))
* update dependency argoproj/argo-cd to v2.5.5 ([#543](https://github.com/keptn/lifecycle-toolkit/issues/543)) ([f1b59ca](https://github.com/keptn/lifecycle-toolkit/commit/f1b59ca4c6fa25a8584db15e9bad50390113fee6))
* update dependency cert-manager/cert-manager to v1.10.1 ([#530](https://github.com/keptn/lifecycle-toolkit/issues/530)) ([cb83e24](https://github.com/keptn/lifecycle-toolkit/commit/cb83e24c24c65dffb3f1923753905c31765020e3))
* update dependency jaegertracing/jaeger to v1.40.0 ([#506](https://github.com/keptn/lifecycle-toolkit/issues/506)) ([ddb9eca](https://github.com/keptn/lifecycle-toolkit/commit/ddb9eca674e7d7af42d99b50d629d21a6824cf09))
* update dependency jaegertracing/jaeger to v1.41.0 ([#565](https://github.com/keptn/lifecycle-toolkit/issues/565)) ([cc5f7ca](https://github.com/keptn/lifecycle-toolkit/commit/cc5f7ca831c719871bd2e4e403cf47bd4792d390))
* update dependency kubernetes-sigs/kustomize to v4 ([#575](https://github.com/keptn/lifecycle-toolkit/issues/575)) ([36a6169](https://github.com/keptn/lifecycle-toolkit/commit/36a61698e2cb1d0c25ac51c8d98d7a40f97b0bc7))
* update dependency kudobuilder/kuttl to v0.14.0 ([#561](https://github.com/keptn/lifecycle-toolkit/issues/561)) ([ff6b95a](https://github.com/keptn/lifecycle-toolkit/commit/ff6b95a84b483eb47c5504a85756338e89318bde))
* update dependency kudobuilder/kuttl to v0.15.0 ([#566](https://github.com/keptn/lifecycle-toolkit/issues/566)) ([9516fcf](https://github.com/keptn/lifecycle-toolkit/commit/9516fcf8776875fd28dea12b2303acf5935295f8))
* update golang docker tag to v1.19.4 ([#495](https://github.com/keptn/lifecycle-toolkit/issues/495)) ([5a74869](https://github.com/keptn/lifecycle-toolkit/commit/5a74869861cf9af29220244c3803003927d1f783))
* update jasonetco/create-an-issue action to v2.9.1 ([#557](https://github.com/keptn/lifecycle-toolkit/issues/557)) ([f6d5934](https://github.com/keptn/lifecycle-toolkit/commit/f6d59345f3954673122097e686e769c3a1ed14e9))
* update kubernetes packages to v0.25.5 (patch) ([#499](https://github.com/keptn/lifecycle-toolkit/issues/499)) ([627b9e1](https://github.com/keptn/lifecycle-toolkit/commit/627b9e163b66121609b13a8d4496cee72fdf8f55))
* update module github.com/magiconair/properties to v1.8.7 ([#503](https://github.com/keptn/lifecycle-toolkit/issues/503)) ([4f87239](https://github.com/keptn/lifecycle-toolkit/commit/4f872397b28063b998091ed43a766b7f9c72fe75))
* update module github.com/onsi/ginkgo/v2 to v2.5.1 ([#384](https://github.com/keptn/lifecycle-toolkit/issues/384)) ([955d41e](https://github.com/keptn/lifecycle-toolkit/commit/955d41e4c1d86596e6b00aee78033f268a5b3d80))
* update module github.com/onsi/ginkgo/v2 to v2.6.1 ([#535](https://github.com/keptn/lifecycle-toolkit/issues/535)) ([e02929a](https://github.com/keptn/lifecycle-toolkit/commit/e02929a0695b54ec455c4e256e140d19831d39b0))
* update module github.com/onsi/gomega to v1.24.2 ([#532](https://github.com/keptn/lifecycle-toolkit/issues/532)) ([2480f21](https://github.com/keptn/lifecycle-toolkit/commit/2480f216a9b619d6268a9edd6e8a9a9366fe58f8))
* update module github.com/prometheus/common to v0.37.1 ([#533](https://github.com/keptn/lifecycle-toolkit/issues/533)) ([b72d52e](https://github.com/keptn/lifecycle-toolkit/commit/b72d52eb1e940247f4accbf7e883edcf965dbd52))
* update module github.com/prometheus/common to v0.39.0 ([#502](https://github.com/keptn/lifecycle-toolkit/issues/502)) ([28ab629](https://github.com/keptn/lifecycle-toolkit/commit/28ab6292d5c06e328f57794e0c309ef0b124f9d5))
* update module github.com/stretchr/testify to v1.8.1 ([#551](https://github.com/keptn/lifecycle-toolkit/issues/551)) ([523bb55](https://github.com/keptn/lifecycle-toolkit/commit/523bb5520e81ce4078188c744b54ca5c9a79212e))
* update module google.golang.org/grpc to v1.51.0 ([#451](https://github.com/keptn/lifecycle-toolkit/issues/451)) ([3828ee5](https://github.com/keptn/lifecycle-toolkit/commit/3828ee58a1737694f93072208f54229d6f0fae8e))
* update module k8s.io/api to v0.25.5 ([#573](https://github.com/keptn/lifecycle-toolkit/issues/573)) ([5159a5e](https://github.com/keptn/lifecycle-toolkit/commit/5159a5e1d9aa774b719e2744b2d264b3559f8467))
* update module k8s.io/component-helpers to v0.25.5 ([#504](https://github.com/keptn/lifecycle-toolkit/issues/504)) ([02b5982](https://github.com/keptn/lifecycle-toolkit/commit/02b59826329a70742dd39edcb9c3a32ba86b0a6b))
* update module k8s.io/kubernetes to v1.25.5 ([#501](https://github.com/keptn/lifecycle-toolkit/issues/501)) ([df8e51c](https://github.com/keptn/lifecycle-toolkit/commit/df8e51c13a5a9b88bd96b67add713251682b274e))
* update opentelemetry-go monorepo to v0.34.0 (minor) ([#498](https://github.com/keptn/lifecycle-toolkit/issues/498)) ([e7db4d0](https://github.com/keptn/lifecycle-toolkit/commit/e7db4d099f1586498f16047c3d0669cdb9a88147))


### Other

* add scarf pixel to markdown files ([#493](https://github.com/keptn/lifecycle-toolkit/issues/493)) ([b05a810](https://github.com/keptn/lifecycle-toolkit/commit/b05a810860d1286abdc72e9d1fe0ef4204453018))
* enhance golangci-lint with code complexity and other measures ([#484](https://github.com/keptn/lifecycle-toolkit/issues/484)) ([1d711d0](https://github.com/keptn/lifecycle-toolkit/commit/1d711d09b181ea72277ac6cb43f1ac605e82955f))
* **operator:** refactor operator and scheduler statuses + add unit tests ([#548](https://github.com/keptn/lifecycle-toolkit/issues/548)) ([c661dc0](https://github.com/keptn/lifecycle-toolkit/commit/c661dc063544835e3854b20c06c22ed536529511))
* unify EvaluationStatus and TaskStatus to single structure ([#569](https://github.com/keptn/lifecycle-toolkit/issues/569)) ([9b31b04](https://github.com/keptn/lifecycle-toolkit/commit/9b31b04546fb600735e18ed3d522b85453bf5be5))
* upgrade examples and tests to v1alpha2 ([#509](https://github.com/keptn/lifecycle-toolkit/issues/509)) ([2a133ea](https://github.com/keptn/lifecycle-toolkit/commit/2a133eaef43fa800fe937f7a5cde115d78e4b5fb))

## [0.4.1](https://github.com/keptn/lifecycle-toolkit/compare/v0.4.0...v0.4.1) (2022-11-30)


### Features

* Move dashboards and fix issues ([#417](https://github.com/keptn/lifecycle-toolkit/issues/417)) ([f6b5bfc](https://github.com/keptn/lifecycle-toolkit/commit/f6b5bfcd0f3254970101c3ac53bdec8d1426b3de))
* **operator:** Copy annotations from parent resource if not defined on pod ([#305](https://github.com/keptn/lifecycle-toolkit/issues/305)) ([c21f015](https://github.com/keptn/lifecycle-toolkit/commit/c21f015a9d4efcc2b59f9c5be41da758dca8e618))
* **operator:** include detailed information about task/evaluation failure in span ([#445](https://github.com/keptn/lifecycle-toolkit/issues/445)) ([94de8d6](https://github.com/keptn/lifecycle-toolkit/commit/94de8d6528a4de293372cffa50f3c12cd24909f5))
* **operator:** refactor existing interfaces ([#419](https://github.com/keptn/lifecycle-toolkit/issues/419)) ([f9c28a8](https://github.com/keptn/lifecycle-toolkit/commit/f9c28a8b677cc82e50201deb743c12458b4dffb4))
* **operator:** Refactor metrics helper functions ([#269](https://github.com/keptn/lifecycle-toolkit/issues/269)) ([b6f3f43](https://github.com/keptn/lifecycle-toolkit/commit/b6f3f43e29737839b25fea16c7b3810f193b313f))
* **operator:** Refactor Task, Evaluation handling + adapt span attributes setting ([#287](https://github.com/keptn/lifecycle-toolkit/issues/287)) ([4d16a77](https://github.com/keptn/lifecycle-toolkit/commit/4d16a779d28738bfbc06a58f0ea2acb0abb08969))
* **operator:** rework Task and Evaluation span structure ([#465](https://github.com/keptn/lifecycle-toolkit/issues/465)) ([e5717c6](https://github.com/keptn/lifecycle-toolkit/commit/e5717c620ce16946cde56e4bbc56e9aa8527b2a8))
* **operator:** rework Workload and Application span structure ([#452](https://github.com/keptn/lifecycle-toolkit/issues/452)) ([9a483ce](https://github.com/keptn/lifecycle-toolkit/commit/9a483ceffc25d9524f894a49603fee07d7032e26))


### Bug Fixes

* adapt name of keptn_app_count metric due to reverted OTel exporter dependency update ([#482](https://github.com/keptn/lifecycle-toolkit/issues/482)) ([97f8e8c](https://github.com/keptn/lifecycle-toolkit/commit/97f8e8cb0b1acbfc868a438dc52dc6dc9c0d2b5e))
* Added back permission to list and watch namespaces ([#404](https://github.com/keptn/lifecycle-toolkit/issues/404)) ([df346f7](https://github.com/keptn/lifecycle-toolkit/commit/df346f7986b1a81797a17356fcca49d7e08062b9))
* Fixed problems in examples ([#378](https://github.com/keptn/lifecycle-toolkit/issues/378)) ([277be10](https://github.com/keptn/lifecycle-toolkit/commit/277be10128a2ed7e8e91181891fdcd5be27978ca))
* **operator:** Also consider StatefulSets/DaemonSets when checking Workload Deployment state ([#406](https://github.com/keptn/lifecycle-toolkit/issues/406)) ([27c189f](https://github.com/keptn/lifecycle-toolkit/commit/27c189f93e363ecb7dde21186207ecb83d82f071))
* **operator:** build env variables are computed during docker build ([#457](https://github.com/keptn/lifecycle-toolkit/issues/457)) ([05ac270](https://github.com/keptn/lifecycle-toolkit/commit/05ac27028fdfe420223882c7b7c231dfb1435079))
* **operator:** cancel pending phases when evaluation fails ([#408](https://github.com/keptn/lifecycle-toolkit/issues/408)) ([7f15baf](https://github.com/keptn/lifecycle-toolkit/commit/7f15baf85bdc7f30537ad9ce5a0c582e65ffb16f))
* **operator:** Changed checks on pod owner replicas ([#412](https://github.com/keptn/lifecycle-toolkit/issues/412)) ([46524a7](https://github.com/keptn/lifecycle-toolkit/commit/46524a72e44afb9089c6939a86481a99f8465da0))
* **operator:** detect Job failure and set Task to failed ([#424](https://github.com/keptn/lifecycle-toolkit/issues/424)) ([19114db](https://github.com/keptn/lifecycle-toolkit/commit/19114db17d5ea01687b184962efcb048c28fdc40))
* **operator:** Do not proceed with WLI if no AppVersion containing it is available ([#377](https://github.com/keptn/lifecycle-toolkit/issues/377)) ([cf74540](https://github.com/keptn/lifecycle-toolkit/commit/cf7454004963ac1975a95a5bd3de2ab3783eb487))
* **operator:** Fixed typo in pre and post deployment checks + sorting the PhaseItem interface functions according to topic ([#405](https://github.com/keptn/lifecycle-toolkit/issues/405)) ([ca8f11d](https://github.com/keptn/lifecycle-toolkit/commit/ca8f11da4bd897dad4ecc0a847745f4c8f0749c5))
* **operator:** increment the correct meter to show deployment count ([#434](https://github.com/keptn/lifecycle-toolkit/issues/434)) ([0287596](https://github.com/keptn/lifecycle-toolkit/commit/028759683af54a1c023f909694d899a5d730b750))
* **operator:** revert broken OTel version ([#447](https://github.com/keptn/lifecycle-toolkit/issues/447)) ([3eb47d0](https://github.com/keptn/lifecycle-toolkit/commit/3eb47d0e08d5fbb400cb68e8c4aecfa49a056ad5))
* **operator:** use correct parent/child span relationship ([#418](https://github.com/keptn/lifecycle-toolkit/issues/418)) ([24efc80](https://github.com/keptn/lifecycle-toolkit/commit/24efc80bcf316aa08ff6d8bc0af963e0657872a6))
* use correct namespace variable in delete cmd ([#446](https://github.com/keptn/lifecycle-toolkit/issues/446)) ([c3b2188](https://github.com/keptn/lifecycle-toolkit/commit/c3b2188f214094bd1d7cf86bc3d7db5a12f33159))


### Dependency Updates

* update denoland/deno docker tag to v1.27.2 ([#354](https://github.com/keptn/lifecycle-toolkit/issues/354)) ([3a37846](https://github.com/keptn/lifecycle-toolkit/commit/3a37846f1d0654798acaf626487d3903d824f7fc))
* update denoland/deno docker tag to v1.28.0 ([#401](https://github.com/keptn/lifecycle-toolkit/issues/401)) ([c4502e1](https://github.com/keptn/lifecycle-toolkit/commit/c4502e1a5327653fad8b9ea9f4daf7b22a1cf739))
* update denoland/deno docker tag to v1.28.1 ([#430](https://github.com/keptn/lifecycle-toolkit/issues/430)) ([fdf3f4b](https://github.com/keptn/lifecycle-toolkit/commit/fdf3f4b3a6471b1631db8f6c7e400b0e262ebadd))
* update dependency argoproj/argo-cd to v2.4.17 ([#435](https://github.com/keptn/lifecycle-toolkit/issues/435)) ([9a4976b](https://github.com/keptn/lifecycle-toolkit/commit/9a4976b4d6e291f9b9a34b314a323587c4535104))
* update dependency argoproj/argo-cd to v2.5.2 ([#438](https://github.com/keptn/lifecycle-toolkit/issues/438)) ([2cf98a2](https://github.com/keptn/lifecycle-toolkit/commit/2cf98a2b9c22e8e6cf6752494a5637eddc74595b))
* update dependency cert-manager/cert-manager to v1.10.0 ([#439](https://github.com/keptn/lifecycle-toolkit/issues/439)) ([bb4e487](https://github.com/keptn/lifecycle-toolkit/commit/bb4e487787636c186ea59ce62e37696ac32ba708))
* update dependency cert-manager/cert-manager to v1.10.1 ([#450](https://github.com/keptn/lifecycle-toolkit/issues/450)) ([8872b3c](https://github.com/keptn/lifecycle-toolkit/commit/8872b3ca944c4e8f647de42b4187f5a418d5247d))
* update dependency jaegertracing/jaeger to v1.38.1 ([#437](https://github.com/keptn/lifecycle-toolkit/issues/437)) ([5bd4e4c](https://github.com/keptn/lifecycle-toolkit/commit/5bd4e4c80d0c59a8a1aa5f22b1b72ced82169178))
* update dependency jaegertracing/jaeger to v1.39.0 ([#440](https://github.com/keptn/lifecycle-toolkit/issues/440)) ([3410b63](https://github.com/keptn/lifecycle-toolkit/commit/3410b635d1bf17884e25b7fdbe56f336c189f246))
* update dependency kubernetes-sigs/controller-tools to v0.10.0 ([#443](https://github.com/keptn/lifecycle-toolkit/issues/443)) ([8c60dc7](https://github.com/keptn/lifecycle-toolkit/commit/8c60dc7059b6e4b10625fdd9634906674837a6ba))
* update dependency kubernetes-sigs/kustomize to v4.5.7 ([#444](https://github.com/keptn/lifecycle-toolkit/issues/444)) ([2d83ce6](https://github.com/keptn/lifecycle-toolkit/commit/2d83ce6ddbac8039e107022dc3aa7c6862faa6dd))
* update ghcr.io/keptn/scheduler docker tag to v202211041667586940 ([#310](https://github.com/keptn/lifecycle-toolkit/issues/310)) ([8d71e29](https://github.com/keptn/lifecycle-toolkit/commit/8d71e297e7d86484c74eddbc810c819a0e3a6b4e))
* update go 1.19 ([#364](https://github.com/keptn/lifecycle-toolkit/issues/364)) ([c72c4bc](https://github.com/keptn/lifecycle-toolkit/commit/c72c4bc8855c362d3bf5e4fe73781c4eaa91364f))
* update helm/kind-action action to v1.4.0 ([#355](https://github.com/keptn/lifecycle-toolkit/issues/355)) ([96cde69](https://github.com/keptn/lifecycle-toolkit/commit/96cde694a8ea4fe20e3b5ea93224671fd36118bf))
* update kubernetes packages to v0.25.3 (minor) ([#263](https://github.com/keptn/lifecycle-toolkit/issues/263)) ([d8cec2f](https://github.com/keptn/lifecycle-toolkit/commit/d8cec2f7f19885bf36484a333ce21710d14a0b2e))
* update kubernetes packages to v0.25.3 (patch) ([#291](https://github.com/keptn/lifecycle-toolkit/issues/291)) ([0a648b1](https://github.com/keptn/lifecycle-toolkit/commit/0a648b1b119eecca0842389a63a98908d9764f8b))
* update kubernetes packages to v0.25.4 (patch) ([#383](https://github.com/keptn/lifecycle-toolkit/issues/383)) ([72088d6](https://github.com/keptn/lifecycle-toolkit/commit/72088d6c91b5f6b0b266627191030cd224b21883))
* update module github.com/magiconair/properties to v1.8.6 ([#331](https://github.com/keptn/lifecycle-toolkit/issues/331)) ([f54665e](https://github.com/keptn/lifecycle-toolkit/commit/f54665e2cae31cd487aafb08690e37f7a88f1d7b))
* update module github.com/prometheus/client_golang to v1.13.1 ([#311](https://github.com/keptn/lifecycle-toolkit/issues/311)) ([1fe4242](https://github.com/keptn/lifecycle-toolkit/commit/1fe42421ca3cad939d75fe2a3069f68aa75306f1))
* update module github.com/prometheus/client_golang to v1.14.0 ([#395](https://github.com/keptn/lifecycle-toolkit/issues/395)) ([39af17b](https://github.com/keptn/lifecycle-toolkit/commit/39af17bb9e7d1827edb27dd1d8130a2152332cde))
* update module google.golang.org/grpc to v1.50.1 ([#274](https://github.com/keptn/lifecycle-toolkit/issues/274)) ([44ac9b4](https://github.com/keptn/lifecycle-toolkit/commit/44ac9b4cf020043b5bee4e4d69ed3a9a27565353))
* update module k8s.io/kubernetes to v1.25.4 ([#399](https://github.com/keptn/lifecycle-toolkit/issues/399)) ([5f47086](https://github.com/keptn/lifecycle-toolkit/commit/5f47086da4c38eb77cd8a009ae8cdb93bbc645b2))
* update module sigs.k8s.io/controller-runtime to v0.13.1 ([#279](https://github.com/keptn/lifecycle-toolkit/issues/279)) ([3afcaad](https://github.com/keptn/lifecycle-toolkit/commit/3afcaad7a560162f154f6002eb381d2df7690de7))
* update module sigs.k8s.io/controller-runtime to v0.13.1 ([#306](https://github.com/keptn/lifecycle-toolkit/issues/306)) ([a3a0600](https://github.com/keptn/lifecycle-toolkit/commit/a3a0600f59983d6f6ab000088dfbff54ff88eb67))


### Docs

* add cert-manager installation instructions to README ([#392](https://github.com/keptn/lifecycle-toolkit/issues/392)) ([58161a1](https://github.com/keptn/lifecycle-toolkit/commit/58161a1c6ecfa0b83534e854ab783cbff48c4bd3))
* adding reference to youtube video ([#407](https://github.com/keptn/lifecycle-toolkit/issues/407)) ([6abcade](https://github.com/keptn/lifecycle-toolkit/commit/6abcaded0427e41fe93e61da4291afa0a49f8c6e))
* fix CRD api version for EvaluationProvider and EvaluationDefinition ([#449](https://github.com/keptn/lifecycle-toolkit/issues/449)) ([d4c6716](https://github.com/keptn/lifecycle-toolkit/commit/d4c6716c86e737cc9c6bdd8f81470821ca948098))


### Other

* add component tests as part of the coverage ([#468](https://github.com/keptn/lifecycle-toolkit/issues/468)) ([d521669](https://github.com/keptn/lifecycle-toolkit/commit/d521669abd2dd868a8e31eb9864bf61e018f2e21))
* add CONTRIBUTING.md file ([#466](https://github.com/keptn/lifecycle-toolkit/issues/466)) ([02c2726](https://github.com/keptn/lifecycle-toolkit/commit/02c272667d5296b1feedb9095c4f5dd72e7c7c10))
* ensures that PR subjects start with lowercase ([#427](https://github.com/keptn/lifecycle-toolkit/issues/427)) ([246f0b6](https://github.com/keptn/lifecycle-toolkit/commit/246f0b6b81849f7c9202a6a17623157d9623f540))
* execute performance tests after all other tests have been executed ([#479](https://github.com/keptn/lifecycle-toolkit/issues/479)) ([145a6ab](https://github.com/keptn/lifecycle-toolkit/commit/145a6abeec614a3e75787c93da98b2ee3dca8ed1))
* **operator:** restructure packages ([#469](https://github.com/keptn/lifecycle-toolkit/issues/469)) ([41f21eb](https://github.com/keptn/lifecycle-toolkit/commit/41f21ebbd7839a64cbdb5c4f49061eab9f66976f))
* **scheduler:** make RealAnna codeowner ([#369](https://github.com/keptn/lifecycle-toolkit/issues/369)) ([aba0a70](https://github.com/keptn/lifecycle-toolkit/commit/aba0a708c3aace2d2309fa571e073d90ad6d6861))

## [0.4.0](https://github.com/keptn/lifecycle-toolkit/compare/v0.3.0...v0.4.0) (2022-11-08)


### ⚠ BREAKING CHANGES

* The lifecycle toolkit now uses keptn-lifecycle-toolkit-system namespace by default (#332)
* Rename to lifecycle toolkit (#286)

### Features

* Add Dashboards for Applications and Workloads ([#219](https://github.com/keptn/lifecycle-toolkit/issues/219)) ([48589e2](https://github.com/keptn/lifecycle-toolkit/commit/48589e2a521df0ff7c607a9fb74f47c06f81d3bf))
* Bootstrap webhook/component/integration/performance tests ([#225](https://github.com/keptn/lifecycle-toolkit/issues/225)) ([dbe08c0](https://github.com/keptn/lifecycle-toolkit/commit/dbe08c0a5947a3fbe42aa94660352c3ef6357f14))
* **operator:** Add additional metrics for Deployment duration and interval ([#220](https://github.com/keptn/lifecycle-toolkit/issues/220)) ([71383c0](https://github.com/keptn/lifecycle-toolkit/commit/71383c0680cd17bec96b01155376cff683034d24))
* **operator:** Add information about current phase in workloadversions and appversions ([#200](https://github.com/keptn/lifecycle-toolkit/issues/200)) ([55fa4e9](https://github.com/keptn/lifecycle-toolkit/commit/55fa4e97c62aec7bd1a45f85d47cfaca48f3dd8f))
* **operator:** Add separate trace for Deployment ([#222](https://github.com/keptn/lifecycle-toolkit/issues/222)) ([6966e3d](https://github.com/keptn/lifecycle-toolkit/commit/6966e3d467e058471f15e90159ed749490bc30b2))
* **operator:** Improve state and phase information ([#211](https://github.com/keptn/lifecycle-toolkit/issues/211)) ([6982074](https://github.com/keptn/lifecycle-toolkit/commit/6982074cae4e8147c4643aae821c284614d542b3))
* **operator:** Use Async Gauges for active KLC Entities ([#206](https://github.com/keptn/lifecycle-toolkit/issues/206)) ([9d61ab2](https://github.com/keptn/lifecycle-toolkit/commit/9d61ab2664d5f3339ed5af4e1303eacf2fc89dec))
* Sign released container images with sigstore/cosign ([#290](https://github.com/keptn/lifecycle-toolkit/issues/290)) ([a8f58a4](https://github.com/keptn/lifecycle-toolkit/commit/a8f58a461b082fd13dc86f700ed01d57075276ca))
* The lifecycle toolkit now uses keptn-lifecycle-toolkit-system namespace by default ([#329](https://github.com/keptn/lifecycle-toolkit/issues/329)) ([ef1a158](https://github.com/keptn/lifecycle-toolkit/commit/ef1a15876958ee8614779a9cd5471a2f4aa528b4))
* The lifecycle toolkit now uses keptn-lifecycle-toolkit-system namespace by default ([#332](https://github.com/keptn/lifecycle-toolkit/issues/332)) ([443be11](https://github.com/keptn/lifecycle-toolkit/commit/443be11bb2d8f650a54aad90f4b040313eee24d8))
* Use debug stages in local docker build make commands ([#234](https://github.com/keptn/lifecycle-toolkit/issues/234)) ([6423834](https://github.com/keptn/lifecycle-toolkit/commit/6423834608ce78ca32d33bf54f27dbbc0ae4c116))


### Bug Fixes

* **operator:** Fix nil pointer exception in case of app not being found ([#233](https://github.com/keptn/lifecycle-toolkit/issues/233)) ([de9a016](https://github.com/keptn/lifecycle-toolkit/commit/de9a01654d7b54809932ef973860ede59f541310))
* **operator:** Fixed starting deployments, when no corresponding app-version is available ([#210](https://github.com/keptn/lifecycle-toolkit/issues/210)) ([3efa13e](https://github.com/keptn/lifecycle-toolkit/commit/3efa13e72b900a11a7dd4f65e0fbaae02211a6e9))
* **operator:** Use correct Span Names ([#327](https://github.com/keptn/lifecycle-toolkit/issues/327)) ([e6a0ea0](https://github.com/keptn/lifecycle-toolkit/commit/e6a0ea038783e1d02a569b3b74d0265de99bea9c))
* **operator:** Use pointer receiver for SpanHandler methods to ensure span map is populated; thread safety via mutex ([#288](https://github.com/keptn/lifecycle-toolkit/issues/288)) ([a127a42](https://github.com/keptn/lifecycle-toolkit/commit/a127a42717068a43c60b4cc30abd56bc1478669c))
* **scheduler:** Fix the status the scheduler is acting on (preDeploymentEvaluationStatus) ([#226](https://github.com/keptn/lifecycle-toolkit/issues/226)) ([1a0dd92](https://github.com/keptn/lifecycle-toolkit/commit/1a0dd929930eb078070fb84b9bab0133ef4bccd9))
* **scheduler:** The client should inherit framework configs ([#309](https://github.com/keptn/lifecycle-toolkit/issues/309)) ([847a460](https://github.com/keptn/lifecycle-toolkit/commit/847a460f7759447213a3e405d743da762e9ed29e))
* Typo in observability example ([#248](https://github.com/keptn/lifecycle-toolkit/issues/248)) ([2f6be5f](https://github.com/keptn/lifecycle-toolkit/commit/2f6be5fe091951231dde005b3b9c99dcf07cab87))


### Docs

* Add KubeCon NA 22 Demo ([#308](https://github.com/keptn/lifecycle-toolkit/issues/308)) ([f0ba5db](https://github.com/keptn/lifecycle-toolkit/commit/f0ba5db31d30e64474bd33d10dd1cdd4878a2dd9))
* Add temporary sub-project logo to the repository ([#207](https://github.com/keptn/lifecycle-toolkit/issues/207)) ([3708cb3](https://github.com/keptn/lifecycle-toolkit/commit/3708cb31dca6d8fb179bf8e46aa422ced3b877ff))
* Fix name of keptnappversions ([#215](https://github.com/keptn/lifecycle-toolkit/issues/215)) ([d6e3e2c](https://github.com/keptn/lifecycle-toolkit/commit/d6e3e2c2859ee1882902c570b7564a999f479f47))
* Update the repository links in README after the org migration ([#208](https://github.com/keptn/lifecycle-toolkit/issues/208)) ([a1ac506](https://github.com/keptn/lifecycle-toolkit/commit/a1ac5060d909e9fbe0d7874aaee20af06805f033))


### Dependency Updates

* update actions/checkout action to v3 ([#282](https://github.com/keptn/lifecycle-toolkit/issues/282)) ([99eae9c](https://github.com/keptn/lifecycle-toolkit/commit/99eae9ce94ebc34ce876bbb5c1d19954f83e36d1))
* update denoland/deno docker tag to v1.27.1 ([#307](https://github.com/keptn/lifecycle-toolkit/issues/307)) ([9061fc5](https://github.com/keptn/lifecycle-toolkit/commit/9061fc5dc366d11c23d6f0122a6fb2cd60b7a35b))
* update golang docker tag to v1.18.8 ([#275](https://github.com/keptn/lifecycle-toolkit/issues/275)) ([c510824](https://github.com/keptn/lifecycle-toolkit/commit/c51082481338edc7405d42baaf15139cb35b51b9))


### Other

* **deps:** Update dependencies ([#265](https://github.com/keptn/lifecycle-toolkit/issues/265)) ([7a87bb8](https://github.com/keptn/lifecycle-toolkit/commit/7a87bb87b697b3052cc4e4cdded3f22cff641ccb))
* Introduce failing observability example for podtatohead ([#204](https://github.com/keptn/lifecycle-toolkit/issues/204)) ([f29910d](https://github.com/keptn/lifecycle-toolkit/commit/f29910d2feb8931cb990794899bea275d47ab7b2))
* **operator:** Add workload and app version to KeptnTask ([#201](https://github.com/keptn/lifecycle-toolkit/issues/201)) ([fde0c67](https://github.com/keptn/lifecycle-toolkit/commit/fde0c67a4dd0b01006d1e0f6b0a240307c07bca4))
* **operator:** Support Progressing state in every phase + refactoring + speed improvements ([#236](https://github.com/keptn/lifecycle-toolkit/issues/236)) ([af1da5d](https://github.com/keptn/lifecycle-toolkit/commit/af1da5d938ce46a3dd6970a467842b01db09c33d))
* Rename to lifecycle toolkit ([#286](https://github.com/keptn/lifecycle-toolkit/issues/286)) ([9177c76](https://github.com/keptn/lifecycle-toolkit/commit/9177c76535b1b9dad9dc64c2d34e5e92819fcd2c))
* Update dependencies and fixes ([#281](https://github.com/keptn/lifecycle-toolkit/issues/281)) ([5f5eda9](https://github.com/keptn/lifecycle-toolkit/commit/5f5eda9c599f421db0c7c94f9f5432945fabea3c))
* update grafana dashboards ([#325](https://github.com/keptn/lifecycle-toolkit/issues/325)) ([0d0f2ab](https://github.com/keptn/lifecycle-toolkit/commit/0d0f2abcd5d3e04383e396209d6495e019eaf6a4))
* Update repo URL everywhere ([#216](https://github.com/keptn/lifecycle-toolkit/issues/216)) ([33d494c](https://github.com/keptn/lifecycle-toolkit/commit/33d494c537ea055d61e6a32d63c7812e0af90575))

## [0.3.0](https://github.com/keptn/lifecycle-toolkit/compare/v0.2.0...v0.3.0) (2022-10-20)


### ⚠ BREAKING CHANGES

* **operator:** Modified behavior of KeptnAppVersion and KeptnWorkloadVersion to support pre and post deployment evaluation checks with Prometheus montoring
* **operator:** now the namespaces have to be annotated/labeled with keptn.sh/lifecycle-toolkit=enabled when the lifecycle controller should be used
* **operator:** Implementation of the KeptnApp CRD and Controller. This modifies the behaviour of the KeptnWorkloadVersion and Keptn MutatingWebhook

### Features

* Namespace keptn-lifecycle-toolkit-system should never call webhook ([#192](https://github.com/keptn/lifecycle-toolkit/issues/192)) ([913a9ff](https://github.com/keptn/lifecycle-toolkit/commit/913a9ffd62f93aa7831b35e29853afff6213a0c9))
* **operator:** add fallback behavior when no keptn annotations are set ([#171](https://github.com/keptn/lifecycle-toolkit/issues/171)) ([b6cc674](https://github.com/keptn/lifecycle-toolkit/commit/b6cc674adb787615fc79dbbc5b10668c367e4736))
* **operator:** Add KeptnApplication controller ([#137](https://github.com/keptn/lifecycle-toolkit/issues/137)) ([271f5a8](https://github.com/keptn/lifecycle-toolkit/commit/271f5a830f216c9f827457d8a391c25d56aed2e3))
* **operator:** Added minimal context information ([#170](https://github.com/keptn/lifecycle-toolkit/issues/170)) ([eebe420](https://github.com/keptn/lifecycle-toolkit/commit/eebe4200aac74a7c2cbc73720d1d9ac6a0c1fc72))
* **operator:** Allow pre- and post-deployment tasks as labels or annotations ([#181](https://github.com/keptn/lifecycle-toolkit/issues/181)) ([4241fe7](https://github.com/keptn/lifecycle-toolkit/commit/4241fe7cfab91aa6d38309eacf5712436a6e8327))
* **operator:** Bootstrap evaluation CRD from app ([#184](https://github.com/keptn/lifecycle-toolkit/issues/184)) ([74c3dbc](https://github.com/keptn/lifecycle-toolkit/commit/74c3dbc7b6d78d8ca7eafbac50abb8c3473701eb))
* **operator:** Bootstrap evaluation CRD from WorkloadVersion ([#188](https://github.com/keptn/lifecycle-toolkit/issues/188)) ([95e206b](https://github.com/keptn/lifecycle-toolkit/commit/95e206b4165b0277f5acbc67fc78a8e28f06741b))
* **operator:** Bootstrap KeptnEvaluationProvider and KeptnEvaluation Definition CRDs ([#165](https://github.com/keptn/lifecycle-toolkit/issues/165)) ([03d2346](https://github.com/keptn/lifecycle-toolkit/commit/03d234610fd8ef9f21e756450c7f503cb236f302))
* **operator:** Fix phase naming ([#197](https://github.com/keptn/lifecycle-toolkit/issues/197)) ([3739127](https://github.com/keptn/lifecycle-toolkit/commit/3739127d2794d75c489a6af04acf57b82920ca46))
* **operator:** Introduce KeptnEvaluation Controller + CRD ([#168](https://github.com/keptn/lifecycle-toolkit/issues/168)) ([1ce044a](https://github.com/keptn/lifecycle-toolkit/commit/1ce044a3470f815597d725d424a5491f828f2c4c))
* **operator:** Introduce Prometheus evaluation ([#183](https://github.com/keptn/lifecycle-toolkit/issues/183)) ([c2ab773](https://github.com/keptn/lifecycle-toolkit/commit/c2ab7733291928eaea5c38287c63e45d12754ba1))
* **operator:** namespace should be annotated when the lifecycle controller is used ([#178](https://github.com/keptn/lifecycle-toolkit/issues/178)) ([fa8b875](https://github.com/keptn/lifecycle-toolkit/commit/fa8b8758ebb5a29064f255a66d9066a863bf0944))


### Docs

* Add documentation for OTel collector as pre-requisite ([#185](https://github.com/keptn/lifecycle-toolkit/issues/185)) ([bc3900c](https://github.com/keptn/lifecycle-toolkit/commit/bc3900ca64f6c7a0ef22ab94a9665aac17a83372))
* Add example for ArgoCD ([#179](https://github.com/keptn/lifecycle-toolkit/issues/179)) ([daf622d](https://github.com/keptn/lifecycle-toolkit/commit/daf622d47068f70539eb5819bc81dfe72e1b105c))
* Add flux example ([#187](https://github.com/keptn/lifecycle-toolkit/issues/187)) ([02cceb3](https://github.com/keptn/lifecycle-toolkit/commit/02cceb37d64c52a12d0779f015cf488b4ad3729f))
* Improve installation steps ([#154](https://github.com/keptn/lifecycle-toolkit/issues/154)) ([d183e4f](https://github.com/keptn/lifecycle-toolkit/commit/d183e4f6b3102e426b9e29d0648cdf0c4c7cc19e))


### Other

* Add Evaluation instructions ([#190](https://github.com/keptn/lifecycle-toolkit/issues/190)) ([6717b89](https://github.com/keptn/lifecycle-toolkit/commit/6717b8931496be4235c3945390be53633ccb9e43))
* Add example Grafana dashboard to observability example ([#199](https://github.com/keptn/lifecycle-toolkit/issues/199)) ([9c20600](https://github.com/keptn/lifecycle-toolkit/commit/9c20600f8a5dd3149f040cf2253cd4b787cc08d3))
* Updated Prometheus Network policy for granting access from lifecycle controller namespace ([#191](https://github.com/keptn/lifecycle-toolkit/issues/191)) ([bd77527](https://github.com/keptn/lifecycle-toolkit/commit/bd775276ad1324278c4bc3c82a9c0352d02bcece))

## [0.2.0](https://github.com/keptn/lifecycle-toolkit/compare/v0.1.0...v0.2.0) (2022-10-12)


### Features

* Added tutorial for setting up observability example ([#145](https://github.com/keptn/lifecycle-toolkit/issues/145)) ([28f5a9c](https://github.com/keptn/lifecycle-toolkit/commit/28f5a9c24d031694e2066318bc85ae6e79dfd095))
* **main:** Make LFC development environment installable with one command ([#138](https://github.com/keptn/lifecycle-toolkit/issues/138)) ([832ca37](https://github.com/keptn/lifecycle-toolkit/commit/832ca37d5a19297a63e17a8d367c126af37275c4))
* **operator:** Add commit hash, buildtime, buildversion to OTel resource attributes ([#121](https://github.com/keptn/lifecycle-toolkit/issues/121)) ([5a2ef61](https://github.com/keptn/lifecycle-toolkit/commit/5a2ef61b965472cfe850672d04b4361f5d48ca0d))
* **operator:** Add Spans for handling webhook requests and inject TraceContext ([#115](https://github.com/keptn/lifecycle-toolkit/issues/115)) ([812f2c5](https://github.com/keptn/lifecycle-toolkit/commit/812f2c5d49314617cb9c7532262e15edecd9f078))
* **operator:** Add support for OTel collector ([#139](https://github.com/keptn/lifecycle-toolkit/issues/139)) ([ac3f0d2](https://github.com/keptn/lifecycle-toolkit/commit/ac3f0d222f43abff7f35f1eb8de5ec80ff7dd8dc))
* **operator:** Added metrics ([#55](https://github.com/keptn/lifecycle-toolkit/issues/55)) ([f8a3cee](https://github.com/keptn/lifecycle-toolkit/commit/f8a3ceea6d1628750e7c3a7c9cd3372642bd0611))
* **operator:** Introduce OTel tracing for Task controller ([#128](https://github.com/keptn/lifecycle-toolkit/issues/128)) ([0baf7a9](https://github.com/keptn/lifecycle-toolkit/commit/0baf7a9d8058877247bc264eb6fdb645b0a77a60))
* **operator:** Introduce OTel tracing for Workload controller ([#125](https://github.com/keptn/lifecycle-toolkit/issues/125)) ([bc03709](https://github.com/keptn/lifecycle-toolkit/commit/bc03709b744d61ad966b5fba9f70dbeaffa10119))
* **operator:** Introduce OTel tracing for WorkloadVersion controller ([#131](https://github.com/keptn/lifecycle-toolkit/issues/131)) ([a195614](https://github.com/keptn/lifecycle-toolkit/commit/a1956141fe80e5b1afd79fb33198313e1dbff7fa))
* **scheduler:** Add OTel Resource Attributes ([#147](https://github.com/keptn/lifecycle-toolkit/issues/147)) ([b952156](https://github.com/keptn/lifecycle-toolkit/commit/b9521568e95e7855ee4fef5d55559376e2d398d9))
* **scheduler:** Add support for OTel collector ([#146](https://github.com/keptn/lifecycle-toolkit/issues/146)) ([9fd210d](https://github.com/keptn/lifecycle-toolkit/commit/9fd210d0355e5d17316f5daa8a8e289a03755d46))
* **scheduler:** Add tracing support ([#129](https://github.com/keptn/lifecycle-toolkit/issues/129)) ([60651d1](https://github.com/keptn/lifecycle-toolkit/commit/60651d15c78f9e0aa786d4dd4836c9ae828b14f3))
* **scheduler:** Background check for pod status in permit plugin ([#124](https://github.com/keptn/lifecycle-toolkit/issues/124)) ([97ceef6](https://github.com/keptn/lifecycle-toolkit/commit/97ceef6938603e315c4e1c8d2bb697aabc3dd7f8))
* **scheduler:** Disable gRPC logs when creating OTLP exporter ([#151](https://github.com/keptn/lifecycle-toolkit/issues/151)) ([d0f69b9](https://github.com/keptn/lifecycle-toolkit/commit/d0f69b9509543a5a11f22e8940a71018509ba048))


### Bug Fixes

* **scheduler:** Create new context when starting background routine for pod checks ([#148](https://github.com/keptn/lifecycle-toolkit/issues/148)) ([543ca87](https://github.com/keptn/lifecycle-toolkit/commit/543ca876b27d90cb906ddb2643112a62dc923f56))
* **scheduler:** Ignoring OTel error logs ([#150](https://github.com/keptn/lifecycle-toolkit/issues/150)) ([0be89a5](https://github.com/keptn/lifecycle-toolkit/commit/0be89a56445a0356275f040dedad8fc8716a0fdd))


### Docs

* Add proper version badge in readme ([#114](https://github.com/keptn/lifecycle-toolkit/issues/114)) ([e4add2d](https://github.com/keptn/lifecycle-toolkit/commit/e4add2de2340f160fe30bd0cd6831107339b175e))
* Improve podtato example with HTTP service lookup ([#113](https://github.com/keptn/lifecycle-toolkit/issues/113)) ([81b1236](https://github.com/keptn/lifecycle-toolkit/commit/81b1236dcff7bd37afd0e39f11638fe01406c7c4))
* Update manifest name in readme ([#111](https://github.com/keptn/lifecycle-toolkit/issues/111)) ([e51dbbc](https://github.com/keptn/lifecycle-toolkit/commit/e51dbbc0198f734fb3905b280bc1ff2e0b24d39e))


### Other

* Updated scheduler readme and developer instructions ([#123](https://github.com/keptn/lifecycle-toolkit/issues/123)) ([9bd5d14](https://github.com/keptn/lifecycle-toolkit/commit/9bd5d1461cdeeca851b6ccb78ee7e6ff0b500c1c))


### Build

* Prepare release ([#149](https://github.com/keptn/lifecycle-toolkit/issues/149)) ([5be4504](https://github.com/keptn/lifecycle-toolkit/commit/5be4504e365b1c89ffc3069871a3f0fc0ecc7482))

## 0.1.0 (2022-10-04)


### Features

* Add scheduler with annotations ([#31](https://github.com/keptn/lifecycle-toolkit/issues/31)) ([9e29019](https://github.com/keptn/lifecycle-toolkit/commit/9e29019c098fd4f1d5e36500bd2c7ef410421aa8))
* Bootstrap Service CR and controller ([#21](https://github.com/keptn/lifecycle-toolkit/issues/21)) ([c714ecc](https://github.com/keptn/lifecycle-toolkit/commit/c714eccc3b9c4d1309036fc9d193da3154b4cac5))
* First draft of a scheduler ([#19](https://github.com/keptn/lifecycle-toolkit/issues/19)) ([1884c86](https://github.com/keptn/lifecycle-toolkit/commit/1884c8678a681ed322a0ef2ea07fad3e24e01237))
* first podtatohead sample deployment manifests ([#45](https://github.com/keptn/lifecycle-toolkit/issues/45)) ([3e92d27](https://github.com/keptn/lifecycle-toolkit/commit/3e92d277ebf1a9063ebcf80f05ebe62958e45cbb))
* First Version of Function Execution ([#35](https://github.com/keptn/lifecycle-toolkit/issues/35)) ([f6badfd](https://github.com/keptn/lifecycle-toolkit/commit/f6badfd19f9f0b15c04364be7b03f524c920a015))
* initial version of function runtime ([#26](https://github.com/keptn/lifecycle-toolkit/issues/26)) ([c8800ee](https://github.com/keptn/lifecycle-toolkit/commit/c8800ee352b5d0d5eccd7338cd4fa6a3ae7d2efa))
* Inject keptn-scheduler when resource contains Keptn annotations ([#18](https://github.com/keptn/lifecycle-toolkit/issues/18)) ([4530e86](https://github.com/keptn/lifecycle-toolkit/commit/4530e8602beb4fc923b767eb586e44752f725400))
* **lfc-scheduler:** Move from Helm to Kustomize ([#53](https://github.com/keptn/lifecycle-toolkit/issues/53)) ([d7ba5f3](https://github.com/keptn/lifecycle-toolkit/commit/d7ba5f35f1b32451f833d9fd53079b4162837bde))
* sample function for deno runtime ([#27](https://github.com/keptn/lifecycle-toolkit/issues/27)) ([2501e46](https://github.com/keptn/lifecycle-toolkit/commit/2501e46a18dfc4ab436669fa7c42c570abad5a52))
* substitute event task ([#43](https://github.com/keptn/lifecycle-toolkit/issues/43)) ([3644a7d](https://github.com/keptn/lifecycle-toolkit/commit/3644a7d9a0d4a565a9d857348a63ed91d8cb8102))
* Switch to distroless-base image ([#46](https://github.com/keptn/lifecycle-toolkit/issues/46)) ([0a735b2](https://github.com/keptn/lifecycle-toolkit/commit/0a735b2ca22a02ca42faf7d091741d39e0f5a547))
* Webhook creates Service, Service creates ServiceRun, ServiceRun creates Event ([#30](https://github.com/keptn/lifecycle-toolkit/issues/30)) ([5ae58c3](https://github.com/keptn/lifecycle-toolkit/commit/5ae58c33abe965e79bb405e74c0f308f1220d4ee))


### Bug Fixes

* Added namespace to task definition for podtato head example ([#72](https://github.com/keptn/lifecycle-toolkit/issues/72)) ([7081f27](https://github.com/keptn/lifecycle-toolkit/commit/7081f2772aee5abec840a58c7ab700603e84cf52))
* Fix CODEOWNERS syntax ([0be5197](https://github.com/keptn/lifecycle-toolkit/commit/0be5197c19ea3066d28fe8e97f274efff21f66ff))
* fixed namespace in scheduler kustomization ([#63](https://github.com/keptn/lifecycle-toolkit/issues/63)) ([237bf4f](https://github.com/keptn/lifecycle-toolkit/commit/237bf4f480161f48aa0c4b5f2afbff433447d2a8))
* Missed error ([#76](https://github.com/keptn/lifecycle-toolkit/issues/76)) ([a59aa15](https://github.com/keptn/lifecycle-toolkit/commit/a59aa1552795bce15e39195af235fd42d1448e61))
* **operator:** Get desired amount of replicas from upper level resource ([#89](https://github.com/keptn/lifecycle-toolkit/issues/89)) ([6767832](https://github.com/keptn/lifecycle-toolkit/commit/67678327c2531c25ea0cdb6f1b805365ae454719))
* **operator:** Update workload if spec changes ([#90](https://github.com/keptn/lifecycle-toolkit/issues/90)) ([ec01ad2](https://github.com/keptn/lifecycle-toolkit/commit/ec01ad2ccd04f0c4e6f9ba47e01c5bada128aa3b))
* **operator:** Update workload instance controller, add example ([#102](https://github.com/keptn/lifecycle-toolkit/issues/102)) ([e679c10](https://github.com/keptn/lifecycle-toolkit/commit/e679c1070f0130bd2d6616bf1856956e64dc0bac))
* query jobs before creating ([#79](https://github.com/keptn/lifecycle-toolkit/issues/79)) ([47f82b8](https://github.com/keptn/lifecycle-toolkit/commit/47f82b891d9d20ade2928faae307009e5c96ae22))
* scheduler config plugin configuration ([#68](https://github.com/keptn/lifecycle-toolkit/issues/68)) ([4c4e3c6](https://github.com/keptn/lifecycle-toolkit/commit/4c4e3c60a0e11267dc69ea7d8470555e3ee4f91e))


### Miscellaneous Chores

* release 0.1.0 ([4c46a42](https://github.com/keptn/lifecycle-toolkit/commit/4c46a4297c540b9da30c5a373624d4b8e8a88231))
* release 0.1.0 ([afa8493](https://github.com/keptn/lifecycle-toolkit/commit/afa849324fa422352ed61faa7f0dc75d74c3c25d))


### Continuous Integration

* Prepare release ([#110](https://github.com/keptn/lifecycle-toolkit/issues/110)) ([9d7644b](https://github.com/keptn/lifecycle-toolkit/commit/9d7644b718e29bd37da398d89dc8b51997667358))

## Changelog
