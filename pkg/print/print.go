package print

import (
	"fmt"
	"io"

	"github.com/philips-labs/fatt/pkg/attestation"
	"github.com/philips-labs/fatt/pkg/oci"
)

// DockerPrinter prints the attestations in OCI format
type DockerPrinter struct {
	w io.Writer
}

var _ attestation.Printer = (*DockerPrinter)(nil)

// NewDockerPrinter creates a new DockerPrinter instance utilizing the given io.Writer
func NewDockerPrinter(w io.Writer) *DockerPrinter {
	return &DockerPrinter{w}
}

// Print prints the attestations to the io.Writer
func (p *DockerPrinter) Print(atts []attestation.Attestation) error {
	for _, att := range atts {
		purl := att.PURL
		ref, err := oci.FromPackageURL(purl)
		if err != nil {
			return err
		}
		if _, err := fmt.Fprintln(p.w, ref.String()); err != nil {
			return err
		}
	}
	return nil
}
