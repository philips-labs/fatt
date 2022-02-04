package options

import "github.com/spf13/cobra"

type RootOptions struct {
	FilePath string
}

var _ CommandFlagger = (*RootOptions)(nil)

func (o *RootOptions) AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&o.FilePath, "file-path", "p", "", "the filepath to find attestation purls (defaults to current working dir)")
}
