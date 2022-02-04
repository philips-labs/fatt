package attestation

import (
	"strings"

	"github.com/package-url/packageurl-go"
)

//go:generate stringer -type=Type
type Type int

const (
	Unknown Type = iota
	SBOM
	Provenance
)

func TypeFromString(text string) Type {
	switch strings.ToLower(text) {
	case "sbom":
		return SBOM
	case "Provenance":
		return Provenance
	default:
		return Unknown
	}
}

func (t *Type) UnmarshalText(text []byte) error {
	*t = TypeFromString(string(text))
	return nil
}

type Attestation struct {
	PURL packageurl.PackageURL
	Type Type
}
