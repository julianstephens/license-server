package model

type Database struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DB       string `mapstructure:"db"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
}

type Server struct {
	Host    string `mapstructure:"host"`
	Port    int    `mapstructure:"port"`
	Env     string `mapstructure:"env"`
	KeyFile string `mapstructure:"keyfile"`
	OutDir  string `mapstructure:"outdir"`
}

type Config struct {
	Database Database `mapstructure:"database"`
	Server   Server   `mapstructure:"server"`
}
