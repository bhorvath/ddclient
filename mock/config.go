package mock

import "github.com/bhorvath/ddclient/config"

// GetAppConfig returns a config.App with test values.
func GetAppConfig() *config.App {
	return &config.App{
		Record: config.Record{
			Domain:		"internet.com",
			Name: 		"test",
			Type:			"A",
		},
		Porkbun: config.Porkbun{
			APIKey:			"api-key",
			SecretKey:	"secret-key",
		},
	}
}

// GetAppArgs returns a config.Args with test values.
func GetAppArgs() *config.Args {
	return &config.Args{
		Record: config.Record{
			Domain:		"internet.com",
			Name: 		"test",
			Type:			"A",
		},
		Porkbun: config.Porkbun{
			APIKey:			"api-key",
			SecretKey:	"secret-key",
		},
	}
}
