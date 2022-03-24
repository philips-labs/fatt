package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/package-url/packageurl-go"
	"github.com/spf13/cobra"

	"github.com/philips-labs/fatt/cmd/fatt/cli/options"
	"github.com/philips-labs/fatt/pkg/attestation"
)

// AttestationsTXT implements the Stringer interface to write the package urls
// to a newline separated file
type AttestationsTXT []*packageurl.PackageURL

func (a AttestationsTXT) String() string {
	var b strings.Builder
	for _, p := range a {
		b.WriteString(p.String() + "\n")
	}
	return b.String()
}

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

			_, cancel := context.WithCancel(cmd.Context())
			defer cancel()

			purls := make(AttestationsTXT, len(po.Attestations))
			for i, att := range po.Attestations {
				r, err := attestation.Publish(po.Repository, po.Version, att)
				if err != nil {
					return err
				}
				fmt.Fprintf(os.Stderr, "cosign upload blob -f %s %s\n", r.AttestationFile, r.OCIRef)
				fmt.Fprintf(os.Stderr, "cosign sign --key %s %s\n", po.KeyRef, r.OCIRef)

				purls[i] = r.PURL
			}

			discoveryOCIRef := fmt.Sprintf("%s:%s.%s", po.Repository, po.Version, "discover")

			fmt.Fprintln(os.Stderr)
			fmt.Fprintln(os.Stderr, "Generating attestations.txt based on uploaded attestations…")
			fmt.Fprintf(os.Stderr, "%s", purls)
			fmt.Fprintln(os.Stderr)
			fmt.Fprintf(os.Stderr, "cosign upload blob -f %s %s\n", "attestations.txt", discoveryOCIRef)
			fmt.Fprintf(os.Stderr, "cosign sign --key %s %s\n", po.KeyRef, discoveryOCIRef)
			return nil
		},
	}

	po.AddFlags(cmd)

	return cmd
}
