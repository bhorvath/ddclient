package config

import "github.com/alexflint/go-arg"

// Args contains configuration arguments which can be set from the command line.
type Args struct {
	Record
	Porkbun
}

// ParseArgs parses and returns command line args.
func ParseArgs() (a *Args) {
	a = &Args{}
	arg.MustParse(a)
	return a
}
