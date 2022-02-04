package options

import "github.com/spf13/cobra"

type RootOptions struct{}

var _ CommandFlagger = (*RootOptions)(nil)

func (o *RootOptions) AddFlags(cmd *cobra.Command) {

}
