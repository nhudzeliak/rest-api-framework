package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strings"

	"github.com/nataliia_hudzeliak/rest-api-framework/app/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// version is a wrapper for version migration.
type version struct {
	ID string `gorm:"column:id; primary_key:yes"`
}

func main() {
	cfg := config.MustConfig()
	logrus.Infof("starting the app with %v env", cfg["envname"])
	defaultDsn := cfg["database.writer"]

	// Parse override flags.
	var dsn string
	flag.StringVar(&dsn, "dsn", "", "override data source name, defaulted to config if not provided")
	flag.Parse()
	if dsn == "" {
		dsn = defaultDsn
	}

	// Read applied migrations.
	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	var appliedMigrations []version
	db.Table("versions").Find(&appliedMigrations)
	appliedMigrationsHash := make(map[string]struct{}, len(appliedMigrations))
	for _, m := range appliedMigrations {
		appliedMigrationsHash[m.ID] = struct{}{}
	}

	// Find un-applied migrations and run them.
	sourceFolder := config.BasePath() + "/db/migrations/"
	files, _ := ioutil.ReadDir(sourceFolder)
	for _, f := range files {
		if strings.HasPrefix(f.Name(), "_") {
			logrus.Infof("skipped migration %v", f.Name())
			continue
		}
		if strings.HasSuffix(f.Name(), ".down.sql") {
			continue
		}
		id := strings.Split(f.Name(), "_")[0]
		if _, ok := appliedMigrationsHash[id]; ok {
			logrus.Infof("skipped migration %v", f.Name())
			continue
		}
		query, _ := os.ReadFile(sourceFolder + f.Name())
		err := db.Exec(string(query)).Error
		if err != nil {
			query, _ = os.ReadFile(sourceFolder + strings.Replace(f.Name(), ".up.sql", ".down.sql", 1))
			err = db.Exec(string(query)).Error
			if err != nil {
				panic(err)
			}
			logrus.Errorf("reverted migration %v", f.Name())
			continue
		}
		logrus.Infof("applied migration %v", f.Name())
		db.Table("versions").Create(&version{ID: id})
	}
}
