package attestation

// Resolver allows to resolve attestations
type Resolver interface {
	Resolve(dir string) ([]Attestation, error)
}

type MultiResolver struct {
	resolvers []Resolver
}

var _ Resolver = (*MultiResolver)(nil)

func NewMultiResolver(resolvers ...Resolver) Resolver {
	return &MultiResolver{
		resolvers: resolvers,
	}
}

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
