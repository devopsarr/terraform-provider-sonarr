# Changelog

## [2.0.0](https://github.com/devopsarr/terraform-provider-sonarr/compare/v1.2.2...v2.0.0) (2022-09-01)


### ⚠ BREAKING CHANGES

* move series to all_series

### Features

* add delay profile data source ([7118dd1](https://github.com/devopsarr/terraform-provider-sonarr/commit/7118dd14e18d124bdb5c58570254a503591f065e))
* add indexer config data source ([049cfe5](https://github.com/devopsarr/terraform-provider-sonarr/commit/049cfe5affd6f54b29094393b98bd289407e141a))
* add indexer data source ([af79355](https://github.com/devopsarr/terraform-provider-sonarr/commit/af79355f18f49ea80bb286c31af6f24ad6944b57))
* add language profile data source ([c980f9d](https://github.com/devopsarr/terraform-provider-sonarr/commit/c980f9d743d8f6c2ee03428282cacbbd1579145d))
* add media management data source ([6ee6c80](https://github.com/devopsarr/terraform-provider-sonarr/commit/6ee6c800357c2a3412c6d48c5706c8a8ee28f66a))
* add naming data source ([c4f5c8e](https://github.com/devopsarr/terraform-provider-sonarr/commit/c4f5c8e9e30866df18629662c9bdf65a62846690))
* add quality profile data soruce ([3f46e2b](https://github.com/devopsarr/terraform-provider-sonarr/commit/3f46e2b9ade3d839b8e91299a47c17f5902c2538))
* add root folder data source ([1958c1f](https://github.com/devopsarr/terraform-provider-sonarr/commit/1958c1f00b9c210103ca6bc406c1013f10a90e02))
* add series data source ([2a7f3ca](https://github.com/devopsarr/terraform-provider-sonarr/commit/2a7f3caf93c31073cc55ebe17cbfc7a0d5e98737))
* add tag data source ([4c260e0](https://github.com/devopsarr/terraform-provider-sonarr/commit/4c260e053b256a45e9b23b64a7a7f2f43593ccfd))
* move series to all_series ([b98d12b](https://github.com/devopsarr/terraform-provider-sonarr/commit/b98d12b38a38ed3002c41b838af1742a0814297e))


### Bug Fixes

* data source read from config ([ec64614](https://github.com/devopsarr/terraform-provider-sonarr/commit/ec646149bf6fcc876aac91c69d9e541d5df6c77d))
* list data source requirements ([c3aa214](https://github.com/devopsarr/terraform-provider-sonarr/commit/c3aa2140131e12a2892f0de5e5414a41f8aafb95))

## [1.2.2](https://github.com/devopsarr/terraform-provider-sonarr/compare/v1.2.1...v1.2.2) (2022-08-30)


### Bug Fixes

* indexer fields type assertion ([d4a746c](https://github.com/devopsarr/terraform-provider-sonarr/commit/d4a746ceaeb88483f21128d8399dfa1adc18ad9c))

## [1.2.1](https://github.com/devopsarr/terraform-provider-sonarr/compare/v1.2.0...v1.2.1) (2022-08-28)


### Bug Fixes

* release gpg dependencies ([b2e7065](https://github.com/devopsarr/terraform-provider-sonarr/commit/b2e7065aa4869219981759cf2af5aa06497d1db8))

## [1.2.1](https://github.com/devopsarr/terraform-provider-sonarr/compare/v1.2.0...v1.2.1) (2022-08-28)


### Bug Fixes

* release gpg dependencies ([0c5a8da](https://github.com/devopsarr/terraform-provider-sonarr/commit/0c5a8da58d37e9443daad242105ca47dbaf82c4e))

## [1.2.0](https://github.com/devopsarr/terraform-provider-sonarr/compare/v1.1.0...v1.2.0) (2022-08-28)


### Features

* add indexer config resource ([375172e](https://github.com/devopsarr/terraform-provider-sonarr/commit/375172e0a572efa3f7293db043bcba7aaa1cacec))
* add indexer resource ([b0602c5](https://github.com/devopsarr/terraform-provider-sonarr/commit/b0602c5b77e7939d6e81cf5ed7212154b0bb932a))
* add indexers data source ([1b2cc84](https://github.com/devopsarr/terraform-provider-sonarr/commit/1b2cc8479b7da70727edc2ba4fcabc289c272934))


### Bug Fixes

* remove set parameter for framework 0.9.0 ([370a452](https://github.com/devopsarr/terraform-provider-sonarr/commit/370a45213b53f87e46516f1404e4310c9ba384ae))
* series and qp datasource set ([3bdac5a](https://github.com/devopsarr/terraform-provider-sonarr/commit/3bdac5a62b55dd04b60cdefe19ec2ae79fa23308))

## [1.1.0](https://github.com/devopsarr/terraform-provider-sonarr/compare/v1.0.0...v1.1.0) (2022-06-06)


### Features

* add media_management resource ([00bfbc5](https://github.com/devopsarr/terraform-provider-sonarr/commit/00bfbc5a9d6d41e6d81b213d3ab228a29602b3dc))
* add naming resource ([79fc714](https://github.com/devopsarr/terraform-provider-sonarr/commit/79fc7142235d79f5cb23d1f5c2566ed8315b4d19))
* add profile validators ([53ce110](https://github.com/devopsarr/terraform-provider-sonarr/commit/53ce1105f5a8ff561039f7e66beb19bd57e425ad))
* add root_folder resources ([6a2d0bd](https://github.com/devopsarr/terraform-provider-sonarr/commit/6a2d0bd4942e9478b58f3b87d2f0f6104e57f184))
* add series resource ([3ec587a](https://github.com/devopsarr/terraform-provider-sonarr/commit/3ec587a31e3c73336ef3ff289ded1644149a23e6))

## 1.0.0 (2022-03-15)


### Features

* add profiles ([3f37682](https://github.com/devopsarr/terraform-provider-sonarr/commit/3f376822edaecbab1129ad202c093eef5d22ff4f))
* add tag resource ([98bb542](https://github.com/devopsarr/terraform-provider-sonarr/commit/98bb5426c912466173e61d1855064110b55ba384))
* align tag to sdk ([987ac70](https://github.com/devopsarr/terraform-provider-sonarr/commit/987ac702ec0f4bf3dd9e0b6a9b3ac119a036f7db))
* first configuration ([b6e0b48](https://github.com/devopsarr/terraform-provider-sonarr/commit/b6e0b486af8a547a2326c3fe8989c0d1abf77a6b))
