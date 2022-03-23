package txt

import (
	"bufio"
	"io"

	"github.com/package-url/packageurl-go"

	"github.com/philips-labs/fatt/pkg/attestation"
)

// Resolver resolves the attestations via txt
type Resolver struct{}

var _ attestation.Resolver = (*Resolver)(nil)

// Resolve resolves the attestations via txt
func (r *Resolver) Resolve(rc io.Reader) ([]attestation.Attestation, error) {
	atts := make([]attestation.Attestation, 0)
	scanner := bufio.NewScanner(rc)

	for scanner.Scan() {
		purl, err := packageurl.FromString(scanner.Text())
		if err != nil {
			return nil, err
		}
		atts = append(atts, attestation.Attestation{
			PURL: purl,
			Type: attestation.SBOM,
		})
	}

	return atts, nil
}
