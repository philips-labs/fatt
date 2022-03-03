package attestation

import (
	"strings"

	"github.com/package-url/packageurl-go"
)

type PackageUrl packageurl.PackageURL

// IsAttestationType indicates if a PURL has the given attestation type.
func (purl PackageUrl) IsAttestationType(aType string) bool {
	if attType, ok := purl.Qualifiers.Map()["attestation_type"]; ok {
		return strings.ToLower(attType) == aType
	}
	return false
}

// PackageTypeOf filters on a specific package type.
// The package "type" or package "protocol" such as maven, npm, nuget, gem, pypi.
func (purl PackageUrl) PackageTypeOf(pType string) bool {
	return purl.Type == pType
}

// NamespaceOf filters on a specific package namespace.
func (purl PackageUrl) NamespaceOf(pNamespace string) bool {
	return purl.Namespace == pNamespace
}

// NameOf filters on a specific package name.
func (purl PackageUrl) NameOf(pName string) bool {
	return purl.Name == pName
}

// VersionOf filters on a specific package version.
func (purl PackageUrl) VersionOf(pVersion string) bool {
	return purl.Version == pVersion
}

// RepositoryUrlOf filters on a specific repository_url qualifier of the package.
func (purl PackageUrl) RepositoryUrlOf(pRepoUrl string) bool {
	if attType, ok := purl.Qualifiers.Map()["repository_url"]; ok {
		return strings.ToLower(attType) == pRepoUrl
	}
	return false
}

// ConvertFiltersToListProjectOptions converts the filter expressions to ListProjectOptions
func ConvertFiltersToFilterPURLOptions(filter string) {

	if strings.Contains(filter, ".RepositoryUrlOf()") {
		purlFilters = append(purlFilters, RepositoryUrlOf)
	} else if strings.Contains(filter, ".VersionOf()") {
		purlFilters = append(purlFilters, VersionOf)
	} else if strings.Contains(filter, ".NameOf()") {
		purlFilters = append(purlFilters, NameOf)
	} else if strings.Contains(filter, ".NamespaceOf()") {
		purlFilters = append(purlFilters, NamespaceOf)
	} else if strings.Contains(filter, ".PackageTypeOf()") {
		purlFilters = append(purlFilters, PackageTypeOf)
	}

	return purlFilters
}
