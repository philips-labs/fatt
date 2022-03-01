package options

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/philips-labs/fatt/pkg/attestation"
	"github.com/philips-labs/fatt/pkg/attestation/resolvers/packagejson"
	"github.com/philips-labs/fatt/pkg/attestation/resolvers/txt"
)

// RootOptions commandline options for the root command
type RootOptions struct {
	FilePath string
	Resolver string
}

var _ CommandFlagger = (*RootOptions)(nil)

// AddFlags implements CommandFlagger to add the RootOptions as flags to the given command
func (o *RootOptions) AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&o.FilePath, "file-path", "p", "", "the filepath to find attestation purls (defaults to current working dir)")
	cmd.PersistentFlags().StringVarP(&o.Resolver, "resolver", "r", "multi", "the resolver to use for finding attestations")
}

// GetResolver returns the resolver based on the resolver commmandline options
func (o *RootOptions) GetResolver() (attestation.Resolver, error) {
	var r attestation.Resolver

	switch strings.ToLower(o.Resolver) {
	case "txt":
		r = &txt.Resolver{}
	case "packagejson":
		r = &packagejson.Resolver{}
	case "multi":
		r = attestation.NewMultiResolver(
			&txt.Resolver{},
			&packagejson.Resolver{},
		)
	default:
		return nil, fmt.Errorf("unsupported resolver, supported resolvers are `txt`, `packagejson`, `multi`")
	}

	return r, nil
}
