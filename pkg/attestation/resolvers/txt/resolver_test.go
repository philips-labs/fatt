package txt_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/philips-labs/fatt/pkg/attestation/resolvers/txt"
)

func TestResolve(t *testing.T) {
	assert := assert.New(t)

	purlsFile := `pkg:docker/philips-labs/fatt@sha256:6cc65b2c82c2baa3391890abb8ab741efbcbc87baff3b06d5797afacb314ddd9?repository_url=ghcr.io&attestation_type=sbom
pkg:docker/philips-labs/fatt@sha256:6cc65b2c82c2baa3391890abb8ab741efbcbc87baff3b06d5797afacb314ddd9?repository_url=ghcr.io&attestation_type=provenance
pkg:docker/philips-labs/fatt@sha256:6cc65b2c82c2baa3391890abb8ab741efbcbc87baff3b06d5797afacb314ddd9?repository_url=ghcr.io
pkg:nuget/philips-labs/fatt@sha256:6cc65b2c82c2baa3391890abb8ab741efbcbc87baff3b06d5797afacb314ddd9?repository_url=nuget.org&attestation_type=provenance
pkg:nuget/philips-labs/fatt@sha256:823413cc65b2c82c2baa3391890abb8ab741e87baff3b06d5797afacb314ddd9?repository_url=nuget.org&attestation_type=sbom`

	r := &txt.Resolver{}
	atts, err := r.Resolve(strings.NewReader(purlsFile))
	assert.NoError(err)
	assert.Len(atts, 5)

	assert.Equal("SBOM", atts[0].Type.String())
	assert.Equal("Provenance", atts[1].Type.String())
	assert.Equal("Unknown", atts[2].Type.String())
	assert.Equal("Provenance", atts[3].Type.String())
	assert.Equal("SBOM", atts[4].Type.String())

	purlsFile = `pkg:docker/philips-labs/fatt@sha256:6cc65b2c82c2baa3391890abb8ab741efbcbc87baff3b06d5797afacb314ddd9?repository_url=ghcr.io&attestation_type=sbom
	ghcr.io/philips-labs/fatt@sha256:6cc65b2c82c2baa3391890abb8ab741efbcbc87baff3b06d5797afacb314ddd9`
	atts, err = r.Resolve(strings.NewReader(purlsFile))
	assert.Error(err)
	assert.EqualError(err, "scheme is missing")
	assert.Len(atts, 0)
}
