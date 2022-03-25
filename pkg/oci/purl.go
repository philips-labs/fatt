package oci

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/package-url/packageurl-go"
)

// FromPackageURL transforms a Package URL to an oci reference
// https://github.com/package-url/purl-spec/blob/master/PURL-TYPES.rst#oci
func FromPackageURL(purl packageurl.PackageURL) (name.Reference, error) {
	ns := purl.Namespace
	v := purl.Version

	if purl.Type == "docker" {
		fmt.Fprintln(os.Stderr, "package url type 'docker' is deprecated, please use package url type 'oci'")
	} else if purl.Type != "oci" {
		return nil, errors.New("unable to handle non oci package url")
	}

	var ociRef string
	if repo, ok := purl.Qualifiers.Map()["repository_url"]; ok {
		ociRef = fmt.Sprintf("%s@%s", repo, v)
		// TODO: Restore this logic when signatures are implemented in publish command.
		// if tag, ok := purl.Qualifiers.Map()["tag"]; ok {
		// 	ociRef = fmt.Sprintf("%s:%s", repo, tag)
		// } else {
		// 	ociRef = fmt.Sprintf("%s@%s", repo, v)
		// }
	} else if strings.Contains(v, "sha") {
		ociRef = fmt.Sprintf("%s/%s@%s", ns, purl.Name, v)
	} else {
		ociRef = fmt.Sprintf("%s/%s:%s", ns, purl.Name, v)
	}

	// removes the name.DefaultRegistry or name.defaultRegistryAlias
	// and the library part (e.g. alpine image)
	r := regexp.MustCompile("^(index.|)docker.io/(library/|)")
	ociRef = r.ReplaceAllString(ociRef, "")

	return name.ParseReference(ociRef)
}

// ToPackageURL transforms an oci reference to Package URL format
// https://github.com/package-url/purl-spec/blob/master/PURL-TYPES.rst#oci
func ToPackageURL(ref, digestRef name.Reference) (*packageurl.PackageURL, error) {
	rs := ref.Context().RepositoryStr()
	ns := rs[:strings.LastIndex(rs, "/")]
	n := rs[strings.LastIndex(rs, "/")+1:]
	v := digestRef.Identifier()

	q := packageurl.QualifiersFromMap(map[string]string{
		"repository_url": fmt.Sprintf("%s/%s", ref.Context().RegistryStr(), ref.Context().RepositoryStr()),
		"tag":            ref.Identifier(),
	})

	return packageurl.NewPackageURL("oci", ns, n, v, q, ""), nil
}
