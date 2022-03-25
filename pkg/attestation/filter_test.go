package attestation_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/philips-labs/fatt/pkg/attestation"
	"github.com/philips-labs/fatt/pkg/attestation/resolvers/txt"
)

func TestReducePurls(t *testing.T) {
	purlsFile := `pkg:docker/philips-labs/fatt@sha256:6cc65b2c82c2baa3391890abb8ab741efbcbc87baff3b06d5797afacb314ddd9?repository_url=ghcr.io%2Fphilips-labs%2Ffatt&tag=v0.1.0
pkg:oci/philips-labs/fatt-cli@sha256:5f74adb7a68a89d96aa3571d93b2e96c68b29891e4c646e5d31f25b3d6eb308a?repository_url=ghcr.io%2Fphilips-labs%2Ffatt-cli&tag=v0.1.0.sbom
pkg:oci/philips-labs/fatt-cli@sha256:bc59775e7ec4983e1b79aa0234ee1020f566eedf01b53e99d44645c12c402826?repository_url=ghcr.io%2Fphilips-labs%2Ffatt-cli&tag=v0.1.0.provenance
pkg:oci/philips-labs/fatt-cli@sha256:5f74adb7a68a89d96aa3571d93b2e96c68b29891e4c646e5d31f25b3d6eb308a?repository_url=index.docker.io%2Fphilipssoftware%2Ffatt-cli&tag=v0.1.0.sbom
pkg:oci/philips-labs/fatt-cli@sha256:bc59775e7ec4983e1b79aa0234ee1020f566eedf01b53e99d44645c12c402826?repository_url=index.docker.io%2Fphilipssoftware%2Ffatt-cli&tag=v0.1.0.provenance`

	testCases := []struct {
		name           string
		filter         string
		resultCount    int
		expectedErrMsg string
	}{
		{name: "empty filter", filter: "", resultCount: 5},
		{name: "sbom filter", filter: `{ .IsAttestationType("sbom") }`, resultCount: 2},
		{name: "provenance filter", filter: `{ .IsAttestationType("provenance") }`, resultCount: 2},
		{name: "docker type filter", filter: `{ .PURL.Type == "docker" }`, resultCount: 1},
		{name: "oci type filter", filter: `{ .PURL.Type == "oci" }`, resultCount: 4},
		{name: "nuget filter", filter: `{ .IsRegistry("index.docker.io") }`, resultCount: 2},
		{name: "nuget sbom filter", filter: `{ .IsRegistry("index.docker.io") && .IsAttestationType("sbom") }`, resultCount: 1},
		{name: "invalid filter", filter: `{ .Version == "v0.3.1" }`, resultCount: 0, expectedErrMsg: "cannot fetch Version from attestation.Attestation (1:25)\n | filter(Attestations, { .Version == \"v0.3.1\" })\n | ........................^"},
		{name: "invalid expression", filter: "{}", resultCount: 0, expectedErrMsg: "unexpected token Bracket(\"}\") (1:23)\n | filter(Attestations, {})\n | ......................^"},
	}

	atts, err := (&txt.Resolver{}).Resolve(strings.NewReader(purlsFile))
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
