package packagejson

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/package-url/packageurl-go"

	"github.com/philips-labs/fatt/pkg/attestation"
)

// PackageJSON structure to retrieve attestations from a package.json
type PackageJSON struct {
	Name         string   `json:"name,omitempty"`
	Attestations []string `json:"attestations,omitempty"`
}

// Resolver resolves the attestations via packagejson
type Resolver struct {
}

var _ attestation.Resolver = (*Resolver)(nil)

// Resolve resolves attestations from the given directory looking at a package.json
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

		if !info.IsDir() && info.Name() == "package.json" {
			file, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			var pJSON PackageJSON
			err = json.Unmarshal(file, &pJSON)
			if err != nil {
				return err
			}

			for _, p := range pJSON.Attestations {
				purl, err := packageurl.FromString(p)
				if err != nil {
					return err
				}
				atts = append(atts, attestation.Attestation{
					PURL: purl,
					Type: attestation.SBOM,
				})
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return atts, nil
}
