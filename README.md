# Fatt

Fatt is a small commandline utility that allows you to fetch attestations for your binaries. It will do so by looking at some locations in the code repository.

> :warning: This project is currently nothing more than a POC.

`fatt` tries to find any [purl][] in your project by looking at predefined fields in the [supported packages](#supported-packages-and-attestations). These fields describe using a [purl][] where to grab the attestation from.

## Fatt Usage

```bash
$ ./bin/fatt --help
Discover and resolve your attestations

Usage:
  fatt [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  list        Lists all attestations
  version     Prints the fatt version

Flags:
  -p, --file-path string   the filepath to find attestation purls (defaults to current working dir)
  -h, --help               help for fatt
  -r, --resolver string    the resolver to use for finding attestations (default "multi")

Use "fatt [command] --help" for more information about a command.
```

### List command: Filter Option

The following attestations fields can be filtered on.

* Type
* Namespace
* Version
* Name
* Scheme

The following functions are available.

* `func (Attestation) IsRegistry(registryURL string) bool`
* `func (Attestation) IsAttestationType(t string) bool`

```bash
$ bin/fatt list -f '{ .IsRegistry("ghcr.io") && .IsAttestationType("sbom") }'
pkg:docker/philips-labs/fatt@sha256:6cc65b2c82c2baa3391890abb8ab741efbcbc87baff3b06d5797afacb314ddd9?repository_url=ghcr.io&attestation_type=sbom
pkg:docker/philips-labs/fatt@sha256:6cc65b2c82c2baa3391890abb8ab741efbcbc87baff3b06d5797afacb314ddd9?repository_url=ghcr.io&attestation_type=sbom
```

## Supported packages and attestations

### NPM

#### SBOM

To fetch an SBOM you can define a [purl][] with `attestation_type`=`sbom` qualifier in `package.json` within a attestations array.

<details>
  <summary>Example cosign stored sbom</summary>

  Using cosign we can leverage any [OCI registry][] to store our attestations.

  ```shell
  $ cosign upload blob -f sbom.spdx.json ghcr.io/philips-labs/fatt:example-sbom-attestation
  Uploading file from [sbom.spdx.json] to [ghcr.io/philips-labs/fatt:example-sbom-attestation] with media type [text/plain]
  File [sbom.spdx.json] is available directly at [ghcr.io/philips-labs/fatt@sha256:6cc65b2c82c2baa3391890abb8ab741efbcbc87baff3b06d5797afacb314ddd9]
  ```

  Now we can use a purl to link this attestation to our Node package.

  ```json
  {
    "name": "@philips-labs/awesome-npm",
    "attestations": [
      "pkg:docker/philips-labs/fatt@sha256:6cc65b2c82c2baa3391890abb8ab741efbcbc87baff3b06d5797afacb314ddd9?repository_url=ghcr.io&attestation_type=sbom",
    ]
  }
  ```

  Using `fatt` we can now scan our project for attestations and fetch them using sget.

  ```shell
  $ attestations="$(bin/fatt list -p examples/awesome-npm -o docker)"
  Fetching attestations for current working directory…
  Found attestations: [{PURL:{Type:docker Namespace:philips-labs Name:fatt Version:sha256:6cc65b2c82c2baa3391890abb8ab741efbcbc87baff3b06d5797afacb314ddd9 Qualifiers:repository_url=ghcr.io&attestation_type=sbom Subpath:} Type:SBOM} {PURL:{Type:docker Namespace:philips-labs Name:fatt Version:sha256:6cc65b2c82c2baa3391890abb8ab741efbcbc87baff3b06d5797afacb314ddd9 Qualifiers:repository_url=ghcr.io&attestation_type=provenance Subpath:} Type:SBOM}]
  Attestation type: sbom
  Attestation type: provenance
  $ while read -r a ; do sget "$a" ; done <<< "$attestations"
  {
    "SPDXID": "SPDXRef-DOCUMENT",
    "name": "ghcr.io/philips-labs/slsa-provenance-v0.7.2",
    "spdxVersion": "SPDX-2.2",
    "creationInfo": {
      "created": "2022-02-25T13:01:35.3837117Z",
      "creators": [
        "Organization: Anchore, Inc",
        "Tool: syft-0.38.0"
      ],
      "licenseListVersion": "3.16"
    },
    …
  ```

</details>

[purl]: https://github.com/package-url/purl-spec "A minimal specification and implementation of purl aka. a Package 'mostly universal' URL."
