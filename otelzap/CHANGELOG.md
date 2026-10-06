# Changelog

## [0.2.3](https://github.com/SpechtLabs/go-otel-utils/compare/otelzap/v0.2.2...otelzap/v0.2.3) (2026-10-06)


### Bug Fixes

* **deps:** break the go 1.27.1 cycle between otelprovider and otelzap ([#59](https://github.com/SpechtLabs/go-otel-utils/issues/59)) ([b285ba3](https://github.com/SpechtLabs/go-otel-utils/commit/b285ba3513abef87692867e45e81401fb8ecf47a))
* **deps:** declare go 1.27 instead of 1.27.1 ([#58](https://github.com/SpechtLabs/go-otel-utils/issues/58)) ([da194c6](https://github.com/SpechtLabs/go-otel-utils/commit/da194c600517ba33b37e7f51248de84b251bf7d2))
* **otelzap:** stop With from adding fields to the shared logger ([#57](https://github.com/SpechtLabs/go-otel-utils/issues/57)) ([a2e25a9](https://github.com/SpechtLabs/go-otel-utils/commit/a2e25a9910027eb38d30498aefd74963666fc486))

## [0.2.2](https://github.com/SpechtLabs/go-otel-utils/compare/otelzap/v0.2.1...otelzap/v0.2.2) (2026-10-05)


### Bug Fixes

* **otelprovider:** move to the stable otel/log v1.47 stack ([#54](https://github.com/SpechtLabs/go-otel-utils/issues/54)) ([3a97ce7](https://github.com/SpechtLabs/go-otel-utils/commit/3a97ce7b91cd9719d4c7889d89892a352ce7e474))

## [0.2.1](https://github.com/SpechtLabs/go-otel-utils/compare/otelzap/v0.2.0...otelzap/v0.2.1) (2026-10-05)


### Bug Fixes

* **otelzap:** stop Attribute and LogValue from panicking on arrays and named slices ([#52](https://github.com/SpechtLabs/go-otel-utils/issues/52)) ([8aa8ae7](https://github.com/SpechtLabs/go-otel-utils/commit/8aa8ae7e28fe4aa9dca9c18f904f726f57fedb24))

## [0.2.0](https://github.com/SpechtLabs/go-otel-utils/compare/otelzap/v0.1.1...otelzap/v0.2.0) (2026-10-05)


### ⚠ BREAKING CHANGES

* **otelzap:** build against the stable otel/log v1.47 API ([#49](https://github.com/SpechtLabs/go-otel-utils/issues/49))

### Bug Fixes

* **otelzap:** build against the stable otel/log v1.47 API ([#49](https://github.com/SpechtLabs/go-otel-utils/issues/49)) ([8ab5ab3](https://github.com/SpechtLabs/go-otel-utils/commit/8ab5ab397c8d645fce06692f07080ff923aa3dbd))

## [0.1.1](https://github.com/SpechtLabs/go-otel-utils/compare/otelzap/v0.1.0...otelzap/v0.1.1) (2026-06-11)


### Bug Fixes

* **deps:** bump otel to 1.44, zap 1.28, smithy 1.27 and inter-module deps to v0.1.0 ([5ca776f](https://github.com/SpechtLabs/go-otel-utils/commit/5ca776f807b57c8fa4133c6c353978b49906d5ab))

## [0.0.16](https://github.com/SpechtLabs/go-otel-utils/compare/otelzap/v0.0.15...otelzap/v0.0.16) (2026-06-11)


### Bug Fixes

* **deps:** update github.com/sierrasoftworks/humane-errors-go digest to 2224f06 ([#17](https://github.com/SpechtLabs/go-otel-utils/issues/17)) ([9283540](https://github.com/SpechtLabs/go-otel-utils/commit/92835409e9305ddbea2c010e3b4ef48a11819f40))
* **deps:** update module github.com/aws/smithy-go to v1.23.0 ([#14](https://github.com/SpechtLabs/go-otel-utils/issues/14)) ([f8fad71](https://github.com/SpechtLabs/go-otel-utils/commit/f8fad71195f16cc5778fe674a0de43dd31d57b51))
* **deps:** update module github.com/spechtlabs/go-otel-utils/otelprovider to v0.0.15 ([#18](https://github.com/SpechtLabs/go-otel-utils/issues/18)) ([6e38af1](https://github.com/SpechtLabs/go-otel-utils/commit/6e38af120822d943d15abd5e44d0772c8c966be2))
