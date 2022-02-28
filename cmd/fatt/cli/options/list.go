package options

import "github.com/spf13/cobra"

// ListOptions commandline options for the list command
type ListOptions struct {
	OutputFormat string
}

var _ CommandFlagger = (*ListOptions)(nil)

// AddFlags implements CommandFlagger to add the RootOptions as flags to the given command
func (o *ListOptions) AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&o.OutputFormat, "output-format", "o", "purl", "output format for the list")
}
