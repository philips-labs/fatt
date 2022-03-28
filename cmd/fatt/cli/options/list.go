package options

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/philips-labs/fatt/pkg/attestation"
	"github.com/philips-labs/fatt/pkg/attestation/discoverers/fs"
	"github.com/philips-labs/fatt/pkg/attestation/discoverers/oci"
	"github.com/philips-labs/fatt/pkg/attestation/resolvers/txt"
	"github.com/philips-labs/fatt/pkg/print"
)

// ListOptions commandline options for the list command
type ListOptions struct {
	*OCIOptions
	FilePath     string
	OutputFormat string
	Filter       string
}

// NewListOptions initializes the ListOptions object
func NewListOptions() *ListOptions {
	return &ListOptions{OCIOptions: &OCIOptions{}}
}

var _ CommandFlagger = (*ListOptions)(nil)

// AddFlags implements CommandFlagger to add the RootOptions as flags to the given command
func (o *ListOptions) AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&o.OutputFormat, "output-format", "o", "purl", "output format for the list")
	cmd.PersistentFlags().StringVarP(&o.Filter, "filter", "f", "", "filter attestations using template expressions")
	o.OCIOptions.AddFlags(cmd)
}

// GetPrinter returns the printer based on the OutputFormat flag
func (o *ListOptions) GetPrinter(w io.Writer) (attestation.Printer, error) {
	var p attestation.Printer

	switch o.OutputFormat {
	case "docker":
		fmt.Fprintln(os.Stderr, "output-format 'docker' is deprecated, please use output-format 'oci'")
		p = print.NewDockerPrinter(w)
	case "oci":
		p = print.NewDockerPrinter(w)
	default:
		p = attestation.NewDefaultPrinter(w)
	}

	return p, nil
}

// GetResolver returns the resolver based on the resolver commmandline options
func (o *ListOptions) GetResolver() (attestation.Resolver, error) {
	return &txt.Resolver{}, nil
}

// GetDiscoverer discovers attestation.txt files from given location
func (o *ListOptions) GetDiscoverer(ctx context.Context) (attestation.Discoverer, error) {
	if _, err := os.Stat(o.FilePath); err == nil {
		return &fs.Discoverer{}, nil
	}
	return oci.NewDiscoverer(o.KeyRef, oci.WithContext(ctx)), nil
}
