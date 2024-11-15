package config

// EnvConfig is the configuration for the application
type EnvConfig struct {
    Host         string `json:"HOST" validate:"required"`
    Port         string `json:"PORT" validate:"required"`
    DatabaseURL  string `json:"DATABASE_URL" validate:"required"`
    JwtSecret    string `json:"JWT_SECRET" validate:"required"`
    SignatureKey string `json:"SIGNATURE_KEY" validate:"required"`
}
