## [2.3.2](https://gitlab.enix.io/products/stx/dothill-provisioner/compare/v2.3.1...v2.3.2) (2021-01-26)


### Bug Fixes

* **controller:** make ControllerPublishVolume idempotent again ([e1eade8](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/e1eade8641de9013d6c3baa6d7a848aece6fe490))
* **node:** check if a volume is mounted somewhere else before ejecting ([8650885](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/8650885b1cfcf113f0d56aec9725d2788c97fb0e)), closes [#72](https://gitlab.enix.io/products/stx/dothill-provisioner/issues/72)

## [2.3.1](https://gitlab.enix.io/products/stx/dothill-provisioner/compare/v2.3.0...v2.3.1) (2021-01-19)


### Bug Fixes

* **controller:** prevent multi-attach by disallowing multi-mapping ([8f058cf](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/8f058cfdb5f7d89217324f7085c7bc6461b13c05))

# [2.3.0](https://gitlab.enix.io/products/stx/dothill-provisioner/compare/v2.2.3...v2.3.0) (2021-01-12)


### Features

* **helm:** add daemonNodeSelector and daemonNodeAffinity ([0f41b0c](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/0f41b0cdebfe9116160bb93f59a09a2f7c4c1357)), closes [#63](https://gitlab.enix.io/products/stx/dothill-provisioner/issues/63)

## [2.2.3](https://gitlab.enix.io/products/stx/dothill-provisioner/compare/v2.2.2...v2.2.3) (2021-01-12)


### Bug Fixes

* **helm:** use a different component label for nodes and controller ([09bba25](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/09bba25422b1d97121168de6f5837b60f554bc38)), closes [#48](https://gitlab.enix.io/products/stx/dothill-provisioner/issues/48)
* **node:** check that a filesystem is healthy before mounting it ([4bae59f](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/4bae59f2074786522cab36022720041cc2eadc31)), closes [#66](https://gitlab.enix.io/products/stx/dothill-provisioner/issues/66)
* **node:** check transport type to use iscsi disks only ([0de5021](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/0de502138dbe2dd4c060c9cb90a5fff1c13d4556)), closes [#68](https://gitlab.enix.io/products/stx/dothill-provisioner/issues/68)
* **node:** fix csi-lib-iscsi commands error reporting ([ddcb3fd](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/ddcb3fd22345cfcb59f528d8eca6366526f227ba)), closes [#64](https://gitlab.enix.io/products/stx/dothill-provisioner/issues/64)
* **node:** support other filesystems than ext4 ([df77fd5](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/df77fd5bf21cb1d8ba4b935d730b5f58ea96ad79)), closes [#67](https://gitlab.enix.io/products/stx/dothill-provisioner/issues/67)

## [2.2.2](https://gitlab.enix.io/products/stx/dothill-provisioner/compare/v2.2.1...v2.2.2) (2021-01-07)


### Bug Fixes

* **node:** fix sda ejection when no device is found ([ecaa987](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/ecaa987960c7d53da8f7e46410e90795f57cf362))

## [2.2.1](https://gitlab.enix.io/products/stx/dothill-provisioner/compare/v2.2.0...v2.2.1) (2020-11-26)


### Bug Fixes

* **controller:** choose lun per initiatorName instead of globally ([72d4210](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/72d4210a334ffb4ccd50da4be95353ccc2d9068b)), closes [#60](https://gitlab.enix.io/products/stx/dothill-provisioner/issues/60)
* **controller:** limit simultaneous CSI calls to 1 by type of call ([abe4dbd](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/abe4dbd723131107e7874292a658845539993bf5))
* **controller:** limit worker-threads to 1 per CSI sidecar ([bf33878](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/bf33878712f8d1624987ba164aad4edbbe616217))
* **controller:** store csi-socket in tmpfs to allow node to be drained ([a9a7329](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/a9a7329f0e5fca4352bd056bc32e36de501eebd1))
* **helm:** fix controller leader election timeout on many new pvc ([a6976f2](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/a6976f2907a02036b8d3ac237de30b214e7dccf8))
* **helm:** WaitForFirstConsumer volume binding mode instead of Immediate ([359aedd](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/359aeddaa33468aec0771d8efbeef8acbe7dfa8d)), closes [#58](https://gitlab.enix.io/products/stx/dothill-provisioner/issues/58)
* version printing ([ac6da3d](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/ac6da3d50f2aada01edeef3ca040b02885a21f22)), closes [#59](https://gitlab.enix.io/products/stx/dothill-provisioner/issues/59)
* **node:** don't disconnect target when multipath is not working ([0d3ed8f](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/0d3ed8f79b3a4b1902167d9a4f8793ee25b604e0))
* **node:** don't give an error when a disk is already unmounted ([f7647ae](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/f7647aec557bfe02a2b8a4b23ae4cfcf9b300cbf))
* **node:** errors when multipath is not used for multiple targets ([c3b7f5b](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/c3b7f5ba2c0d7c7d61272bb8cb60d9d16a19ae6b))
* **node:** fix choose LUN with high value existing LUNs ([55b4aa3](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/55b4aa3159d5ee9502d564ce061d942b44da52f6))
* **node:** handle gracefully already existing socket ([71fcb4d](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/71fcb4d36119a68387e71f7ce84be511796c13ee)), closes [#47](https://gitlab.enix.io/products/stx/dothill-provisioner/issues/47)
* **node:** limit NodePublishVolume concurrency to 1, fail otherwise ([3cbc6a9](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/3cbc6a9e91d05aa71e98bb80f564830c4137cdf5))
* **node:** persist multipathd configuration ([0c2e14f](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/0c2e14fa54d6fe18890b977ef4a456bb29c86981))
* **node:** remove mutex publish/unpublish mutex ([40f3518](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/40f3518ebee2892180751b93790ee66a72573024))
* **node:** upgrade csi-lib-iscsi to fix devices ejection ([89feac7](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/89feac7f7039e4541a0bd2543a2a41a91b3c1c93))

# [2.2.0](https://gitlab.enix.io/products/stx/dothill-provisioner/compare/v2.1.0...v2.2.0) (2020-11-09)


### Bug Fixes

* **node:** don't rescan iSCSI sessions after removing device ([84653ac](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/84653ac8746a59af046910e7f63ca1881b9c20df)), closes [#52](https://gitlab.enix.io/products/stx/dothill-provisioner/issues/52)
* **node:** upgrade csi-lib-iscsi to disconnect volumes properly ([4a8b64e](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/4a8b64ea4e977e21372c49c56a37a42c854639d9)), closes [#32](https://gitlab.enix.io/products/stx/dothill-provisioner/issues/32)


### Features

* **node:** add support for multipathd ([480a6ed](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/480a6ed54b8afacf82f5ac4c9d2419e5701454a2)), closes [#33](https://gitlab.enix.io/products/stx/dothill-provisioner/issues/33)

# [2.1.0](https://gitlab.enix.io/products/stx/dothill-provisioner/compare/v2.0.2...v2.1.0) (2020-10-28)


### Bug Fixes

* **node:** add missing discovery flag in iscsi connector ([9db3dc5](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/9db3dc50d93ce9be42da4f0aa3d0dade5189e0a7))
* **testing:** fix compilation ([cbdeffd](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/cbdeffde1fdffb08e21abc49e8c5baca2abd3c25))


### Features

* **example:** add 2 examples ([82bfc10](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/82bfc10becc3813b6887438c5858268eac0dff7d))
* **node:** add CSI liveness probe with iscsi_tcp kernel mod check ([08e84b3](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/08e84b30781b4fe41e8b25ec0a706397e65446d3)), closes [#50](https://gitlab.enix.io/products/stx/dothill-provisioner/issues/50)
* **node:** containerize iscsid  ([a90d306](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/a90d306b239e26838bf957e6eb72b9e0ca6e8372)), closes [#44](https://gitlab.enix.io/products/stx/dothill-provisioner/issues/44) [#49](https://gitlab.enix.io/products/stx/dothill-provisioner/issues/49)
* **node:** generate unique iqn once per node ([cd98706](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/cd987064003d112c6253061aadc5634a721a88ea)), closes [#23](https://gitlab.enix.io/products/stx/dothill-provisioner/issues/23)
* **node:** load iscsi_tcp kernel module in initContainer ([4fe865d](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/4fe865d143d3721ad198e848b9279e4a5aafc634))

## [2.0.2](https://gitlab.enix.io/products/stx/dothill-provisioner/compare/v2.0.1...v2.0.2) (2020-10-14)


### Bug Fixes

* **ci:** push docker images tagged with the "v" prefix ([d3a2bfc](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/d3a2bfcc92de000b7033836944bb073661b1ad06))

## [2.0.1](https://gitlab.enix.io/products/stx/dothill-provisioner/compare/v2.0.0...v2.0.1) (2020-10-14)


### Bug Fixes

* **ci:** prevent dirty build on release ([57deab8](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/57deab8f9ebaaaf9c811ea3e786cc52c6fa9d702))
* **helm:** convert name to kebab-case to fix GitHub release name ([26655cb](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/26655cb73d493b5ef4e99b4d4507cb5a413da62f))

# [2.0.0](https://gitlab.enix.io/products/stx/dothill-provisioner/compare/v1.4.0...v2.0.0) (2020-10-14)


### Bug Fixes

* **ci:** fix chart version ([352a224](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/352a22410c8bb8a1e1f6005d1559d5094792d7d7))


### BREAKING CHANGES

* **ci:** complete project rewrite as a CSI plugin

# [1.4.0](https://gitlab.enix.io/products/stx/dothill-provisioner/compare/v1.3.0...v1.4.0) (2020-10-14)


### Bug Fixes

* **arch:** fix inexistant go client version ([2faaa85](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/2faaa85adf0efdd84c544b5585e9907355b821e2))
* **build:** fix dependencies for production build ([47dcfcd](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/47dcfcda74191b2ef7c6c337db23346adb8823c5))
* **build:** fix failing build after go mod migration ([08d7911](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/08d79115e1b71267cb5fe426dcfc6f582ad3ea4b))
* **build:** fix klog dep version ([aa192c3](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/aa192c37469d53a28ac7b8a12060e520040660a9))
* **ci:** add commit sha  to un-versionned images ([2887fc3](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/2887fc3441ef9a51803fe589c6f6ce3d7e233a37))
* **ci:** add gcc to build image ([2a84e5f](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/2a84e5ffe741690158e41745ef7bbea6d73d71f8))
* **ci:** add git to build imahge ([65d57d0](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/65d57d0912e1594017d920f9f9e4bd97f1879fb5))
* **ci:** don't create an image for every commit ([d7d5fe5](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/d7d5fe5903c57db329d9e25de1a3fc195396a8b1))
* **ci:** don't skip build and push on changelog update ([6b831b0](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/6b831b0c705d3447b9f8664eb936876ef28ffd02))
* **ci:** fix dependencies error ([4ce6ed0](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/4ce6ed0b2e0dbe716d40789f045fcef7caf69909))
* **ci:** fix jobs dependencies ([4fdd899](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/4fdd899fa6ab0177e79ab5e5f18ed914fca0489c))
* **ci:** remove sanity tests ([716527a](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/716527ae01c588a9ff0a4a7c76872b7b4a79798b))
* **controller:** fail on publish volume when host doesn't exists ([07203cb](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/07203cbc9d1b64f41246c7e160bf3cbeedeba3f8))
* **controller:** fail on validate volume caps when volume doesn't exists ([ba8bed0](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/ba8bed01cb54553ca2a5b878efb9330ee3f49714))
* **controller:** fix CreateVolume with an existing name for sanity tests ([c6085e1](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/c6085e13eef771733f913667116f5ff68b62d26e))
* **controller:** fix LUN selection with no volume ([8c44200](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/8c442002202a4d39e9097739fd3972ac8da8397e))
* **controller:** minor fixes ([c763c2e](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/c763c2ee1de4c93e8e10b59132410c15c7584c80))
* **controller:** revert fail on publish volume when host doesn't exists ([33fbc15](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/33fbc15e8db4eda1b5038722df7ad926165750d6))
* **controller:** unmap volume from given node id instead of all nodes ([51e719e](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/51e719ea50a9c6de47a3c1fb0c7101ca8918f8f8))
* **deploy:** consistent plugin name accross configs ([ac5bff5](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/ac5bff51be0974136800c72a4ee44dd59ac5a495))
* **deploy:** replace default namespace by kube-system ([f068850](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/f068850cf0fc3b18fbf2d7808fd88c9a23ec245e))
* **deploy:** set fixed versions on external dependencies ([cf53ffd](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/cf53ffdcec5bc5b667cc92073c1e88c29174b360))
* **deploy:** update deploy files to properly deploy drivers ([207fcab](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/207fcabd181309630c3928f222a49a9e0fc0c2ba))
* **deploy:** working k8s deployment ([cc94561](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/cc945616d025a734211bde6961e1e456b435613e))
* **example:** fix example storage class name ([7bfa944](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/7bfa94451186ecb76f2e9d83469652e07133eb93))
* **helm:** fix cluster role binding service account namespace ([ee624f7](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/ee624f70007cae527fb8c01177fc679ab6ca1559))
* **helm:** handle psp admission controller ([f0882d3](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/f0882d3b76a2f64e446b600abf818d1355486ad8))
* **helm:** remove imagePullSecret as all images used are public ([069010b](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/069010b97dbb8b3841c22a51fa582f8d3d33636a))
* **node:** check if device is actually using multipath instead of just assuming it ([10b92b8](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/10b92b8292220a138e5b23e22728e6c9894d8c1b))
* **node:** detect file system on alpine and debian ([8d22bf9](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/8d22bf9f728255e227e19eb3791586bafc8a6aaa))
* **node:** minor mount fixes ([75637bb](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/75637bb77a74d896be1e166ed2a3e8cb4568bb63))
* **node:** rescan iscsi sessions after volume unmount ([08512c4](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/08512c4a3b810d1eaf232ba6fe22e940957c6227))
* **provisioner:** add klog flag missing after update ([160a3cb](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/160a3cb2fb0a0bc6a27a8b22d14635c9bfd6770e))
* **resizer:** remove useless code ([292d3b9](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/292d3b922c55a8bfcf8ad0b28d4eff0d89167cf3))
* **resizer:** send size difference to api instead of new size ([09e003c](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/09e003c7390aee9129ef68cc506b16f9a7d74bc3))
* **resizer:** tell k8s to resize the fs after expanding ([470d6bc](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/470d6bca981b410267bc400187ff36441ec058d6))
* **testing:** fix exit code and initiator name file ([d37fafc](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/d37fafcca4833dc0ce062022730a22991f412709))
* **testing:** improve sanity requirements error handling ([959e4ec](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/959e4ec07b5ead410f23d27f903057efa30d7a77))


### Features

* **build:** optimize docker build ([aa13e86](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/aa13e8654157bd90a5dbec705c6fa3267f9e2197))
* **ci:** deploy helm chart on GitHub and image on docker hub ([bc8c81a](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/bc8c81a6f55cedbf6e862223e8058b3c5a4a30dd))
* **ci:** generate CHANGELOG.md ([657685a](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/657685a9d1d696469009a898b671f18ce71bb5d5))
* **controller:** implement validate volume capabilities ([1ecd5bc](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/1ecd5bc87480212572877af24cb39328d65a86f7))
* **controller:** implement volume attach routine ([fa16310](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/fa16310d45d63ffdb3d791c1041c0aeabbf1e9ae))
* **controller:** make controller unpublish idempotent ([de0ebbc](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/de0ebbc99e01f032778250b5b95a7305eb401d45))
* **controller:** migrate provisioning code ([81c50d6](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/81c50d6c4283bd145fe0c3145c2ce1968808980e))
* **deploy:** add an example of values.yaml ([06d2dcf](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/06d2dcf11ec1b0b856c20a4b0d87ae6dd345a683))
* **deploy:** add external attacher container into deployment ([528d529](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/528d529b4288290f638b4d67b08f677bf1df344f))
* **deploy:** build a rc on master and trigger releases manually only ([4a8092d](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/4a8092dbd665a2b5294421a7b52c8bb2ae1e0d97))
* **driver:** graceful stop & unbind on signal ([7a641c7](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/7a641c754f2ca923c1e3843e3c826ef56898d4eb))
* **driver:** properly handle concurrent RPCs ([38342b4](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/38342b4fcbdba5e1344c994d9884b4cca2d13d1e))
* **driver:** working provision/deletion with multipath ([e79b930](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/e79b930684af507a3fe236e77030987cdd5e110b))
* **helm:** initial chart ([80e79ad](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/80e79ad2f16692e1bb7486407cac5a5e6e7099e4))
* **helm:** split extraArgs config for node and controller ([13d7aae](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/13d7aaeb3edcf02e7742ba845c3e2e546a647189))
* **helm:** template kubelet path with default value ([1ba0e17](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/1ba0e178a2d8bc2c0ffc51f797269e03a2c6867b))
* **helm:** template storage classes and their secrets ([61a9773](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/61a9773028b52cefd6bdf5420367ede3b0199840))
* **node:** allow to use a custom path for /var/lib/kubelet ([6d3d793](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/6d3d7934322113efd2ec021ab50370ca668ddf07))
* **node:** codebase for node plugin ([9401af4](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/9401af47dbb3bd7025b4ad69d41687fc8fc84d36))
* **node:** format disks only if needed ([cce9571](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/cce9571edf81ec60298863a6b6cd80aa76c111f7))
* **node:** implement info calls ([e4f84ca](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/e4f84ca783b923a239ba98d390b659f0c1f28d87))
* **node:** mount published volumes ([8f0cbb1](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/8f0cbb12dbd9ed3610ebb51ef05bf01bd1904d4c))
* **resizer:** run preflight checks to ensure resizing is possible ([43dbc04](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/43dbc045f8ea634090198e8274d667073952a684))
* **testing:** properly implement k8s sanity tests & fix most of them ([aa658d6](https://gitlab.enix.io/products/stx/dothill-provisioner/commit/aa658d642eee3a94c2c2f8caef654b7941e2f134))
