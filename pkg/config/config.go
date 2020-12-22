package config

import "fmt"

type MysqlConfig struct {
	Host     string `default:"localhost"`
	Port     string `default:"3306"`
	User     string `default:"root"`
	Password string `default:"root"`
	Database string `default:"biubiu"`
	Disable  bool   `default:"false"`
}

func (m *MysqlConfig) GetUrl() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", m.User, m.Password, m.Host, m.Port, m.Database)
}
