package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres" // :(
	dbname   = "homework"
)

// Connect to postgresql database
func Connect() (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable", host, port, dbname, user, password))
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
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
