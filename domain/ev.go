package domain

type EV struct {
	ID            int
	StateOfCharge float64 // kWh
	TargetCharge  float64 // kWh
	MaxChargeRate float64 // kW
	Deadline      int     // Minutes until required departure
}
