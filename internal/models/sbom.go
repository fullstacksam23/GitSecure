package models

type SBOMResponse struct {
	SBOM SPDXDocument `json:"sbom"`
}

type SPDXDocument struct {
	SPDXVersion string        `json:"spdxVersion"`
	Name        string        `json:"name"`
	Packages    []SPDXPackage `json:"packages"`
}

type SPDXPackage struct {
	Name         string        `json:"name"`
	SPDXID       string        `json:"SPDXID"`
	VersionInfo  string        `json:"versionInfo"`
	ExternalRefs []ExternalRef `json:"externalRefs"`
}

type ExternalRef struct {
	ReferenceCategory string `json:"referenceCategory"`
	ReferenceType     string `json:"referenceType"`
	ReferenceLocator  string `json:"referenceLocator"`
}
