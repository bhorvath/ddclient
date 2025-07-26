package config

// Record specifies configurable options pertaining to the DNS record with which we will interact with.
type Record struct {
	Domain string `help:"the domain of the record"`
	Type   string `help:"the type of the record"`
	Name   string `help:"the name of the record"`
}
