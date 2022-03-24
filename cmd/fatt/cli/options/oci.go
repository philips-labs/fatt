package options

import (
	"github.com/spf13/cobra"
)

// OCIOptions commandline options to fetch from oci registry
type OCIOptions struct {
	KeyRef string
}

var _ CommandFlagger = (*OCIOptions)(nil)

// AddFlags implements CommandFlagger to add the RootOptions as flags to the given command
func (o *OCIOptions) AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&o.KeyRef, "key", "", "", "path to the public key file, URL, or KMS URI")
}
