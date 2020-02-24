//https://github.com/microsoft/sql-server-samples/blob/master/samples/tutorials/go/crud.go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
)

func TestConnection() {
	condb := GetConnection()

	var (
		ID   int
		Name string
	)
	rows, err := condb.Query("select * from dbo.TestTable")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		err := rows.Scan(&ID, &Name)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(ID)
		log.Println(Name)
	}

	defer condb.Close()
}

func CreateMasterRecord(db *sql.DB, title string, company string, place string) (int64, error) {
	tsql := fmt.Sprintf("INSERT INTO dbo.Jobs (Title,Company,Place) VALUES (N'%s',N'%s',N'%s') SELECT SCOPE_IDENTITY();",
		strings.Trim(title, " "), strings.Trim(company, " "), strings.Trim(place, " "))
	var id int64
	err := db.QueryRow(tsql).Scan(&id)

	return int64(id), err
}

func CreateDetailRecord(db *sql.DB, jobId int, key string, value string) (int64, error) {
	tsql := fmt.Sprintf("INSERT INTO dbo.JobDetails (JobId,[key],[value]) VALUES (%d,N'%s',N'%s');",
		jobId, strings.Trim(key, " "), strings.Trim(value, " "))
	var id int64
	err := db.QueryRow(tsql).Scan(&id)

	return id, err
}

func GetConnection() *sql.DB {
	condb, errdb := sql.Open("mssql", "server=.;port=1433;database=JobDb;")
	if errdb != nil {
		fmt.Println(" Error open db:", errdb.Error())
	}
	return condb
}
