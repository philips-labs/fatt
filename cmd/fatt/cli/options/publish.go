package options

import (
	"github.com/spf13/cobra"
)

// PublishOptions commandline options for the list command
type PublishOptions struct {
	*OCIOptions
	Version      string
	Attestations []string
}

// NewPublishOptions initializes the ListOptions object
func NewPublishOptions() *PublishOptions {
	return &PublishOptions{OCIOptions: &OCIOptions{}}
}

var _ CommandFlagger = (*PublishOptions)(nil)

// AddFlags implements CommandFlagger to add the RootOptions as flags to the given command
func (o *PublishOptions) AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&o.Version, "version", "", "", "the version to publish the attestations for.")
	o.OCIOptions.AddFlags(cmd)
}
