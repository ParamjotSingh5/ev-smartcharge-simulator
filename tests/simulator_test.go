package tests

import (
	"testing"

	"github.com/ParamjotSingh5/EV-charging-station-simulator/application"
	"github.com/ParamjotSingh5/EV-charging-station-simulator/domain"
)

func TestStepSimulation(t *testing.T) {
	ports := []*domain.ChargingPort{
		{ID: 1},
		{ID: 2},
	}
	station := &domain.ChargingStation{
		Ports:         ports,
		TotalCapacity: 40.0,
	}
	sim := application.NewSimulator(station, 30) // 30-min timestep

	ev := &domain.EV{ID: 1, StateOfCharge: 0, TargetCharge: 30, MaxChargeRate: 22, Deadline: 120}
	ok := sim.AddEV(0, ev)
	if !ok {
		t.Fatal("failed to add EV")
	}
	sim.Step() // First 30 minutes, should deliver up to max charge allowed

	if ev.StateOfCharge <= 0 {
		t.Error("State of charge did not increase after Step")
	}
	// Add more assertions as needed, including edge cases (full battery, power above EV max, etc.)

	if ev.StateOfCharge > ev.TargetCharge {
		t.Error("State of charge exceeded target charge")
	}

	if !station.Ports[0].Occupied {
		t.Error("Port should still be occupied")
	}

	// Simulate until deadline
	for sim.CurrentTime < ev.Deadline {
		sim.Step()
	}

	if station.Ports[0].Occupied {
		t.Error("Port should be free after deadline")
	}
	if station.Ports[0].EV != nil {
		t.Error("EV should be removed from port after deadline")
	}
	if ev.StateOfCharge < ev.TargetCharge {
		t.Error("EV should have reached target charge before deadline")
	}
}

func TestEqualShareStrategy(t *testing.T) {
	ev1 := &domain.EV{ID: 1, MaxChargeRate: 11}
	ev2 := &domain.EV{ID: 2, MaxChargeRate: 22}
	ports := []*domain.ChargingPort{
		{ID: 1, Occupied: true, EV: ev1},
		{ID: 2, Occupied: true, EV: ev2},
	}
	strat := &application.EqualSharingStrategy{}
	allocations := strat.AssignPower(ports, 40)
	if allocations[0] != 11 {
		t.Errorf("Expected EV1 to be capped at 11 kW, got %v", allocations[0])
	}
	if allocations[1] != 20 {
		t.Errorf("Expected EV2 to get 20 kW, got %v", allocations[1])
	}

	t.Logf("Allocations: EV1: %v kW, EV2: %v kW", allocations[0], allocations[1])
}

func TestEarliestDeadlineStrategy(t *testing.T) {
	ev1 := &domain.EV{ID: 1, MaxChargeRate: 11, Deadline: 30}
	ev2 := &domain.EV{ID: 2, MaxChargeRate: 22, Deadline: 60}
	ports := []*domain.ChargingPort{
		{ID: 1, Occupied: true, EV: ev1},
		{ID: 2, Occupied: true, EV: ev2},
	}
	strat := &application.EarliestDeadlineFirstStrategy{}
	allocations := strat.AssignPower(ports, 20)
	// EV1 should be prioritized
	if allocations[0] != 11 {
		t.Errorf("Expected EV1 to get 11 kW, got %v", allocations[0])
	}
	if allocations[1] != 9 {
		t.Errorf("Expected EV2 to get remaining 9 kW, got %v", allocations[1])
	}
}
