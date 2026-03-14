package main

import "testing"

// Compile-time contract: newFXApp must remain the bootstrap factory shape used by run().
var _ func() lifecycleApp = newFXApp

func TestNewFXAppFactoryContract(t *testing.T) {
	t.Parallel()
}
