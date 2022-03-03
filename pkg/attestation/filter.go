package attestation

import (
	"fmt"
	"strings"

	"github.com/antonmedv/expr"
)

type AttestationFilterEnv struct {
	Attestations []Attestation
}

// IsAttestationType indicates if a PURL has the given attestation type.
func (a Attestation) IsAttestationType(t string) bool {
	if attType, ok := a.PURL.Qualifiers.Map()["attestation_type"]; ok {
		return strings.ToLower(attType) == t
	}
	return false
}

// RepositoryUrlOf filters on a specific repository_url qualifier of the package.
func (a Attestation) RepositoryUrlOf(registryUrl string) bool {
	if attType, ok := a.PURL.Qualifiers.Map()["repository_url"]; ok {
		return strings.ToLower(attType) == registryUrl
	}
	return false
}

// Reduce filters the Attestations based on the given filter
func Reduce(atts []Attestation, filter string) ([]Attestation, error) {
	if strings.TrimSpace(filter) == "" {
		return atts, nil
	}
	program, err := expr.Compile(fmt.Sprintf("filter(Attestations, %s)", filter))
	if err != nil {
		return nil, err
	}

	result, err := expr.Run(program, AttestationFilterEnv{atts})
	if err != nil {
		return nil, err
	}

	var filteredResults []Attestation
	for _, a := range result.([]interface{}) {
		filteredResults = append(filteredResults, a.(Attestation))
	}
	return filteredResults, nil
}
