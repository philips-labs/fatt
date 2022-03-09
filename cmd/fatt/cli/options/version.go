package options

import (
	"github.com/spf13/cobra"
)

// VersionOptions commandline options for the version command
type VersionOptions struct {
	OutputFormat string
}

var _ CommandFlagger = (*VersionOptions)(nil)

// AddFlags implements CommandFlagger to add the RootOptions as flags to the given command
func (o *VersionOptions) AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&o.OutputFormat, "output-format", "o", "", "output format for the version command, valid option is: json")
}
