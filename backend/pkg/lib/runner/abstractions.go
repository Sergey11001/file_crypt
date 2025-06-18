package runner

import (
	"context"
)

// Labeler returns label of the application.
type Labeler interface {
	Label() string
}

// Module of the application.
type Module interface {

	// Name returns the module name.
	Name() string

	// Run runs the module.
	Run(ctx context.Context)
}
