package attestation

// Resolver allows to resolve attestations
type Resolver interface {
	Resolve(dir string) ([]Attestation, error)
}
