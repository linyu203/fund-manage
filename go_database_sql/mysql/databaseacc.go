package main

import (
    "database/sql"
    "fmt"
    "os"
    "github.com/go-sql-driver/mysql"
)

// db is the global database connection pool.
var db *sql.DB

type Fund struct {
    Name string
    Count int
    Description string
    Creation  mysql.NullTime
}

type Description string

type Bond struct {
    Parsekey string
    Creation mysql.NullTime
}

func InitTable() {
    var err error
    db, err = initSocketConnectionPool()
    if err != nil {
	fmt.Println("Init Socket error %v", err)
	return 
    }


    // Create the funds and bonds tables if it does not already exist.
    if _, err = db.Exec(`CREATE TABLE IF NOT EXISTS funds
    ( name varchar(31) NOT NULL, 
      description varchar(101) NOT NULL, 
      creation datetime NOT NULL, 
      PRIMARY KEY(name) );`); err != nil {
        fmt.Println("DB.Exec: unable to create table: %v", err)
    }
    if _, err = db.Exec(`CREATE TABLE IF NOT EXISTS bonds
    ( fundName varchar(31) NOT NULL, 
      parsekey varchar(41) NOT NULL, 
      creation datetime NOT NULL );`); err != nil {
        fmt.Println("DB.Exec: unable to create table: %v", err)
    }
}

func CreateFund(name, description string) error {
    sqlInsert := "INSERT INTO funds (name, description, creation) VALUES (?, ?, NOW())"
    if name != "" && description != "" {
        if _, err := db.Exec(sqlInsert, name, description); err != nil {
            return fmt.Errorf("DB.Exec: %v", err)
        } else {
            fmt.Printf("Create fund successfully fund name %s!, decription %s\n", description)
            return nil
        }
    }
    return fmt.Errorf("fund name and description should not be empty!");
    // [END cloud_sql_mysql_databasesql_connection]
}

func InsertBond(fundName, parsekey string) error {
    sqlInsert := "INSERT INTO bonds (fundName, parsekey, creation) VALUES (?, ?, NOW())"
    if fundName != "" && parsekey != "" {
        if _, err := db.Exec(sqlInsert, fundName, parsekey); err != nil {
            return fmt.Errorf("DB.Exec: %v", err)
        } else {
            fmt.Printf("insert bond %s into fund %s successfully.\n", parsekey, fundName)
            return nil
        }
    }
    return fmt.Errorf("fund name and parsekey should not be empty!");
    // [END cloud_sql_mysql_databasesql_connection]
}

func RemoveBond(fundName, parsekey string) error {
    sqlDelete := "delete from bonds where fundName=? and parsekey=?"
    if fundName != "" && parsekey != "" {
        if _, err := db.Exec(sqlDelete, fundName, parsekey); err != nil {
            return fmt.Errorf("DB.Exec: %v", err)
        } else {
            fmt.Printf("Delete bond %s from fund %s successfully.\n", parsekey, fundName)
            return nil
        }
    }
    return fmt.Errorf("fund name and parsekey should not be empty!");
    // [END cloud_sql_mysql_databasesql_connection]
}

func RemoveFund(fundName string) error {
    sqlDeleteBonds := "delete from bonds where fundName=?"
    sqlDeleteFund := "delete from funds where name=?"
    if fundName != ""{
        if _, err := db.Exec(sqlDeleteBonds, fundName); err != nil {
            return fmt.Errorf("DB.Exec: %v", err)
        } else if _, err := db.Exec(sqlDeleteFund, fundName); err != nil{
            return fmt.Errorf("DB.Exec: %v", err)
        } else {
            fmt.Printf("Delete fund %s successfully.\n", fundName)
            return nil
        }
    }
    return fmt.Errorf("the fund name should not be empty!");
    // [END cloud_sql_mysql_databasesql_connection]
}

func GetAllfunds() ([]Fund, error){
    var funds []Fund
    rows,err := db.Query(`
        select f.name, count(b.parsekey) as count, f.description, f.creation
        from funds as f
        left join bonds as b on b.fundName=f.name
        group by f.name
        `)
    if err != nil {
        return funds, fmt.Errorf("DB.Query: %v", err)
    }
    defer rows.Close()
    for rows.Next() {
        nextfund := Fund{}
        err := rows.Scan(&nextfund.Name, &nextfund.Count, &nextfund.Description, &nextfund.Creation)
        if err != nil {
            return funds, fmt.Errorf("Rows.Scan: %v", err)
        }
        funds = append(funds, nextfund)
    }
    return funds, nil
}

func Getbonds(fundName string) ([]Bond, error){
    var bonds []Bond
    rows,err := db.Query(`select parsekey, creation from bonds where fundName = ?`, fundName)
    if err != nil {
        return bonds, fmt.Errorf("DB.Query: %v", err)
    }
    defer rows.Close()
    for rows.Next() {
        nextbond := Bond{}
        err := rows.Scan(&nextbond.Parsekey, &nextbond.Creation)
        if err != nil {
            return bonds, fmt.Errorf("Rows.Scan: %v", err)
        }
        bonds = append(bonds, nextbond)
    }
    return bonds, nil
}

func GetDescription(fundName string) (string, error){
    var description string
    err := db.QueryRow(`select description from funds where name=?`, fundName).Scan(&description)
    if err != nil {
        return description, fmt.Errorf("DB.Query: %v", err)
    }
    return description, nil
}

// mustGetEnv is a helper function for getting environment variables.
// Displays a warning if the environment variable is not set.
func mustGetenv(k string) string {
    v := os.Getenv(k)
    if v == "" {
        fmt.Printf("Warning: %s environment variable not set.\n", k)
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

func removeEmptyFunds() (error){
    rows,err := db.Query(`
        select tb.name from (
        select f.name, count(b.parsekey) as count, f.description, f.creation
        from funds as f
        left join bonds as b on b.fundName=f.name
        group by f.name
) as tb
where tb.count <= 1
        `)
    if err != nil {
        return fmt.Errorf("DB.Query: %v", err)
    }
    defer rows.Close()
    var name string
    for rows.Next() {
        err := rows.Scan(&name)
        if err != nil {
            return fmt.Errorf("Rows.Scan: %v", err)
        }
	fmt.Println(name)
	if err := RemoveFund(name); err != nil{
	    fmt.Printf("Error Revome fund %s %v",name, err);
	}
    }
    return  nil
}



/*  Test 
func main(){
    InitTable();
    CreateFund("testttFund", "A test fund")
    InsertBond("testttFund", "IBM 4.6 04/20/2020 Corp")
    funds,err := GetAllfunds();
    if err != nil {
	fmt.Printf("error: %v", err)
	return 
    }
    for idx,fund := range(funds) {
	fmt.Println( idx, fund);
    }

}*/
