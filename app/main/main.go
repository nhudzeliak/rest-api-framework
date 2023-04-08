package main

import (
	"flag"
	"strconv"

	"github.com/nataliia_hudzeliak/rest-api-framework/app/config"
	"github.com/nataliia_hudzeliak/rest-api-framework/app/controllers"
	"github.com/nataliia_hudzeliak/rest-api-framework/app/services/api"

	"github.com/sirupsen/logrus"
)

func main() {
	// Load configs.
	cfg := config.MustConfig()
	logrus.Infof("starting the app with %v env", cfg["envname"])

	// Parse override flags.
	var host string
	var port int
	flag.StringVar(&host, "host", "", "override host, defaulted to config if not provided")
	flag.IntVar(&port, "port", 0, "override port, defaulted to config if not provided")
	flag.Parse()
	if host == "" {
		host = cfg["self.host"]
	}
	if port == 0 {
		port, _ = strconv.Atoi(cfg["self.port"])
	}
	logrus.Infof("starting the app at %v:%v", host, port)

	// Instantiate the api server.
	controllers.MustInitialize()
	server := api.NewServer(host, port, controllers.Controllers)
	err := server.Run()
	if err != nil {
		logrus.WithError(err).Fatal("server existed with error")
	}
	logrus.Info("server shut down gracefully")
}
