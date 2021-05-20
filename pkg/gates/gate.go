package gates

import "context"

// GateType represents when the gate is supposed to be executed
type GateType int

const (
	// PreDelete evaluates before image deletion
	PreDelete GateType = iota
	// PreImport evaluates before image import
	PreImport
)

// Gate is an object that evaluates if something has happened
type Gate interface {
	// GetGateType returns the type of gate
	GetGateType() GateType
	// Evaluate checks to see if the gate is open or closed. True = open, False = closed
	Evaluate(ctx context.Context) (bool, error)
}
