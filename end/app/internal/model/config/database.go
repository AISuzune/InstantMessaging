package config

import "fmt"

type Database struct {
	Mysql *Mysql `mapstructure:"mysql" yaml:"mysql"`
	Redis *Redis `mapstructure:"redis" yaml:"redis"`
}

type Mysql struct {
	Addr     string `mapstructure:"addr" yaml:"addr"`
	Port     string `mapstructure:"port" yaml:"port"`
	Db       string `mapstructure:"db" yaml:"db"`
	Username string `mapstructure:"username" yaml:"username"`
	Password string `mapstructure:"password" yaml:"password"`
	Charset  string `mapstructure:"charset" yaml:"charset"`
}

func (m *Mysql) GetDsn() string {
	return fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		m.Username,
		m.Password,
		m.Addr,
		m.Port,
		m.Db,
		m.Charset)
}

type Redis struct {
	Addr     string `mapstructure:"addr" yaml:"addr"`
	Port     string `mapstructure:"port" yaml:"port"`
	Password string `mapstructure:"password" yaml:"password"`
	Db       int    `mapstructure:"db" yaml:"db"`
}
