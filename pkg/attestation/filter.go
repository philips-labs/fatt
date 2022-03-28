package attestation

import (
	"fmt"
	"strings"

	"github.com/antonmedv/expr"
)

// FilteredEnv holds the context to apply filter expressions.
type FilteredEnv struct {
	Attestations []Attestation
}

// IsAttestationType indicates if a PURL has the given attestation type.
func (a Attestation) IsAttestationType(t string) bool {
	if tag, ok := a.PURL.Qualifiers.Map()["tag"]; ok {
		return strings.HasSuffix(tag, "."+t)
	}
	return false
}

// IsRegistry filters on a specific repository_url qualifier of the package.
func (a Attestation) IsRegistry(registryURL string) bool {
	if repoURL, ok := a.PURL.Qualifiers.Map()["repository_url"]; ok {
		return strings.HasPrefix(repoURL, registryURL+"/")
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

	result, err := expr.Run(program, FilteredEnv{atts})
	if err != nil {
		return nil, err
	}

	var filteredResults []Attestation
	for _, a := range result.([]interface{}) {
		filteredResults = append(filteredResults, a.(Attestation))
	}
	return filteredResults, nil
}
