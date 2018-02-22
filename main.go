package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"flag"
	"regexp"
)

var path string

func main() {
	dataSource := flag.String("data-source", "", "Example: root:root@/htmlacademy")
	pathSource := flag.String("path", "", "Example: /var/migrations")

	flag.Parse()

	if *pathSource != "" {
		path = *pathSource
	}

	if *dataSource == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	db, err := sql.Open("mysql", *dataSource)

	defer db.Close()

	if err != nil {
		panic(err)
	}

	for _, table := range GetTables(db) {
		WriteMigration(table, GetCreateStatement(db, table))
	}
}

func GetFilePath(fileName string) string {
	if path != "" {
		return path +fileName+ ".sql"
	}

	return  "schema/"+fileName+ ".sql"
}

func WriteMigration(fileName string, migration string) {
	_, err := os.Stat(GetFilePath(fileName))

	if os.IsNotExist(err) {
		file, err := os.Create(GetFilePath(fileName))
		defer file.Close()

		file.WriteString(migration)

		if err != nil {
			panic(err)
		}
	} else {
		var file, _ = os.OpenFile(GetFilePath(fileName), os.O_RDWR, 0644)
		defer file.Close()

		file.WriteString(migration)

		file.Sync()
	}
}

func GetCreateStatement(db *sql.DB, table string) string {
	var tableName string
	var createStatement string

	err := db.QueryRow("SHOW CREATE TABLE " + table).Scan(&tableName, &createStatement)

	if err != nil {
		panic(err)
	}

	reg := regexp.MustCompile(`AUTO_INCREMENT=\d+`)
	res := reg.ReplaceAllString(createStatement, "")

	return  res + "\n"
}

func GetTables(db *sql.DB) []string  {
	results, err := db.Query("SHOW TABLES")
	tables := make([]string, 0)

	if err != nil {
		panic(err)
	}

	for results.Next() {
		var tableName string

		results.Scan(&tableName)

		tables = append(tables, tableName)
	}

	return  tables
}