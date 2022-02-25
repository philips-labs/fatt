package oci

import (
	"fmt"
	"strings"

	"github.com/package-url/packageurl-go"
)

func ImageURLFromPURL(purl packageurl.PackageURL) string {
	ns := purl.Namespace
	v := purl.Version

	if repo, ok := purl.Qualifiers.Map()["repository_url"]; ok {
		if strings.Contains(v, "sha") {
			return fmt.Sprintf("%s/%s/%s@%s", repo, ns, purl.Name, v)
		}
		return fmt.Sprintf("%s/%s/%s:%s", repo, ns, purl.Name, v)
	}

	if strings.Contains(v, "sha") {
		return fmt.Sprintf("%s/%s@%s", ns, purl.Name, v)
	}
	return fmt.Sprintf("%s/%s:%s", ns, purl.Name, v)
}
