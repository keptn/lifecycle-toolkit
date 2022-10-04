# Changelog

## [0.1.5](https://github.com/mowies/lifecycle-controller/compare/v0.1.4...v0.1.5) (2022-10-04)


### Features

* a new feature ([d5093f4](https://github.com/mowies/lifecycle-controller/commit/d5093f4ba1d36c64d729bade416e1b1ce9e33d84))

## [0.1.4](https://github.com/mowies/lifecycle-controller/compare/v0.1.3...v0.1.4) (2022-10-04)


### Bug Fixes

* Add release write permissions ([754d1f7](https://github.com/mowies/lifecycle-controller/commit/754d1f75be2743741445890e9ec797d85a343bd5))

## [0.1.3](https://github.com/mowies/lifecycle-controller/compare/v0.1.2...v0.1.3) (2022-10-04)


### Bug Fixes

* Add missing login setup for GHCR ([15237eb](https://github.com/mowies/lifecycle-controller/commit/15237eb357d2ce356b1eaf1ae443c056bd3377b4))

## [0.1.2](https://github.com/mowies/lifecycle-controller/compare/v0.1.1...v0.1.2) (2022-10-04)


### Bug Fixes

* Add correct github token permissions in release pipeline ([5d7a333](https://github.com/mowies/lifecycle-controller/commit/5d7a3338138795adb51e99db199fde59b5772ce1))

## [0.1.1](https://github.com/mowies/lifecycle-controller/compare/v0.1.0...v0.1.1) (2022-10-04)


### Features

* Use personal GHCR ([3810946](https://github.com/mowies/lifecycle-controller/commit/38109465109eb6fea8b4de607af84ff41c380119))

## [0.1.0](https://github.com/mowies/lifecycle-controller/compare/v0.1.0...v0.1.0) (2022-10-04)


### Features

* Add scheduler with annotations ([#31](https://github.com/mowies/lifecycle-controller/issues/31)) ([9e29019](https://github.com/mowies/lifecycle-controller/commit/9e29019c098fd4f1d5e36500bd2c7ef410421aa8))
* Bootstrap Service CR and controller ([#21](https://github.com/mowies/lifecycle-controller/issues/21)) ([c714ecc](https://github.com/mowies/lifecycle-controller/commit/c714eccc3b9c4d1309036fc9d193da3154b4cac5))
* First draft of a scheduler ([#19](https://github.com/mowies/lifecycle-controller/issues/19)) ([1884c86](https://github.com/mowies/lifecycle-controller/commit/1884c8678a681ed322a0ef2ea07fad3e24e01237))
* first podtatohead sample deployment manifests ([#45](https://github.com/mowies/lifecycle-controller/issues/45)) ([3e92d27](https://github.com/mowies/lifecycle-controller/commit/3e92d277ebf1a9063ebcf80f05ebe62958e45cbb))
* First Version of Function Execution ([#35](https://github.com/mowies/lifecycle-controller/issues/35)) ([f6badfd](https://github.com/mowies/lifecycle-controller/commit/f6badfd19f9f0b15c04364be7b03f524c920a015))
* initial version of function runtime ([#26](https://github.com/mowies/lifecycle-controller/issues/26)) ([c8800ee](https://github.com/mowies/lifecycle-controller/commit/c8800ee352b5d0d5eccd7338cd4fa6a3ae7d2efa))
* Inject keptn-scheduler when resource contains Keptn annotations ([#18](https://github.com/mowies/lifecycle-controller/issues/18)) ([4530e86](https://github.com/mowies/lifecycle-controller/commit/4530e8602beb4fc923b767eb586e44752f725400))
* **lfc-scheduler:** Move from Helm to Kustomize ([#53](https://github.com/mowies/lifecycle-controller/issues/53)) ([d7ba5f3](https://github.com/mowies/lifecycle-controller/commit/d7ba5f35f1b32451f833d9fd53079b4162837bde))
* sample function for deno runtime ([#27](https://github.com/mowies/lifecycle-controller/issues/27)) ([2501e46](https://github.com/mowies/lifecycle-controller/commit/2501e46a18dfc4ab436669fa7c42c570abad5a52))
* substitute event task ([#43](https://github.com/mowies/lifecycle-controller/issues/43)) ([3644a7d](https://github.com/mowies/lifecycle-controller/commit/3644a7d9a0d4a565a9d857348a63ed91d8cb8102))
* Switch to distroless-base image ([#46](https://github.com/mowies/lifecycle-controller/issues/46)) ([0a735b2](https://github.com/mowies/lifecycle-controller/commit/0a735b2ca22a02ca42faf7d091741d39e0f5a547))
* Webhook creates Service, Service creates ServiceRun, ServiceRun creates Event ([#30](https://github.com/mowies/lifecycle-controller/issues/30)) ([5ae58c3](https://github.com/mowies/lifecycle-controller/commit/5ae58c33abe965e79bb405e74c0f308f1220d4ee))


### Bug Fixes

* Added namespace to task definition for podtato head example ([#72](https://github.com/mowies/lifecycle-controller/issues/72)) ([7081f27](https://github.com/mowies/lifecycle-controller/commit/7081f2772aee5abec840a58c7ab700603e84cf52))
* Fix CODEOWNERS syntax ([0be5197](https://github.com/mowies/lifecycle-controller/commit/0be5197c19ea3066d28fe8e97f274efff21f66ff))
* fixed namespace in scheduler kustomization ([#63](https://github.com/mowies/lifecycle-controller/issues/63)) ([237bf4f](https://github.com/mowies/lifecycle-controller/commit/237bf4f480161f48aa0c4b5f2afbff433447d2a8))
* Missed error ([#76](https://github.com/mowies/lifecycle-controller/issues/76)) ([a59aa15](https://github.com/mowies/lifecycle-controller/commit/a59aa1552795bce15e39195af235fd42d1448e61))
* query jobs before creating ([#79](https://github.com/mowies/lifecycle-controller/issues/79)) ([47f82b8](https://github.com/mowies/lifecycle-controller/commit/47f82b891d9d20ade2928faae307009e5c96ae22))
* scheduler config plugin configuration ([#68](https://github.com/mowies/lifecycle-controller/issues/68)) ([4c4e3c6](https://github.com/mowies/lifecycle-controller/commit/4c4e3c60a0e11267dc69ea7d8470555e3ee4f91e))


### Miscellaneous Chores

* release 0.1.0 ([4c46a42](https://github.com/mowies/lifecycle-controller/commit/4c46a4297c540b9da30c5a373624d4b8e8a88231))
* release 0.1.0 ([afa8493](https://github.com/mowies/lifecycle-controller/commit/afa849324fa422352ed61faa7f0dc75d74c3c25d))

## Changelog
