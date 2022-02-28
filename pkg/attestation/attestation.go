package attestation

import (
	"strings"

	"github.com/package-url/packageurl-go"
)

//go:generate stringer -type=Type
// Type enumeration of attestation types
type Type int

const (
	// Unknown unknown attestation
	Unknown Type = iota
	// SBOM an sbom attestation
	SBOM
	// Provenance a provenance attestation
	Provenance
)

// TypeFromString takes an attestation type as string and return the Type
func TypeFromString(text string) Type {
	switch strings.ToLower(text) {
	case "sbom":
		return SBOM
	case "provenance":
		return Provenance
	default:
		return Unknown
	}
}

// UnmarshalText unmarshals the type from a text form
func (t *Type) UnmarshalText(text []byte) error {
	*t = TypeFromString(string(text))
	return nil
}

// Attestation an attestation url
type Attestation struct {
	PURL packageurl.PackageURL
	Type Type
}
