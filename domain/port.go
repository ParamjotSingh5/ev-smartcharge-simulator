package domain

type ChargingPort struct {
	ID       int
	Occupied bool
	EV       *EV
}
