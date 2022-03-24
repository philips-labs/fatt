package oci_test

import (
	"testing"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/package-url/packageurl-go"
	"github.com/stretchr/testify/assert"

	"github.com/philips-labs/fatt/pkg/oci"
)

func TestToPackageUrl(t *testing.T) {
	assert := assert.New(t)

	ociRef, err := name.ParseReference("ghcr.io/philips-labs/fatt:v0.1.0.sbom")
	assert.NoError(err)

	// TODO: to comply with the purl spec need to have digest as version
	expectedPURL, err := packageurl.FromString("pkg:oci/philips-labs/fatt@v0.1.0.sbom?repository_url=ghcr.io/philips-labs/fatt&tag=v0.1.0.sbom")
	assert.NoError(err)

	purl, err := oci.ToPackageURL(ociRef)
	assert.NoError(err)
	assertPURL(t, expectedPURL, *purl)

	ociRef, err = name.ParseReference("ghcr.io/philips-labs/with-long-repo/fatt:v0.1.0.provenance")
	assert.NoError(err)

	// TODO: to comply with the purl spec need to have digest as version
	expectedPURL, err = packageurl.FromString("pkg:oci/philips-labs/with-long-repo/fatt@v0.1.0.provenance?repository_url=ghcr.io/philips-labs/with-long-repo/fatt&tag=v0.1.0.provenance")
	assert.NoError(err)

	purl, err = oci.ToPackageURL(ociRef)
	assert.NoError(err)
	assertPURL(t, expectedPURL, *purl)

	ociRef, err = name.ParseReference("philipssoftware/fatt:v0.1.0.sbom")
	assert.NoError(err)

	// TODO: to comply with the purl spec need to have digest as version
	expectedPURL, err = packageurl.FromString("pkg:oci/philipssoftware/fatt@v0.1.0.sbom?repository_url=index.docker.io/philipssoftware/fatt&tag=v0.1.0.sbom")
	assert.NoError(err)

	purl, err = oci.ToPackageURL(ociRef)
	assert.NoError(err)
	assertPURL(t, expectedPURL, *purl)
}

func TestFromPackageURL(t *testing.T) {
	assert := assert.New(t)

	purl, err := packageurl.FromString("pkg:oci/philips-labs/slsa-provenance@sha256:e3378aef23821fd6e210229e5b98b5bead2858581b2d590d9e3b49d53c3f71e7?repository_url=ghcr.io/philips-labs/slsa-provenance&tag=v0.7.2")
	assert.NoError(err)

	expectedOCIRef, err := name.ParseReference("ghcr.io/philips-labs/slsa-provenance:v0.7.2")
	assert.NoError(err)

	ociRef, err := oci.FromPackageURL(purl)
	assert.NoError(err)
	assert.Equal(expectedOCIRef.String(), ociRef.String())

	purl, err = packageurl.FromString("pkg:oci/philips-labs/slsa-provenance@sha256:e3378aef23821fd6e210229e5b98b5bead2858581b2d590d9e3b49d53c3f71e7?repository_url=ghcr.io/philips-labs/slsa-provenance")
	assert.NoError(err)

	expectedOCIRef, err = name.ParseReference("ghcr.io/philips-labs/slsa-provenance@sha256:e3378aef23821fd6e210229e5b98b5bead2858581b2d590d9e3b49d53c3f71e7")
	assert.NoError(err)

	ociRef, err = oci.FromPackageURL(purl)
	assert.NoError(err)
	assert.Equal(expectedOCIRef.String(), ociRef.String())

	purl, err = packageurl.FromString("pkg:oci/philipssoftware/slsa-provenance@sha256:e3378aef23821fd6e210229e5b98b5bead2858581b2d590d9e3b49d53c3f71e7?repository_url=index.docker.io/philipssoftware/slsa-provenance")
	assert.NoError(err)

	expectedOCIRef, err = name.ParseReference("philipssoftware/slsa-provenance@sha256:e3378aef23821fd6e210229e5b98b5bead2858581b2d590d9e3b49d53c3f71e7")
	assert.NoError(err)

	ociRef, err = oci.FromPackageURL(purl)
	assert.NoError(err)
	assert.Equal(expectedOCIRef.String(), ociRef.String())

	purl, err = packageurl.FromString("pkg:oci/philipssoftware/slsa-provenance@sha256:e3378aef23821fd6e210229e5b98b5bead2858581b2d590d9e3b49d53c3f71e7")
	assert.NoError(err)

	expectedOCIRef, err = name.ParseReference("philipssoftware/slsa-provenance@sha256:e3378aef23821fd6e210229e5b98b5bead2858581b2d590d9e3b49d53c3f71e7")
	assert.NoError(err)

	ociRef, err = oci.FromPackageURL(purl)
	assert.NoError(err)
	assert.Equal(expectedOCIRef.String(), ociRef.String())

	purl, err = packageurl.FromString("pkg:oci/library/alpine@sha256:ceeae2849a425ef1a7e591d8288f1a58cdf1f4e8d9da7510e29ea829e61cf512?repository_url=docker.io/library/alpine&tag=latest")
	assert.NoError(err)

	expectedOCIRef, err = name.ParseReference("alpine:latest")
	assert.NoError(err)

	ociRef, err = oci.FromPackageURL(purl)
	assert.NoError(err)
	assert.Equal(expectedOCIRef.String(), ociRef.String())
}

func assertPURL(t *testing.T, expected, actual packageurl.PackageURL) {
	assert := assert.New(t)

	assert.Equal(expected.Type, actual.Type)
	assert.Equal(expected.Namespace, actual.Namespace)
	assert.Equal(expected.Name, actual.Name)
	assert.Equal(expected.Version, actual.Version)
	assert.Equal(expected.Qualifiers, actual.Qualifiers)
	assert.Equal(expected.Subpath, actual.Subpath)
}
