# Fatt

Fatt is a small commandline utility that allows you to fetch attestations for your binaries or other non OCI package types. This can be done by searching the given file path for an `attestations.txt` file or by providing a reference to an OCI blob that contains this `attestations.txt` contents.

> :warning: This project is a POC.

`fatt` also offers a feature to generate and publish this `attestations.txt` to an OCI registry. Within the `attestations.txt` references to your attestations are stored in [purl][] format. Next to that these attestations can also be printed in oci format so you can use them in shell scripts with for example `sget`. See the usage and examples.

## Fatt Usage

```bash
$ /fatt --help
Discover and resolve your attestations

Usage:
  fatt [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  list        Lists all attestations
  publish     Publishes given attestations to an OCI registry
  version     Prints the fatt version

Flags:
  -h, --help   help for fatt

Use "fatt [command] --help" for more information about a command.
```

### List filter options

Filters use the Go template language.

The following fields are supported.

* Type
* PURL.Type
* PURL.Namespace
* PURL.Name
* PURL.Version
* PURL.Subpath
* PURL.Qualifiers

The following functions are available.

* `func (Attestation) IsRegistry(registryURL string) bool`
* `func (Attestation) IsAttestationType(t string) bool`

```bash
$ fatt list -f '{ .IsRegistry("ghcr.io") && .IsAttestationType("sbom") }'
pkg:docker/philips-labs/fatt@sha256:6cc65b2c82c2baa3391890abb8ab741efbcbc87baff3b06d5797afacb314ddd9?repository_url=ghcr.io&attestation_type=sbom
pkg:docker/philips-labs/fatt@sha256:6cc65b2c82c2baa3391890abb8ab741efbcbc87baff3b06d5797afacb314ddd9?repository_url=ghcr.io&attestation_type=sbom
```

## OCI as blob storage

Using OCI as a blob storage we can leverage this to distribute and discover attestations for software assets that are not OCI/Docker images. For this ecosystem we already have nice integration using [cosign](https://github.com/sigstore/cosign). With `fatt` we expand these capabilities to also support other ecosystems like NPM, Nuget, Gradle. Based on a publishing convention we can make it clear for our package consumers where they can retrieve our `build provenance` and `SBOM`.

Any attestations published to an OCI registry should be captured in a `attestations.txt`. Fatt publish command automates the manual steps to upload the attestations using cosign and then generates an `attestations.txt` using PURL references of the published attestations. Then it also automatically publishes these `attestations.txt`. Manually you would perform following steps:

```shell
$ tree .
.
├── provenance.att
└── sbom-spdx.json

$ cosign upload blob -f provenance.att ghcr.io/philips-labs/fatt-attestations-example:v0.1.0.provenance
$ cosign sign --key cosign.key ghcr.io/philips-labs/fatt-attestations-example:v0.1.0.provenance
$ cosign upload blob -f sbom-spdx.json ghcr.io/philips-labs/fatt-attestations-example:v0.1.0.sbom
$ cosign sign --key cosign.key ghcr.io/philips-labs/fatt-attestations-example:v0.1.0.sbom

# Now you would manually fill out an attestations.txt using these image digests following this spec
# https://github.com/package-url/purl-spec/blob/master/PURL-TYPES.rst#oci
$ tree
.
├── attestations.txt
├── provenance.att
└── sbom-spdx.json

$ cosign upload blob -f attestations.txt ghcr.io/philips-labs/fatt-attestations-example:v0.1.0.discover
$ cosign sign --key cosign.key ghcr.io/philips-labs/fatt-attestations-example:v0.1.0.discover
# Now we have the discovery file stored in an OCI registry we can look it up
$ fatt list --key cosign.pub ghcr.io/philips-labs/fatt-attestations-example:v0.1.0.discover
```

See below an example on how `fatt publish` automates this whole flow for you. (:warning: NOTE that signing currently is not yet build in. [#20](https://github.com/philips-labs/fatt/issues/20)) To identify the attestation type we define this in a `URL Scheme` format (see the example below).

```shell
$ tree .
.
├── provenance.att
└── sbom-spdx.json

$ fatt publish --repository ghcr.io/philips-labs/fatt-attestations-example --version v0.1.0 sbom://sbom-spdx.json provenance://provenance.att
Publishing attestations…
Uploading file from [sbom-spdx.json] to [ghcr.io/philips-labs/fatt-attestations-example:v0.1.0.sbom] with media type [text/plain]
File [sbom-spdx.json] is available directly at [ghcr.io/v2/philips-labs/fatt-attestations-example/blobs/sha256:877084e55eb2896eb3d159df7483862e8f7470469d9ac732a54da2298bcf456c]
Uploading file from [provenance.att] to [ghcr.io/philips-labs/fatt-attestations-example:v0.1.0.provenance] with media type [text/plain]
File [provenance.att] is available directly at [ghcr.io/v2/philips-labs/fatt-attestations-example/blobs/sha256:a167d9ca71c4fda26e092eaa0a1d5242389b2f202ca822dff8f088faf8cce00e]

Generating attestations.txt based on uploaded attestations…
Uploading file from [attestations.txt] to [ghcr.io/philips-labs/fatt-attestations-example:v0.1.0.discovery] with media type [text/plain]
File [attestations.txt] is available directly at [ghcr.io/v2/philips-labs/fatt-attestations-example/blobs/sha256:2106cfd71501952197e00e1099b515fbcbe4dd852c7bf2bd4a87fa58d3bae0d7]
$ cosign sign --key cosign.key ghcr.io/philips-labs/fatt-attestations-example:v0.1.0.sbom
$ cosign sign --key cosign.key ghcr.io/philips-labs/fatt-attestations-example:v0.1.0.provenance
$ cosign sign --key cosign.key ghcr.io/philips-labs/fatt-attestations-example:v0.1.0.discover
$ fatt list "ghcr.io/philips-labs/fatt-attestations-example:v0.1.0.discovery" > attestations.txt
Fetching attestations from ghcr.io/philips-labs/fatt-attestations-example@sha256:2106cfd71501952197e00e1099b515fbcbe4dd852c7bf2bd4a87fa58d3bae0d7…

$ cat attestations.txt
pkg:oci/philips-labs/fatt-attestations-example@sha256:d17ece80fca09d53d5d23c54900697870fa1dc2c9161097c22d59b3775b88cc0?repository_url=ghcr.io%2Fphilips-labs%2Ffatt-attestations-example&tag=v0.1.0.sbom
pkg:oci/philips-labs/fatt-attestations-example@sha256:f25d28beea7c81af4160a32256831380d7173449cfc49dde70bcca1b697f9c7e?repository_url=ghcr.io%2Fphilips-labs%2Ffatt-attestations-example&tag=v0.1.0.provenance

$ fatt list --output-format oci .
Fetching attestations from ghcr.io/philips-internal/slsa-workflow-examples/fatt-test@sha256:03a820de6d554b6d77c8323c459f8f17ec6bb970f448b551d5a912a201e532e1…
ghcr.io/philips-labs/fatt-attestations-example@sha256:d17ece80fca09d53d5d23c54900697870fa1dc2c9161097c22d59b3775b88cc0
ghcr.io/philips-labs/fatt-attestations-example@sha256:f25d28beea7c81af4160a32256831380d7173449cfc49dde70bcca1b697f9c7e
```

[purl]: https://github.com/package-url/purl-spec "A minimal specification and implementation of purl aka. a Package 'mostly universal' URL."
