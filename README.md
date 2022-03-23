# Fatt

Fatt is a small commandline utility that allows you to fetch attestations for your binaries. It will do so by looking at some locations in the code repository.

> :warning: This project is currently nothing more than a POC.

`fatt` tries to find any [purl][] in your project by searching the given path recursively for `attestations.txt`. Within an `attestations.txt` you can describe where your project stores attestations using [purl][] format.

In addition `fatt` allows to fetch these attestations from an OCI registry. It assumes that given location contains an uploaded blob using cosign containing the contents of `attestations.txt`.

## Fatt Usage

```bash
$ ./bin/fatt list --help
Lists all attestations

Usage:
  fatt list <discovery-path> [flags]

Flags:
  -f, --filter string          filter attestations using template expressions
  -h, --help                   help for list
      --key string             path to the public key file, URL, or KMS URI
  -o, --output-format string   output format for the list (default "purl")
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
$ bin/fatt list -f '{ .IsRegistry("ghcr.io") && .IsAttestationType("sbom") }'
pkg:docker/philips-labs/fatt@sha256:6cc65b2c82c2baa3391890abb8ab741efbcbc87baff3b06d5797afacb314ddd9?repository_url=ghcr.io&attestation_type=sbom
pkg:docker/philips-labs/fatt@sha256:6cc65b2c82c2baa3391890abb8ab741efbcbc87baff3b06d5797afacb314ddd9?repository_url=ghcr.io&attestation_type=sbom
```

## OCI

Any attestations published to an OCI registry should be captured in a `attestations.txt`. To distribute this discovery filewe can publish this as well to an OCI registry.

We can do this utilizing cosign. (we might add integration later to reduce the manual steps).

<details>
  <summary>Store attestations.txt using cosign.</summary>

  Using cosign we can leverage any [OCI registry][] to store our attestations. Once we stored the attestations we can capture that in an `attestations.txt` using [purl][] format. This `attestations.txt` we can also store in the [OCI registry][].

  ```shell
  $ cosign upload blob -f attestations.txt ghcr.io/philips-labs/fatt:attestations
  Uploading file from [attestations.txt] to [ghcr.io/philips-labs/fatt:attestations] with media type [text/plain]
  File [attestations.txt] is available directly at [ghcr.io/philips-labs/fatt@sha256:6cc65b2c82c2baa3391890abb8ab741efbcbc87baff3b06d5797afacb314ddd9]
  $ cosign sign --key cosign.key ghcr.io/philips-labs/fatt:attestations
  ```

  Using `fatt` we can now list the attestations stored in the OCI registry. `fatt` utilizes `sget` to fetch the `attestations.txt` and verify the signature. As we captured our attestations in PURL format we can also translate the attestations to docker format so we can also utilize sget to fetch the attestations themself.

  ```shell
  $ attestations="$(bin/fatt list --key cosign.pub -o docker ghcr.io/philips-labs/fatt:attestations)"
  Fetching attestations from ghcr.io/philips-labs/fatt:attestations…

  Verification for ghcr.io/philips-internal/attestations/slsa-workflow-examples/awesome-node-cli --
  The following checks were performed on each of these signatures:
    - The cosign claims were validated
    - The signatures were verified against the specified public key

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
