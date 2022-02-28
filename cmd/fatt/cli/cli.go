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
	"github.com/philips-labs/fatt/pkg/oci"
)

const (
	cliName = "fatt"
)

var (
	ro = &options.RootOptions{}
)

// New create a new instance of the fatt commandline interface
func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   cliName,
		Short: "Fetches an attestation",
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
			default:
				fmt.Fprintln(os.Stderr, "unsupported resolver, supported resolvers are `txt` and `packagejson`.")
			}

			atts, err := r.Resolve(ro.FilePath)
			if err != nil {
				return fmt.Errorf("failed to resolve attestations: %w", err)
			}

			fmt.Fprintf(os.Stderr, "Found attestations: %+v\n", atts)
			for _, att := range atts {
				purl := att.PURL
				if attType, ok := purl.Qualifiers.Map()["attestation_type"]; ok {
					fmt.Fprintf(os.Stderr, "Attestation type: %s\n", attType)
				}
				switch att.PURL.Type {
				case "docker":
					fmt.Fprintln(os.Stdout, oci.ImageURLFromPURL(purl))
				default:
					fmt.Fprintln(os.Stderr, "Unsupported purl type")
				}
			}

			return nil
		},
	}

	ro.AddFlags(cmd)

	return cmd
}
