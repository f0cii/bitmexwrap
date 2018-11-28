package bitmexwrap

type Balance struct {
	Currency  string
	Available float64
	Frozen    float64
	Balance   float64
}
