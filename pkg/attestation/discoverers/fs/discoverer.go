package fs

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/philips-labs/fatt/pkg/attestation"
)

// Discoverer discovers an attestations.txt from the filesystem
type Discoverer struct{}

var _ attestation.Discoverer = (*Discoverer)(nil)

// Discover discovers an attestations.txt
func (r *Discoverer) Discover(dir string) (io.ReadCloser, error) {
	if _, err := os.Stat(dir); err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
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

			_, err = io.Copy(buf, file)
			if err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return io.NopCloser(buf), nil
}
