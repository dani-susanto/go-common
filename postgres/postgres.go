package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/XSAM/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

func New(
	username string,
	password string,
	host string,
	port string,
	database string,
) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		username,
		password,
		host,
		port,
		database,
	)

	db, err := otelsql.Open(
		"postgres",
		dsn,
		otelsql.WithAttributes(semconv.DBSystemPostgreSQL),
		otelsql.WithSpanOptions(
			otelsql.SpanOptions{
				OmitConnResetSession: true,
				OmitRows:             true,
				Ping:                 false,
			},
		),
	)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}
