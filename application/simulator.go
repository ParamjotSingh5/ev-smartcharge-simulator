package application

import (
	"github.com/ParamjotSingh5/EV-charging-station-simulator/domain"
)

type Simulator struct {
	Station      *domain.ChargingStation
	CurrentTime  int // minute count
	TimeStepMins int
}

func NewSimulator(station *domain.ChargingStation, timeStepMins int) *Simulator {
	return &Simulator{
		Station:      station,
		CurrentTime:  0,
		TimeStepMins: timeStepMins,
	}
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

// AddEV connects an EV to a specified port if it's available.
func (s *Simulator) AddEV(portIdx int, ev *domain.EV) bool {
	if s.Station.Ports[portIdx].Occupied {
		return false
	}
	s.Station.Ports[portIdx].EV = ev
	s.Station.Ports[portIdx].Occupied = true
	return true
}
