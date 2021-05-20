package gates

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/relativitydev/containership/pkg/prisma"
)

const retryWait = 5 * time.Second
const retry = 4

// prismaGate object
type prismaGate struct {
	client     *prisma.Client
	Registry   string
	Repository string
	Tag        string
}

// prismaGateImageMetadata holds the data the gate needs
type prismaGateImageMetadata struct {
	registry   string
	repository string
	tag        string
}

// newPrismaGate creates a new gate object
func newPrismaGate(metadata map[string]string) (Gate, error) {
	parsedMetadata, err := parsePrismaGateImageMetadata(metadata)
	if err != nil {
		return nil, fmt.Errorf("unable to get prisma metadata: %v", err)
	}

	client, err := prisma.NewClient(parsedMetadata.registry)
	if err != nil {
		return nil, fmt.Errorf("failed to get registry config: %v", err)
	}

	return &prismaGate{
		client:     client,
		Registry:   parsedMetadata.registry,
		Repository: parsedMetadata.repository,
		Tag:        parsedMetadata.tag,
	}, nil
}

// GetGateType returns the type of gate
func (s *prismaGate) GetGateType() GateType {
	return PreImport
}

// Evaluate determines if this image has any security vulnerabilities
func (s *prismaGate) Evaluate(ctx context.Context) (bool, error) {
	reports, err := s.scanImage(ctx)
	if err != nil {
		return false, errors.Wrap(err, "Failed to scan image")
	}

	var allReports []prisma.ScanReport = *reports

	var rpt prisma.ScanReport = allReports[0]

	return !s.hasVulnerabilities(rpt), nil
}

// parsePrismaGateImageMetadata parses metadata
func parsePrismaGateImageMetadata(metadata map[string]string) (*prismaGateImageMetadata, error) {
	meta := prismaGateImageMetadata{}

	if val, ok := metadata["registry"]; ok && val != "" {
		meta.registry = val
	} else {
		return nil, fmt.Errorf("no registry string given")
	}

	if val, ok := metadata["repository"]; ok && val != "" {
		meta.repository = val
	} else {
		return nil, fmt.Errorf("no repository string given")
	}

	if val, ok := metadata["tag"]; ok && val != "" {
		meta.tag = val
	} else {
		return nil, fmt.Errorf("no tag string given")
	}

	return &meta, nil
}

func (s *prismaGate) scanImage(ctx context.Context) (*[]prisma.ScanReport, error) {
	for i := 0; i < retry; i++ {
		resp, err := s.client.Get(ctx, fmt.Sprintf("%s/%s:%s", s.Registry, s.Repository, s.Tag))
		if err != nil {
			return nil, errors.Wrap(err, "Failed to get image scan report")
		}

		if len(*resp) == 0 {
			_, err := s.client.ScanImage(ctx, s.Registry, s.Repository, s.Tag)
			if err != nil {
				return nil, fmt.Errorf("unable to scan image: %s/%s:%s :: %s", s.Registry, s.Repository, s.Tag, err)
			}

			time.Sleep(retryWait)

			continue
		}

		return resp, nil
	}

	return nil, fmt.Errorf("image scan failed")
}

func (s *prismaGate) hasVulnerabilities(rpt prisma.ScanReport) bool {
	vulnerabilityLevel, defined := os.LookupEnv("PRISMA_VULNERABILITY_LEVEL")

	if defined {
		switch strings.ToLower(vulnerabilityLevel) {
		case "critical":
			return rpt.VulnerabilityDistribution.Critical > 0
		case "high":
			return (rpt.VulnerabilityDistribution.Critical + rpt.VulnerabilityDistribution.High) > 0
		case "medium":
			return (rpt.VulnerabilityDistribution.Critical + rpt.VulnerabilityDistribution.High + rpt.VulnerabilityDistribution.Medium) > 0
		case "low":
			return rpt.VulnerabilitiesCount > 0
		}
	}

	return rpt.VulnerabilitiesCount > 0
}
