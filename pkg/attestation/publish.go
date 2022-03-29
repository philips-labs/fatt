package attestation

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/package-url/packageurl-go"
	cremote "github.com/sigstore/cosign/pkg/cosign/remote"

	"github.com/philips-labs/fatt/pkg/oci"
)

// File extends cremote.File by adding Scheme
type File interface {
	cremote.File
	Scheme() string
}

// PublishResult captures the result after publishing the attestations
type PublishResult struct {
	OCIRef string
	PURL   *packageurl.PackageURL
}

type file struct {
	scheme string
	path   string
}

// Scheme implements AttestationFile
func (f file) Scheme() string {
	return f.scheme
}

// Contents implements AttestationFile
func (f file) Contents() ([]byte, error) {
	return os.ReadFile(f.path)
}

// Path implements AttestationFile
func (f file) Path() string {
	return f.path
}

// Platform implements AttestationFile
func (file) Platform() *v1.Platform {
	return nil
}

// String implements AttestationFile
func (f file) String() string {
	return f.scheme + "://" + f.path
}

var _ File = file{}

// ParseFileRef parses a file reference as a AttestationFile
func ParseFileRef(fileRef string) (File, error) {
	refParts := strings.Split(fileRef, "://")

	if len(refParts) != 2 {
		return nil, errors.New("could not parse attestation scheme, use <scheme>://<file> format")
	}

	return &file{
		scheme: refParts[0],
		path:   refParts[1],
	}, nil
}

// Publish publishes the attestations to an oci repository
func Publish(ctx context.Context, repository, tagPrefix string, version string, att cremote.File) (*PublishResult, error) {
	t, err := getType(att.String())
	if err != nil {
		return nil, err
	}

	ociRef := fmt.Sprintf("%s:%s-%s.%s", repository, tagPrefix, version, t)

	ref, err := name.ParseReference(ociRef)
	if err != nil {
		return nil, err
	}

	digestRef, err := uploadBlob(ctx, att, ref)
	if err != nil {
		return nil, err
	}

	purl, err := oci.ToPackageURL(ref, digestRef)
	if err != nil {
		return nil, err
	}

	return &PublishResult{
		OCIRef: ociRef,
		PURL:   purl,
	}, nil
}

func getType(att string) (string, error) {
	prov := "provenance"
	sbom := "sbom"
	discovery := "discovery"

	switch {
	case strings.HasPrefix(att, prov+"://"):
		return prov, nil
	case strings.HasPrefix(att, sbom+"://"):
		return sbom, nil
	case strings.HasPrefix(att, discovery+"://"):
		return discovery, nil
	default:
		return "", errors.New("currently only sbom:// and provenance:// schemes are supported")
	}
}

func uploadBlob(ctx context.Context, file cremote.File, ref name.Reference) (name.Digest, error) {
	mt := cremote.DefaultMediaTypeGetter
	opts := []remote.Option{
		remote.WithAuthFromKeychain(authn.DefaultKeychain),
		remote.WithContext(ctx),
	}

	return cremote.UploadFiles(ref, []cremote.File{file}, mt, opts...)
}
