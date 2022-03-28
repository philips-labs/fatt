# fatt/installer-action

This action enables you to install `fatt` on your runner. During installation the integrity of `fatt` is verified based on it's cosign signature.

For a quick start guide on the usage of `fatt`, please refer to https://github.com/philips-labs/fatt#quick-start.
For available `cosign` releases, see https://github.com/sigstore/cosign/releases.

## Usage

This action currently supports GitHub-provided Linux, macOS and Windows runners.

Add the following entry to your Github workflow YAML file:

```yaml
uses: philips-labs/fatt/installer-action@main
with:
  fatt-release: 'v0.2.0' # optional
  install-path: '.fatt/bin' #optional
```

Now you can use fatt to `list` and `publish` attestations. Keep in mind to add `packages: write` permission if you want to publish the attestations to oci.

## Example

```yaml
jobs:
  publish-attestations:
    runs-on: ubuntu-20.04

    permissions:
      packages: write

    env:
      PACKAGE: ghcr.io/philips-labs/fatt/attestations-example
      PACKAGE_VERSION: v0.2.0

    steps:
      - name: Install cosign
        uses: sigstore/cosign-installer@v2.1.0
        with:
          cosign-release: v1.6.0

      - name: Install fatt
        uses: philips-labs/fatt/installer-action@main
        with:
          fatt-release: v0.2.0
          install-dir: .fatt/bin

      - name: Generate SBOM
        run: echo 'We could have generated a real sbom using syft here…' > sbom-spdx.json

      - name: Generate provenance
        run: echo 'We can use philips-labs/slsa-provenance-action to generate provenance…' > provenance.att

      - name: Login to ghcr.io
        uses: docker/login-action@dd4fa0671be5250ee6f50aedf4cb05514abda2c7 #v1.14.1
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: Install signing key
        run: echo "${{ secrets.COSIGN_PRIVATE_KEY }}" > cosign.key

      - name: Publish attestations
        run: |
          fatt publish \
            --repository "${PACKAGE}" \
            --version "${PACKAGE_VERSION}" \
            "sbom://sbom-spdx.json" "provenance://provenance.att"

      - name: Sign attestations and discovery
        run: |
          cosign sign --key cosign.key "${PACKAGE}:${PACKAGE_VERSION}.provenance"
          cosign sign --key cosign.key "${PACKAGE}:${PACKAGE_VERSION}.sbom"
          cosign sign --key cosign.key "${PACKAGE}:${PACKAGE_VERSION}.discovery"
        env:
          COSIGN_PASSWORD: ${{ secrets.COSIGN_PASSWORD }}

      - name: Discover attestations
        run: |
          echo "${{ secrets.COSIGN_PUBLIC_KEY }}" > cosign.pub
          fatt list --key cosign.pub "${PACKAGE}:${PACKAGE_VERSION}.discovery"

      - name: Cleanup signing key
        if: ${{ always() }}
        run: rm cosign.key
```
