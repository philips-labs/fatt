package attestation_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/philips-labs/fatt/pkg/attestation"
	"github.com/philips-labs/fatt/pkg/attestation/resolvers/txt"
)

func TestReducePurls(t *testing.T) {

	purlsFile := `pkg:docker/philips-labs/fatt@sha256:6cc65b2c82c2baa3391890abb8ab741efbcbc87baff3b06d5797afacb314ddd9?repository_url=ghcr.io&attestation_type=sbom
pkg:docker/philips-labs/fatt@sha256:6cc65b2c82c2baa3391890abb8ab741efbcbc87baff3b06d5797afacb314ddd9?repository_url=ghcr.io&attestation_type=provenance
pkg:docker/philips-labs/fatt@sha256:6cc65b2c82c2baa3391890abb8ab741efbcbc87baff3b06d5797afacb314ddd9?repository_url=ghcr.io
pkg:nuget/philips-labs/fatt@sha256:6cc65b2c82c2baa3391890abb8ab741efbcbc87baff3b06d5797afacb314ddd9?repository_url=nuget.org&attestation_type=provenance
pkg:nuget/philips-labs/fatt@sha256:823413cc65b2c82c2baa3391890abb8ab741e87baff3b06d5797afacb314ddd9?repository_url=nuget.org&attestation_type=sbom`

	testCases := []struct {
		name           string
		filter         string
		resultCount    int
		expectedErrMsg string
	}{
		{name: "sbom filter", filter: `{ .IsAttestationType("sbom") }`, resultCount: 2},
		{name: "provenance filter", filter: `{ .IsAttestationType("provenance") }`, resultCount: 2},
		{name: "empty filter", filter: "", resultCount: 5},
		{name: "docker type filter", filter: `{ .PURL.Type == "docker" }`, resultCount: 3},
		{name: "nuget filter", filter: `{ .IsRegistry("nuget.org") }`, resultCount: 2},
		{name: "nuget sbom filter", filter: `{ .IsRegistry("nuget.org") && .IsAttestationType("sbom") }`, resultCount: 1},
		{name: "invalid filter", filter: `{ .Version == "v0.3.1" }`, resultCount: 0, expectedErrMsg: "cannot fetch Version from attestation.Attestation (1:25)\n | filter(Attestations, { .Version == \"v0.3.1\" })\n | ........................^"},
		{name: "invalid expression", filter: "{}", resultCount: 0, expectedErrMsg: "unexpected token Bracket(\"}\") (1:23)\n | filter(Attestations, {})\n | ......................^"},
	}

	atts, err := txt.ReadAttestations(strings.NewReader(purlsFile))
	assert.NoError(t, err)

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			assert := assert.New(tt)
			results, err := attestation.Reduce(atts, tc.filter)

			if tc.expectedErrMsg != "" {
				assert.Error(err)
				assert.EqualError(err, tc.expectedErrMsg)
				assert.Nil(results)
			} else {
				assert.NoError(err)
				assert.Len(results, tc.resultCount)
			}
		})
	}
}
