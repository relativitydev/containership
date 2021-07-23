/*
 This is the basic contract of a gate. Implement and register your own gate to make it accessible for consumption.
*/

package gates

import "context"

// Lifecycle represents when the gate is supposed to be executed
type Lifecycle int

const (
	// Delete evaluates before image deletion
	Delete Lifecycle = iota
	// Promotion evaluates before image promotion
	Promotion
)

// Gate is an object that evaluates if something has happened
type Gate interface {
	// GetLifecycle returns the type of gate
	GetLifecycle() Lifecycle
	// Evaluate checks to see if the gate is open or closed. True = open, False = closed
	Evaluate(ctx context.Context) (bool, error)
}
