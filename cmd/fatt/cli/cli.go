package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/philips-labs/fatt/cmd/fatt/cli/options"
	"github.com/philips-labs/fatt/pkg/resolver"
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
			fmt.Println("Fetching attestations for current working directory…")

			if ro.FilePath == "" {
				d, err := os.Getwd()
				if err != nil {
					return err
				}
				ro.FilePath = d
			}

			atts, err := resolver.Resolve(ro.FilePath)
			if err != nil {
				return fmt.Errorf("failed to resolve attestations: %w", err)
			}

			for _, att := range atts {
				fmt.Println(att)
			}

			return nil
		},
	}

	ro.AddFlags(cmd)

	return cmd
}
