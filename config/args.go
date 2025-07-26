package config

import "github.com/alexflint/go-arg"

// Args contains configuration arguments which can be set from the command line.
type Args struct {
	Record
	Porkbun
	ConfigFilePath string `arg:"--config" help:"config file to use"`
	Save           bool   `help:"save configs to file (no other action is taken)"`
}

// ParseArgs parses and returns command line args.
func ParseArgs() (a *Args) {
	a = &Args{}
	arg.MustParse(a)
	return a
}
