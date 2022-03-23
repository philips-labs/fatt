package attestation

import "io"

// Discoverer allows to discover attestations
type Discoverer interface {
	Discover(string) (io.ReadCloser, error)
}
