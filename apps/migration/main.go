package main

import (
	"app/config"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"text/template"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
)

var templateContent = `
-- +migrate Up

-- +migrate Down
`
var tpl = template.Must(template.New("new_migration").Parse(templateContent))

func init() {
	godotenv.Load()
}

func main() {
	config := config.InitConfig()

	if len(os.Args) < 2 {
		fmt.Println("Missing parameter, provide action!")
		return
	}

	postgresDsn := config.GetPostgresDSN()
	db, err := sql.Open("postgres", postgresDsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	migrations := &migrate.FileMigrationSource{
		Dir: "migration",
	}
	migrate.SetTable("migrations")

	switch os.Args[1] {
	case "new":
		if len(os.Args) < 3 {
			log.Println("Name is required for new action!")
			return
		}
		name := os.Args[2]

		if _, err := os.Stat("./migration"); os.IsNotExist(err) {
			panic(err)
		}

		fileName := fmt.Sprintf("%s-%s.sql", time.Now().Format("20060102150405"), strings.TrimSpace(name))
		pathName := path.Join("./migration", fileName)
		f, err := os.Create(pathName)
		if err != nil {
			panic(err)
		}
		defer func() { _ = f.Close() }()

		if err := tpl.Execute(f, nil); err != nil {
			panic(err)
		}

		fmt.Printf("Created migration %s\n", pathName)
	case "up":
		m, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Applied %d migrations!\n", m)
	case "down":
		m, err := migrate.ExecMax(db, "postgres", migrations, migrate.Down, 1)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Applied %d migrations!\n", m)
	case "status":
		m, err := migrate.GetMigrationRecords(db, "postgres")
		if err != nil {
			panic(err)
		}
		for _, record := range m {
			fmt.Printf("%s - %s\n", record.Id, record.AppliedAt)
		}
	}
}
