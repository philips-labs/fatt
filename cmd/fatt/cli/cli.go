package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/philips-labs/fatt/cmd/fatt/cli/options"
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

			return nil
		},
	}

	ro.AddFlags(cmd)

	return cmd
}
