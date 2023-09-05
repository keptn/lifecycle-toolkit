# Changelog

## 1.0.0 (2023-09-05)


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
