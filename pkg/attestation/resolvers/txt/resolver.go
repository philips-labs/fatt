package txt

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/package-url/packageurl-go"

	"github.com/philips-labs/fatt/pkg/attestation"
)

// Resolver resolves the attestations via txt
type Resolver struct{}

var _ attestation.Resolver = (*Resolver)(nil)

// Resolve resolves the attestations via txt
func (r *Resolver) Resolve(dir string) ([]attestation.Attestation, error) {
	if _, err := os.Stat(dir); err != nil {
		return nil, err
	}

	atts := make([]attestation.Attestation, 0)

	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return nil
		}

		if !info.IsDir() && info.Name() == "attestations.txt" {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			prls, err := ReadAttestations(file)
			if err != nil {
				return err
			}
			atts = append(atts, prls...)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return atts, nil
}

// ReadAttestations reads an io.reader and scans for attestations.
func ReadAttestations(r io.Reader) ([]attestation.Attestation, error) {
	atts := make([]attestation.Attestation, 0)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		purl, err := packageurl.FromString(scanner.Text())
		if err != nil {
			return nil, err
		}
		atts = append(atts, attestation.Attestation{
			PURL: purl,
			Type: attestation.SBOM,
		})
	}
	return atts, nil
}
