package db

import (
	"database/sql"
	"fmt"
	"io/ioutil"

	_ "github.com/lib/pq"
	"gopkg.in/yaml.v2"
)

type DBConf struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBname   string `yaml:"dbname`
}

func getConfig(configPath string) (*DBConf, error) {
	dbconf := DBConf{}
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &dbconf)
	if err != nil {
		return nil, err
	}
	return &dbconf, nil
}

// Connect to postgresql database
func Connect(dbConfPath string) (*sql.DB, error) {
	conf, err := getConfig(dbConfPath)
	if err != nil {
		return nil, fmt.Errorf("Could not open database config file, %s", err)
	}
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		conf.Host, conf.Port, conf.DBname, conf.User, conf.Password))
	if err != nil {
		return nil, fmt.Errorf("Unable connect to database %s", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Unable connect to database %s", err)
	}
	return db, nil
}

// GetCPUQuery represents a sql query to calculate max and min cpu usage
func GetCPUQuery(db *sql.DB) (*sql.Stmt, error) {
	stmt, err := db.Prepare(`
		SELECT
			max(usage) AS max_cpu_usage, min(usage) AS min_cpu_usage
		FROM cpu_usage
		WHERE
			host=$1 AND ts >= $2 AND ts <= $3
		GROUP BY date_trunc('minute', ts)
	`)
	return stmt, err
}
