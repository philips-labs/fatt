package attestation

// Resolver allows to resolve attestations
type Resolver interface {
	Resolve(dir string) ([]Attestation, error)
}

// MultiResolver resolves attestations via the given resolvers
type MultiResolver struct {
	resolvers []Resolver
}

var _ Resolver = (*MultiResolver)(nil)

// NewMultiResolver creates a new MultiResolver using the given resolvers
func NewMultiResolver(resolvers ...Resolver) Resolver {
	return &MultiResolver{
		resolvers: resolvers,
	}
}

// Resolve resolves the attestations
func (r *MultiResolver) Resolve(dir string) ([]Attestation, error) {
	attestations := make([]Attestation, 0)

	for _, rr := range r.resolvers {
		a, err := rr.Resolve(dir)
		if err != nil {
			return nil, err
		}
		attestations = append(attestations, a...)
	}

	return attestations, nil
}
