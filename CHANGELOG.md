# Changelog

## [0.8.2](https://github.com/keptn/lifecycle-toolkit/compare/klt-v0.8.1...klt-v0.8.2) (2023-09-06)


### Features

* add `aggregation` field in `KeptnMetric` ([#1780](https://github.com/keptn/lifecycle-toolkit/issues/1780)) ([c0b66ea](https://github.com/keptn/lifecycle-toolkit/commit/c0b66eae296e0502608dd66c5fe7eb8f890497e6))
* add `step` field in `KeptnMetric` ([#1755](https://github.com/keptn/lifecycle-toolkit/issues/1755)) ([03ca7dd](https://github.com/keptn/lifecycle-toolkit/commit/03ca7ddde4ce787d0bfddaba2bb3f7b422ff5d6a))
* add cloud events support ([#1843](https://github.com/keptn/lifecycle-toolkit/issues/1843)) ([5b47120](https://github.com/keptn/lifecycle-toolkit/commit/5b471203e412a919903876212ac45c04f180e482))
* add grafana labels to work with kube-prometheus-stack ([#1757](https://github.com/keptn/lifecycle-toolkit/issues/1757)) ([3b7d5ed](https://github.com/keptn/lifecycle-toolkit/commit/3b7d5ed9bd4f09ff49a84e34e3b708edbbff12d8))
* add monitor action to all KLT workflows ([#1923](https://github.com/keptn/lifecycle-toolkit/issues/1923)) ([ee0a0f3](https://github.com/keptn/lifecycle-toolkit/commit/ee0a0f36a42178353202fe0ca407f385cd82b5b5))
* **lifecycle-operator:** clean up KeptnTask API by removing duplicated attributes ([#1965](https://github.com/keptn/lifecycle-toolkit/issues/1965)) ([257b220](https://github.com/keptn/lifecycle-toolkit/commit/257b220a6171ccc82d1b471002b6cf773ec9bd09))
* metrics-operator monorepo setup ([#1791](https://github.com/keptn/lifecycle-toolkit/issues/1791)) ([51445eb](https://github.com/keptn/lifecycle-toolkit/commit/51445ebd24b0914d34b0339ab05ec939440aa4a3))
* **metrics-operator:** add analysis controller ([#1875](https://github.com/keptn/lifecycle-toolkit/issues/1875)) ([017e08b](https://github.com/keptn/lifecycle-toolkit/commit/017e08b0a65679ca417e363f2223b7f4fef3bc55))
* **metrics-operator:** add Analysis CRD ([#1839](https://github.com/keptn/lifecycle-toolkit/issues/1839)) ([9521a16](https://github.com/keptn/lifecycle-toolkit/commit/9521a16ce4946d3169993780f2d2a4f3a75d0445))
* **metrics-operator:** add AnalysisDefinition CRD ([#1823](https://github.com/keptn/lifecycle-toolkit/issues/1823)) ([adf4621](https://github.com/keptn/lifecycle-toolkit/commit/adf4621c2e8147bc0e4ee7f1719859007290c978))
* **metrics-operator:** add AnalysisValueTemplate CRD  ([#1822](https://github.com/keptn/lifecycle-toolkit/issues/1822)) ([f25b24d](https://github.com/keptn/lifecycle-toolkit/commit/f25b24dfef07a600c0fbcd4bdb540efe58cff387))
* **metrics-operator:** introduce range operators in AnalysisDefinition ([#1976](https://github.com/keptn/lifecycle-toolkit/issues/1976)) ([7fb8952](https://github.com/keptn/lifecycle-toolkit/commit/7fb8952c514909ce2c0202e01f1cf501de2c8d55))
* monorepo setup for lifecycle-operator, scheduler and runtimes ([#1857](https://github.com/keptn/lifecycle-toolkit/issues/1857)) ([84e243a](https://github.com/keptn/lifecycle-toolkit/commit/84e243a213ffba86eddd51ccc4bf4dbd61140069))
* update stability of Certificate Manager ([#1733](https://github.com/keptn/lifecycle-toolkit/issues/1733)) ([e83d2ae](https://github.com/keptn/lifecycle-toolkit/commit/e83d2ae4a4724d4da8fee63f23a0c063275fac91))


### Bug Fixes

* add missing cert-injection annotations to helm-chart test result ([#1873](https://github.com/keptn/lifecycle-toolkit/issues/1873)) ([56d6598](https://github.com/keptn/lifecycle-toolkit/commit/56d659812c2fdeda1137d9f9f04701cb1775e434))
* admit pod without creating KLT resources if owner of the pod is not supported ([#1752](https://github.com/keptn/lifecycle-toolkit/issues/1752)) ([f47ca50](https://github.com/keptn/lifecycle-toolkit/commit/f47ca50c1e34548466f9368c260c533e293754ad))
* bump KLT version in helm values ([#1697](https://github.com/keptn/lifecycle-toolkit/issues/1697)) ([342d9d1](https://github.com/keptn/lifecycle-toolkit/commit/342d9d14fff2ad7b387d2adff0ef7331d43f48ff))
* fix Go badge ([#1983](https://github.com/keptn/lifecycle-toolkit/issues/1983)) ([c989a6c](https://github.com/keptn/lifecycle-toolkit/commit/c989a6c5a0d0498d9157380018c3c19013a9cd30))
* **operator:** sanitize app name annotation from uppercase to lowercase ([#1793](https://github.com/keptn/lifecycle-toolkit/issues/1793)) ([0986360](https://github.com/keptn/lifecycle-toolkit/commit/0986360ddaed1ee1c56552f46b98964809dd19c1))
* remove klt-cert-manager from version bumps during KLT release ([#1783](https://github.com/keptn/lifecycle-toolkit/issues/1783)) ([a53e8e0](https://github.com/keptn/lifecycle-toolkit/commit/a53e8e0991e46193d7d0e10325ceae98934b7acd))
* take last element in tag as Workload version number ([#1726](https://github.com/keptn/lifecycle-toolkit/issues/1726)) ([dc3ade0](https://github.com/keptn/lifecycle-toolkit/commit/dc3ade0af20e8ac424bfb860bd8d871d76b81119))
* update DOCKER CMD on docs/Makefile ([#1745](https://github.com/keptn/lifecycle-toolkit/issues/1745)) ([a9ac9f6](https://github.com/keptn/lifecycle-toolkit/commit/a9ac9f6cba77fd12154bccc8e8647b6dd1a8fba0))


### Other

* add metrics-operator back to renovate ([#2047](https://github.com/keptn/lifecycle-toolkit/issues/2047)) ([e5a92c1](https://github.com/keptn/lifecycle-toolkit/commit/e5a92c1c8405d478018006f018a13fa81cb729b4))
* add status field docs to all CRDs ([#1807](https://github.com/keptn/lifecycle-toolkit/issues/1807)) ([650ecba](https://github.com/keptn/lifecycle-toolkit/commit/650ecba95624ed3dc2bd61bf1f86578f450223a5))
* cleanup unused env variables in Makefile ([#1913](https://github.com/keptn/lifecycle-toolkit/issues/1913)) ([1ddd089](https://github.com/keptn/lifecycle-toolkit/commit/1ddd089556aa381b7bf85764e0d6dba7789cf99c))
* create pull request template ([#1936](https://github.com/keptn/lifecycle-toolkit/issues/1936)) ([a3f366d](https://github.com/keptn/lifecycle-toolkit/commit/a3f366d4fe4aa2e6d425526c9726dbe7a084e6ac))
* fix minor security issues ([#1728](https://github.com/keptn/lifecycle-toolkit/issues/1728)) ([ea73cd9](https://github.com/keptn/lifecycle-toolkit/commit/ea73cd983102632fb162e1b4c8ae56687b288b25))
* improved example on app yamls ([#1821](https://github.com/keptn/lifecycle-toolkit/issues/1821)) ([584138f](https://github.com/keptn/lifecycle-toolkit/commit/584138fe9a37db2e6f9518906de483b6fde0f504))
* **main:** release lifecycle-operator-and-scheduler libraries ([#1979](https://github.com/keptn/lifecycle-toolkit/issues/1979)) ([12d0f40](https://github.com/keptn/lifecycle-toolkit/commit/12d0f40725e466825c4a0d483fa344e5888b03ae))
* more renaming ([#1830](https://github.com/keptn/lifecycle-toolkit/issues/1830)) ([f2d5bdd](https://github.com/keptn/lifecycle-toolkit/commit/f2d5bdd5700fef1289f67763bf361af7f2bacbd7))
* move from continuous helmify to custom chart ([#1840](https://github.com/keptn/lifecycle-toolkit/issues/1840)) ([b8d6241](https://github.com/keptn/lifecycle-toolkit/commit/b8d6241687ac9a02ac09a97c10312dde957634ae))
* **operator:** remove dependency on metrics-operator ([#1715](https://github.com/keptn/lifecycle-toolkit/issues/1715)) ([8e2aa3b](https://github.com/keptn/lifecycle-toolkit/commit/8e2aa3b37cb074d32623289702bbda81119b5784))
* **operator:** standardize k8s Events on lifecycle path ([#1692](https://github.com/keptn/lifecycle-toolkit/issues/1692)) ([92730ad](https://github.com/keptn/lifecycle-toolkit/commit/92730ad5bdcca5328b1f5c04636ae057c1d923e5))
* **operator:** unexport EventSender in BuilderOptions ([#1698](https://github.com/keptn/lifecycle-toolkit/issues/1698)) ([c7e7335](https://github.com/keptn/lifecycle-toolkit/commit/c7e7335a680d3ad7ce76dfdce59e6b64a5d9f41a))
* promote Release Lifecycle to beta ([#1833](https://github.com/keptn/lifecycle-toolkit/issues/1833)) ([ee90157](https://github.com/keptn/lifecycle-toolkit/commit/ee90157b1d9180526a3483482f1f4d5275178fc8))
* release cert-manager 1.0.0 ([#1619](https://github.com/keptn/lifecycle-toolkit/issues/1619)) ([5a11d9a](https://github.com/keptn/lifecycle-toolkit/commit/5a11d9a1165f894e7d7e3efd0f29d90d03f7eb36))
* release cert-manager 1.1.0 ([#1972](https://github.com/keptn/lifecycle-toolkit/issues/1972)) ([bb133cf](https://github.com/keptn/lifecycle-toolkit/commit/bb133cfd2ac3207e8a4006eb7a9390dc58737465))
* release cert-manager 1.1.0 ([#1993](https://github.com/keptn/lifecycle-toolkit/issues/1993)) ([a8c22f7](https://github.com/keptn/lifecycle-toolkit/commit/a8c22f779eafd68ea12c97c808ad2041fc89acbf))
* release cert-manager 1.1.0 ([#1998](https://github.com/keptn/lifecycle-toolkit/issues/1998)) ([5fbee38](https://github.com/keptn/lifecycle-toolkit/commit/5fbee380634244e32ac0bb90b0cf4b74ee72b661))
* release deno-runtime 1.0.0 ([#1975](https://github.com/keptn/lifecycle-toolkit/issues/1975)) ([8df9ca4](https://github.com/keptn/lifecycle-toolkit/commit/8df9ca4840201106361b8ab8d2e7765d946b5ed2))
* release deno-runtime 1.0.1 ([#1990](https://github.com/keptn/lifecycle-toolkit/issues/1990)) ([4e088c5](https://github.com/keptn/lifecycle-toolkit/commit/4e088c535645815dcb3fb58f1c09dc67b97e0c02))
* release lifecycle-operator 0.8.2 ([#2033](https://github.com/keptn/lifecycle-toolkit/issues/2033)) ([17ef13a](https://github.com/keptn/lifecycle-toolkit/commit/17ef13a2fe2782c7499d923d1a91fd05f6503047))
* release metrics-operator 0.8.2 ([#2030](https://github.com/keptn/lifecycle-toolkit/issues/2030)) ([c523cb0](https://github.com/keptn/lifecycle-toolkit/commit/c523cb017f3ff9d6adb65b9a0e7d9e63712e6b3e))
* release python-runtime 1.0.0 ([#1969](https://github.com/keptn/lifecycle-toolkit/issues/1969)) ([9a995c4](https://github.com/keptn/lifecycle-toolkit/commit/9a995c447e65a4a96d4d3dca53f40a0c1c383b70))
* release scheduler 0.8.2 ([#2032](https://github.com/keptn/lifecycle-toolkit/issues/2032)) ([cb4d2b1](https://github.com/keptn/lifecycle-toolkit/commit/cb4d2b14a7a772572b505fa844db6f08a45db291))
* release scheduler 0.8.2 ([#2043](https://github.com/keptn/lifecycle-toolkit/issues/2043)) ([621c59d](https://github.com/keptn/lifecycle-toolkit/commit/621c59d26c22af492ea3fbc947071c6a07c0ffbd))
* remove cert manager from renovate ignores ([#1962](https://github.com/keptn/lifecycle-toolkit/issues/1962)) ([972b3bb](https://github.com/keptn/lifecycle-toolkit/commit/972b3bbcec735c3361fb8e236b1e2c61edaebbb4))
* remove generation of KLT manifest releases ([#1850](https://github.com/keptn/lifecycle-toolkit/issues/1850)) ([36b2dff](https://github.com/keptn/lifecycle-toolkit/commit/36b2dff6c8ec26aa4552542e236b3d790d65da57))
* remove helm chart generation from CI ([#1856](https://github.com/keptn/lifecycle-toolkit/issues/1856)) ([768b3e9](https://github.com/keptn/lifecycle-toolkit/commit/768b3e9d6d3a2254be13402615ee1ff06c734214))
* remove monitor action ([#1989](https://github.com/keptn/lifecycle-toolkit/issues/1989)) ([b0b37ea](https://github.com/keptn/lifecycle-toolkit/commit/b0b37ea56c5b36298834cb8fc1b6fbc8ac345906))
* rename operator folder to lifecycle-operator ([#1819](https://github.com/keptn/lifecycle-toolkit/issues/1819)) ([97a2d25](https://github.com/keptn/lifecycle-toolkit/commit/97a2d25919c0a02165dd0dc6c7c82d57ad200139))
* rename sonar settings ([#1831](https://github.com/keptn/lifecycle-toolkit/issues/1831)) ([952712f](https://github.com/keptn/lifecycle-toolkit/commit/952712f2c1c70dfa412bcb8adfdbe96fd392ad7d))
* simplify dashboard installation ([#1815](https://github.com/keptn/lifecycle-toolkit/issues/1815)) ([7c4e7bc](https://github.com/keptn/lifecycle-toolkit/commit/7c4e7bc2b5f7679daf6e2a42f0b72eda6a8b2ea3))
* support external cert-manager ([#1864](https://github.com/keptn/lifecycle-toolkit/issues/1864)) ([50dac48](https://github.com/keptn/lifecycle-toolkit/commit/50dac48ddba2181acb28aa714ba7e6605b038df5))


### Docs

* adapt DORA metrics kubectl command ([#1865](https://github.com/keptn/lifecycle-toolkit/issues/1865)) ([05d2f51](https://github.com/keptn/lifecycle-toolkit/commit/05d2f51e30b82be89752de3a29fe2743c9dc9154))
* add info to implementing/tasks about sequential execution ([#1950](https://github.com/keptn/lifecycle-toolkit/issues/1950)) ([61f92c3](https://github.com/keptn/lifecycle-toolkit/commit/61f92c335b26e79cf058db39c9b31fb44088951c))
* added troubleshooting page ([#1860](https://github.com/keptn/lifecycle-toolkit/issues/1860)) ([43f439b](https://github.com/keptn/lifecycle-toolkit/commit/43f439b09878fbe06bbae2714e033389a367f702))
* change releases link to get started in mainpage header menu ([#1738](https://github.com/keptn/lifecycle-toolkit/issues/1738)) ([1f9ea33](https://github.com/keptn/lifecycle-toolkit/commit/1f9ea332ad5c2d8375f8d3300927c5ad31ccbf0a))
* clean up KLT README file ([#1685](https://github.com/keptn/lifecycle-toolkit/issues/1685)) ([5204457](https://github.com/keptn/lifecycle-toolkit/commit/5204457870bbd2ee6e450b244043fe51aea3212a))
* conceptual mapping of Keptn v1 components to KLT components ([#1628](https://github.com/keptn/lifecycle-toolkit/issues/1628)) ([fa2b54d](https://github.com/keptn/lifecycle-toolkit/commit/fa2b54da52c58a3d9a12fd6cbb39a73e0e58e072))
* contrib: DCO ([#1886](https://github.com/keptn/lifecycle-toolkit/issues/1886)) ([2fdd9cb](https://github.com/keptn/lifecycle-toolkit/commit/2fdd9cb7f51c7ae210b8a8d2fc4da93cfafd95ec))
* contrib: set up dev environment ([#1888](https://github.com/keptn/lifecycle-toolkit/issues/1888)) ([2778c21](https://github.com/keptn/lifecycle-toolkit/commit/2778c211d5af5c7e4877445ac3ee768843a24c15))
* contribution guidelines ([#1906](https://github.com/keptn/lifecycle-toolkit/issues/1906)) ([25bf6ce](https://github.com/keptn/lifecycle-toolkit/commit/25bf6ce1a32b15c870c3dadff861999eef8a4896))
* create crd-template.md ([#1885](https://github.com/keptn/lifecycle-toolkit/issues/1885)) ([06defd3](https://github.com/keptn/lifecycle-toolkit/commit/06defd31b6d1c805a3ff67cfd9800a3e888337e8))
* deployment flow architecture ([#1470](https://github.com/keptn/lifecycle-toolkit/issues/1470)) ([6fe5192](https://github.com/keptn/lifecycle-toolkit/commit/6fe5192293175bf6946f74bba0e9057e704ab8c7))
* document `timeframe` feature for `KeptnMetric` ([#1703](https://github.com/keptn/lifecycle-toolkit/issues/1703)) ([077f0d5](https://github.com/keptn/lifecycle-toolkit/commit/077f0d5d0a49bc5b1f0e800274343660b8218c65))
* edits to metrics-operator architecture page ([#1679](https://github.com/keptn/lifecycle-toolkit/issues/1679)) ([7eb8afe](https://github.com/keptn/lifecycle-toolkit/commit/7eb8afecde7b5d7d8f96cb3c1f2370eb60410a7a))
* excercises text changed to exercises ([#1693](https://github.com/keptn/lifecycle-toolkit/issues/1693)) ([df4cda6](https://github.com/keptn/lifecycle-toolkit/commit/df4cda68c97c464983a482faeaabe91e38d05637))
* fix branding homepage ([#2041](https://github.com/keptn/lifecycle-toolkit/issues/2041)) ([e91137e](https://github.com/keptn/lifecycle-toolkit/commit/e91137ee91e08bb59b91d8f1dad1bcbc977e2a39))
* fix broken links ([#1921](https://github.com/keptn/lifecycle-toolkit/issues/1921)) ([44074a2](https://github.com/keptn/lifecycle-toolkit/commit/44074a29391ab157b0279f2728b6232c15a8c13a))
* fix link to deno web page ([#1908](https://github.com/keptn/lifecycle-toolkit/issues/1908)) ([d63182f](https://github.com/keptn/lifecycle-toolkit/commit/d63182f1af83f282676bee97301e80e1f97fdea6))
* fix links in dev builds ([#1722](https://github.com/keptn/lifecycle-toolkit/issues/1722)) ([a35ed45](https://github.com/keptn/lifecycle-toolkit/commit/a35ed45fabd2314922fc9631ef585bdb10c3c295))
* fix typo ([#1706](https://github.com/keptn/lifecycle-toolkit/issues/1706)) ([3690cd3](https://github.com/keptn/lifecycle-toolkit/commit/3690cd3aea124b48873168b2a14d31cd67840a8b))
* fix typo ([#1754](https://github.com/keptn/lifecycle-toolkit/issues/1754)) ([9ebdcec](https://github.com/keptn/lifecycle-toolkit/commit/9ebdcec83ab3c9416c5ae2a597352cfc34271661))
* fix typos and grammar issues ([#1925](https://github.com/keptn/lifecycle-toolkit/issues/1925)) ([5570d55](https://github.com/keptn/lifecycle-toolkit/commit/5570d555bfc4bbdcbfc66b2725d5352090e5b937))
* fixed typo ([#1799](https://github.com/keptn/lifecycle-toolkit/issues/1799)) ([b9393be](https://github.com/keptn/lifecycle-toolkit/commit/b9393be558da31b5d0ee4541fc28942b85752f07))
* fixed typo Troubleshoort to Troubleshoot ([#1776](https://github.com/keptn/lifecycle-toolkit/issues/1776)) ([eb7c9b2](https://github.com/keptn/lifecycle-toolkit/commit/eb7c9b20b7b71ed598a0ac8681335d91c4a621db))
* git contributing guide ([#1892](https://github.com/keptn/lifecycle-toolkit/issues/1892)) ([30393c6](https://github.com/keptn/lifecycle-toolkit/commit/30393c6e6d1127cde12812810afc1bb8f97a1498))
* how to use GitHub Codespaces to contribute to Keptn ([#1977](https://github.com/keptn/lifecycle-toolkit/issues/1977)) ([c96cb72](https://github.com/keptn/lifecycle-toolkit/commit/c96cb727931abdeeb9bcbbc2e551dc5722a9ed88))
* implement KLT -&gt; Keptn name change ([#2001](https://github.com/keptn/lifecycle-toolkit/issues/2001)) ([440c308](https://github.com/keptn/lifecycle-toolkit/commit/440c3082e5400f89d791724651984ba2bc0a4724))
* implement KLT -&gt; Keptn name change for tasks page ([#2016](https://github.com/keptn/lifecycle-toolkit/issues/2016)) ([8516716](https://github.com/keptn/lifecycle-toolkit/commit/8516716c8c7a7a265df3a9d65b2291707c96cd76))
* improve "Intro to Keptn" page ([#2040](https://github.com/keptn/lifecycle-toolkit/issues/2040)) ([af4f417](https://github.com/keptn/lifecycle-toolkit/commit/af4f417ea0455d2dc26c00eafc9445fe1a208d49))
* improve how-to and ref info for apps ([#1868](https://github.com/keptn/lifecycle-toolkit/issues/1868)) ([7139ffd](https://github.com/keptn/lifecycle-toolkit/commit/7139ffd67115210c4fc760987748c01853b94eab))
* keptn Scheduler architecture documentation ([#1777](https://github.com/keptn/lifecycle-toolkit/issues/1777)) ([ce96200](https://github.com/keptn/lifecycle-toolkit/commit/ce96200b9bfed62062b199845104c4493b3a2627))
* manifests ref section, edits to ref section intros ([#1800](https://github.com/keptn/lifecycle-toolkit/issues/1800)) ([604876f](https://github.com/keptn/lifecycle-toolkit/commit/604876f6a7b51b7185db18c5daabe0d23a2b5eaf))
* migrate quality gates ([#1708](https://github.com/keptn/lifecycle-toolkit/issues/1708)) ([0cec244](https://github.com/keptn/lifecycle-toolkit/commit/0cec244f1ed3da8e10312dc6e820b3129d161d9f))
* restructure migration guide ([#1838](https://github.com/keptn/lifecycle-toolkit/issues/1838)) ([8eb05c7](https://github.com/keptn/lifecycle-toolkit/commit/8eb05c7e75b7d2377835b6915df15bfa7a19c756))
* rewrite README file ([#1862](https://github.com/keptn/lifecycle-toolkit/issues/1862)) ([e304969](https://github.com/keptn/lifecycle-toolkit/commit/e3049690da849cce749fab9bb8c1d4e0edd6f473))
* simplify installation page ([#1751](https://github.com/keptn/lifecycle-toolkit/issues/1751)) ([a2144f5](https://github.com/keptn/lifecycle-toolkit/commit/a2144f501f8b8a062f77b711e09bdfbaba3bbf6e))
* small improvements ([#1951](https://github.com/keptn/lifecycle-toolkit/issues/1951)) ([6273709](https://github.com/keptn/lifecycle-toolkit/commit/627370907807d1805dcff7e224a2edafe30c8383))
* tutorial Run Standalone tasks for non-k8s apps ([#1947](https://github.com/keptn/lifecycle-toolkit/issues/1947)) ([a8d6902](https://github.com/keptn/lifecycle-toolkit/commit/a8d6902c7e283bf6a9dec1526b35fa57e1d05fc1))
* update auto generated docs to include `spec.range.step` in `KeptnMetric` ([#1806](https://github.com/keptn/lifecycle-toolkit/issues/1806)) ([8a90145](https://github.com/keptn/lifecycle-toolkit/commit/8a90145bc024b8046423a3847664cb16f74e53a5))
* update docs/content/en/contribute/docs/local-building/index.md ([#1753](https://github.com/keptn/lifecycle-toolkit/issues/1753)) ([14494c5](https://github.com/keptn/lifecycle-toolkit/commit/14494c5a66e2d4e26ba40299b25424d6127516fa))
* update the doc source file structure page ([#1984](https://github.com/keptn/lifecycle-toolkit/issues/1984)) ([124e243](https://github.com/keptn/lifecycle-toolkit/commit/124e2433888e76277fca640ba56f5a19bbebdfbe))


### Dependency Updates

* update actions/setup-node action to v3.7.0 ([#1713](https://github.com/keptn/lifecycle-toolkit/issues/1713)) ([7a610ef](https://github.com/keptn/lifecycle-toolkit/commit/7a610ef50c54ac590a072741e76eeca4c6b27cc4))
* update actions/setup-node action to v3.8.1 ([#1912](https://github.com/keptn/lifecycle-toolkit/issues/1912)) ([642842a](https://github.com/keptn/lifecycle-toolkit/commit/642842abec87475ccbbda035ef992051c42071ef))
* update aquasecurity/trivy-action action to v0.12.0 ([#2010](https://github.com/keptn/lifecycle-toolkit/issues/2010)) ([093670c](https://github.com/keptn/lifecycle-toolkit/commit/093670caded5e7d78de5923eed2b795b71ca0eac))
* update curlimages/curl docker tag to v8.2.1 ([#1792](https://github.com/keptn/lifecycle-toolkit/issues/1792)) ([88a54f9](https://github.com/keptn/lifecycle-toolkit/commit/88a54f97c1573038f8c1e762e2cffe40de513a7e))
* update denoland/deno docker tag to alpine-1.36.1 ([#1924](https://github.com/keptn/lifecycle-toolkit/issues/1924)) ([4031ec0](https://github.com/keptn/lifecycle-toolkit/commit/4031ec01b44777c86f9d623dcfd8195177be01fa))
* update dependency argoproj/argo-cd to v2.7.10 ([#1795](https://github.com/keptn/lifecycle-toolkit/issues/1795)) ([3936cf0](https://github.com/keptn/lifecycle-toolkit/commit/3936cf0581eeb27cb409055b0be49127ff4c7c6d))
* update dependency argoproj/argo-cd to v2.7.11 ([#1877](https://github.com/keptn/lifecycle-toolkit/issues/1877)) ([72fba14](https://github.com/keptn/lifecycle-toolkit/commit/72fba143c16f341bf01899a4e92a21255aaac855))
* update dependency argoproj/argo-cd to v2.7.8 ([#1763](https://github.com/keptn/lifecycle-toolkit/issues/1763)) ([b168ef5](https://github.com/keptn/lifecycle-toolkit/commit/b168ef58a99bf7993ab153a9361e79f31ddeffd0))
* update dependency argoproj/argo-cd to v2.8.0 ([#1927](https://github.com/keptn/lifecycle-toolkit/issues/1927)) ([2a6bc6a](https://github.com/keptn/lifecycle-toolkit/commit/2a6bc6a26de4e61a7945b225cc079960f0d7381c))
* update dependency argoproj/argo-cd to v2.8.2 ([#1956](https://github.com/keptn/lifecycle-toolkit/issues/1956)) ([04456d5](https://github.com/keptn/lifecycle-toolkit/commit/04456d5a75dd855e675bd7375e360612e2a24057))
* update dependency autoprefixer to v10.4.15 ([#1909](https://github.com/keptn/lifecycle-toolkit/issues/1909)) ([8dbec2d](https://github.com/keptn/lifecycle-toolkit/commit/8dbec2d6116fb20bac86162aaea2b75c24eb96be))
* update dependency bitnami-labs/readme-generator-for-helm to v2.5.1 ([#1849](https://github.com/keptn/lifecycle-toolkit/issues/1849)) ([48236c9](https://github.com/keptn/lifecycle-toolkit/commit/48236c954a5e97df03d774415443d1dea30eab88))
* update dependency cloud-bulldozer/kube-burner to v1.7.4 ([#1797](https://github.com/keptn/lifecycle-toolkit/issues/1797)) ([69606e6](https://github.com/keptn/lifecycle-toolkit/commit/69606e60e3478049e3533a2711a0ee1fb4de0d8e))
* update dependency cloud-bulldozer/kube-burner to v1.7.5 ([#1910](https://github.com/keptn/lifecycle-toolkit/issues/1910)) ([29a82be](https://github.com/keptn/lifecycle-toolkit/commit/29a82be86a945825f0b526f149c223ad7652163f))
* update dependency cloud-bulldozer/kube-burner to v1.7.6 ([#2000](https://github.com/keptn/lifecycle-toolkit/issues/2000)) ([ca1fe57](https://github.com/keptn/lifecycle-toolkit/commit/ca1fe57fecd5978d63e900a0cb060f29ba7fbe74))
* update dependency golangci/golangci-lint to v1.54.1 ([#1928](https://github.com/keptn/lifecycle-toolkit/issues/1928)) ([cc36f34](https://github.com/keptn/lifecycle-toolkit/commit/cc36f34096b892611cb42e073cd6f9cc01c5365f))
* update dependency golangci/golangci-lint to v1.54.2 ([#1937](https://github.com/keptn/lifecycle-toolkit/issues/1937)) ([db5bcbf](https://github.com/keptn/lifecycle-toolkit/commit/db5bcbf33477b11dea602e2050bc9c366c654a6f))
* update dependency helm/helm to v3.12.2 ([#1764](https://github.com/keptn/lifecycle-toolkit/issues/1764)) ([8216e6b](https://github.com/keptn/lifecycle-toolkit/commit/8216e6b65aac53f670aec5f383f3edbfbcbd526b))
* update dependency jaegertracing/jaeger to v1.48.0 ([#1542](https://github.com/keptn/lifecycle-toolkit/issues/1542)) ([0596eb5](https://github.com/keptn/lifecycle-toolkit/commit/0596eb51012ff912af9d6e26a5341b0eb707ea06))
* update dependency jaegertracing/jaeger-operator to v1.47.0 ([#1638](https://github.com/keptn/lifecycle-toolkit/issues/1638)) ([6bb371e](https://github.com/keptn/lifecycle-toolkit/commit/6bb371eb15d60f5d1a2c5e6d175e78a2c3650489))
* update dependency jaegertracing/jaeger-operator to v1.48.0 ([#2018](https://github.com/keptn/lifecycle-toolkit/issues/2018)) ([db781c7](https://github.com/keptn/lifecycle-toolkit/commit/db781c737d39b79e8a745e245b327d1d6220d98f))
* update dependency kubernetes-sigs/controller-tools to v0.12.1 ([#1765](https://github.com/keptn/lifecycle-toolkit/issues/1765)) ([ba79a32](https://github.com/keptn/lifecycle-toolkit/commit/ba79a32ef6acc9de8fb5d618b9ede7d6f96ce15e))
* update dependency kubernetes-sigs/controller-tools to v0.13.0 ([#1930](https://github.com/keptn/lifecycle-toolkit/issues/1930)) ([8b34b63](https://github.com/keptn/lifecycle-toolkit/commit/8b34b63404d0339633ef41ff1cf2005deae8d2b7))
* update dependency kubernetes-sigs/kustomize to v5.1.1 ([#1853](https://github.com/keptn/lifecycle-toolkit/issues/1853)) ([354ab3f](https://github.com/keptn/lifecycle-toolkit/commit/354ab3f980c2569e17a0354ece417df40317d120))
* update ghcr.io/keptn/certificate-operator docker tag to v1.1.0 ([#1964](https://github.com/keptn/lifecycle-toolkit/issues/1964)) ([cdc07ae](https://github.com/keptn/lifecycle-toolkit/commit/cdc07ae2ef2e3164d5da74dd777a975e03ff8ee2))
* update ghcr.io/keptn/deno-runtime docker tag to v1 ([#1988](https://github.com/keptn/lifecycle-toolkit/issues/1988)) ([c6958fd](https://github.com/keptn/lifecycle-toolkit/commit/c6958fd15b84813167d231b745d53748d812e69c))
* update ghcr.io/keptn/python-runtime docker tag to v1 ([#1985](https://github.com/keptn/lifecycle-toolkit/issues/1985)) ([121de2e](https://github.com/keptn/lifecycle-toolkit/commit/121de2eb8881064f390b5abaef7e6c5f5031c5ec))
* update github.com/keptn/lifecycle-toolkit/klt-cert-manager digest to 0b618c4 ([#1654](https://github.com/keptn/lifecycle-toolkit/issues/1654)) ([c749313](https://github.com/keptn/lifecycle-toolkit/commit/c749313bfad7bd98b8d0ae7cc6dd2ea56f23e041))
* update github.com/keptn/lifecycle-toolkit/klt-cert-manager digest to cba2de5 ([#1762](https://github.com/keptn/lifecycle-toolkit/issues/1762)) ([b77bcea](https://github.com/keptn/lifecycle-toolkit/commit/b77bceae39d6e4372b879afa326e7658d4ccdd89))
* update helm/kind-action action to v1.8.0 ([#1714](https://github.com/keptn/lifecycle-toolkit/issues/1714)) ([af76757](https://github.com/keptn/lifecycle-toolkit/commit/af7675775ae679ca077cf6281210cf7ec88a768f))
* update keptn components to v0.8.2 ([#2048](https://github.com/keptn/lifecycle-toolkit/issues/2048)) ([49a884b](https://github.com/keptn/lifecycle-toolkit/commit/49a884b70bbe19ea1523d553195db879a71acc4b))
* update keptn/docs-tooling action to v0.1.4 ([#1781](https://github.com/keptn/lifecycle-toolkit/issues/1781)) ([bba98c2](https://github.com/keptn/lifecycle-toolkit/commit/bba98c2f2c0eb49fdf15d53b825e78cb319f96b5))
* update kubernetes packages (patch) ([#1786](https://github.com/keptn/lifecycle-toolkit/issues/1786)) ([cba2de5](https://github.com/keptn/lifecycle-toolkit/commit/cba2de5a5cd04c094131552aaf92c2b85ac23d21))
* update module github.com/imdario/mergo to v1 ([#1664](https://github.com/keptn/lifecycle-toolkit/issues/1664)) ([3c009d0](https://github.com/keptn/lifecycle-toolkit/commit/3c009d07c379e30489072744c2ceef10edd30923))
* update module github.com/onsi/gomega to v1.27.9 ([#1787](https://github.com/keptn/lifecycle-toolkit/issues/1787)) ([90b6ce9](https://github.com/keptn/lifecycle-toolkit/commit/90b6ce92253f52a43f3c13dddaa918ca73b515d0))
* update module golang.org/x/net to v0.12.0 ([#1662](https://github.com/keptn/lifecycle-toolkit/issues/1662)) ([49318bf](https://github.com/keptn/lifecycle-toolkit/commit/49318bfc40497a120304de9d831dfe033259220f))
* update module google.golang.org/grpc to v1.56.2 ([#1663](https://github.com/keptn/lifecycle-toolkit/issues/1663)) ([0b618c4](https://github.com/keptn/lifecycle-toolkit/commit/0b618c4bf15209fbb81ec7c05f1d05543bdfd1cf))
* update otel/opentelemetry-collector docker tag to v0.81.0 ([#1188](https://github.com/keptn/lifecycle-toolkit/issues/1188)) ([cbfdc84](https://github.com/keptn/lifecycle-toolkit/commit/cbfdc8499cce26b4d66e1d2179302e1ebfc14141))
* update otel/opentelemetry-collector docker tag to v0.84.0 ([#1916](https://github.com/keptn/lifecycle-toolkit/issues/1916)) ([7e4bab4](https://github.com/keptn/lifecycle-toolkit/commit/7e4bab4a89ae9fc984083840bd719f12ecd995e0))
* update sigstore/cosign-installer action to v3.1.2 ([#2009](https://github.com/keptn/lifecycle-toolkit/issues/2009)) ([044d3b5](https://github.com/keptn/lifecycle-toolkit/commit/044d3b52338f027c2095d3949e4f21a2848016bf))

## [0.8.1](https://github.com/keptn/lifecycle-toolkit/compare/klt-v0.8.0...klt-v0.8.1) (2023-07-07)


### Features

* add support for timeframe in `KeptnMetric` ([#1471](https://github.com/keptn/lifecycle-toolkit/issues/1471)) ([4d9ceb7](https://github.com/keptn/lifecycle-toolkit/commit/4d9ceb7a630ad8007f2583040f24ec514870006d))
* cert-manager monorepo setup ([#1528](https://github.com/keptn/lifecycle-toolkit/issues/1528)) ([0156f15](https://github.com/keptn/lifecycle-toolkit/commit/0156f15779f0ec756c801b51a84c3e9cdd0edeb7))
* update Prometheus API to query metrics over a range ([#1587](https://github.com/keptn/lifecycle-toolkit/issues/1587)) ([47a3e06](https://github.com/keptn/lifecycle-toolkit/commit/47a3e067896224e8fc82ba67ec69d4094a685a35))


### Bug Fixes

* **examples:** add new task definitions to kustomize base ([#1674](https://github.com/keptn/lifecycle-toolkit/issues/1674)) ([adba1ec](https://github.com/keptn/lifecycle-toolkit/commit/adba1ec386eba555594498c3c17f65150743a10c))
* **helm-chart:** propagate labels for validation webhook ([#1678](https://github.com/keptn/lifecycle-toolkit/issues/1678)) ([5602bd1](https://github.com/keptn/lifecycle-toolkit/commit/5602bd108b9a326466defe43717ab7b1b2c52469))
* **operator:** avoid multiple creations of the same KeptnTask ([#1676](https://github.com/keptn/lifecycle-toolkit/issues/1676)) ([78ba574](https://github.com/keptn/lifecycle-toolkit/commit/78ba57415e7a60696c71e60b5e2c0ba1fe8c89db))
* **operator:** ensure that generated resource names contain no unallowed character ([#1661](https://github.com/keptn/lifecycle-toolkit/issues/1661)) ([59db60f](https://github.com/keptn/lifecycle-toolkit/commit/59db60f67d295728e5df4a1a132ad7f48ff847eb))
* **operator:** make sure there is exactly one job per task execution ([#1672](https://github.com/keptn/lifecycle-toolkit/issues/1672)) ([b68ba87](https://github.com/keptn/lifecycle-toolkit/commit/b68ba874c516a02cbf9bc431f3d58b82ce34b564))
* **operator:** parse flags so they can be configured ([#1649](https://github.com/keptn/lifecycle-toolkit/issues/1649)) ([4243085](https://github.com/keptn/lifecycle-toolkit/commit/424308509b2556377b4683b836d88a7797870dc3))
* **operator:** provide the right app version for single-service applications ([c7d35b8](https://github.com/keptn/lifecycle-toolkit/commit/c7d35b8f1f26e73a7846ae4d2a39dff3b0cc87f7))
* **python-runtime:** install curl to execute scripts referenced via url ([#1681](https://github.com/keptn/lifecycle-toolkit/issues/1681)) ([ac0d515](https://github.com/keptn/lifecycle-toolkit/commit/ac0d51532182bbdb559b72319eed05f2b12fa1d0))


### Dependency Updates

* update anchore/sbom-action action to v0.14.3 ([#1626](https://github.com/keptn/lifecycle-toolkit/issues/1626)) ([2a1026c](https://github.com/keptn/lifecycle-toolkit/commit/2a1026cdc312b36198841bdd40264f4352064e51))
* update busybox docker tag to v1.36.1 ([#1595](https://github.com/keptn/lifecycle-toolkit/issues/1595)) ([6770912](https://github.com/keptn/lifecycle-toolkit/commit/6770912c5c7fb3562db51835edb468cc6866e85c))
* update dependency argoproj/argo-cd to v2.7.6 ([#1596](https://github.com/keptn/lifecycle-toolkit/issues/1596)) ([1c77c81](https://github.com/keptn/lifecycle-toolkit/commit/1c77c815abd06f29f91aafee440890bd36efea28))
* update dependency golangci/golangci-lint to v1.53.3 ([#1606](https://github.com/keptn/lifecycle-toolkit/issues/1606)) ([227800a](https://github.com/keptn/lifecycle-toolkit/commit/227800a882028ee857480c70a567aeffccc5a60e))
* update dependency helm/helm to v3.12.1 ([#1607](https://github.com/keptn/lifecycle-toolkit/issues/1607)) ([ac93ba4](https://github.com/keptn/lifecycle-toolkit/commit/ac93ba4976998b4ebccc964dc553ebb34304ec89))
* update dependency kubernetes-sigs/kustomize to v5.1.0 ([#1655](https://github.com/keptn/lifecycle-toolkit/issues/1655)) ([791e211](https://github.com/keptn/lifecycle-toolkit/commit/791e211d515e7eb74d8d347420500323b1cc7cef))
* update github.com/keptn/lifecycle-toolkit/klt-cert-manager digest to 1c77c81 ([#1593](https://github.com/keptn/lifecycle-toolkit/issues/1593)) ([472eac0](https://github.com/keptn/lifecycle-toolkit/commit/472eac0c8e2bd26eecf20e047123689761c914e7))
* update github.com/keptn/lifecycle-toolkit/klt-cert-manager digest to 4ad9bbf ([#1631](https://github.com/keptn/lifecycle-toolkit/issues/1631)) ([9060ae1](https://github.com/keptn/lifecycle-toolkit/commit/9060ae1d7ee738c39e14fd6cb1691bfd3a386bc4))
* update github.com/keptn/lifecycle-toolkit/metrics-operator digest to 472eac0 ([#1594](https://github.com/keptn/lifecycle-toolkit/issues/1594)) ([7087bb6](https://github.com/keptn/lifecycle-toolkit/commit/7087bb6972fcddfc1ebc5a1d3bc5b2c01246ddfd))
* update kubernetes packages (patch) ([#1634](https://github.com/keptn/lifecycle-toolkit/issues/1634)) ([4ad9bbf](https://github.com/keptn/lifecycle-toolkit/commit/4ad9bbf74960f6d744d97e43b0e4e6130892ab36))
* update module github.com/datadog/datadog-api-client-go/v2 to v2.14.0 ([#1656](https://github.com/keptn/lifecycle-toolkit/issues/1656)) ([2b1a1e9](https://github.com/keptn/lifecycle-toolkit/commit/2b1a1e9cc07133f575c258510f9bac05f13cd603))
* update module github.com/onsi/gomega to v1.27.8 ([#1552](https://github.com/keptn/lifecycle-toolkit/issues/1552)) ([fe9e7ec](https://github.com/keptn/lifecycle-toolkit/commit/fe9e7ecc0bababfbf0d4fe5a35c6c735d4c28c76))
* update module github.com/prometheus/client_golang to v1.16.0 ([#1657](https://github.com/keptn/lifecycle-toolkit/issues/1657)) ([c2e56c5](https://github.com/keptn/lifecycle-toolkit/commit/c2e56c574dd660f4b8b32ede487a21b892a67daf))
* update sigstore/cosign-installer action to v3.1.0 ([#1627](https://github.com/keptn/lifecycle-toolkit/issues/1627)) ([a23ba71](https://github.com/keptn/lifecycle-toolkit/commit/a23ba710d00139e2d82d1676c23ef4a993d76e0e))
* update sigstore/cosign-installer action to v3.1.1 ([#1644](https://github.com/keptn/lifecycle-toolkit/issues/1644)) ([c93a496](https://github.com/keptn/lifecycle-toolkit/commit/c93a496330f720568e0df8266151d0b58915dc55))


### Docs

* `deno` rather than `function` for deno-runtime runner ([#1611](https://github.com/keptn/lifecycle-toolkit/issues/1611)) ([72f5b82](https://github.com/keptn/lifecycle-toolkit/commit/72f5b82c6e4cf20114b5d716bda7fab99fe72ae9))
* add optional field in secretKeyRef ([#1590](https://github.com/keptn/lifecycle-toolkit/issues/1590)) ([d0d5bcb](https://github.com/keptn/lifecycle-toolkit/commit/d0d5bcbaa97a019444ef3a78d5967ef178ed4fd3))
* explain namespaces for metrics and evaluations ([#1641](https://github.com/keptn/lifecycle-toolkit/issues/1641)) ([72f7038](https://github.com/keptn/lifecycle-toolkit/commit/72f7038e55719dd6b7fba32ee8645aabc36d69e5))
* fix &lt;/details&gt; ending under watch pods in _index.md ([#1636](https://github.com/keptn/lifecycle-toolkit/issues/1636)) ([7274cf7](https://github.com/keptn/lifecycle-toolkit/commit/7274cf7d27e2aa35901b692643ae3145b2610ee0))
* fix installation steps in E2E example ([#1645](https://github.com/keptn/lifecycle-toolkit/issues/1645)) ([d6f4307](https://github.com/keptn/lifecycle-toolkit/commit/d6f4307effd89717c3abb339a85850de8d1a838d))
* minor improvement for docs contribution guide ([#1613](https://github.com/keptn/lifecycle-toolkit/issues/1613)) ([6eac2bf](https://github.com/keptn/lifecycle-toolkit/commit/6eac2bf404ab1357903f8df990540379773eb11c))
* move Dora metrics info to "Implementing" section ([#1639](https://github.com/keptn/lifecycle-toolkit/issues/1639)) ([55ee941](https://github.com/keptn/lifecycle-toolkit/commit/55ee94122c2147472f3984d475ea780fa3259789))
* re-generate CRD docs ([#1612](https://github.com/keptn/lifecycle-toolkit/issues/1612)) ([36845da](https://github.com/keptn/lifecycle-toolkit/commit/36845da8397d9067192498bb563f45d42ef60079))
* set up Migration Guide ([#1506](https://github.com/keptn/lifecycle-toolkit/issues/1506)) ([c2e9f4a](https://github.com/keptn/lifecycle-toolkit/commit/c2e9f4a29890ab384769ca4c679b60efdbda3f57))


### Other

* bump release-please base commit ([#1621](https://github.com/keptn/lifecycle-toolkit/issues/1621)) ([c11bba3](https://github.com/keptn/lifecycle-toolkit/commit/c11bba3cd9625a819871bc66ee678ab56d93c773))
* **operator:** print trace in the logs only if the collector is not enabled ([c7d35b8](https://github.com/keptn/lifecycle-toolkit/commit/c7d35b8f1f26e73a7846ae4d2a39dff3b0cc87f7))
* **operator:** refactor k8s Event sending mechanism ([#1687](https://github.com/keptn/lifecycle-toolkit/issues/1687)) ([20839af](https://github.com/keptn/lifecycle-toolkit/commit/20839affec4fd637fe9f5d658be1549b2eaf8b5e))
* replace stale bot with GH actions workflow ([#1629](https://github.com/keptn/lifecycle-toolkit/issues/1629)) ([351c092](https://github.com/keptn/lifecycle-toolkit/commit/351c09296ac4f9beeab6ef2cad457e9689ddc0a5))

## [0.8.0](https://github.com/keptn/lifecycle-toolkit/compare/v0.7.1...v0.8.0) (2023-06-21)


### ⚠ BREAKING CHANGES

* **operator:** support python-runtime runner for KeptnTasks
* **operator:** support container-runtime runner for KeptnTasks

### Features

* add python-runtime ([#1496](https://github.com/keptn/lifecycle-toolkit/issues/1496)) ([76a4bd9](https://github.com/keptn/lifecycle-toolkit/commit/76a4bd92607d05c16c63ccc4c1dd91e35cb4d6b0))
* add validating webhook for KeptnTaskDefinition ([#1514](https://github.com/keptn/lifecycle-toolkit/issues/1514)) ([d55a7ef](https://github.com/keptn/lifecycle-toolkit/commit/d55a7ef421f15305ef711babe9ba0b682dbf65ef))
* **cert-manager:** additional options for targeting WebhookConfigurations and CRDs ([#1276](https://github.com/keptn/lifecycle-toolkit/issues/1276)) ([dadd70b](https://github.com/keptn/lifecycle-toolkit/commit/dadd70b2d63d458c881bda2d6f5755ac257f32e2))
* **metrics-operator:** introduce ErrMsg field into KeptnMetric status ([#1365](https://github.com/keptn/lifecycle-toolkit/issues/1365)) ([092d284](https://github.com/keptn/lifecycle-toolkit/commit/092d28499a74d0ac11c69400bc9454ee2285366d))
* **operator:** adapt TaskDefinition validation webhook to consider python and deno runtime ([#1534](https://github.com/keptn/lifecycle-toolkit/issues/1534)) ([59cdfc8](https://github.com/keptn/lifecycle-toolkit/commit/59cdfc882b7075244a91f0aaf6633d843d9e0099))
* **operator:** introduce fallback search to KLT default namespace when KeptnEvaluationDefinition is not found ([#1359](https://github.com/keptn/lifecycle-toolkit/issues/1359)) ([d5ddf26](https://github.com/keptn/lifecycle-toolkit/commit/d5ddf266a737d3a69d5919f4231a03732c59694f))
* **operator:** support container-runtime runner for KeptnTasks ([02ce860](https://github.com/keptn/lifecycle-toolkit/commit/02ce86023b3db175481b859f379cb4298d03566a))
* **operator:** support python-runtime runner for KeptnTasks ([b79f7c4](https://github.com/keptn/lifecycle-toolkit/commit/b79f7c421d15a8456d50754b631bef61fc0c8dd8))
* **operator:** trim KeptnAppVersion name that exceed max limit ([#1296](https://github.com/keptn/lifecycle-toolkit/issues/1296)) ([0bf2f9e](https://github.com/keptn/lifecycle-toolkit/commit/0bf2f9e78f6a65d79eed0135e49289816e9a2533))


### Bug Fixes

* added the missing link  ([#1537](https://github.com/keptn/lifecycle-toolkit/issues/1537)) ([27fb2c2](https://github.com/keptn/lifecycle-toolkit/commit/27fb2c23bdda13efa047d4101ea8db6595d936a1))
* **cert-manager:** avoid index-out-of-bounds error when updating webhook configs ([#1497](https://github.com/keptn/lifecycle-toolkit/issues/1497)) ([0f28b8c](https://github.com/keptn/lifecycle-toolkit/commit/0f28b8c2b5854944d9b0e72fee97a6bc91f39bea))
* **helm-chart:** fix Python runtime version number ([#1586](https://github.com/keptn/lifecycle-toolkit/issues/1586)) ([05572ec](https://github.com/keptn/lifecycle-toolkit/commit/05572ec6a1a6dd7a56c786333af046eea19929bd))
* **metrics-operator:** improve error handling in metrics providers ([#1466](https://github.com/keptn/lifecycle-toolkit/issues/1466)) ([9801e5d](https://github.com/keptn/lifecycle-toolkit/commit/9801e5dfe9e17fc6c30ef832d97439955964fdcc))
* **metrics-operator:** introduce IsStatusSet method to KeptnMetric ([#1427](https://github.com/keptn/lifecycle-toolkit/issues/1427)) ([24a60f5](https://github.com/keptn/lifecycle-toolkit/commit/24a60f5e6f8f3a383dfce554d644bfd974c4b5fd))
* **operator:** use new RuntimeSpec instead of FunctionSpec ([#1529](https://github.com/keptn/lifecycle-toolkit/issues/1529)) ([6189317](https://github.com/keptn/lifecycle-toolkit/commit/61893175eb73d3f559c7086c09c8939485dd23d3))
* remove scarf redirect from containers images ([#1443](https://github.com/keptn/lifecycle-toolkit/issues/1443)) ([a20b2e7](https://github.com/keptn/lifecycle-toolkit/commit/a20b2e707fd2c0bb03b661c6a6cca272eb088ee1))
* restore go files ([#1371](https://github.com/keptn/lifecycle-toolkit/issues/1371)) ([9a4a6fd](https://github.com/keptn/lifecycle-toolkit/commit/9a4a6fd026bbdbfe986449373bad2b116c34b3d4))


### Other

* add example for python task definition ([#1554](https://github.com/keptn/lifecycle-toolkit/issues/1554)) ([908b03d](https://github.com/keptn/lifecycle-toolkit/commit/908b03d7ae6d4e6e5136257611f0c402f32e08f0))
* bump up helm chart version ([#1351](https://github.com/keptn/lifecycle-toolkit/issues/1351)) ([737d478](https://github.com/keptn/lifecycle-toolkit/commit/737d4782c6c90c77930f58590dcd1098b68b6ef1))
* **cert-manager:** updated readme of cert-manager ([#1393](https://github.com/keptn/lifecycle-toolkit/issues/1393)) ([12fcca8](https://github.com/keptn/lifecycle-toolkit/commit/12fcca82b904ba2f2ca72347cb176771d84d63df))
* minor refactoring of the evaluation controller ([#1356](https://github.com/keptn/lifecycle-toolkit/issues/1356)) ([4398e96](https://github.com/keptn/lifecycle-toolkit/commit/4398e9677ca60a4dd10bd7198479796f36f26026))
* **operator:** bump OTel dependencies to the latest version ([#1419](https://github.com/keptn/lifecycle-toolkit/issues/1419)) ([a7475c2](https://github.com/keptn/lifecycle-toolkit/commit/a7475c2ae13479fed55fa4a322af3c2a47649fa1))
* **operator:** explicitly define ImagePullPolicy of Job container to IfNotPresent ([#1509](https://github.com/keptn/lifecycle-toolkit/issues/1509)) ([bb916f3](https://github.com/keptn/lifecycle-toolkit/commit/bb916f3e3875ec4b3c3e5efbb9f1a65be2a58196))
* **operator:** make use of status.jobName when searching for job in KeptnTask controller ([#1436](https://github.com/keptn/lifecycle-toolkit/issues/1436)) ([28dd6b7](https://github.com/keptn/lifecycle-toolkit/commit/28dd6b77c4cacd038539e30ac8275d6f63d39155))
* **operator:** refactor KeptnTask controller logic ([#1536](https://github.com/keptn/lifecycle-toolkit/issues/1536)) ([ed85fc9](https://github.com/keptn/lifecycle-toolkit/commit/ed85fc972cd676ca0be05dae10aabf90e969d503))
* **operator:** refactor keptntaskcontroller to use builder interface ([#1450](https://github.com/keptn/lifecycle-toolkit/issues/1450)) ([a3f5e5b](https://github.com/keptn/lifecycle-toolkit/commit/a3f5e5bc3fc8bd8073d264c30c39a38fc09d364e))
* **operator:** use List() when fetching KeptnWorkloadVersions for KeptnAppVersion ([#1456](https://github.com/keptn/lifecycle-toolkit/issues/1456)) ([ecd8c48](https://github.com/keptn/lifecycle-toolkit/commit/ecd8c487b22b11bea0646a3c5b2a1f9a22c80d2f))
* remove code duplication ([#1372](https://github.com/keptn/lifecycle-toolkit/issues/1372)) ([da66c80](https://github.com/keptn/lifecycle-toolkit/commit/da66c80653b4a992fd94e49b067f4a21bdf3978b))
* remove decoder injector interface from webhook ([#1563](https://github.com/keptn/lifecycle-toolkit/issues/1563)) ([7850766](https://github.com/keptn/lifecycle-toolkit/commit/785076613942995fdda8882dcbb74b4b1963675e))
* remove space in python sample folder ([#1550](https://github.com/keptn/lifecycle-toolkit/issues/1550)) ([53443ac](https://github.com/keptn/lifecycle-toolkit/commit/53443ac4042e5f78e402df7378770fd163b5167c))
* standardize generation of resource names ([#1472](https://github.com/keptn/lifecycle-toolkit/issues/1472)) ([f7abcb0](https://github.com/keptn/lifecycle-toolkit/commit/f7abcb096838c0071c07b95bf6ff938de9be4975))
* use cert-manager library in lifecycle-operator and metrics-operator to reduce code duplication ([#1379](https://github.com/keptn/lifecycle-toolkit/issues/1379)) ([831fc46](https://github.com/keptn/lifecycle-toolkit/commit/831fc46d9e4ebb059473f137ef6c012373c6179c))
* website edit links should point to page ([#1566](https://github.com/keptn/lifecycle-toolkit/issues/1566)) ([8b62f33](https://github.com/keptn/lifecycle-toolkit/commit/8b62f3355afc178f15416b930ea1c77844f2b02c))


### Dependency Updates

* update anchore/sbom-action action to v0.14.2 ([#1401](https://github.com/keptn/lifecycle-toolkit/issues/1401)) ([9085785](https://github.com/keptn/lifecycle-toolkit/commit/9085785b669bbc5bdd418afa6e9bd2f81c788653))
* update aquasecurity/trivy-action action to v0.11.0 ([#1531](https://github.com/keptn/lifecycle-toolkit/issues/1531)) ([66c9505](https://github.com/keptn/lifecycle-toolkit/commit/66c95058235deddac86daaa6bf897ae680181052))
* update aquasecurity/trivy-action action to v0.11.2 ([#1551](https://github.com/keptn/lifecycle-toolkit/issues/1551)) ([2d588db](https://github.com/keptn/lifecycle-toolkit/commit/2d588dbfd2405b9ff09b68a38f6b3ff5e6779d63))
* update busybox docker tag to v1.36.1 ([#1437](https://github.com/keptn/lifecycle-toolkit/issues/1437)) ([9ba5cae](https://github.com/keptn/lifecycle-toolkit/commit/9ba5cae8b3f0be8b380e28883530f97db76773dc))
* update checkmarx/kics-github-action action to v1.7.0 ([#1435](https://github.com/keptn/lifecycle-toolkit/issues/1435)) ([f9d609c](https://github.com/keptn/lifecycle-toolkit/commit/f9d609c47545c8fa772329056606891534a6eed6))
* update curlimages/curl docker tag to v8.1.0 ([#1439](https://github.com/keptn/lifecycle-toolkit/issues/1439)) ([9e90d17](https://github.com/keptn/lifecycle-toolkit/commit/9e90d17c211709b357b40ef8f0843a9e1bf0364f))
* update curlimages/curl docker tag to v8.1.1 ([#1455](https://github.com/keptn/lifecycle-toolkit/issues/1455)) ([d1279a9](https://github.com/keptn/lifecycle-toolkit/commit/d1279a9fe0f09177449b20d4b3fc8f0f3c10d81a))
* update curlimages/curl docker tag to v8.1.2 ([#1530](https://github.com/keptn/lifecycle-toolkit/issues/1530)) ([ef3e89e](https://github.com/keptn/lifecycle-toolkit/commit/ef3e89e148afd1d5dd54ea47f7973ac8532dd203))
* update dependency argoproj/argo-cd to v2.7.1 ([#1374](https://github.com/keptn/lifecycle-toolkit/issues/1374)) ([9b9a973](https://github.com/keptn/lifecycle-toolkit/commit/9b9a973a95ca59e91927695f92c0d56389be3f4f))
* update dependency argoproj/argo-cd to v2.7.2 ([#1423](https://github.com/keptn/lifecycle-toolkit/issues/1423)) ([e381f7f](https://github.com/keptn/lifecycle-toolkit/commit/e381f7fc6d79703b9f32dbf49331247107597b20))
* update dependency argoproj/argo-cd to v2.7.3 ([#1512](https://github.com/keptn/lifecycle-toolkit/issues/1512)) ([6146e79](https://github.com/keptn/lifecycle-toolkit/commit/6146e79a62d12de37fec9218c5f6f7acd2255b82))
* update dependency argoproj/argo-cd to v2.7.4 ([#1541](https://github.com/keptn/lifecycle-toolkit/issues/1541)) ([712bd9a](https://github.com/keptn/lifecycle-toolkit/commit/712bd9aa8ee08c6de61a6cb63d7633b3ef00bf38))
* update dependency autoprefixer to v10.4.14 ([#1560](https://github.com/keptn/lifecycle-toolkit/issues/1560)) ([a07261e](https://github.com/keptn/lifecycle-toolkit/commit/a07261ee6b3600f2106b547576e51f93c327655b))
* update dependency golangci/golangci-lint to v1.53.2 ([#1538](https://github.com/keptn/lifecycle-toolkit/issues/1538)) ([e387822](https://github.com/keptn/lifecycle-toolkit/commit/e3878222a7994a253b12923e852c0e25c49756af))
* update dependency helm/helm to v3.12.0 ([#1440](https://github.com/keptn/lifecycle-toolkit/issues/1440)) ([aff755d](https://github.com/keptn/lifecycle-toolkit/commit/aff755d4539310787d47a448f1fc6600ffd04c33))
* update dependency jaegertracing/jaeger to v1.45.0 ([#1407](https://github.com/keptn/lifecycle-toolkit/issues/1407)) ([dab62de](https://github.com/keptn/lifecycle-toolkit/commit/dab62dea8b5c0a45baba9d9b3e33c1e6f6f640e4))
* update dependency jaegertracing/jaeger-operator to v1.44.0 ([#1258](https://github.com/keptn/lifecycle-toolkit/issues/1258)) ([dab73fb](https://github.com/keptn/lifecycle-toolkit/commit/dab73fb94f85022d84436453a84bf19f7f95cc5c))
* update dependency jaegertracing/jaeger-operator to v1.45.0 ([#1478](https://github.com/keptn/lifecycle-toolkit/issues/1478)) ([7bc4feb](https://github.com/keptn/lifecycle-toolkit/commit/7bc4feb66cd303235a407f9341cd723562262dca))
* update dependency kubernetes-sigs/controller-tools to v0.12.0 ([#1383](https://github.com/keptn/lifecycle-toolkit/issues/1383)) ([0a6b7e7](https://github.com/keptn/lifecycle-toolkit/commit/0a6b7e795a9e58425ab6baacf38a82a22dbbc0c8))
* update dependency kubernetes-sigs/kustomize to v5.0.3 ([#1402](https://github.com/keptn/lifecycle-toolkit/issues/1402)) ([fad37af](https://github.com/keptn/lifecycle-toolkit/commit/fad37afd14cd781c9561d43dd2f1af6824973693))
* update github.com/keptn/lifecycle-toolkit/klt-cert-manager digest to 65b4139 ([#1429](https://github.com/keptn/lifecycle-toolkit/issues/1429)) ([57fdcdd](https://github.com/keptn/lifecycle-toolkit/commit/57fdcddcf73c71dde07641cb13f9c7b16cff6cf5))
* update github.com/keptn/lifecycle-toolkit/klt-cert-manager digest to 7c4d2ab ([#1510](https://github.com/keptn/lifecycle-toolkit/issues/1510)) ([36d0c23](https://github.com/keptn/lifecycle-toolkit/commit/36d0c23ea66646b873e5c4b571b6c1c4b3102dfd))
* update github.com/keptn/lifecycle-toolkit/klt-cert-manager digest to 9eafb78 ([#1454](https://github.com/keptn/lifecycle-toolkit/issues/1454)) ([b66ad6f](https://github.com/keptn/lifecycle-toolkit/commit/b66ad6fab019640380a11acf837b3589605e6219))
* update github.com/keptn/lifecycle-toolkit/klt-cert-manager digest to e381f7f ([#1422](https://github.com/keptn/lifecycle-toolkit/issues/1422)) ([daedf87](https://github.com/keptn/lifecycle-toolkit/commit/daedf878eeb8d00d717e5746b46dd651c6fba8de))
* update github.com/keptn/lifecycle-toolkit/metrics-operator digest to 36d0c23 ([#1511](https://github.com/keptn/lifecycle-toolkit/issues/1511)) ([3b59742](https://github.com/keptn/lifecycle-toolkit/commit/3b59742d3ad2793a628a9679ba4e7c3bb4a8c488))
* update github.com/keptn/lifecycle-toolkit/metrics-operator digest to 57fdcdd ([#1430](https://github.com/keptn/lifecycle-toolkit/issues/1430)) ([54a9384](https://github.com/keptn/lifecycle-toolkit/commit/54a93840094e7dd3c9799e5e0a8ae889d51bb2ac))
* update github.com/keptn/lifecycle-toolkit/metrics-operator digest to bb916f3 ([#1463](https://github.com/keptn/lifecycle-toolkit/issues/1463)) ([4292570](https://github.com/keptn/lifecycle-toolkit/commit/4292570ec3256b9aa2291f5abc5769ef22e3cdf2))
* update github.com/keptn/lifecycle-toolkit/metrics-operator digest to e381f7f ([#1268](https://github.com/keptn/lifecycle-toolkit/issues/1268)) ([f0f7edf](https://github.com/keptn/lifecycle-toolkit/commit/f0f7edf7041d438b8d8804ad9341ef878ed625de))
* update golang docker tag to v1.20.4 ([#1346](https://github.com/keptn/lifecycle-toolkit/issues/1346)) ([8fedf0f](https://github.com/keptn/lifecycle-toolkit/commit/8fedf0f11c6f4e55e1ac47ab8e80705e189ffff8))
* update helm/kind-action action to v1.7.0 ([#1479](https://github.com/keptn/lifecycle-toolkit/issues/1479)) ([fb22707](https://github.com/keptn/lifecycle-toolkit/commit/fb22707ee148bb5c3ae2a724ca1bde3c0d5df929))
* update kubernetes packages (patch) ([#1432](https://github.com/keptn/lifecycle-toolkit/issues/1432)) ([7f5b3ab](https://github.com/keptn/lifecycle-toolkit/commit/7f5b3abb87ebe8e0c040e415f69ed12e25ebb7fd))
* update module github.com/argoproj/argo-rollouts to v1.5.0 ([#1408](https://github.com/keptn/lifecycle-toolkit/issues/1408)) ([2f75e73](https://github.com/keptn/lifecycle-toolkit/commit/2f75e739ca8fa218b3d15ccf657a5d85530eecf5))
* update module github.com/argoproj/argo-rollouts to v1.5.1 ([#1513](https://github.com/keptn/lifecycle-toolkit/issues/1513)) ([de95b50](https://github.com/keptn/lifecycle-toolkit/commit/de95b50ae10da0e0e6c40eb5545fd98bc0f5ffcb))
* update module github.com/benbjohnson/clock to v1.3.4 ([#1403](https://github.com/keptn/lifecycle-toolkit/issues/1403)) ([f88dfd5](https://github.com/keptn/lifecycle-toolkit/commit/f88dfd5c0d7a544d94054cce6693ebd3d88f0a9f))
* update module github.com/benbjohnson/clock to v1.3.5 ([#1464](https://github.com/keptn/lifecycle-toolkit/issues/1464)) ([abf10bf](https://github.com/keptn/lifecycle-toolkit/commit/abf10bfaf033a1f57b8f65d3c7127dd962926ed4))
* update module github.com/datadog/datadog-api-client-go/v2 to v2.13.0 ([#1519](https://github.com/keptn/lifecycle-toolkit/issues/1519)) ([d774568](https://github.com/keptn/lifecycle-toolkit/commit/d77456836530c9a456f03a8f4330082514077388))
* update module github.com/imdario/mergo to v0.3.16 ([#1482](https://github.com/keptn/lifecycle-toolkit/issues/1482)) ([9eafb78](https://github.com/keptn/lifecycle-toolkit/commit/9eafb78b51d60d13af44cad281c8c631b02773c3))
* update module github.com/onsi/ginkgo/v2 to v2.11.0 ([#1553](https://github.com/keptn/lifecycle-toolkit/issues/1553)) ([7c4d2ab](https://github.com/keptn/lifecycle-toolkit/commit/7c4d2abcb84e672d5c800fed0ed59a0672d112c3))
* update module github.com/onsi/ginkgo/v2 to v2.9.4 ([#1384](https://github.com/keptn/lifecycle-toolkit/issues/1384)) ([2ed8dd7](https://github.com/keptn/lifecycle-toolkit/commit/2ed8dd7a7d62a44bab22cc1da11f80e02fd129f8))
* update module github.com/onsi/ginkgo/v2 to v2.9.5 ([#1433](https://github.com/keptn/lifecycle-toolkit/issues/1433)) ([fcdd9fe](https://github.com/keptn/lifecycle-toolkit/commit/fcdd9fea3860ba9d5ec52b3733a48258df4e8549))
* update module github.com/onsi/ginkgo/v2 to v2.9.7 ([#1517](https://github.com/keptn/lifecycle-toolkit/issues/1517)) ([225c04b](https://github.com/keptn/lifecycle-toolkit/commit/225c04b6e9c5a848dc5e2f37239c77d893bec539))
* update module github.com/onsi/gomega to v1.27.7 ([#1473](https://github.com/keptn/lifecycle-toolkit/issues/1473)) ([50f7415](https://github.com/keptn/lifecycle-toolkit/commit/50f7415a832f2cc4e90db5faf51f17cf471558cc))
* update module github.com/open-feature/go-sdk to v1.4.0 ([#1516](https://github.com/keptn/lifecycle-toolkit/issues/1516)) ([a2ef768](https://github.com/keptn/lifecycle-toolkit/commit/a2ef76846ed9853a5fe2bbe88d0ae35cf1fb82f8))
* update module github.com/prometheus/client_golang to v1.15.1 ([#1386](https://github.com/keptn/lifecycle-toolkit/issues/1386)) ([8b73046](https://github.com/keptn/lifecycle-toolkit/commit/8b730461a9892f5ab06e51ee9519ec6fa7d83125))
* update module github.com/prometheus/common to v0.44.0 ([#1452](https://github.com/keptn/lifecycle-toolkit/issues/1452)) ([af22685](https://github.com/keptn/lifecycle-toolkit/commit/af2268566b74b251da17dec5576af3cd03159482))
* update module github.com/stretchr/testify to v1.8.3 ([#1434](https://github.com/keptn/lifecycle-toolkit/issues/1434)) ([65b4139](https://github.com/keptn/lifecycle-toolkit/commit/65b41399b2e0d5c4109af484a80d4bb2c56f9215))
* update module github.com/stretchr/testify to v1.8.4 ([#1515](https://github.com/keptn/lifecycle-toolkit/issues/1515)) ([c732492](https://github.com/keptn/lifecycle-toolkit/commit/c732492a85397f4b986383b5d37b6d73e1dada5b))
* update module golang.org/x/net to v0.10.0 ([#1453](https://github.com/keptn/lifecycle-toolkit/issues/1453)) ([65a3e4b](https://github.com/keptn/lifecycle-toolkit/commit/65a3e4b402d693a64dc9be452aead9c4773d6945))
* update module google.golang.org/grpc to v1.54.1 ([#1404](https://github.com/keptn/lifecycle-toolkit/issues/1404)) ([a5d9b19](https://github.com/keptn/lifecycle-toolkit/commit/a5d9b19901f673768cef63dcc1606aafbc5a1b51))
* update module google.golang.org/grpc to v1.55.0 ([#1480](https://github.com/keptn/lifecycle-toolkit/issues/1480)) ([d5a8e7c](https://github.com/keptn/lifecycle-toolkit/commit/d5a8e7cbf0095119b646f23b74891dcb231e2e0c))
* update module k8s.io/klog/v2 to v2.100.1 ([#1324](https://github.com/keptn/lifecycle-toolkit/issues/1324)) ([6524d58](https://github.com/keptn/lifecycle-toolkit/commit/6524d583dc9d99bd67b9a599f48f78b6d89a3877))
* update module k8s.io/kubernetes to v1.25.10 ([#1475](https://github.com/keptn/lifecycle-toolkit/issues/1475)) ([e65715c](https://github.com/keptn/lifecycle-toolkit/commit/e65715cfe98eebfcdee599253a1e63e482773f4d))
* update octokit/request-action action to v2.1.8 ([#1524](https://github.com/keptn/lifecycle-toolkit/issues/1524)) ([dcc66a0](https://github.com/keptn/lifecycle-toolkit/commit/dcc66a03a1eacfc3e6ac9a094cb430ba7d38314c))
* update octokit/request-action action to v2.1.9 ([#1533](https://github.com/keptn/lifecycle-toolkit/issues/1533)) ([bba7339](https://github.com/keptn/lifecycle-toolkit/commit/bba73398c9e2b2f076618f8161f0d4fe74df8207))
* update sigstore/cosign-installer action to v3.0.3 ([#1308](https://github.com/keptn/lifecycle-toolkit/issues/1308)) ([ac98fe5](https://github.com/keptn/lifecycle-toolkit/commit/ac98fe566f2652eebdd6e578a6f3491df9e471d1))
* update sigstore/cosign-installer action to v3.0.5 ([#1438](https://github.com/keptn/lifecycle-toolkit/issues/1438)) ([1fba2b4](https://github.com/keptn/lifecycle-toolkit/commit/1fba2b4985a424c728ca02747c56a343fcf3fdbe))


### Docs

* add cluster requirements ([#1364](https://github.com/keptn/lifecycle-toolkit/issues/1364)) ([e06b01e](https://github.com/keptn/lifecycle-toolkit/commit/e06b01e4b3723b16b8479f3b22fa3021e8dead55))
* add content to implementing/otel page ([#1492](https://github.com/keptn/lifecycle-toolkit/issues/1492)) ([452c3a9](https://github.com/keptn/lifecycle-toolkit/commit/452c3a917cfe2679553fbecc46c133da4887dc2d))
* add docs for Python runtime ([#1549](https://github.com/keptn/lifecycle-toolkit/issues/1549)) ([2e53fda](https://github.com/keptn/lifecycle-toolkit/commit/2e53fda17183bd1fcffa0d8aaa03c305088465b4))
* add info about automatic application discovery ([#1353](https://github.com/keptn/lifecycle-toolkit/issues/1353)) ([d42d023](https://github.com/keptn/lifecycle-toolkit/commit/d42d023d1d431782deb0c3ef8fa20fa2f2375ad3))
* added comments to document the meaning of CRD properties ([#1360](https://github.com/keptn/lifecycle-toolkit/issues/1360)) ([a8bc440](https://github.com/keptn/lifecycle-toolkit/commit/a8bc440a4f15f624455c513373033c78a31a53b5))
* content for KeptnTaskDefinition ref and tasks guide ([#1392](https://github.com/keptn/lifecycle-toolkit/issues/1392)) ([13b0495](https://github.com/keptn/lifecycle-toolkit/commit/13b04956a02a0384bfc1ad6b043e901613d1d5b2))
* create "observability" getting started guide ([#1376](https://github.com/keptn/lifecycle-toolkit/issues/1376)) ([4815986](https://github.com/keptn/lifecycle-toolkit/commit/48159862bf89b8cc1692500af7b05487a6cc03cb))
* create keptn metrics getting started ([#1375](https://github.com/keptn/lifecycle-toolkit/issues/1375)) ([8de6d8f](https://github.com/keptn/lifecycle-toolkit/commit/8de6d8f8ca34c576466d9cc8b32d1d3865123ad8))
* create KeptnApp reference page ([#1391](https://github.com/keptn/lifecycle-toolkit/issues/1391)) ([4aa141a](https://github.com/keptn/lifecycle-toolkit/commit/4aa141a069b8b7d25c508ff92309ad460120beb4))
* create KeptnConfig yaml ref page ([#1369](https://github.com/keptn/lifecycle-toolkit/issues/1369)) ([e40292c](https://github.com/keptn/lifecycle-toolkit/commit/e40292ce995070a492187f5dcc7db363e03eb260))
* create pre/post-deploy getting started ([#1362](https://github.com/keptn/lifecycle-toolkit/issues/1362)) ([d602115](https://github.com/keptn/lifecycle-toolkit/commit/d602115ef4b158b04be3c630ca45c6e4f39fc0f3))
* delete obsolete pages ([#1520](https://github.com/keptn/lifecycle-toolkit/issues/1520)) ([96e69c2](https://github.com/keptn/lifecycle-toolkit/commit/96e69c2a10e9a0446b1992a880ec74400bbe0cfe))
* document container-runtime and python-runtime runners ([#1579](https://github.com/keptn/lifecycle-toolkit/issues/1579)) ([3834b70](https://github.com/keptn/lifecycle-toolkit/commit/3834b709da4552a07c8acd8699d1e8a583a621bd))
* enhance install page ([#1399](https://github.com/keptn/lifecycle-toolkit/issues/1399)) ([025709e](https://github.com/keptn/lifecycle-toolkit/commit/025709e3147abef79d2ddbecb795db0c5e8bf2a8))
* final polish of getting started guides ([#1449](https://github.com/keptn/lifecycle-toolkit/issues/1449)) ([30e6647](https://github.com/keptn/lifecycle-toolkit/commit/30e664703c3b42aa5c2049535d528f69cbcfe4b4))
* fix edit links ([#1583](https://github.com/keptn/lifecycle-toolkit/issues/1583)) ([1384679](https://github.com/keptn/lifecycle-toolkit/commit/138467970e966c90a1bd8fa428b4c5efe0a9bd50))
* fix getting started guides ([#1447](https://github.com/keptn/lifecycle-toolkit/issues/1447)) ([6035e55](https://github.com/keptn/lifecycle-toolkit/commit/6035e552d3f46e2553603f711db008784ff99d0e))
* fix link to v1 docs ([#1461](https://github.com/keptn/lifecycle-toolkit/issues/1461)) ([a7f54ad](https://github.com/keptn/lifecycle-toolkit/commit/a7f54ad22c6b13f6bee7dcced2e3ab71bb1be365))
* fix markdown links ([#1414](https://github.com/keptn/lifecycle-toolkit/issues/1414)) ([b2392c1](https://github.com/keptn/lifecycle-toolkit/commit/b2392c1d6a81df92adf6228167a52233eb1757ae))
* fix readme links to point to website instead of repo files ([#1344](https://github.com/keptn/lifecycle-toolkit/issues/1344)) ([e5f0344](https://github.com/keptn/lifecycle-toolkit/commit/e5f034425dd1012c899d60fb2ff9d755853924aa))
* fix typo ([#1578](https://github.com/keptn/lifecycle-toolkit/issues/1578)) ([fe5bbea](https://github.com/keptn/lifecycle-toolkit/commit/fe5bbea0773cf212ce0d9bc72c09182ad5bb8916))
* fix typos ([#1562](https://github.com/keptn/lifecycle-toolkit/issues/1562)) ([be47052](https://github.com/keptn/lifecycle-toolkit/commit/be47052fe76312aab413d29e5132cc88b1fdc8d6))
* improve list on install landing page ([#1400](https://github.com/keptn/lifecycle-toolkit/issues/1400)) ([3d23e29](https://github.com/keptn/lifecycle-toolkit/commit/3d23e29b82d1296627900850b19af7ea2eb30d87))
* mention Prometheus in intro ([#1405](https://github.com/keptn/lifecycle-toolkit/issues/1405)) ([2c51231](https://github.com/keptn/lifecycle-toolkit/commit/2c51231fd700009c1588259de1974e1dfa80e8b8))
* metrics & evaluation ref and guides ([#1385](https://github.com/keptn/lifecycle-toolkit/issues/1385)) ([7712bfa](https://github.com/keptn/lifecycle-toolkit/commit/7712bfae84a21adaf6341ca02ec3589d0459854f))
* misspelled file name, misordered pages ([#1363](https://github.com/keptn/lifecycle-toolkit/issues/1363)) ([be3c2f1](https://github.com/keptn/lifecycle-toolkit/commit/be3c2f1b469a15292bbd698af2200c0a4fb4002e))
* refactor contributing guide - general guidelines ([#1411](https://github.com/keptn/lifecycle-toolkit/issues/1411)) ([7170eec](https://github.com/keptn/lifecycle-toolkit/commit/7170eec563efcaa149dd9bdfb8153a26f846638d))
* refactor contributing guide - linter requirements ([#1412](https://github.com/keptn/lifecycle-toolkit/issues/1412)) ([2ccdec7](https://github.com/keptn/lifecycle-toolkit/commit/2ccdec7e3aafbfda8bb223441cefefbf092084ea))
* refactor contributing guide - local building ([#1484](https://github.com/keptn/lifecycle-toolkit/issues/1484)) ([751552c](https://github.com/keptn/lifecycle-toolkit/commit/751552cb7d5c4acea0b8eb4d55171aabca0675d9))
* refactor contributing guide - source file structure ([#1523](https://github.com/keptn/lifecycle-toolkit/issues/1523)) ([c97b4b9](https://github.com/keptn/lifecycle-toolkit/commit/c97b4b9cc5a02b2ff39a3fb23aa79f69667dd27b))
* regenerate CRD docs ([#1507](https://github.com/keptn/lifecycle-toolkit/issues/1507)) ([672e281](https://github.com/keptn/lifecycle-toolkit/commit/672e281f1b44a7e83449c32e08a9de4a44c8d287))
* remove old "Tasks" section from docs ([#1572](https://github.com/keptn/lifecycle-toolkit/issues/1572)) ([8f0f4f0](https://github.com/keptn/lifecycle-toolkit/commit/8f0f4f0eb4276ba399e1ab77915ef41e690a09db))
* small edit of original Getting Started guide ([#1367](https://github.com/keptn/lifecycle-toolkit/issues/1367)) ([0fd922a](https://github.com/keptn/lifecycle-toolkit/commit/0fd922ad161fb30dcf834d87aec225f48d619f4d))
* update KLT intro page ([#1495](https://github.com/keptn/lifecycle-toolkit/issues/1495)) ([d1db5d2](https://github.com/keptn/lifecycle-toolkit/commit/d1db5d2f054ed56911bc3e91507f9d539dc5829e))
* updated the misspelled word ([#1544](https://github.com/keptn/lifecycle-toolkit/issues/1544)) ([0cb7f14](https://github.com/keptn/lifecycle-toolkit/commit/0cb7f14ae0f6f7844051219e823261a99c6dfe54))

## [0.7.1](https://github.com/keptn/lifecycle-toolkit/compare/v0.7.0...v0.7.1) (2023-05-03)


### Features

* add support for multiple metrics providers ([#1193](https://github.com/keptn/lifecycle-toolkit/issues/1193)) ([3c465d0](https://github.com/keptn/lifecycle-toolkit/commit/3c465d07044b0317cbb6e462004dff9cf8f1d533))
* datadog metric provider for KLT ([#948](https://github.com/keptn/lifecycle-toolkit/issues/948)) ([597a23f](https://github.com/keptn/lifecycle-toolkit/commit/597a23f93433ce56aac7000cf1806dd79f67b3f6))
* improve API reference generation script with path extension ([#1271](https://github.com/keptn/lifecycle-toolkit/issues/1271)) ([74fa4f5](https://github.com/keptn/lifecycle-toolkit/commit/74fa4f56471853e564a91d49132b4f7ce2367f44))
* make examples resource footprint smaller, fix bugs ([#1171](https://github.com/keptn/lifecycle-toolkit/issues/1171)) ([8b165d3](https://github.com/keptn/lifecycle-toolkit/commit/8b165d3bf63a63b452ac2f1423166978b80facc9))
* **operator:** add information about evaluation target in status ([#1341](https://github.com/keptn/lifecycle-toolkit/issues/1341)) ([cc03a85](https://github.com/keptn/lifecycle-toolkit/commit/cc03a8513d469200becb371844e15fd4f832371c))
* **operator:** additional parameters for KeptnTask to support retry logic ([#1084](https://github.com/keptn/lifecycle-toolkit/issues/1084)) ([eed5568](https://github.com/keptn/lifecycle-toolkit/commit/eed5568e51f381f62e2e6db3fddc13a610bcd5e0))
* **operator:** bootstrapped KeptnAppCreationRequest CRD ([#1134](https://github.com/keptn/lifecycle-toolkit/issues/1134)) ([6b58da3](https://github.com/keptn/lifecycle-toolkit/commit/6b58da3c907af591633052ba3d7fc49a2b801ebc))
* **operator:** consider corner cases in KACR controller ([#1270](https://github.com/keptn/lifecycle-toolkit/issues/1270)) ([b3b7010](https://github.com/keptn/lifecycle-toolkit/commit/b3b70109a9125ef4ea017a5f3d25d02146438a46))
* **operator:** create KeptnAppCreationRequest in pod webhook ([#1277](https://github.com/keptn/lifecycle-toolkit/issues/1277)) ([da942c2](https://github.com/keptn/lifecycle-toolkit/commit/da942c2f12fe4a8d5fd89ec0615228d05064b183))
* **operator:** implement KeptnAppCreationRequest controller ([#1191](https://github.com/keptn/lifecycle-toolkit/issues/1191)) ([79afd83](https://github.com/keptn/lifecycle-toolkit/commit/79afd83476baa567285bddc3fc4bc40a76783e67))
* **operator:** introduce fallback search to KLT default namespace when KeptnTaskDefinition is not found ([#1340](https://github.com/keptn/lifecycle-toolkit/issues/1340)) ([6794fe2](https://github.com/keptn/lifecycle-toolkit/commit/6794fe2d7e3334deb17dd13e5580bfc358edb57c))
* **operator:** introduce retry logic for KeptnTasks ([#1088](https://github.com/keptn/lifecycle-toolkit/issues/1088)) ([e49b5a3](https://github.com/keptn/lifecycle-toolkit/commit/e49b5a3f133b5ba5f1ceaad53c12899415ea58b2))
* **operator:** polish `KeptnConfig` and use Env Var for initial configuration ([#1097](https://github.com/keptn/lifecycle-toolkit/issues/1097)) ([559acee](https://github.com/keptn/lifecycle-toolkit/commit/559acee5059016b96703fb9f6f8d842d3c392c29))
* **operator:** propagate KeptnTaskDefinition labels and annotations to Job Pods ([#1283](https://github.com/keptn/lifecycle-toolkit/issues/1283)) ([83be9d9](https://github.com/keptn/lifecycle-toolkit/commit/83be9d98381f7a53c0de324cd868fd03635b52ef))
* **operator:** support Argo Rollout resources ([#879](https://github.com/keptn/lifecycle-toolkit/issues/879)) ([c2b0fa3](https://github.com/keptn/lifecycle-toolkit/commit/c2b0fa35f875d250564f1a75acab6752e65b504d))
* use smaller distroless images for released containers ([#1092](https://github.com/keptn/lifecycle-toolkit/issues/1092)) ([8a7a6af](https://github.com/keptn/lifecycle-toolkit/commit/8a7a6af9f44c3a3f88b0a2f2331e3e820741d26f))


### Bug Fixes

* adapt mapping for community files ([#1215](https://github.com/keptn/lifecycle-toolkit/issues/1215)) ([99ef223](https://github.com/keptn/lifecycle-toolkit/commit/99ef2235648e1fa97d9ae30c7df4551dbb7bcf94))
* add missing control-plane label into lifecycle operator service ([#1148](https://github.com/keptn/lifecycle-toolkit/issues/1148)) ([df04fbe](https://github.com/keptn/lifecycle-toolkit/commit/df04fbe5512d1f9e0c8d8e81a253a8c4892e1dec))
* fix examples restart make command, reduce prometheus resources ([#1158](https://github.com/keptn/lifecycle-toolkit/issues/1158)) ([06b10a8](https://github.com/keptn/lifecycle-toolkit/commit/06b10a82fd3e8c942ef7592f919acb60552c4ae4))
* fix examples, update podtatohead ([#1098](https://github.com/keptn/lifecycle-toolkit/issues/1098)) ([f581ed5](https://github.com/keptn/lifecycle-toolkit/commit/f581ed500f55da2be69a10ac67da5d8717ac3104))
* fix kubecon examples ([#1225](https://github.com/keptn/lifecycle-toolkit/issues/1225)) ([a47fe1d](https://github.com/keptn/lifecycle-toolkit/commit/a47fe1d10d433a121381d0fdd2a9def087f14046))
* fix kubecon examples ([#1226](https://github.com/keptn/lifecycle-toolkit/issues/1226)) ([5fb61ba](https://github.com/keptn/lifecycle-toolkit/commit/5fb61ba00c57bfa4d062d137e2d89781b9d274ea))
* fix metrics demo setup ([#1207](https://github.com/keptn/lifecycle-toolkit/issues/1207)) ([b261172](https://github.com/keptn/lifecycle-toolkit/commit/b261172ff2e2921923ca7d6bb6519a55182bacdf))
* generate missing CRD docs, fix validation pipeline ([#1086](https://github.com/keptn/lifecycle-toolkit/issues/1086)) ([71e9073](https://github.com/keptn/lifecycle-toolkit/commit/71e9073288c55da4ccbc51c1beb01a2d00b0921a))
* helm generation checker pipeline ([#1209](https://github.com/keptn/lifecycle-toolkit/issues/1209)) ([72396cd](https://github.com/keptn/lifecycle-toolkit/commit/72396cda0e8b02913f060a6e99e782be2fab4e85))
* **helm-chart:** fix missing values in the KLT helm chart ([#1082](https://github.com/keptn/lifecycle-toolkit/issues/1082)) ([52311c1](https://github.com/keptn/lifecycle-toolkit/commit/52311c1a1ee023ff5d271456c1fe09737ba94d60))
* **metrics-operator:** normalize Dynatrace URL ([#1145](https://github.com/keptn/lifecycle-toolkit/issues/1145)) ([b33b4f4](https://github.com/keptn/lifecycle-toolkit/commit/b33b4f49320bf75c5098190d0b36ab3c49be9b45))
* move prometheus install into make file ([#1093](https://github.com/keptn/lifecycle-toolkit/issues/1093)) ([f6f44e4](https://github.com/keptn/lifecycle-toolkit/commit/f6f44e4c5d880f00605e2923292c039503bb7903))
* **operator:** fix otel collector URL setup  ([#1262](https://github.com/keptn/lifecycle-toolkit/issues/1262)) ([c3754b7](https://github.com/keptn/lifecycle-toolkit/commit/c3754b755146e27bedd73a18d5f66d9c01a46677))
* **operator:** look up latest AppVersion based on creation timestamp ([#1186](https://github.com/keptn/lifecycle-toolkit/issues/1186)) ([45a96e7](https://github.com/keptn/lifecycle-toolkit/commit/45a96e7fdf464d1248c3674787966d2e5ae50828))
* removed failure branch ([#1175](https://github.com/keptn/lifecycle-toolkit/issues/1175)) ([66df012](https://github.com/keptn/lifecycle-toolkit/commit/66df01257a1abce26e4b3577527a8fbe651358d6))
* security pipeline ([#1333](https://github.com/keptn/lifecycle-toolkit/issues/1333)) ([79e475c](https://github.com/keptn/lifecycle-toolkit/commit/79e475ce1f8d7f05d88365d8711fb491eae0c374))
* use correct control-plane label for metrics-operator ([#1147](https://github.com/keptn/lifecycle-toolkit/issues/1147)) ([1035183](https://github.com/keptn/lifecycle-toolkit/commit/10351834d70b000b8296bea65b80fc9e22e54cee))
* use custom k8s label to inject certificates where needed ([#1288](https://github.com/keptn/lifecycle-toolkit/issues/1288)) ([8fe5df3](https://github.com/keptn/lifecycle-toolkit/commit/8fe5df34e2e4f5ef544a1952040ec1b170148d7a))
* use hash as revision instead of generation number ([#1243](https://github.com/keptn/lifecycle-toolkit/issues/1243)) ([2ad5d81](https://github.com/keptn/lifecycle-toolkit/commit/2ad5d811921834c7049e76879cbf61c819f1a39d))


### Dependency Updates

* bump denoland/deno to 1.32.5 ([#1329](https://github.com/keptn/lifecycle-toolkit/issues/1329)) ([73f0af0](https://github.com/keptn/lifecycle-toolkit/commit/73f0af062832dc0d86297fb2305c287450c3bc05))
* remove github.com/open-feature/flagd ([#1110](https://github.com/keptn/lifecycle-toolkit/issues/1110)) ([e118851](https://github.com/keptn/lifecycle-toolkit/commit/e11885180bcab7dac93409fd5868328c6dade508))
* update actions/setup-go action to v4 ([#1051](https://github.com/keptn/lifecycle-toolkit/issues/1051)) ([8b470d4](https://github.com/keptn/lifecycle-toolkit/commit/8b470d4c8e7285b481df5ae2e1e4674413caaaab))
* update amannn/action-semantic-pull-request action to v5.2.0 ([#1102](https://github.com/keptn/lifecycle-toolkit/issues/1102)) ([c57b1fe](https://github.com/keptn/lifecycle-toolkit/commit/c57b1febc0501648c5c2d3c94c4434536ece2871))
* update anchore/sbom-action action to v0.13.4 ([#1101](https://github.com/keptn/lifecycle-toolkit/issues/1101)) ([4c9a1aa](https://github.com/keptn/lifecycle-toolkit/commit/4c9a1aabe6550d21d80c8696421203896847297e))
* update anchore/sbom-action action to v0.14.1 ([#1187](https://github.com/keptn/lifecycle-toolkit/issues/1187)) ([21e72a3](https://github.com/keptn/lifecycle-toolkit/commit/21e72a3a38dafda79d2cdb9f129bf48382696821))
* update aquasecurity/trivy-action action to v0.10.0 ([#1255](https://github.com/keptn/lifecycle-toolkit/issues/1255)) ([1ff448c](https://github.com/keptn/lifecycle-toolkit/commit/1ff448cafb76c385ef358dd36392b99371d0561e))
* update curlimages/curl docker tag to v8 ([#1116](https://github.com/keptn/lifecycle-toolkit/issues/1116)) ([05bf675](https://github.com/keptn/lifecycle-toolkit/commit/05bf6750563d8956dd29a8964dea1a79136db810))
* update dawidd6/action-download-artifact action to v2.26.1 ([#1189](https://github.com/keptn/lifecycle-toolkit/issues/1189)) ([1053717](https://github.com/keptn/lifecycle-toolkit/commit/10537174792d1b3514336452bca6dbba3ef49de5))
* update dawidd6/action-download-artifact action to v2.27.0 ([#1256](https://github.com/keptn/lifecycle-toolkit/issues/1256)) ([dc3e9b2](https://github.com/keptn/lifecycle-toolkit/commit/dc3e9b21e897d10e33a4ef406f357d44675cfacf))
* update dependency argoproj/argo-cd to v2.6.6 ([#1039](https://github.com/keptn/lifecycle-toolkit/issues/1039)) ([fb0f7a3](https://github.com/keptn/lifecycle-toolkit/commit/fb0f7a39bad1bbeda96210bd198d3e0ca0b6cb86))
* update dependency argoproj/argo-cd to v2.6.7 ([#1121](https://github.com/keptn/lifecycle-toolkit/issues/1121)) ([97c4b58](https://github.com/keptn/lifecycle-toolkit/commit/97c4b5823398310c11f50ff94d6c4cb4a29526a2))
* update dependency golangci/golangci-lint to v1.52.0 ([#1103](https://github.com/keptn/lifecycle-toolkit/issues/1103)) ([2b28b4f](https://github.com/keptn/lifecycle-toolkit/commit/2b28b4f58fb14421a704225cf437d08cba6b27b6))
* update dependency golangci/golangci-lint to v1.52.1 ([#1108](https://github.com/keptn/lifecycle-toolkit/issues/1108)) ([f5fb9ea](https://github.com/keptn/lifecycle-toolkit/commit/f5fb9ead9a53ad9fbd458df8c8b6aa6bfaa720da))
* update dependency golangci/golangci-lint to v1.52.2 ([#1142](https://github.com/keptn/lifecycle-toolkit/issues/1142)) ([1071f02](https://github.com/keptn/lifecycle-toolkit/commit/1071f0297b59bbd3a755e0a21c8e894bbfde1907))
* update dependency helm/helm to v3.11.2 ([#1050](https://github.com/keptn/lifecycle-toolkit/issues/1050)) ([2669e1d](https://github.com/keptn/lifecycle-toolkit/commit/2669e1d4760ed89797e312b7160d76d76f2171e8))
* update dependency helm/helm to v3.11.3 ([#1234](https://github.com/keptn/lifecycle-toolkit/issues/1234)) ([13c8fd8](https://github.com/keptn/lifecycle-toolkit/commit/13c8fd894b4603319591da52ef812214b00c782a))
* update dependency jaegertracing/jaeger to v1.43.0 ([#794](https://github.com/keptn/lifecycle-toolkit/issues/794)) ([abd4e09](https://github.com/keptn/lifecycle-toolkit/commit/abd4e0977fbc60638e32a19704580a667e0de282))
* update dependency jaegertracing/jaeger to v1.44.0 ([#1229](https://github.com/keptn/lifecycle-toolkit/issues/1229)) ([1257f0b](https://github.com/keptn/lifecycle-toolkit/commit/1257f0b8d4d83bffeca14139054b50e3a6b20324))
* update dependency jaegertracing/jaeger-operator to v1.43.0 ([#1152](https://github.com/keptn/lifecycle-toolkit/issues/1152)) ([9890213](https://github.com/keptn/lifecycle-toolkit/commit/9890213a643aac0977f677ac181d1d39e77fa5b5))
* update dependency kubernetes-sigs/controller-tools to v0.11.4 ([#1280](https://github.com/keptn/lifecycle-toolkit/issues/1280)) ([cfeec33](https://github.com/keptn/lifecycle-toolkit/commit/cfeec33a26207a21d95f21bcafd7e7aa55881b7c))
* update dependency kubernetes-sigs/kustomize to v5 ([#769](https://github.com/keptn/lifecycle-toolkit/issues/769)) ([33107ac](https://github.com/keptn/lifecycle-toolkit/commit/33107ac07a737053f60b449100c3d86b5cc910b7))
* update ghcr.io/podtato-head/entry docker tag to v0.2.8 ([#1211](https://github.com/keptn/lifecycle-toolkit/issues/1211)) ([d8f56b1](https://github.com/keptn/lifecycle-toolkit/commit/d8f56b12883a1904f7f6f600710a61d36e4753cd))
* update ghcr.io/podtato-head/hat docker tag to v0.2.8 ([#1212](https://github.com/keptn/lifecycle-toolkit/issues/1212)) ([ff09fbc](https://github.com/keptn/lifecycle-toolkit/commit/ff09fbc69c65209d26c04bc5418281fad5f438b6))
* update ghcr.io/podtato-head/left-arm docker tag to v0.2.8 ([#1217](https://github.com/keptn/lifecycle-toolkit/issues/1217)) ([549e76d](https://github.com/keptn/lifecycle-toolkit/commit/549e76d698eeccada0aae74c1c267b8c98f6b727))
* update ghcr.io/podtato-head/left-leg docker tag to v0.2.8 ([#1218](https://github.com/keptn/lifecycle-toolkit/issues/1218)) ([dd15d4a](https://github.com/keptn/lifecycle-toolkit/commit/dd15d4a0e0e4b4986d8c3d0860662bfa8ad1110e))
* update ghcr.io/podtato-head/right-arm docker tag to v0.2.8 ([#1219](https://github.com/keptn/lifecycle-toolkit/issues/1219)) ([48f6030](https://github.com/keptn/lifecycle-toolkit/commit/48f603067056d82c9c2fc5369872e5ed6572510d))
* update ghcr.io/podtato-head/right-leg docker tag to v0.2.8 ([#1220](https://github.com/keptn/lifecycle-toolkit/issues/1220)) ([3a4be7f](https://github.com/keptn/lifecycle-toolkit/commit/3a4be7f3048672994c5b4943844a2aa42f8954b1))
* update github.com/keptn/lifecycle-toolkit/metrics-operator digest to 6b58da3 ([#1141](https://github.com/keptn/lifecycle-toolkit/issues/1141)) ([3859059](https://github.com/keptn/lifecycle-toolkit/commit/385905913f867c88482e5b1a40643a8957a8e022))
* update github.com/keptn/lifecycle-toolkit/metrics-operator digest to 720e9e9 ([#1035](https://github.com/keptn/lifecycle-toolkit/issues/1035)) ([8a77f00](https://github.com/keptn/lifecycle-toolkit/commit/8a77f0004bc643a45890a888f3b7942f0c4ef794))
* update github.com/keptn/lifecycle-toolkit/metrics-operator digest to b32d753 ([#1164](https://github.com/keptn/lifecycle-toolkit/issues/1164)) ([4480444](https://github.com/keptn/lifecycle-toolkit/commit/448044446a315973cbf80fb4d9a24e58bb797eb5))
* update github.com/keptn/lifecycle-toolkit/metrics-operator digest to dd15d4a ([#1182](https://github.com/keptn/lifecycle-toolkit/issues/1182)) ([87b170f](https://github.com/keptn/lifecycle-toolkit/commit/87b170f4ac7a31812d235e1587d92789b9448da5))
* update github.com/keptn/lifecycle-toolkit/metrics-operator digest to f5fb9ea ([#1107](https://github.com/keptn/lifecycle-toolkit/issues/1107)) ([65f6a83](https://github.com/keptn/lifecycle-toolkit/commit/65f6a832a60a63226170b342d79ca187bddb6cec))
* update golang docker tag to v1.20.2 ([#1036](https://github.com/keptn/lifecycle-toolkit/issues/1036)) ([720e9e9](https://github.com/keptn/lifecycle-toolkit/commit/720e9e9b7040a1048d3f6fc86917bd678436c437))
* update golang docker tag to v1.20.3 ([#1183](https://github.com/keptn/lifecycle-toolkit/issues/1183)) ([f9a1bc7](https://github.com/keptn/lifecycle-toolkit/commit/f9a1bc7f2c4797d61d0afa75844052f945b15acc))
* update kubernetes packages (patch) ([#1228](https://github.com/keptn/lifecycle-toolkit/issues/1228)) ([ec1ece4](https://github.com/keptn/lifecycle-toolkit/commit/ec1ece41fcb797df3b02f2dd5fd53484b126e10e))
* update kubernetes packages to v0.26.3 (patch) ([#1072](https://github.com/keptn/lifecycle-toolkit/issues/1072)) ([a6459f8](https://github.com/keptn/lifecycle-toolkit/commit/a6459f8c2cd98fac9897ca40e299790ba35cd569))
* update module github.com/benbjohnson/clock to v1.3.1 ([#1257](https://github.com/keptn/lifecycle-toolkit/issues/1257)) ([e644597](https://github.com/keptn/lifecycle-toolkit/commit/e644597b17fb10fbae9d25a60b26487d96224841))
* update module github.com/benbjohnson/clock to v1.3.3 ([#1293](https://github.com/keptn/lifecycle-toolkit/issues/1293)) ([b7b2383](https://github.com/keptn/lifecycle-toolkit/commit/b7b23833faccd963c66a36f3c0b3a1af05d4d05c))
* update module github.com/datadog/datadog-api-client-go/v2 to v2.11.0 ([#1109](https://github.com/keptn/lifecycle-toolkit/issues/1109)) ([fbc021e](https://github.com/keptn/lifecycle-toolkit/commit/fbc021eadaf2442bfa67bc32947193b7fb6c20ae))
* update module github.com/datadog/datadog-api-client-go/v2 to v2.12.0 ([#1259](https://github.com/keptn/lifecycle-toolkit/issues/1259)) ([db347de](https://github.com/keptn/lifecycle-toolkit/commit/db347dec20302fae2688df613ca81526c95a0484))
* update module github.com/go-logr/logr to v1.2.4 ([#1153](https://github.com/keptn/lifecycle-toolkit/issues/1153)) ([c1ecfd0](https://github.com/keptn/lifecycle-toolkit/commit/c1ecfd03e39fc72f8304c78f80e71941aa0162e9))
* update module github.com/imdario/mergo to v0.3.14 ([#1073](https://github.com/keptn/lifecycle-toolkit/issues/1073)) ([ad408fd](https://github.com/keptn/lifecycle-toolkit/commit/ad408fd3b2dbd8480baad0c730b3f209c0d3b503))
* update module github.com/imdario/mergo to v0.3.15 ([#1132](https://github.com/keptn/lifecycle-toolkit/issues/1132)) ([17baf34](https://github.com/keptn/lifecycle-toolkit/commit/17baf3499a8681a850dfa7d88d049c26b86f8784))
* update module github.com/onsi/gomega to v1.27.4 ([#967](https://github.com/keptn/lifecycle-toolkit/issues/967)) ([502189a](https://github.com/keptn/lifecycle-toolkit/commit/502189ad73ed9abc9bedb89bd701258ac729024e))
* update module github.com/onsi/gomega to v1.27.5 ([#1133](https://github.com/keptn/lifecycle-toolkit/issues/1133)) ([7d0cf4b](https://github.com/keptn/lifecycle-toolkit/commit/7d0cf4b87d3c897ad348e9d1d31284710f99fbd0))
* update module github.com/onsi/gomega to v1.27.6 ([#1166](https://github.com/keptn/lifecycle-toolkit/issues/1166)) ([ab3a091](https://github.com/keptn/lifecycle-toolkit/commit/ab3a091748b4d3cdccad4854edfaf01798387376))
* update module github.com/prometheus/client_golang to v1.15.0 ([#1236](https://github.com/keptn/lifecycle-toolkit/issues/1236)) ([80b46c2](https://github.com/keptn/lifecycle-toolkit/commit/80b46c285c811feb126192471e8b7077246a3500))
* update module github.com/prometheus/common to v0.42.0 ([#1111](https://github.com/keptn/lifecycle-toolkit/issues/1111)) ([7ac89de](https://github.com/keptn/lifecycle-toolkit/commit/7ac89de96afd21e63c5d77465880f9428a822213))
* update module github.com/spf13/afero to v1.9.5 ([#1037](https://github.com/keptn/lifecycle-toolkit/issues/1037)) ([108e2a5](https://github.com/keptn/lifecycle-toolkit/commit/108e2a50764677dfe2bf6568c35bc2187ddcc206))
* update module github.com/stretchr/testify to v1.8.2 ([#937](https://github.com/keptn/lifecycle-toolkit/issues/937)) ([ddd3732](https://github.com/keptn/lifecycle-toolkit/commit/ddd3732423cffa1d11bf7b0cef86a8229c3216e2))
* update module golang.org/x/net to v0.9.0 ([#1298](https://github.com/keptn/lifecycle-toolkit/issues/1298)) ([ba7b679](https://github.com/keptn/lifecycle-toolkit/commit/ba7b679b0781de5558777fc93b8e9deb4ff6406a))
* update module google.golang.org/grpc to v1.53.0 ([#817](https://github.com/keptn/lifecycle-toolkit/issues/817)) ([f5a3493](https://github.com/keptn/lifecycle-toolkit/commit/f5a3493545f391112f341b5c54b6bdf442d8179b))
* update module google.golang.org/grpc to v1.54.0 ([#1112](https://github.com/keptn/lifecycle-toolkit/issues/1112)) ([ad2dc51](https://github.com/keptn/lifecycle-toolkit/commit/ad2dc511b9dcd3a8bc3bcd17e6344c9251a17a39))
* update module k8s.io/component-helpers to v0.25.9 ([#1235](https://github.com/keptn/lifecycle-toolkit/issues/1235)) ([16b9a2b](https://github.com/keptn/lifecycle-toolkit/commit/16b9a2baefc87c695028993160aca4e5bbf5145d))
* update module k8s.io/kubernetes to v1.25.8 ([#938](https://github.com/keptn/lifecycle-toolkit/issues/938)) ([65b854a](https://github.com/keptn/lifecycle-toolkit/commit/65b854ac9d02057e2bb99c270f629a44d67f258d))
* update module sigs.k8s.io/controller-runtime to v0.14.5 ([#1038](https://github.com/keptn/lifecycle-toolkit/issues/1038)) ([1be4f11](https://github.com/keptn/lifecycle-toolkit/commit/1be4f11872a634a037ed60cdf07ecf4a58c3b2c0))
* update module sigs.k8s.io/controller-runtime to v0.14.6 ([#1160](https://github.com/keptn/lifecycle-toolkit/issues/1160)) ([5f0071d](https://github.com/keptn/lifecycle-toolkit/commit/5f0071d114e28863192427e33ac5daa412418995))
* update peter-evans/create-pull-request action to v5 ([#1190](https://github.com/keptn/lifecycle-toolkit/issues/1190)) ([6c205b1](https://github.com/keptn/lifecycle-toolkit/commit/6c205b1b75ba6ba3591379244a87e5fc5eabc8a2))
* update sigstore/cosign-installer action to v3.0.2 ([#1198](https://github.com/keptn/lifecycle-toolkit/issues/1198)) ([31c657a](https://github.com/keptn/lifecycle-toolkit/commit/31c657afeac15de38a561d5a73b19d5013edc33c))


### Other

* adapt CODEOWNERS to new team structure ([#1250](https://github.com/keptn/lifecycle-toolkit/issues/1250)) ([0f11b85](https://github.com/keptn/lifecycle-toolkit/commit/0f11b85474e54a1d7811bb5df90498e40e13ecda))
* bump go to 1.20 ([#1294](https://github.com/keptn/lifecycle-toolkit/issues/1294)) ([0a6ac23](https://github.com/keptn/lifecycle-toolkit/commit/0a6ac23eb77e2b5f8bdd20028309254dda2c9d1d))
* bump GO_VERSION to 1.20 in pipelines ([#1326](https://github.com/keptn/lifecycle-toolkit/issues/1326)) ([7e8079e](https://github.com/keptn/lifecycle-toolkit/commit/7e8079ecfcd612bc31f581e66f76e4c895283efc))
* **cert-manager:** reduce secret permissions ([#1295](https://github.com/keptn/lifecycle-toolkit/issues/1295)) ([bd8de3b](https://github.com/keptn/lifecycle-toolkit/commit/bd8de3b6461fcd599b58461ffbf42ff2e087951e))
* fix failing component test ([#1282](https://github.com/keptn/lifecycle-toolkit/issues/1282)) ([00fd1f3](https://github.com/keptn/lifecycle-toolkit/commit/00fd1f3f2f6bd8e4d547d7e6b14cbd9fe9e14d42))
* improve CRD docs generation script output ([#1157](https://github.com/keptn/lifecycle-toolkit/issues/1157)) ([b27adf1](https://github.com/keptn/lifecycle-toolkit/commit/b27adf1edc297dd7723998fcf4960929b2b2952d))
* **metrics-operator:** add configuration parameters for container securityContext ([#1290](https://github.com/keptn/lifecycle-toolkit/issues/1290)) ([27439ff](https://github.com/keptn/lifecycle-toolkit/commit/27439ff3ad9fff30341fb987fd06ecb3aaef0d1d))
* **metrics-operator:** restrict custom metrics ClusterRole privileges ([#1330](https://github.com/keptn/lifecycle-toolkit/issues/1330)) ([6f59a6c](https://github.com/keptn/lifecycle-toolkit/commit/6f59a6c0c75d79a54a874fcda64181180723551b))
* **operator:** read-only RBAC for KeptnConfig controller ([#1096](https://github.com/keptn/lifecycle-toolkit/issues/1096)) ([ea91ff3](https://github.com/keptn/lifecycle-toolkit/commit/ea91ff36dfbe13811143031462407fecf7791596))
* refactor and add unit tests to watcher ([#1253](https://github.com/keptn/lifecycle-toolkit/issues/1253)) ([4b40b7e](https://github.com/keptn/lifecycle-toolkit/commit/4b40b7ecccecc3bccaccc0532be1cd16d9c7ba6a))
* remove cert-manager leftovers ([#1216](https://github.com/keptn/lifecycle-toolkit/issues/1216)) ([1c58ba8](https://github.com/keptn/lifecycle-toolkit/commit/1c58ba8edc2320d58096dea1f8beecfdc8c949b5))
* reorder integration test execution ([#1264](https://github.com/keptn/lifecycle-toolkit/issues/1264)) ([71f2f78](https://github.com/keptn/lifecycle-toolkit/commit/71f2f787072627c807c2802bc57d016ee5b51d1b))
* revert test makefile changes ([#1281](https://github.com/keptn/lifecycle-toolkit/issues/1281)) ([2261a4a](https://github.com/keptn/lifecycle-toolkit/commit/2261a4ab055095ad6d92b991e5520ff76ad9ba86))
* set up YAML linter rules, fix YAML files accordingly ([#1174](https://github.com/keptn/lifecycle-toolkit/issues/1174)) ([86fbb75](https://github.com/keptn/lifecycle-toolkit/commit/86fbb757e55cb959e8ccb2ddc62bfcba55271452))
* stop pushing dev container images to GHCR ([#1192](https://github.com/keptn/lifecycle-toolkit/issues/1192)) ([fa53443](https://github.com/keptn/lifecycle-toolkit/commit/fa53443c704193db029c92ff4ad08fc2bfaaa24d))


### Docs

* add better overview KeptnApp to readme ([#1254](https://github.com/keptn/lifecycle-toolkit/issues/1254)) ([497e57e](https://github.com/keptn/lifecycle-toolkit/commit/497e57e6ddeac0735871367042e96bca06b73356))
* add community files to webpage ([#1077](https://github.com/keptn/lifecycle-toolkit/issues/1077)) ([ed3836a](https://github.com/keptn/lifecycle-toolkit/commit/ed3836aa6d75b186b09cb7c8ecf7c049f58af999))
* add metrics-operator architecture ([#1151](https://github.com/keptn/lifecycle-toolkit/issues/1151)) ([80d0045](https://github.com/keptn/lifecycle-toolkit/commit/80d0045daccc292cad745aebbc2feba0c1e55cbd))
* added example for autoscaling using KeptnMetric ([#1173](https://github.com/keptn/lifecycle-toolkit/issues/1173)) ([98dd248](https://github.com/keptn/lifecycle-toolkit/commit/98dd248f1fd22964fbae80a1de113f76fff3e55e))
* adding KLT runtime info for local development ([#1246](https://github.com/keptn/lifecycle-toolkit/issues/1246)) ([c8131b6](https://github.com/keptn/lifecycle-toolkit/commit/c8131b63b55fc11085a7ab83b555b077656d9a7a))
* change linting CLI and add custom rules ([#1031](https://github.com/keptn/lifecycle-toolkit/issues/1031)) ([acf5f91](https://github.com/keptn/lifecycle-toolkit/commit/acf5f91b2bc600d0008e64ffe602da8147134330))
* cleanup after theme migration ([#1045](https://github.com/keptn/lifecycle-toolkit/issues/1045)) ([0125462](https://github.com/keptn/lifecycle-toolkit/commit/01254620b2735605d103f418e14487c43e6d3a1e))
* describe automatic application discovery ([#1304](https://github.com/keptn/lifecycle-toolkit/issues/1304)) ([d576a33](https://github.com/keptn/lifecycle-toolkit/commit/d576a3323ff32a28ff1a8a207d811a01e1865d08))
* fix missing code fence ([#1343](https://github.com/keptn/lifecycle-toolkit/issues/1343)) ([2576a98](https://github.com/keptn/lifecycle-toolkit/commit/2576a98cc252e66fc8e3ea3c0532e363f173fc61))
* fix typo ([#1252](https://github.com/keptn/lifecycle-toolkit/issues/1252)) ([4a96b06](https://github.com/keptn/lifecycle-toolkit/commit/4a96b0637f01d597fb0ff2631b5784267386862e))
* fix typo in the getting started docs ([#1204](https://github.com/keptn/lifecycle-toolkit/issues/1204)) ([c9b1a42](https://github.com/keptn/lifecycle-toolkit/commit/c9b1a42df8a8da1e158fd7a0c1d2a92726b7ca08))
* improve docs for KeptnEvaluationDefinition ([#1335](https://github.com/keptn/lifecycle-toolkit/issues/1335)) ([d9e0aac](https://github.com/keptn/lifecycle-toolkit/commit/d9e0aac74c4b8fc4a3debd3f08d22155190d9f5f))
* improve headline of Getting Started subsection ([#1350](https://github.com/keptn/lifecycle-toolkit/issues/1350)) ([a3ef431](https://github.com/keptn/lifecycle-toolkit/commit/a3ef431c5c41978a85df97828f9f297648f2bfdc))
* improve landing page, based on slides ([#1272](https://github.com/keptn/lifecycle-toolkit/issues/1272)) ([117cda4](https://github.com/keptn/lifecycle-toolkit/commit/117cda489d60a5056e1341ec6f2ee4daf2172dec))
* improve Notes within the documentation ([#962](https://github.com/keptn/lifecycle-toolkit/issues/962)) ([4e69699](https://github.com/keptn/lifecycle-toolkit/commit/4e69699d12d82cac1c8516c42dee8e131d993b4c))
* improve rendering of links for local markdown files ([#1177](https://github.com/keptn/lifecycle-toolkit/issues/1177)) ([070bbee](https://github.com/keptn/lifecycle-toolkit/commit/070bbee37ee44be3f8f2c7f03078c3276ee317e3))
* make cert-manager and manifest installation a detail ([#1099](https://github.com/keptn/lifecycle-toolkit/issues/1099)) ([66b3f01](https://github.com/keptn/lifecycle-toolkit/commit/66b3f013bba551be7471088a53b117242a34b543))
* modify footer ([#1163](https://github.com/keptn/lifecycle-toolkit/issues/1163)) ([ee4ffcf](https://github.com/keptn/lifecycle-toolkit/commit/ee4ffcfd13469f3fb733e44319508b0c91bd765d))
* **operator:** adjust docs comment ([#1126](https://github.com/keptn/lifecycle-toolkit/issues/1126)) ([4078fad](https://github.com/keptn/lifecycle-toolkit/commit/4078fada7a1c8a5a392aa11fe0f161acf190478d))
* replace cert-manager ([#1210](https://github.com/keptn/lifecycle-toolkit/issues/1210)) ([a84cbc7](https://github.com/keptn/lifecycle-toolkit/commit/a84cbc7d67e8608f4c6638ab8d8efc1424f46933))
* set up get-started directory tree ([#1303](https://github.com/keptn/lifecycle-toolkit/issues/1303)) ([57b6574](https://github.com/keptn/lifecycle-toolkit/commit/57b6574c76d3b5de1c7e361e5c4e64c449894d7a))
* set up structure for yaml ref pages, guide section ([#1184](https://github.com/keptn/lifecycle-toolkit/issues/1184)) ([c164595](https://github.com/keptn/lifecycle-toolkit/commit/c164595d7682257e6c78f219cb9bc1e2b0f2fb82))
* set up top-level Installation section ([#1162](https://github.com/keptn/lifecycle-toolkit/issues/1162)) ([2c62593](https://github.com/keptn/lifecycle-toolkit/commit/2c6259366a5ff545404dad2c692b7f1db524de83))
* technologies to get familiar before working with KLT ([#1060](https://github.com/keptn/lifecycle-toolkit/issues/1060)) ([58e8a4c](https://github.com/keptn/lifecycle-toolkit/commit/58e8a4cac7b02a1b3fc37b5b448cb5df45a4c484))
* update API reference docs pages ([#1273](https://github.com/keptn/lifecycle-toolkit/issues/1273)) ([706292a](https://github.com/keptn/lifecycle-toolkit/commit/706292aeb3d8280a7eea37ec089e5b7e83e6076e))
* update docs for multi metrics provider support, fix API reference generator ([#1251](https://github.com/keptn/lifecycle-toolkit/issues/1251)) ([1dfd653](https://github.com/keptn/lifecycle-toolkit/commit/1dfd653eba98056cbee207d32ef5aa5b567bfb10))
* update KeptnConfig docs to include KeptnAppCreationRequestTimeout ([#1348](https://github.com/keptn/lifecycle-toolkit/issues/1348)) ([117c263](https://github.com/keptn/lifecycle-toolkit/commit/117c263f0e4da43f32c5b3603c792de9a1badf66))
* update KeptnTaskDefinition to include fallback search to default KLT namespace ([#1349](https://github.com/keptn/lifecycle-toolkit/issues/1349)) ([2f5587e](https://github.com/keptn/lifecycle-toolkit/commit/2f5587ed2407abe053ed6d8fb49ef7a7c1123eb4))
* update list of videos about KLT ([#1105](https://github.com/keptn/lifecycle-toolkit/issues/1105)) ([ade49e1](https://github.com/keptn/lifecycle-toolkit/commit/ade49e1fd8d4cb9b811fb4b4be871c1907f124ad))

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

* Evaluation and Task statuses in KeptnWorkloadInstance/KeptnAppVersion use the same structure
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

* **dashboards:** use fixed color mode for succeeded AppVersion/WorkloadInstance tiles ([#515](https://github.com/keptn/lifecycle-toolkit/issues/515)) ([8cdb23e](https://github.com/keptn/lifecycle-toolkit/commit/8cdb23ee61cc7ee22be2b7326bbf202ee3ddf09f))


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
* **operator:** Add information about current phase in workloadinstances and appversions ([#200](https://github.com/keptn/lifecycle-toolkit/issues/200)) ([55fa4e9](https://github.com/keptn/lifecycle-toolkit/commit/55fa4e97c62aec7bd1a45f85d47cfaca48f3dd8f))
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

* **operator:** Modified behavior of KeptnAppVersion and KeptnWorkloadInstance to support pre and post deployment evaluation checks with Prometheus montoring
* **operator:** now the namespaces have to be annotated/labeled with keptn.sh/lifecycle-toolkit=enabled when the lifecycle controller should be used
* **operator:** Implementation of the KeptnApp CRD and Controller. This modifies the behaviour of the KeptnWorkloadInstance and Keptn MutatingWebhook

### Features

* Namespace keptn-lifecycle-toolkit-system should never call webhook ([#192](https://github.com/keptn/lifecycle-toolkit/issues/192)) ([913a9ff](https://github.com/keptn/lifecycle-toolkit/commit/913a9ffd62f93aa7831b35e29853afff6213a0c9))
* **operator:** add fallback behavior when no keptn annotations are set ([#171](https://github.com/keptn/lifecycle-toolkit/issues/171)) ([b6cc674](https://github.com/keptn/lifecycle-toolkit/commit/b6cc674adb787615fc79dbbc5b10668c367e4736))
* **operator:** Add KeptnApplication controller ([#137](https://github.com/keptn/lifecycle-toolkit/issues/137)) ([271f5a8](https://github.com/keptn/lifecycle-toolkit/commit/271f5a830f216c9f827457d8a391c25d56aed2e3))
* **operator:** Added minimal context information ([#170](https://github.com/keptn/lifecycle-toolkit/issues/170)) ([eebe420](https://github.com/keptn/lifecycle-toolkit/commit/eebe4200aac74a7c2cbc73720d1d9ac6a0c1fc72))
* **operator:** Allow pre- and post-deployment tasks as labels or annotations ([#181](https://github.com/keptn/lifecycle-toolkit/issues/181)) ([4241fe7](https://github.com/keptn/lifecycle-toolkit/commit/4241fe7cfab91aa6d38309eacf5712436a6e8327))
* **operator:** Bootstrap evaluation CRD from app ([#184](https://github.com/keptn/lifecycle-toolkit/issues/184)) ([74c3dbc](https://github.com/keptn/lifecycle-toolkit/commit/74c3dbc7b6d78d8ca7eafbac50abb8c3473701eb))
* **operator:** Bootstrap evaluation CRD from WorkloadInstance ([#188](https://github.com/keptn/lifecycle-toolkit/issues/188)) ([95e206b](https://github.com/keptn/lifecycle-toolkit/commit/95e206b4165b0277f5acbc67fc78a8e28f06741b))
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
* **operator:** Introduce OTel tracing for WorkloadInstance controller ([#131](https://github.com/keptn/lifecycle-toolkit/issues/131)) ([a195614](https://github.com/keptn/lifecycle-toolkit/commit/a1956141fe80e5b1afd79fb33198313e1dbff7fa))
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
