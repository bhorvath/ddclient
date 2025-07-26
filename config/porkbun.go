package config

// Porkbun specifies options pertaining to the Porkbun API.
type Porkbun struct {
	APIKey    string `help:"Porkbun API key"`
	SecretKey string `help:"Porkbun secret key"`
}
