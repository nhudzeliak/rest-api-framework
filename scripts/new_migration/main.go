package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/nataliia_hudzeliak/rest-api-framework/app/config"
)

func main() {
	// Parse flags.
	var name string
	flag.StringVar(&name, "name", "", "migration name")
	flag.Parse()
	if name == "" {
		panic("migration name can't be empty")
	}

	// Connect to the database.
	fullPath := config.BasePath() + fmt.Sprintf("/db/migrations/%v_%v", time.Now().UTC().UnixNano(), name)
	os.Create(fullPath + ".up.sql")
	os.Create(fullPath + ".down.sql")
}
