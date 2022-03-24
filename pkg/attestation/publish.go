package attestation

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/package-url/packageurl-go"

	"github.com/philips-labs/fatt/pkg/oci"
)

// PublishResult captures the result after publishing the attestations
type PublishResult struct {
	AttestationFile string
	OCIRef          string
	PURL            *packageurl.PackageURL
}

// Publish publishes the attestations to an oci repository
func Publish(repository, version, attestationRef string) (*PublishResult, error) {
	t, err := getType(attestationRef)
	if err != nil {
		return nil, err
	}

	ociRef := fmt.Sprintf("%s:%s.%s", repository, version, t)

	ref, err := name.ParseReference(ociRef)
	if err != nil {
		return nil, err
	}

	purl, err := oci.ToPackageURL(ref)
	if err != nil {
		return nil, err
	}

	fileName := strings.Split(attestationRef, "://")[1]

	return &PublishResult{
		AttestationFile: fileName,
		OCIRef:          ociRef,
		PURL:            purl,
	}, nil
}

func getType(attestationRef string) (string, error) {
	prov := "provenance"
	sbom := "sbom"

	switch {
	case strings.HasPrefix(attestationRef, prov+"://"):
		return prov, nil
	case strings.HasPrefix(attestationRef, sbom+"://"):
		return sbom, nil
	default:
		return "", errors.New("could not parse attestation scheme, use <scheme>://<file> format")
	}
}
