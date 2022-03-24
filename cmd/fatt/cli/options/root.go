package options

import (
	"github.com/spf13/cobra"
)

// RootOptions commandline options for the root command
type RootOptions struct {
}

// NewRootOptions initializes the RootOptions object
func NewRootOptions() *RootOptions {
	return &RootOptions{}
}

var _ CommandFlagger = (*RootOptions)(nil)

// AddFlags implements CommandFlagger to add the RootOptions as flags to the given command
func (o *RootOptions) AddFlags(cmd *cobra.Command) {

}
