package oci

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/sigstore/cosign/v2/cmd/cosign/cli/fulcio"
	"github.com/sigstore/cosign/v2/cmd/cosign/cli/options"
	"github.com/sigstore/cosign/v2/pkg/cosign"
	ociremote "github.com/sigstore/cosign/v2/pkg/oci/remote"
	sigs "github.com/sigstore/cosign/v2/pkg/signature"

	"github.com/philips-labs/fatt/pkg/attestation"
)

// DiscoverOption allows to configure the Discoverer
type DiscoverOption func(*Discoverer)

// WithContext sets the given context on the Discoverer
func WithContext(ctx context.Context) DiscoverOption {
	return func(d *Discoverer) {
		d.context = ctx
	}
}

// Discoverer discovers the attestations from an oci reference
type Discoverer struct {
	keyRef  string
	context context.Context
}

var _ attestation.Discoverer = (*Discoverer)(nil)

// NewDiscoverer creates a new instance of a Discoverer
func NewDiscoverer(keyRef string, options ...DiscoverOption) *Discoverer {
	d := &Discoverer{keyRef: keyRef}
	for _, o := range options {
		o(d)
	}

	if d.context == nil {
		d.context = context.Background()
	}

	return d
}

// Discover discovers an attestations.txt from an oci registry
func (r *Discoverer) Discover(blobRef string) (io.Reader, error) {
	// Mojority of implementation taken from https://github.com/sigstore/cosign/blob/main/pkg/sget/sget.go

	buf := &bytes.Buffer{}

	ref, err := name.ParseReference(blobRef)
	if err != nil {
		return nil, err
	}

	opts := []remote.Option{
		remote.WithAuthFromKeychain(authn.DefaultKeychain),
		remote.WithContext(r.context),
	}

	co := &cosign.CheckOpts{
		ClaimVerifier:      cosign.SimpleClaimVerifier,
		RegistryClientOpts: []ociremote.Option{ociremote.WithRemoteOptions(opts...)},
	}
	if _, ok := ref.(name.Tag); ok {
		if r.keyRef == "" && !options.EnableExperimental() {
			return nil, errors.New("public key must be specified when fetching by tag, you must fetch by digest or supply a public key")
		}
	}

	ref, err = ociremote.ResolveDigest(ref, co.RegistryClientOpts...)
	if err != nil {
		return nil, err
	}

	if r.keyRef != "" {
		pub, err := sigs.LoadPublicKey(r.context, r.keyRef)
		if err != nil {
			return nil, err
		}
		co.SigVerifier = pub
	}

	if co.SigVerifier != nil || options.EnableExperimental() {
		co.RootCerts, err = fulcio.GetRoots()
		if err != nil {
			return nil, err
		}

		fmt.Fprintf(os.Stderr, "Verifying signature for %s…\n", ref)
		_, _, err := cosign.VerifyImageSignatures(r.context, ref, co)
		if err != nil {
			return nil, err
		}
	}

	img, err := remote.Image(ref, opts...)
	if err != nil {
		return nil, err
	}
	layers, err := img.Layers()
	if err != nil {
		return nil, err
	}
	if len(layers) != 1 {
		return nil, errors.New("invalid artifact")
	}
	rc, err := layers[0].Compressed()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	_, err = io.Copy(buf, rc)

	return buf, err
}
