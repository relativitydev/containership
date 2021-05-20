/*
Note: We are using Prisma to determine if the image is deployed rather than direct K8s API calls because we need to check the image in all clusters. Authenticating and querying each cluster individually would be difficult, but Prisma already has the information we need.
*/

package gates

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/relativitydev/containership/pkg/prisma"
)

// deployedImageGate object
type deployedImageGate struct {
	client     *prisma.Client
	Registry   string
	Repository string
	Tag        string
}

// DeployedImageGateMetadata holds the data the gate needs
type deployedImageGateMetadata struct {
	registry   string
	repository string
	tag        string
}

// newDeployedImageGate creates a new gate object
func newDeployedImageGate(metadata map[string]string) (Gate, error) {
	parsedMetadata, err := parseDeployedImageMetadata(metadata)
	if err != nil {
		return nil, fmt.Errorf("unable to get deployed image gate metadata: %s", err)
	}

	client, err := prisma.NewClient(parsedMetadata.registry)
	if err != nil {
		return nil, fmt.Errorf("Failed to get registry config: %s", err)
	}

	return &deployedImageGate{
		client:     client,
		Registry:   parsedMetadata.registry,
		Repository: parsedMetadata.repository,
		Tag:        parsedMetadata.tag,
	}, nil
}

// GetGateType returns the type of gate
func (s *deployedImageGate) GetGateType() GateType {
	return PreDelete
}

// Evaluate determines if this image is deployed on any hosts
func (s *deployedImageGate) Evaluate(ctx context.Context) (bool, error) {
	reports, err := s.client.GetImage(ctx, fmt.Sprintf("%s/%s:%s", s.Registry, s.Repository, s.Tag))
	if err != nil {
		return false, errors.Wrap(err, "failed to get image report from prisma")
	}

	var allReports []prisma.ScanReport = *reports
	if len(allReports) != 0 {
		var rpt prisma.ScanReport = allReports[0]
		return len(rpt.Hosts) == 0, nil
	}

	// no report found means it isn't deployed
	return true, nil
}

// parseMetadata parses metadata
func parseDeployedImageMetadata(metadata map[string]string) (*deployedImageGateMetadata, error) {
	meta := deployedImageGateMetadata{}

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
