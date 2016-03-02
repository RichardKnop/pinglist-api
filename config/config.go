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

// SessionConfig stores session configuration for the web app
type SessionConfig struct {
	Secret string
	Path   string
	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'.
	// MaxAge>0 means Max-Age attribute present and given in seconds.
	MaxAge int
	// When you tag a cookie with the HttpOnly flag, it tells the browser that
	// this particular cookie should only be accessed by the server.
	// Any attempt to access the cookie from client script is strictly forbidden.
	HTTPOnly bool
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
	Scheme string
	Host   string
}

// Config stores all configuration options
type Config struct {
	Database      DatabaseConfig
	Oauth         OauthConfig
	Session       SessionConfig
	Facebook      FacebookConfig
	Sendgrid      SendgridConfig
	Stripe        StripeConfig
	Web           WebConfig
	IsDevelopment bool
}
