package attestation

import (
	"context"
	"testing"

	v1 "github.com/google/go-containerregistry/pkg/v1"
	cremote "github.com/sigstore/cosign/pkg/cosign/remote"
	"github.com/stretchr/testify/assert"
)

type mockFile struct {
	scheme string
	path   string
}

// Scheme implements AttestationFile
func (f mockFile) Scheme() string {
	return f.scheme
}

// Contents implements AttestationFile
func (f mockFile) Contents() ([]byte, error) {
	return make([]byte, 0), nil
}

// Path implements AttestationFile
func (f mockFile) Path() string {
	return f.path
}

// Platform implements AttestationFile
func (mockFile) Platform() *v1.Platform {
	return nil
}

// String implements AttestationFile
func (f mockFile) String() string {
	return f.scheme + "://" + f.path
}

var _ cremote.File = mockFile{}

func TestBuildOciRef(t *testing.T) {

	testCases := []struct {
		name           string
		repo           string
		version        string
		tagPrefix      string
		attType        string
		expectedResult string
		expectedErrMsg string
	}{
		{
			name:           "with tag-prefix parameter",
			repo:           "ghcr.io/philips-labs/fatt-attestations-example",
			version:        "v0.1.0",
			tagPrefix:      "test-application",
			attType:        "sbom",
			expectedResult: "ghcr.io/philips-labs/fatt-attestations-example:test-application-v0.1.0.sbom",
		},
		{
			name:           "without tag-prefix parameter",
			repo:           "ghcr.io/philips-labs/fatt-attestations-example",
			version:        "v0.1.0",
			tagPrefix:      "",
			attType:        "sbom",
			expectedResult: "ghcr.io/philips-labs/fatt-attestations-example:v0.1.0.sbom",
		},
		{
			name:           "without version parameter",
			repo:           "ghcr.io/philips-labs/fatt-attestations-example",
			version:        "",
			tagPrefix:      "test-application",
			attType:        "sbom",
			expectedErrMsg: "version is required",
		},
		{
			name:           "without attestation type parameter",
			repo:           "ghcr.io/philips-labs/fatt-attestations-example",
			version:        "v1.0.0",
			tagPrefix:      "test-application",
			expectedErrMsg: "attestation type is required",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			assert := assert.New(tt)
			results, err := buildOciRef(tc.repo, tc.tagPrefix, tc.version, tc.attType)
			if tc.expectedErrMsg != "" {
				assert.Error(err)
				assert.EqualError(err, tc.expectedErrMsg)
				assert.Nil(results)
			} else {
				assert.NoError(err)
				assert.Equal(results.String(), tc.expectedResult)
			}
		})
	}
}
func TestPublish(t *testing.T) {
	sbom := mockFile{
		scheme: "sbom://",
		path:   "something.txt",
	}
	invalidFile := file{
		scheme: "stuff://",
		path:   "stuff-spdx.json",
	}

	testCases := []struct {
		name           string
		repo           string
		version        string
		tagPrefix      string
		file           cremote.File
		expectedErrMsg string
	}{
		{
			name:           "empty args",
			repo:           "",
			version:        "",
			tagPrefix:      "",
			expectedErrMsg: "repository is required",
		},
		{
			name:           "without version but with tag-prefix parameter",
			repo:           "ghcr.io/philips-labs/fatt-attestations-example",
			version:        "",
			tagPrefix:      "test-application",
			file:           sbom,
			expectedErrMsg: "version is required",
		},
		{
			name:           "without file parameter",
			repo:           "ghcr.io/philips-labs/fatt-attestations-example",
			version:        "v0.1.0",
			tagPrefix:      "test-application",
			expectedErrMsg: "attestation file is required",
		},
		{
			name:           "with invalid file parameter",
			repo:           "ghcr.io/philips-labs/fatt-attestations-example",
			version:        "v0.1.0",
			tagPrefix:      "test-application",
			file:           invalidFile,
			expectedErrMsg: "currently only sbom:// and provenance:// schemes are supported",
		},
	}
	ctx := context.Background()
	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			assert := assert.New(tt)
			results, err := Publish(ctx, tc.repo, tc.tagPrefix, tc.version, tc.file)
			if tc.expectedErrMsg != "" {
				assert.Error(err)
				assert.EqualError(err, tc.expectedErrMsg)
				assert.Nil(results)
			} else {
				assert.NoError(err)
			}
		})
	}
}
