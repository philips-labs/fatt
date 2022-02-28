package options

import "github.com/spf13/cobra"

// RootOptions commandline options for the root command
type RootOptions struct {
	FilePath string
}

var _ CommandFlagger = (*RootOptions)(nil)

// AddFlags implements CommandFlagger to add the RootOptions as flags to the given command
func (o *RootOptions) AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&o.FilePath, "file-path", "p", "", "the filepath to find attestation purls (defaults to current working dir)")
}
