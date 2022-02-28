package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/philips-labs/fatt/cmd/fatt/cli/options"
	"github.com/philips-labs/fatt/pkg/attestation"
	"github.com/philips-labs/fatt/pkg/oci"
)

var (
	lo = &options.ListOptions{}
)

// NewListCommand creates a new instance of a list command
func NewListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all attestations",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(os.Stderr, "Fetching attestations for current working directory…")

			if ro.FilePath == "" {
				d, err := os.Getwd()
				if err != nil {
					return err
				}
				ro.FilePath = d
			}

			r, err := ro.GetResolver()
			if err != nil {
				return err
			}

			atts, err := r.Resolve(ro.FilePath)
			if err != nil {
				return fmt.Errorf("failed to resolve attestations: %w", err)
			}

			var p attestation.Printer
			switch lo.OutputFormat {
			case "docker":
				p = oci.NewDockerPrinter(os.Stdout)
			default:
				p = attestation.NewDefaultPrinter(os.Stdout)
			}
			return p.Print(atts)
		},
	}

	lo.AddFlags(cmd)

	return cmd
}
