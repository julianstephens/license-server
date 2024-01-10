package model

type DatabaseCfg struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DB       string `mapstructure:"dbname"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
}

type ServerCfg struct {
	Host                 string `mapstructure:"host"`
	Port                 int    `mapstructure:"port"`
	Env                  string `mapstructure:"env"`
	KeyFile              string `mapstructure:"keyfile"`
	TokenFile            string `mapstructure:"tokenfile"`
	LogFile              string `mapstructure:"logfile"`
	OutDir               string `mapstructure:"outdir"`
	LicenseLength        int    `mapstructure:"licenselength"`
	DefaultLicenseLength int    `mapstructure:"default_license_term"`
	MaxOfflineDuration   int    `mapstructure:"maximum_Offline_duration"`
}

type AuthCfg struct {
	Host       string `mapstructure:"host"`
	Port       int    `mapstructure:"port"`
	JwksUrl    string `mapstructure:"jwks_url"`
	TenantID   string `mapstructure:"tenant_id"`
	ResourceID string `mapstructure:"resource_id"`
	ClientID   string `mapstructure:"client_id"`
}

type AppCfg struct {
	Name     string `mapstructure:"name"`
	Filename string `mapstructure:"filename"`
	Version  string `mapstructure:"version"`
}

type Config struct {
	Database DatabaseCfg `mapstructure:"database"`
	Server   ServerCfg   `mapstructure:"server"`
	Auth     AuthCfg     `mapstructure:"auth"`
	App      AppCfg      `mapstructure:"app"`
}
