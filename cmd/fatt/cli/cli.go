package cli

import (
	"github.com/spf13/cobra"

	"github.com/philips-labs/fatt/cmd/fatt/cli/options"
)

const (
	cliName = "fatt"
)

var (
	ro = options.NewRootOptions()
)

// New create a new instance of the fatt commandline interface
func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   cliName,
		Short: "Discover and resolve your attestations",
	}

	cmd.AddCommand(NewListCommand())
	cmd.AddCommand(NewPublishCommand())
	cmd.AddCommand(NewVersionCommand())

	ro.AddFlags(cmd)

	return cmd
}
