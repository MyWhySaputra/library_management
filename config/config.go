package config

import (
    "database/sql"
    _ "github.com/lib/pq"
)

func ConnectDB(connectionString string) (*sql.DB, error) {
    return sql.Open("postgres", connectionString)
}
