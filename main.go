package main

import (
	"log"
	"os"

	"golang_strarter_kit_2025/app/database"
	"golang_strarter_kit_2025/bootstrap"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "Golang Starter Kit",
		Usage: "CLI tool for managing migrations",
		Commands: []*cli.Command{
			database.MigrationCmd, // make:migration
		},
	}

	// Jalankan CLI jika ada args
	if len(os.Args) > 1 {
		if err := app.Run(os.Args); err != nil {
			log.Fatal(err)
		}
		return
	}

	// Jalankan server biasa jika tidak ada CLI arg
	bootstrap.Init()
}
