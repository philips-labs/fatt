package attestation

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/package-url/packageurl-go"
	cremote "github.com/sigstore/cosign/pkg/cosign/remote"

	"github.com/philips-labs/fatt/pkg/oci"
)

// PublishResult captures the result after publishing the attestations
type PublishResult struct {
	AttestationFile string
	OCIRef          string
	PURL            *packageurl.PackageURL
}

// Publish publishes the attestations to an oci repository
func Publish(ctx context.Context, repository, version, attestationRef string) (*PublishResult, error) {
	t, err := getType(attestationRef)
	if err != nil {
		return nil, err
	}

	ociRef := fmt.Sprintf("%s:%s.%s", repository, version, t)

	ref, err := name.ParseReference(ociRef)
	if err != nil {
		return nil, err
	}

	fileName := strings.Split(attestationRef, "://")[1]
	digestRef, err := uploadBlob(ctx, fileName, ref)
	if err != nil {
		return nil, err
	}

	purl, err := oci.ToPackageURL(ref, digestRef)
	if err != nil {
		return nil, err
	}

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

func uploadBlob(ctx context.Context, fileName string, ref name.Reference) (name.Digest, error) {
	mt := cremote.DefaultMediaTypeGetter
	opts := []remote.Option{
		remote.WithAuthFromKeychain(authn.DefaultKeychain),
		remote.WithContext(ctx),
	}

	file := cremote.FileFromFlag(fileName)
	return cremote.UploadFiles(ref, []cremote.File{file}, mt, opts...)
}
