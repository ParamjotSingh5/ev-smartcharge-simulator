package tests

import (
	"testing"

	"github.com/ParamjotSingh5/EV-charging-station-simulator/application"
	"github.com/ParamjotSingh5/EV-charging-station-simulator/domain"
)

// Define EqualShareStrategy to conform to V2GChargingStrategy
type EqualShareStrategy struct{}

func (e *EqualShareStrategy) AssignPower(ports []*domain.ChargingPort, stationCapacity float64, gridNeeds float64) []application.PowerAssignment {
	n := len(ports)
	assignments := make([]application.PowerAssignment, n)
	if n == 0 {
		return assignments
	}
	share := stationCapacity / float64(n)
	for i, port := range ports {
		assign := share
		if share > port.EV.MaxChargeRate {
			assign = port.EV.MaxChargeRate
		}
		assignments[i] = application.PowerAssignment{Power: assign, PortIndex: i}
	}
	return assignments
}

func TestEqualShareStrategy_NoV2G(t *testing.T) {
	ev1 := &domain.EV{ID: 1, MaxChargeRate: 11, StateOfCharge: 0, TargetCharge: 30, V2GCapable: false}
	ev2 := &domain.EV{ID: 2, MaxChargeRate: 22, StateOfCharge: 0, TargetCharge: 30, V2GCapable: false}
	ports := []*domain.ChargingPort{
		{ID: 1, Occupied: true, EV: ev1},
		{ID: 2, Occupied: true, EV: ev2},
	}
	strat := &EqualShareStrategy{}
	station := &domain.ChargingStation{Ports: ports, TotalCapacity: 40}
	sim := application.NewSimulator(station, 60) // 60 min step
	sim.SetStrategy(strat)
	sim.StepWithStrategy(0) // gridNeeds = 0, normal charging

	if ev1.StateOfCharge != 11 {
		t.Errorf("EV1 should receive 11 kWh (max rate), got %v", ev1.StateOfCharge)
	}
	if ev2.StateOfCharge != 20 {
		t.Errorf("EV2 should receive 20 kWh (split of remaining), got %v", ev2.StateOfCharge)
	}
}
