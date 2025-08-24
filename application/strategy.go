package application

import "github.com/ParamjotSingh5/EV-charging-station-simulator/domain"

type ChargingStrategy interface {
	// Distribute available station power among occupied ports
	AssignPower(ports []*domain.ChargingPort, stationCapacity float64) []float64
}

// Equal Power Sharing - splits capacity evenly
type EqualSharingStrategy struct{}

func (eps *EqualSharingStrategy) AssignPower(ports []*domain.ChargingPort, stationCapacity float64) []float64 {
	n := len(ports)
	allocation := make([]float64, n)
	if n == 0 {
		return allocation
	}
	powerPerPort := stationCapacity / float64(n)
	for i, port := range ports {
		if powerPerPort > port.EV.MaxChargeRate {
			allocation[i] = port.EV.MaxChargeRate
		} else {
			allocation[i] = powerPerPort
		}
	}
	return allocation
}

type EarliestDeadlineFirstStrategy struct{}

func (edf *EarliestDeadlineFirstStrategy) AssignPower(ports []*domain.ChargingPort, stationCapacity float64) []float64 {
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
	allocation := make([]float64, len(ports))
	remaining := stationCapacity
	for _, port := range sorted {
		idx := findPortIndex(ports, port)
		need := port.EV.MaxChargeRate
		if need > remaining {
			allocation[idx] = remaining
			break
		}
		allocation[idx] = need
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

func (fcfs *FirstComeFirstServeStrategy) AssignPower(ports []*domain.ChargingPort, stationCapacity float64) []float64 {
	allocation := make([]float64, len(ports))
	remaining := stationCapacity
	for i, port := range ports {
		need := port.EV.MaxChargeRate
		if need > remaining {
			allocation[i] = remaining
			break
		}
		allocation[i] = need
		remaining -= need
	}
	return allocation
}
