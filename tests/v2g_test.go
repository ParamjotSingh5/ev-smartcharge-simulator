package tests

import (
	"testing"

	"github.com/ParamjotSingh5/EV-charging-station-simulator/application"
	"github.com/ParamjotSingh5/EV-charging-station-simulator/domain"
)

// Ensure the EqualShareWithV2G strategy conforms to V2GChargingStrategy
func TestEqualShareWithV2G_Discharge(t *testing.T) {
	ev := &domain.EV{
		ID: 1, StateOfCharge: 40,
		TargetCharge:  20,
		MaxChargeRate: 10, MaxDischargeRate: 5,
		Deadline: 60, V2GCapable: true,
	}
	port := &domain.ChargingPort{ID: 1, Occupied: true, EV: ev}
	strat := &application.EqualShareWithV2G{}
	station := &domain.ChargingStation{Ports: []*domain.ChargingPort{port}, TotalCapacity: 10}
	sim := application.NewSimulator(station, 60)
	sim.SetStrategy(strat)
	sim.StepWithStrategy(-5) // grid requests 5kW discharge for 60 minutes

	// Check SOC decreased properly
	expectedSOC := 20.0 // Can't discharge past TargetCharge per safety spec
	if ev.StateOfCharge != expectedSOC {
		t.Errorf("Expected SOC %v after V2G discharge, got %v", expectedSOC, ev.StateOfCharge)
	}

	// Should not discharge past TargetCharge
	sim.StepWithStrategy(-10) // Request far more than possible
	if ev.StateOfCharge < 20 {
		t.Errorf("EV over-discharged below target: %v", ev.StateOfCharge)
	}
}
