package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"pantrypal/backend/internal/platform/db"
)

func main() {
	dbPath := flag.String("db", "../database/sqlite/pantrypal.db", "path to sqlite db file")
	migrationPath := flag.String("migration", "./migrations/001_init_schema.sql", "path to migration sql file")
	seedPath := flag.String("seed", "", "optional path to seed sql file")
	reset := flag.Bool("reset", false, "delete db file before applying migration + seed")
	flag.Parse()

	if *reset {
		_ = os.Remove(*dbPath)
	}

	conn, err := db.Open(*dbPath)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer conn.Close()

	if err := db.ApplySQLFile(conn, *migrationPath); err != nil {
		log.Fatal(err)
	}
	if strings.TrimSpace(*seedPath) != "" {
		if err := db.ApplySQLFile(conn, *seedPath); err != nil {
			log.Fatal(err)
		}
	}

	log.Printf("database bootstrapped at %s", *dbPath)
}
