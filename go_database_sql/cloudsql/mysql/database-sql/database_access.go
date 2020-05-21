// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Sample database-sql demonstrates connection to a Cloud SQL instance from App Engine
// standard. The application is a Golang version of the "Tabs vs Spaces" web
// app presented at Cloud Next '19 as seen in this video:
// https://www.youtube.com/watch?v=qVgzP3PsXFw&t=1833s
package main

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/go-sql-driver/mysql"
)

// db is the global database connection pool.
var db *sql.DB

// parsedTemplate is the global parsed HTML template.
// vote struct contains a single row from the votes table in the database.
// Each vote includes a candidate ("TABS" or "SPACES") and a timestamp.
type vote struct {
	Candidate string
	VoteTime  mysql.NullTime
}

// voteDiff is used to provide a string representation of the current voting
// margin, such as "1 vote" (singular) or "2 votes" (plural).
type voteDiff int


// templateData struct is used to pass data to the HTML template.
type templateData struct {
	TabsCount   uint
	SpacesCount uint
	VoteMargin  string
	RecentVotes []vote
}

func InitTable() {
	var err error

	db, err = initSocketConnectionPool()
	if err != nil {
		log.Fatalf("initSocketConnectionPool: unable to connect: %s", err)
	}

	// Create the funds and bonds tables if it does not already exist.
	if _, err = db.Exec(`CREATE TABLE IF NOT EXISTS funds
	( name varchar(31) NOT NULL, 
	  description varchar(101) NOT NULL, 
	  creation datetime NOT NULL, 
	  RPIMARY KEY (name) );`); err != nil {
		log.Fatalf("DB.Exec: unable to create table: %s", err)
	}
	if _, err = db.Exec(`CREATE TABLE IF NOT EXISTS bonds
	( fundName varchar(31) NOT NULL, 
	  parsekey varchar(41) NOT NULL, 
	  creation datetime NOT NULL );`); err != nil {
		log.Fatalf("DB.Exec: unable to create table: %s", err)
	}

}

// recentVotes returns a slice of the last 5 votes cast.
func recentVotes() ([]vote, error) {
	var votes []vote
	rows, err := db.Query(`SELECT candidate, time_cast FROM votes ORDER BY time_cast DESC LIMIT 5`)
	if err != nil {
		return votes, fmt.Errorf("DB.Query: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		nextVote := vote{}
		err := rows.Scan(&nextVote.Candidate, &nextVote.VoteTime)
		if err != nil {
			return votes, fmt.Errorf("Rows.Scan: %v", err)
		}
		votes = append(votes, nextVote)
	}
	return votes, nil
}

// currentTotals returns a templateData structure for populating the web page.
func currentTotals() (templateData, error) {

	// get total votes for each candidate
	var tabVotes, spaceVotes uint
	err := db.QueryRow(`SELECT count(vote_id) FROM votes WHERE candidate='TABS'`).Scan(&tabVotes)
	if err != nil {
		return templateData{}, fmt.Errorf("DB.QueryRow: %v", err)
	}
	err = db.QueryRow(`SELECT count(vote_id) FROM votes WHERE candidate='SPACES'`).Scan(&spaceVotes)
	if err != nil {
		return templateData{}, fmt.Errorf("DB.QueryRow: %v", err)
	}

	var voteDiffStr string = voteDiff(int(math.Abs(float64(tabVotes) - float64(spaceVotes)))).String()

	latestVotesCast, err := recentVotes()
	if err != nil {
		return templateData{}, fmt.Errorf("recentVotes: %v", err)
	}
	return templateData{tabVotes, spaceVotes, voteDiffStr, latestVotesCast}, nil

}

// showTotals renders an HTML template showing the current vote totals.
func showTotals(w http.ResponseWriter, r *http.Request) error {

	totals, err := currentTotals()
	if err != nil {
		return fmt.Errorf("currentTotals: %v", err)
	}
	err = parsedTemplate.Execute(w, totals)
	if err != nil {
		return fmt.Errorf("Template.Execute: %v", err)
	}
	return nil
}

// saveVote saves a vote passed as http.Request form data.
func saveVote(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return fmt.Errorf("Request.ParseForm: %v", err)
	}

	var team string
	if teamprop, ok := r.Form["team"]; ok {
		team = teamprop[0]
	} else {
		return fmt.Errorf("team property missing from form submission")
	}

	// [START cloud_sql_mysql_databasesql_connection]
	sqlInsert := "INSERT INTO votes (candidate) VALUES (?)"
	if team == "TABS" || team == "SPACES" {
		if _, err := db.Exec(sqlInsert, team); err != nil {
			fmt.Fprintf(w, "unable to save vote: %s", err)
			return fmt.Errorf("DB.Exec: %v", err)
		} else {
			fmt.Fprintf(w, "Vote successfully cast for %s!\n", team)
		}
	}
	return nil
	// [END cloud_sql_mysql_databasesql_connection]
}

// mustGetEnv is a helper function for getting environment variables.
// Displays a warning if the environment variable is not set.
func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Printf("Warning: %s environment variable not set.\n", k)
	}
	return v
}

// initSocketConnectionPool initializes a Unix socket connection pool for
// a Cloud SQL instance of MySQL.
func initSocketConnectionPool() (*sql.DB, error) {
	// [START cloud_sql_mysql_databasesql_create_socket]
	var (
		dbUser                 = mustGetenv("DB_USER")
		dbPwd                  = mustGetenv("DB_PASS")
		instanceConnectionName = mustGetenv("CLOUD_SQL_CONNECTION_NAME")
		dbName                 = mustGetenv("DB_NAME")
	)

	var dbURI string
	dbURI = fmt.Sprintf("%s:%s@unix(/cloudsql/%s)/%s", dbUser, dbPwd, instanceConnectionName, dbName)

	// dbPool is the pool of database connections.
	dbPool, err := sql.Open("mysql", dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %v", err)
	}

	// [START_EXCLUDE]
	configureConnectionPool(dbPool)
	// [END_EXCLUDE]

	return dbPool, nil
	// [END cloud_sql_mysql_databasesql_create_socket]
}


// configureConnectionPool sets database connection pool properties.
// For more information, see https://golang.org/pkg/database/sql
func configureConnectionPool(dbPool *sql.DB) {
	// [START cloud_sql_mysql_databasesql_limit]

	// Set maximum number of connections in idle connection pool.
	dbPool.SetMaxIdleConns(5)

	// Set maximum number of open connections to the database.
	dbPool.SetMaxOpenConns(7)

	// [END cloud_sql_mysql_databasesql_limit]

	// [START cloud_sql_mysql_databasesql_lifetime]

	// Set Maximum time (in seconds) that a connection can remain open.
	dbPool.SetConnMaxLifetime(1800)

	// [END cloud_sql_mysql_databasesql_lifetime]
}
