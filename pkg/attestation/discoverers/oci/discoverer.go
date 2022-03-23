package oci

import (
	"bytes"
	"context"
	"io"

	"github.com/sigstore/cosign/pkg/sget"

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
	wc := &bytes.Buffer{}
	err := sget.New(blobRef, r.keyRef, wc).Do(r.context)

	return wc, err
}
