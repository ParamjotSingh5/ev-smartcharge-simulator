package application

import (
	"github.com/ParamjotSingh5/EV-charging-station-simulator/domain"
)

type Simulator struct {
	Station      *domain.ChargingStation
	CurrentTime  int // minute count
	TimeStepMins int
	Strategy     ChargingStrategy
}

func NewSimulator(station *domain.ChargingStation, timeStepMins int) *Simulator {
	return &Simulator{
		Station:      station,
		CurrentTime:  0,
		TimeStepMins: timeStepMins,
	}
}

func (s *Simulator) SetStrategy(strategy ChargingStrategy) {
	s.Strategy = strategy
}

// Step advances the simulation by one time step, distributing available power among connected EVs.
// It updates each EV's state of charge based on its max charge rate, target charge, and the station's total capacity.
func (s *Simulator) Step() {
	s.CurrentTime += s.TimeStepMins
	totalPower := s.Station.TotalCapacity
	activeEVs := []*domain.EV{}

	// Collect all connected EVs
	for _, port := range s.Station.Ports {
		if port.Occupied && port.EV != nil {
			if port.EV.Deadline <= s.CurrentTime {
				port.Occupied = false
				port.EV = nil
				continue
			}
			if port.EV.StateOfCharge >= port.EV.TargetCharge {
				port.Occupied = false
				port.EV = nil
				continue
			}
			activeEVs = append(activeEVs, port.EV)
		}
	}

	if len(activeEVs) == 0 {
		return // No EVs to charge
	}

	powerPerEV := totalPower / float64(len(activeEVs))

	for _, ev := range activeEVs {
		if ev.StateOfCharge >= ev.TargetCharge {
			continue // Skip if already at or above target
		}

		maxPossibleCharge := ev.MaxChargeRate * float64(s.TimeStepMins) / 60.0 // kWh for this timestep
		requiredCharge := ev.TargetCharge - ev.StateOfCharge
		actualCharge := min(maxPossibleCharge, requiredCharge, powerPerEV*float64(s.TimeStepMins)/60.0)

		ev.StateOfCharge += actualCharge
		if ev.StateOfCharge > ev.TargetCharge {
			ev.StateOfCharge = ev.TargetCharge // Cap at target charge
		}
	}
}

func (s *Simulator) StepWithStrategy(gridNeeds float64) {
	// gridNeeds: negative for grid needing power (V2G), positive for normal charging

	var connectedPorts []*domain.ChargingPort
	for _, port := range s.Station.Ports {
		if port.Occupied && port.EV != nil {
			connectedPorts = append(connectedPorts, port)
		}
	}
	if len(connectedPorts) == 0 {
		s.CurrentTime += s.TimeStepMins
		return
	}

	var assignments []PowerAssignment
	if s.Strategy != nil {
		assignments = s.Strategy.AssignPower(connectedPorts, s.Station.TotalCapacity, gridNeeds)
	} else {
		// fallback if not set
	}

	for i, port := range connectedPorts {
		assign := assignments[i]
		// Charging (assign.Power > 0), or Discharging (assign.Power < 0)
		energyDelivered := assign.Power * float64(s.TimeStepMins) / 60.0 // kWh (neg or pos)
		port.EV.StateOfCharge += energyDelivered
		// Clamp values
		if port.EV.StateOfCharge > port.EV.TargetCharge && assign.Power < 0 {
			port.EV.StateOfCharge = port.EV.TargetCharge
		}
		if port.EV.StateOfCharge > port.EV.TargetCharge && assign.Power > 0 {
			port.EV.StateOfCharge = port.EV.TargetCharge
		}
		if port.EV.StateOfCharge < 0 {
			port.EV.StateOfCharge = 0
		}
		// handle EV departure if charge/discharge is done
	}
	s.CurrentTime += s.TimeStepMins
}

// AddEV connects an EV to a specified port if it's available.
func (s *Simulator) AddEV(portIdx int, ev *domain.EV) bool {
	if s.Station.Ports[portIdx].Occupied {
		return false
	}
	s.Station.Ports[portIdx].EV = ev
	s.Station.Ports[portIdx].Occupied = true
	return true
}
