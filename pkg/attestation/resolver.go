package attestation

import (
	"io"
)

// Resolver allows to resolve attestations
type Resolver interface {
	Resolve(io.Reader) ([]Attestation, error)
}
