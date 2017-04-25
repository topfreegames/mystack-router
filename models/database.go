// mystack-router api
// https://github.com/topfreegames/mystack-router
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/jmoiron/sqlx"
)

//GetDB Connection using the given properties
func GetDB(
	host, user, dbName, password string,
	port, connectionTimeoutMS int,
) (*sqlx.DB, error) {
	if connectionTimeoutMS <= 0 {
		connectionTimeoutMS = 100
	}
	connStr := fmt.Sprintf(
		"host=%s user=%s port=%d sslmode=disable dbname=%s connect_timeout=2",
		host, user, port, dbName,
	)
	if password != "" {
		connStr += fmt.Sprintf(" password=%s", password)
	}

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, err
	}

	shouldPing(db.DB, time.Duration(connectionTimeoutMS)*time.Millisecond)

	return db, nil
}

//ShouldPing the database
func shouldPing(db *sql.DB, timeout time.Duration) error {
	var err error
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = timeout
	ticker := backoff.NewTicker(b)

	// Ticks will continue to arrive when the previous operation is still running,
	// so operations that take a while to fail could run in quick succession.
	for range ticker.C {
		if err = db.Ping(); err != nil {
			continue
		}

		ticker.Stop()
		return nil
	}

	return fmt.Errorf("could not ping database")
}
