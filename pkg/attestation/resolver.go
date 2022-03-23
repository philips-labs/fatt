package attestation

import (
	"io"
)

// Resolver allows to resolve attestations
type Resolver interface {
	Resolve(io.ReadCloser) ([]Attestation, error)
}
