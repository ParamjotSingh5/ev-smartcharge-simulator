package application

import "github.com/ParamjotSingh5/EV-charging-station-simulator/domain"

// direction: +ve for charging, -ve for discharging
type PowerAssignment struct {
	Power     float64 // kW assigned (neg for feeding grid)
	PortIndex int
}

type ChargingStrategy interface {
	// Distribute available station power among occupied ports
	AssignPower(ports []*domain.ChargingPort, stationCapacity float64, gridNeeds float64) []PowerAssignment
}

type EqualShareWithV2G struct{}

func (e *EqualShareWithV2G) AssignPower(ports []*domain.ChargingPort, stationCapacity float64, gridNeeds float64) []PowerAssignment {
	n := len(ports)
	assigned := make([]PowerAssignment, n)
	// gridNeeds: positive means grid wants to absorb, negative means grid requests extra energy

	// Simple: try to discharge from all V2G EVs to meet gridNeeds (if negative), otherwise charge as normal
	totalDischargePossible := 0.0
	for _, port := range ports {
		if port.EV.V2GCapable && port.EV.StateOfCharge > port.EV.TargetCharge {
			avail := port.EV.StateOfCharge - port.EV.TargetCharge
			if avail > port.EV.MaxDischargeRate {
				avail = port.EV.MaxDischargeRate
			}
			totalDischargePossible += avail
		}
	}
	// Discharge mode: grid needs power
	if gridNeeds < 0 && totalDischargePossible > 0 {
		needed := -gridNeeds
		for i, port := range ports {
			if port.EV.V2GCapable && port.EV.StateOfCharge > port.EV.TargetCharge {
				avail := port.EV.StateOfCharge - port.EV.TargetCharge
				if avail > port.EV.MaxDischargeRate {
					avail = port.EV.MaxDischargeRate
				}
				power := avail
				if power > needed {
					power = needed
				}
				assigned[i] = PowerAssignment{Power: -power, PortIndex: i} // negative means to grid
				needed -= power
				if needed <= 0 {
					break
				}
			}
		}
	} else {
		// Normal charging
		share := stationCapacity / float64(n)
		for i, port := range ports {
			limit := port.EV.MaxChargeRate
			if share > limit {
				assigned[i] = PowerAssignment{Power: limit, PortIndex: i}
			} else {
				assigned[i] = PowerAssignment{Power: share, PortIndex: i}
			}
		}
	}
	return assigned
}

// Equal Power Sharing - splits capacity evenly
type EqualSharingStrategy struct{}

func (eps *EqualSharingStrategy) AssignPower(ports []*domain.ChargingPort, stationCapacity float64, gridNeeds float64) []PowerAssignment {
	n := len(ports)
	allocation := make([]PowerAssignment, n)
	if n == 0 {
		return allocation
	}
	powerPerPort := stationCapacity / float64(n)
	for i, port := range ports {
		allocation[i].PortIndex = i
		if powerPerPort > port.EV.MaxChargeRate {
			allocation[i].Power = port.EV.MaxChargeRate
		} else {
			allocation[i].Power = powerPerPort
		}
	}
	return allocation
}

type EarliestDeadlineFirstStrategy struct{}

func (edf *EarliestDeadlineFirstStrategy) AssignPower(ports []*domain.ChargingPort, stationCapacity float64, gridNeeds float64) []PowerAssignment {
	// Sort by EV deadline ascending
	sorted := make([]*domain.ChargingPort, len(ports))
	copy(sorted, ports)
	// Simple bubble sort (could use sort.Slice for bigger projects)
	for i := range sorted {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i].EV.Deadline > sorted[j].EV.Deadline {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	allocation := make([]PowerAssignment, len(ports))
	remaining := stationCapacity
	for _, port := range sorted {
		idx := findPortIndex(ports, port)
		allocation[idx].PortIndex = idx
		need := port.EV.MaxChargeRate
		if need > remaining {
			allocation[idx].Power = remaining
			break
		}
		allocation[idx].Power = need
		remaining -= need
	}
	return allocation
}

func findPortIndex(ports []*domain.ChargingPort, target *domain.ChargingPort) int {
	for i, port := range ports {
		if port == target {
			return i
		}
	}
	return -1
}

type FirstComeFirstServeStrategy struct{}

func (fcfs *FirstComeFirstServeStrategy) AssignPower(ports []*domain.ChargingPort, stationCapacity float64, gridNeeds float64) []PowerAssignment {
	allocation := make([]PowerAssignment, len(ports))
	remaining := stationCapacity
	for i, port := range ports {
		need := port.EV.MaxChargeRate
		allocation[i].PortIndex = i
		if need > remaining {
			allocation[i].Power = remaining
			break
		}
		allocation[i].Power = need
		remaining -= need
	}
	return allocation
}
