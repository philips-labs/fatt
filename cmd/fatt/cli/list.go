package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/philips-labs/fatt/cmd/fatt/cli/options"
	"github.com/philips-labs/fatt/pkg/attestation"
	"github.com/philips-labs/fatt/pkg/attestation/resolvers/packagejson"
	"github.com/philips-labs/fatt/pkg/attestation/resolvers/txt"
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

			var r attestation.Resolver
			switch strings.ToLower(ro.Resolver) {
			case "txt":
				r = &txt.Resolver{}
			case "packagejson":
				r = &packagejson.Resolver{}
			case "multi":
				r = attestation.NewMultiResolver(
					&txt.Resolver{},
					&packagejson.Resolver{},
				)
			default:
				fmt.Fprintln(os.Stderr, "unsupported resolver, supported resolvers are `txt`, `packagejson`, `multi`.")
			}

			atts, err := r.Resolve(ro.FilePath)
			if err != nil {
				return fmt.Errorf("failed to resolve attestations: %w", err)
			}

			p := attestation.NewDefaultPrinter(os.Stdout)
			return p.Print(atts)
		},
	}

	lo.AddFlags(cmd)

	return cmd
}
