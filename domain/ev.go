package domain

type EV struct {
	ID               int
	StateOfCharge    float64 // kWh
	TargetCharge     float64 // kWh
	MaxChargeRate    float64 // kW
	MaxDischargeRate float64 // kW, max rate for V2G
	Deadline         int     // Minutes until required departure
	V2GCapable       bool    // Vehicle-to-Grid capability
}
