package txt

import (
	"bufio"
	"io"
	"strings"

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
			Type: getType(purl),
		})
	}

	return atts, nil
}

func getType(p packageurl.PackageURL) attestation.Type {
	if attType, ok := p.Qualifiers.Map()["attestation_type"]; ok {
		switch strings.ToLower(attType) {
		case "provenance":
			return attestation.Provenance
		case "sbom":
			return attestation.SBOM
		default:
			return attestation.Unknown
		}
	}

	return attestation.Unknown
}
