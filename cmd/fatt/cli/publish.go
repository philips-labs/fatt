package cli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"

	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/package-url/packageurl-go"
	"github.com/spf13/cobra"

	"github.com/philips-labs/fatt/cmd/fatt/cli/options"
	"github.com/philips-labs/fatt/pkg/attestation"
)

// AttestationsTXT implements the Stringer interface to write the package urls
// to a newline separated file
type AttestationsTXT []*packageurl.PackageURL

// Scheme implements remote.File
func (a AttestationsTXT) Scheme() string {
	return "discovery"
}

// Contents implements remote.File
func (a AttestationsTXT) Contents() ([]byte, error) {
	var b bytes.Buffer
	for _, p := range a {
		b.WriteString(p.String() + "\n")
	}
	return b.Bytes(), nil
}

// Path implements remote.File
func (AttestationsTXT) Path() string {
	return "attestations.txt"
}

// Platform implements remote.File
func (AttestationsTXT) Platform() *v1.Platform {
	return nil
}

// Platform implements remote.File
func (a AttestationsTXT) String() string {
	return a.Scheme() + "://" + a.Path()
}

var _ attestation.File = AttestationsTXT{}

// NewPublishCommand creates a new instance of a publish command
func NewPublishCommand() *cobra.Command {
	po := options.NewPublishOptions()
	cmd := &cobra.Command{
		Use:   "publish <attestations...>",
		Short: "Publishes given attestations to an OCI registry",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("please provide any attestion files to be published")
			}

			po.Attestations = args

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(os.Stderr, "Publishing attestations…")

			ctx, cancel := context.WithCancel(cmd.Context())
			defer cancel()

			purls := make(AttestationsTXT, len(po.Attestations))
			for i, att := range po.Attestations {
				att, err := attestation.ParseFileRef(att)
				if err != nil {
					return err
				}
				r, err := attestation.Publish(ctx, po.Repository, po.TagPrefix, po.Version, att)
				if err != nil {
					return err
				}

				purls[i] = r.PURL
			}

			fmt.Fprintln(os.Stderr)
			fmt.Fprintln(os.Stderr, "Generating attestations.txt based on uploaded attestations…")
			_, err := attestation.Publish(ctx, po.Repository, po.TagPrefix, po.Version, purls)
			if err != nil {
				return err
			}
			fmt.Fprintln(os.Stderr)
			return nil
		},
	}

	po.AddFlags(cmd)

	return cmd
}
