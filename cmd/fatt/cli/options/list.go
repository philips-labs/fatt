package options

import (
	"io"

	"github.com/spf13/cobra"

	"github.com/philips-labs/fatt/pkg/attestation"
	"github.com/philips-labs/fatt/pkg/oci"
)

// ListOptions commandline options for the list command
type ListOptions struct {
	OutputFormat string
	Filter       string
}

var _ CommandFlagger = (*ListOptions)(nil)

// AddFlags implements CommandFlagger to add the RootOptions as flags to the given command
func (o *ListOptions) AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&o.OutputFormat, "output-format", "o", "purl", "output format for the list")
	cmd.PersistentFlags().StringVarP(&o.Filter, "filter", "f", "", "filter attestations using template expressions")
}

// GetPrinter returns the printer based on the OutputFormat flag
func (o *ListOptions) GetPrinter(w io.Writer) (attestation.Printer, error) {
	var p attestation.Printer

	switch o.OutputFormat {
	case "docker":
		p = oci.NewDockerPrinter(w)
	default:
		p = attestation.NewDefaultPrinter(w)
	}

	return p, nil
}
