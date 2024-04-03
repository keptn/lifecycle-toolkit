# Changelog

## [2.0.3](https://github.com/keptn/lifecycle-toolkit/compare/deno-runtime-v2.0.2...deno-runtime-v2.0.3) (2024-03-19)


### Dependency Updates

* update denoland/deno docker tag to alpine-1.41.1 ([#3165](https://github.com/keptn/lifecycle-toolkit/issues/3165)) ([036a1d4](https://github.com/keptn/lifecycle-toolkit/commit/036a1d45e1b851bd9b2a55648f95c47367638d30))
* update denoland/deno docker tag to alpine-1.41.2 ([#3262](https://github.com/keptn/lifecycle-toolkit/issues/3262)) ([53a32b2](https://github.com/keptn/lifecycle-toolkit/commit/53a32b2d65cacec787ffe653b81df68d87fd70d4))
* update denoland/deno docker tag to alpine-1.41.3 ([#3279](https://github.com/keptn/lifecycle-toolkit/issues/3279)) ([8a1eebf](https://github.com/keptn/lifecycle-toolkit/commit/8a1eebf2242118f3e582605494b0fd917641dda6))

## [2.0.2](https://github.com/keptn/lifecycle-toolkit/compare/deno-runtime-v2.0.1...deno-runtime-v2.0.2) (2024-02-29)


### Dependency Updates

* update denoland/deno docker tag to alpine-1.40.3 ([#2943](https://github.com/keptn/lifecycle-toolkit/issues/2943)) ([fde53ce](https://github.com/keptn/lifecycle-toolkit/commit/fde53ce523438fc6b040d9df7951aa1ec04a82f4))
* update denoland/deno docker tag to alpine-1.40.4 ([#3031](https://github.com/keptn/lifecycle-toolkit/issues/3031)) ([839c61e](https://github.com/keptn/lifecycle-toolkit/commit/839c61ec6a34b9b2a44f65cc2f0231c38f1d6f30))
* update denoland/deno docker tag to alpine-1.40.5 ([#3060](https://github.com/keptn/lifecycle-toolkit/issues/3060)) ([4b25727](https://github.com/keptn/lifecycle-toolkit/commit/4b25727e1a75e99f3b3b709aeeb3c49e5845c0b3))
* update denoland/deno docker tag to alpine-1.41.0 ([#3126](https://github.com/keptn/lifecycle-toolkit/issues/3126)) ([a60e4c6](https://github.com/keptn/lifecycle-toolkit/commit/a60e4c6da63f36ac06cbdaf4994bbae8c062ac18))

## [2.0.1](https://github.com/keptn/lifecycle-toolkit/compare/deno-runtime-v2.0.0...deno-runtime-v2.0.1) (2024-02-06)


### Other

* revert "deps(deno-runtime): update libcrypto3 and libssl3" ([7e18270](https://github.com/keptn/lifecycle-toolkit/commit/7e1827088848dc486afb007c354155d2f9a5ed5c))

## [2.0.0](https://github.com/keptn/lifecycle-toolkit/compare/deno-runtime-v1.0.2...deno-runtime-v2.0.0) (2024-02-06)


### ⚠ BREAKING CHANGES

* **lifecycle-operator:** The environment variable giving deno and python runtime access to context information has been renamed from `CONTEXT` to `KEPTN_CONTEXT`

### Features

* **lifecycle-operator:** rename CONTEXT to KEPTN_CONTEXT in task runtimes ([#2521](https://github.com/keptn/lifecycle-toolkit/issues/2521)) ([a7322bd](https://github.com/keptn/lifecycle-toolkit/commit/a7322bd9266fa1589d77b06675d70d1a9e6c29ac))


### Other

* add config for spell checker action, fix typos ([#2443](https://github.com/keptn/lifecycle-toolkit/issues/2443)) ([eac178f](https://github.com/keptn/lifecycle-toolkit/commit/eac178f650962208449553086d54d26d27fa4da3))
* **deno-runtime:** add read/write permissions to deno runtime image ([#2618](https://github.com/keptn/lifecycle-toolkit/issues/2618)) ([8425f50](https://github.com/keptn/lifecycle-toolkit/commit/8425f50bf745282e78ed0bba7300810d82bd84c9))
* rename Keptn default namespace to 'keptn-system' ([#2565](https://github.com/keptn/lifecycle-toolkit/issues/2565)) ([aec1148](https://github.com/keptn/lifecycle-toolkit/commit/aec11489451ab1b0bcd69a6b90b0d45f69c5df7c))


### Docs

* mention `KEPTN_CONTEXT` env var in runtime readmes files ([#2588](https://github.com/keptn/lifecycle-toolkit/issues/2588)) ([dfefc90](https://github.com/keptn/lifecycle-toolkit/commit/dfefc90e9e5075ef130e3962b1ded983b2b213f4))
* remove Scarf transparent pixels ([#2590](https://github.com/keptn/lifecycle-toolkit/issues/2590)) ([95851fa](https://github.com/keptn/lifecycle-toolkit/commit/95851fa52cb3a6565a4b52ae0e8b00dcc9861a3b))


### Dependency Updates

* **deno-runtime:** update libcrypto3 and libssl3 ([#2953](https://github.com/keptn/lifecycle-toolkit/issues/2953)) ([882b442](https://github.com/keptn/lifecycle-toolkit/commit/882b44222fee306704674a91875ffdf1ccc7a3af))
* update denoland/deno docker tag to alpine-1.38.0 ([#2413](https://github.com/keptn/lifecycle-toolkit/issues/2413)) ([14f3cc3](https://github.com/keptn/lifecycle-toolkit/commit/14f3cc3b191403b759909f7433596aa3713a095e))
* update denoland/deno docker tag to alpine-1.38.1 ([#2474](https://github.com/keptn/lifecycle-toolkit/issues/2474)) ([04248eb](https://github.com/keptn/lifecycle-toolkit/commit/04248eb7a7d7307334733834d680670e6f006b53))
* update denoland/deno docker tag to alpine-1.38.2 ([#2528](https://github.com/keptn/lifecycle-toolkit/issues/2528)) ([964fd25](https://github.com/keptn/lifecycle-toolkit/commit/964fd259c5e56652d87af32b83cceec7cae77cac))
* update denoland/deno docker tag to alpine-1.38.3 ([#2567](https://github.com/keptn/lifecycle-toolkit/issues/2567)) ([f1b969e](https://github.com/keptn/lifecycle-toolkit/commit/f1b969e66b64c03dfe70ae4960daa78246e26a61))
* update denoland/deno docker tag to alpine-1.38.4 ([#2625](https://github.com/keptn/lifecycle-toolkit/issues/2625)) ([32220d9](https://github.com/keptn/lifecycle-toolkit/commit/32220d90fd6a4bc645965c8ebc26f8580c643a16))
* update denoland/deno docker tag to alpine-1.38.5 ([#2648](https://github.com/keptn/lifecycle-toolkit/issues/2648)) ([a3d77a3](https://github.com/keptn/lifecycle-toolkit/commit/a3d77a377dc7869655b901c697ee62c3e7ae86e3))
* update denoland/deno docker tag to alpine-1.39.0 ([#2684](https://github.com/keptn/lifecycle-toolkit/issues/2684)) ([5456786](https://github.com/keptn/lifecycle-toolkit/commit/5456786ed7128766016b065e262fad1a4a4f3ee6))
* update denoland/deno docker tag to alpine-1.39.2 ([#2721](https://github.com/keptn/lifecycle-toolkit/issues/2721)) ([8e1e088](https://github.com/keptn/lifecycle-toolkit/commit/8e1e088ba31de2444a306be5750fa5972e224b28))
* update denoland/deno docker tag to alpine-1.39.4 ([#2795](https://github.com/keptn/lifecycle-toolkit/issues/2795)) ([5a5d5f4](https://github.com/keptn/lifecycle-toolkit/commit/5a5d5f4fc05e5f10f069db71137751a59febcb20))
* update denoland/deno docker tag to alpine-1.40.2 ([#2886](https://github.com/keptn/lifecycle-toolkit/issues/2886)) ([cad0c57](https://github.com/keptn/lifecycle-toolkit/commit/cad0c578950fc6dbd0102764508a05d2a86d749f))

## [1.0.2](https://github.com/keptn/lifecycle-toolkit/compare/deno-runtime-v1.0.1...deno-runtime-v1.0.2) (2023-10-30)


### Docs

* implement KLT -&gt; Keptn name change ([#2001](https://github.com/keptn/lifecycle-toolkit/issues/2001)) ([440c308](https://github.com/keptn/lifecycle-toolkit/commit/440c3082e5400f89d791724651984ba2bc0a4724))


### Dependency Updates

* update denoland/deno docker tag to alpine-1.36.4 ([#2012](https://github.com/keptn/lifecycle-toolkit/issues/2012)) ([f2f3162](https://github.com/keptn/lifecycle-toolkit/commit/f2f316271d86209da124ea3554fa2e821d79e953))
* update denoland/deno docker tag to alpine-1.37.0 ([#2157](https://github.com/keptn/lifecycle-toolkit/issues/2157)) ([0f863d0](https://github.com/keptn/lifecycle-toolkit/commit/0f863d03c46a16ee7e105335ae610f3c4776d4f8))
* update denoland/deno docker tag to alpine-1.37.1 ([#2218](https://github.com/keptn/lifecycle-toolkit/issues/2218)) ([21652a8](https://github.com/keptn/lifecycle-toolkit/commit/21652a8bf5a10eae55d4c1fd81e270ee581eb4a1))
* update denoland/deno docker tag to alpine-1.37.2 ([#2280](https://github.com/keptn/lifecycle-toolkit/issues/2280)) ([b6f5c18](https://github.com/keptn/lifecycle-toolkit/commit/b6f5c184bf5dccc26003e63ba11edce80f10eb66))
* update dependency autoprefixer to v10.4.15 ([#1909](https://github.com/keptn/lifecycle-toolkit/issues/1909)) ([8dbec2d](https://github.com/keptn/lifecycle-toolkit/commit/8dbec2d6116fb20bac86162aaea2b75c24eb96be))

## [1.0.1](https://github.com/keptn/lifecycle-toolkit/compare/deno-runtime-v1.0.0...deno-runtime-v1.0.1) (2023-08-30)


### Dependency Updates

* update denoland/deno docker tag to alpine-1.36.3 ([#1944](https://github.com/keptn/lifecycle-toolkit/issues/1944)) ([da95b40](https://github.com/keptn/lifecycle-toolkit/commit/da95b4025775b399084b2937b4ea0c0c360ec86c))

## 1.0.0 (2023-08-29)


### Features

* monorepo setup for lifecycle-operator, scheduler and runtimes ([#1857](https://github.com/keptn/lifecycle-toolkit/issues/1857)) ([84e243a](https://github.com/keptn/lifecycle-toolkit/commit/84e243a213ffba86eddd51ccc4bf4dbd61140069))
