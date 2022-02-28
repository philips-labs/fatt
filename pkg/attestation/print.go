package attestation

import (
	"fmt"
	"io"
)

// Printer allows to print results
type Printer interface {
	Print(atts []Attestation) error
}

// DefaultPrinter prints the attestation purls
type DefaultPrinter struct {
	w io.Writer
}

var _ Printer = (*DefaultPrinter)(nil)

// NewDefaultPrinter creates a new DefaultPrinter instance utilizing the given io.Writer
func NewDefaultPrinter(w io.Writer) *DefaultPrinter {
	return &DefaultPrinter{w}
}

// Print prints the attestations to the io.Writer
func (p *DefaultPrinter) Print(atts []Attestation) error {
	for _, att := range atts {
		if err := p.PrintAttestation(att); err != nil {
			return err
		}
	}

	return nil
}

// PrintAttestation prints a single attestation to the io.Writer
func (p *DefaultPrinter) PrintAttestation(a Attestation) error {
	_, err := fmt.Fprintln(p.w, a.PURL.String())
	return err
}
