package config

// Porkbun specifies options pertaining to the Porkbun API.
type Porkbun struct {
	APIKey    string `arg:"required" help:"Porkbun API key"`
	SecretKey string `arg:"required" help:"Porkbun secret key"`
}
