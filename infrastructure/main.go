package main

import (
	"fmt"

	"github.com/ParamjotSingh5/EV-charging-station-simulator/application"
	"github.com/ParamjotSingh5/EV-charging-station-simulator/domain"
)

func main() {
	ports := []*domain.ChargingPort{
		{ID: 1}, {ID: 2}, {ID: 3},
	}
	station := &domain.ChargingStation{
		Ports:         ports,
		TotalCapacity: 50.0,
	}
	sim := application.NewSimulator(station, 10) // 10-minute steps

	ev1 := &domain.EV{ID: 1, StateOfCharge: 10, TargetCharge: 50, MaxChargeRate: 22, Deadline: 60}
	ev2 := &domain.EV{ID: 2, StateOfCharge: 20, TargetCharge: 60, MaxChargeRate: 11, Deadline: 80}
	sim.AddEV(0, ev1)
	sim.AddEV(1, ev2)

	for i := 0; i < 10; i++ {
		sim.Step()
		fmt.Printf("T+%d min: EV1: %.2f kWh, EV2: %.2f kWh\n", sim.CurrentTime, ev1.StateOfCharge, ev2.StateOfCharge)
	}
}
