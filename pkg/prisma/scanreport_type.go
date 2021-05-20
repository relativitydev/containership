// Code generated … DO NOT EDIT.
// https://mholt.github.io/json-to-go/

package prisma

import (
	"time"
)

// ScanReport is the response object the Prisma returns from /registry
type ScanReport struct {
	ID                        string                    `json:"_id"`
	Type                      string                    `json:"type"`
	Hostname                  string                    `json:"hostname"`
	ScanTime                  time.Time                 `json:"scanTime"`
	Binaries                  []Binaries                `json:"binaries"`
	Secrets                   []string                  `json:"Secrets"`
	StartupBinaries           []StartupBinaries         `json:"startupBinaries"`
	OsDistro                  string                    `json:"osDistro"`
	OsDistroRelease           string                    `json:"osDistroRelease"`
	Distro                    string                    `json:"distro"`
	Packages                  []Packages                `json:"packages"`
	Files                     []interface{}             `json:"files"`
	PackageManager            bool                      `json:"packageManager"`
	Image                     Image                     `json:"image"`
	History                   []History                 `json:"history"`
	SHA                       string                    `json:"id"`
	ComplianceIssues          interface{}               `json:"complianceIssues"`
	AllCompliance             AllCompliance             `json:"allCompliance"`
	Vulnerabilities           []Vulnerabilities         `json:"vulnerabilities"`
	RepoTag                   RepoTag                   `json:"repoTag"`
	Tags                      []Tags                    `json:"tags"`
	RepoDigests               []string                  `json:"repoDigests"`
	CreationTime              time.Time                 `json:"creationTime"`
	VulnerabilitiesCount      int                       `json:"vulnerabilitiesCount"`
	ComplianceIssuesCount     int                       `json:"complianceIssuesCount"`
	VulnerabilityDistribution VulnerabilityDistribution `json:"vulnerabilityDistribution"`
	ComplianceDistribution    ComplianceDistribution    `json:"complianceDistribution"`
	VulnerabilityRiskScore    int                       `json:"vulnerabilityRiskScore"`
	ComplianceRiskScore       int                       `json:"complianceRiskScore"`
	Layers                    []string                  `json:"layers"`
	TopLayer                  string                    `json:"topLayer"`
	RiskFactors               RiskFactors               `json:"riskFactors"`
	Labels                    []string                  `json:"labels"`
	InstalledProducts         InstalledProducts         `json:"installedProducts"`
	ScanVersion               string                    `json:"scanVersion"`
	FirstScanTime             time.Time                 `json:"firstScanTime"`
	CloudMetadata             CloudMetadata             `json:"cloudMetadata"`
	Instances                 []Instances               `json:"instances"`
	Hosts                     map[string]interface{}    `json:"hosts"`
	Err                       string                    `json:"err"`
	Collections               []string                  `json:"collections"`
	ScanID                    int                       `json:"scanID"`
	TrustStatus               string                    `json:"trustStatus"`
}

type Binaries struct {
	Name       string `json:"name"`
	Path       string `json:"path"`
	Md5        string `json:"md5"`
	CveCount   int    `json:"cveCount"`
	LayerTime  int    `json:"layerTime"`
	Version    string `json:"version,omitempty"`
	MissingPkg bool   `json:"missingPkg,omitempty"`
}

type StartupBinaries struct {
	Name      string `json:"name"`
	Path      string `json:"path"`
	Md5       string `json:"md5"`
	CveCount  int    `json:"cveCount"`
	LayerTime int    `json:"layerTime"`
	Version   string `json:"version,omitempty"`
}

type Pkgs struct {
	Version    string   `json:"version"`
	Name       string   `json:"name"`
	CveCount   int      `json:"cveCount"`
	License    string   `json:"license"`
	LayerTime  int      `json:"layerTime"`
	BinaryPkgs []string `json:"binaryPkgs,omitempty"`
}

type Packages struct {
	PkgsType string `json:"pkgsType"`
	Pkgs     []Pkgs `json:"pkgs"`
}

type Image struct {
	Created time.Time `json:"created"`
}

type History struct {
	Created     int      `json:"created"`
	Instruction string   `json:"instruction"`
	SizeBytes   int      `json:"sizeBytes,omitempty"`
	ID          string   `json:"id"`
	Tags        []string `json:"tags,omitempty"`
}

type AllCompliance struct {
	Enabled bool `json:"enabled"`
}

type AttackComplexityLow struct {
}

type AttackVectorNetwork struct {
}

type HasFix struct {
}

type HighSeverity struct {
}

type RecentVulnerability struct {
}

type RiskFactors struct {
}

type CriticalSeverity struct {
}

type DoS struct {
}

type RemoteExecution struct {
}

type Vulnerabilities struct {
	Text            string      `json:"text"`
	ID              int         `json:"id"`
	Severity        string      `json:"severity"`
	Cvss            float64     `json:"cvss"`
	Status          string      `json:"status"`
	Cve             string      `json:"cve"`
	Cause           string      `json:"cause"`
	Description     string      `json:"description"`
	Title           string      `json:"title"`
	VecStr          string      `json:"vecStr"`
	Exploit         string      `json:"exploit"`
	RiskFactors     RiskFactors `json:"riskFactors,omitempty"`
	Link            string      `json:"link"`
	Type            string      `json:"type"`
	PackageName     string      `json:"packageName"`
	PackageVersion  string      `json:"packageVersion"`
	LayerTime       int         `json:"layerTime"`
	Templates       interface{} `json:"templates"`
	Twistlock       bool        `json:"twistlock"`
	Published       int         `json:"published"`
	ApplicableRules []string    `json:"applicableRules"`
	Discovered      time.Time   `json:"discovered"`
	Block           bool        `json:"block"`
}

type RepoTag struct {
	Registry string `json:"registry"`
	Repo     string `json:"repo"`
	Tag      string `json:"tag"`
}

type Tags struct {
	Registry string `json:"registry"`
	Repo     string `json:"repo"`
	Tag      string `json:"tag"`
}

type VulnerabilityDistribution struct {
	Critical int `json:"critical"`
	High     int `json:"high"`
	Medium   int `json:"medium"`
	Low      int `json:"low"`
	Total    int `json:"total"`
}

type ComplianceDistribution struct {
	Critical int `json:"critical"`
	High     int `json:"high"`
	Medium   int `json:"medium"`
	Low      int `json:"low"`
	Total    int `json:"total"`
}

type InstalledProducts struct {
	Docker            string `json:"docker"`
	HasPackageManager bool   `json:"hasPackageManager"`
}

type CloudMetadata struct {
}

type Instances struct {
	Image    string    `json:"image"`
	Host     string    `json:"host"`
	Modified time.Time `json:"modified"`
	Tag      string    `json:"tag"`
	Repo     string    `json:"repo"`
	Registry string    `json:"registry"`
}
