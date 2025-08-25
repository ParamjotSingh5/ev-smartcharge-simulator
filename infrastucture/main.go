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
		TotalCapacity: 20.0,
	}
	sim := application.NewSimulator(station, 10) // 10-minute steps
	sim.SetStrategy(&application.EqualShareWithV2G{})

	ev1 := &domain.EV{ID: 1, StateOfCharge: 10, TargetCharge: 50, MaxChargeRate: 22, Deadline: 40}
	ev2 := &domain.EV{ID: 2, StateOfCharge: 20, TargetCharge: 60, MaxChargeRate: 11, Deadline: 80}
	ev3 := &domain.EV{ID: 3, StateOfCharge: 30, TargetCharge: 70, MaxChargeRate: 7, Deadline: 20}

	sim.AddEV(1, ev2)
	sim.AddEV(2, ev3)
	sim.AddEV(0, ev1)

	for i := 0; i < 10; i++ {
		sim.StepWithStrategy(-10)
		fmt.Printf("T+%d min: EV1: %.2f kWh, EV2: %.2f kWh, EV3: %.2F kWh \n", sim.CurrentTime, ev1.StateOfCharge, ev2.StateOfCharge, ev3.StateOfCharge)
	}
}
