package gates

import (
	"fmt"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/go-logr/logr"
)

// GateHandler encapsulates the logic of calling the correct gate
type GateHandler struct {
	logger logr.Logger
}

// NewGateHandler creates a GateHandler object
func NewGateHandler() *GateHandler {
	return &GateHandler{
		logger: ctrl.Log.WithName("gates/gate_handler"),
	}
}

// GetGate creates the specified type of Gate object
func (h *GateHandler) GetGate(name string, gateMetadata map[string]string) (Gate, error) {
	switch name {
	case "prisma":
		return newPrismaGate(gateMetadata)
	case "deployedImage":
		return newDeployedImageGate(gateMetadata)
	default:
		return nil, fmt.Errorf("No gate found for type: %s", name)
	}
}
