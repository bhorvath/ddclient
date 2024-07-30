package config

// Record specifies configurable options pertaining to the DNS record with which we will interact with.
type Record struct {
	Domain string `arg:"required" help:"the domain of the record"`
	Type   string `arg:"required" help:"the type of the record"`
	Name   string `arg:"required" help:"the name of the record"`
}
