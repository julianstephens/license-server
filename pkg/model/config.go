package model

type Database struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DB       string `mapstructure:"dbname"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
}

type Server struct {
	Host                 string `mapstructure:"host"`
	Port                 int    `mapstructure:"port"`
	Env                  string `mapstructure:"env"`
	KeyFile              string `mapstructure:"keyfile"`
	LogFile              string `mapstructure:"logfile"`
	OutDir               string `mapstructure:"outdir"`
	LicenseLength        int    `mapstructure:"licenselength"`
	DefaultLicenseLength int    `mapstructure:"default_license_term"`
	MaxOfflineDuration   int    `mapstructure:"maximum_Offline_duration"`
}

type Config struct {
	Database Database `mapstructure:"database"`
	Server   Server   `mapstructure:"server"`
}
