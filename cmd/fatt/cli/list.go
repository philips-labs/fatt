package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/philips-labs/fatt/cmd/fatt/cli/options"
	"github.com/philips-labs/fatt/pkg/attestation"
)

// NewListCommand creates a new instance of a list command
func NewListCommand() *cobra.Command {
	lo := options.NewListOptions()

	cmd := &cobra.Command{
		Use:   "list <discovery-path>",
		Short: "Lists all attestations",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				d, err := os.Getwd()
				if err != nil {
					return err
				}
				lo.FilePath = d
			} else {
				lo.FilePath = args[0]
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(os.Stderr, "Fetching attestations from %s…\n", lo.FilePath)
			ctx, cancel := context.WithCancel(cmd.Context())
			defer cancel()

			d, err := lo.GetDiscoverer(ctx)
			if err != nil {
				return err
			}

			attReader, err := d.Discover(lo.FilePath)
			if err != nil {
				return err
			}

			r, err := lo.GetResolver()
			if err != nil {
				return err
			}

			atts, err := r.Resolve(attReader)
			if err != nil {
				return fmt.Errorf("failed to resolve attestations: %w", err)
			}

			if lo.Filter != "" {
				atts, err = attestation.Reduce(atts, lo.Filter)
				if err != nil {
					return err
				}
			}

			p, err := lo.GetPrinter(os.Stdout)
			if err != nil {
				return err
			}

			return p.Print(atts)
		},
	}

	lo.AddFlags(cmd)

	return cmd
}
