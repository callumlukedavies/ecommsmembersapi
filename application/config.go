package application

import (
	"os"
	"strconv"

	"github.com/go-sql-driver/mysql"
)

type Config struct {
	ServerPort  uint16
	MySQLConfig mysql.Config
}

func LoadConfig() Config {
	cfg := Config{
		ServerPort: 8080,
		MySQLConfig: mysql.Config{
			User:                 "root",
			Passwd:               "Orangeseven7byte!",
			Net:                  "tcp",
			Addr:                 "127.0.0.1:3306",
			DBName:               "itemsdb",
			AllowNativePasswords: true,
		},
	}

	if serverPort, exists := os.LookupEnv("PORT"); exists {
		if portAddr, err := strconv.ParseUint(serverPort, 10, 16); err != nil {
			cfg.ServerPort = uint16(portAddr)
		}
	}

	if mysqlPort, exists := os.LookupEnv("SQL_PORT"); exists {
		cfg.MySQLConfig.Addr = mysqlPort
	}

	return cfg
}
