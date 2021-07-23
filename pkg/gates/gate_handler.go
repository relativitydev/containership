package gates

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/relativitydev/containership/api/v1beta2"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

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
func (h *GateHandler) GetGate(gateTypeName string) (Gate, error) {

	_, err := h.getGateCustomResource(gateTypeName)
	if err != nil {
		return nil, err
	}

	switch gateTypeName {
	case "prisma-promotion":
		return nil, nil
	case "prisma-deletion":
		return nil, nil
	default:
		return nil, fmt.Errorf("no gate found for type: %s", gateTypeName)
	}
}

func (h *GateHandler) getGateCustomResource(gateTypeName string) (map[string]string, error) {
	cl, err := client.New(config.GetConfigOrDie(), client.Options{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create client when looking for gate CR")
	}

	gateList := &v1beta2.GateList{}

	err = cl.List(context.Background(), gateList)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list gates to get metadata")
	}

	for _, gate := range gateList.Items {
		if gate.Spec.Type == gateTypeName {
			return gate.Spec.Metadata, nil
		}
	}

	// If no gate type match is found that is okay. Maybe not every gate needs metadata?
	return nil, nil
}
