package config

type Config struct {
	DataBase *Database `mapstructure:"database"  yaml:"database"`
	Logger   *Logger   `mapstructure:"logger" yaml:"logger"`
	Server   *Server   `mapstructure:"server"  yaml:"server"`
	Cors     CORS      `mapstructure:"cors" yaml:"cors"`
}
