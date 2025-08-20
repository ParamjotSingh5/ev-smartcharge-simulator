package domain

type ChargingStation struct {
	Ports         []*ChargingPort
	TotalCapacity float64 // kW total
}
