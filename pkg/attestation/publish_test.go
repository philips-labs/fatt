package attestation

import (
	"testing"

	"github.com/package-url/packageurl-go"
	"github.com/stretchr/testify/assert"
)

func TestToPackageUrl(t *testing.T) {
	assert := assert.New(t)

	ociRef := "ghcr.io/philips-labs/fatt:v0.1.0.sbom"
	// TODO: to comply with the purl spec need to have digest as version
	expectedPURL, err := packageurl.FromString("pkg:oci/philips-labs/fatt@v0.1.0.sbom?repository_url=ghcr.io/philips-labs/fatt&tag=v0.1.0.sbom")
	assert.NoError(err)

	purl, err := toPackageURL(ociRef)
	assert.NoError(err)
	assertPURL(t, expectedPURL, *purl)

	ociRef = "ghcr.io/philips-labs/with-long-repo/fatt:v0.1.0.provenance"
	// TODO: to comply with the purl spec need to have digest as version
	expectedPURL, err = packageurl.FromString("pkg:oci/philips-labs/with-long-repo/fatt@v0.1.0.provenance?repository_url=ghcr.io/philips-labs/with-long-repo/fatt&tag=v0.1.0.provenance")
	assert.NoError(err)

	purl, err = toPackageURL(ociRef)
	assert.NoError(err)
	assertPURL(t, expectedPURL, *purl)
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
