package main_test

import (
	"testing"

	"mobilehound/simul"
)

func TestSimulation(t *testing.T) {
	simul.Start("randhound.toml")
}
