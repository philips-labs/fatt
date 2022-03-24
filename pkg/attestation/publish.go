package attestation

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/package-url/packageurl-go"
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
	purl, err := toPackageURL(ociRef)
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

// https://github.com/package-url/purl-spec/blob/master/PURL-TYPES.rst#oci
func toPackageURL(ociRef string) (*packageurl.PackageURL, error) {
	ref, err := name.ParseReference(ociRef)
	if err != nil {
		return nil, err
	}

	rs := ref.Context().RepositoryStr()
	ns := rs[:strings.LastIndex(rs, "/")]
	n := rs[strings.LastIndex(rs, "/")+1:]
	v := ref.Identifier() //TODO get digest to comply with purl spec

	q := packageurl.QualifiersFromMap(map[string]string{
		"repository_url": fmt.Sprintf("%s/%s", ref.Context().RegistryStr(), ref.Context().RepositoryStr()),
		"tag":            ref.Identifier(),
	})

	return packageurl.NewPackageURL("oci", ns, n, v, q, ""), nil
}
