package config

// DatabaseConfig stores database connection options
type DatabaseConfig struct {
	Type         string
	Host         string
	Port         int
	User         string
	Password     string
	DatabaseName string
	MaxIdleConns int
	MaxOpenConns int
}

// OauthConfig stores oauth service configuration options
type OauthConfig struct {
	AccessTokenLifetime  int
	RefreshTokenLifetime int
	AuthCodeLifetime     int
}

// AWSConfig stores AWS related configuration
type AWSConfig struct {
	Region                     string
	AssetsBucket               string
	APNSPlatformApplicationARN string
	GCMPlatformApplicationARN  string
}

// FacebookConfig stores Facebook app config
type FacebookConfig struct {
	AppID     string
	AppSecret string
}

// SendgridConfig stores sengrid configuration options
type SendgridConfig struct {
	APIKey string
}

// StripeConfig stores stripe configuration options
type StripeConfig struct {
	SecretKey      string
	PublishableKey string
}

// WebConfig stores web related config like scheme and host
type WebConfig struct {
	AppScheme string
	AppHost   string
}

// PinglistConfig stores app specific config
type PinglistConfig struct {
	PasswordResetLifetime int
	ContactEmail          string
}

// Config stores all configuration options
type Config struct {
	Database      DatabaseConfig
	Oauth         OauthConfig
	AWS           AWSConfig
	Facebook      FacebookConfig
	Sendgrid      SendgridConfig
	Stripe        StripeConfig
	Web           WebConfig
	Pinglist      PinglistConfig
	IsDevelopment bool
}
