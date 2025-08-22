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
