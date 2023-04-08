package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
)

var (
	basePath  string
	apiClient *APIClient
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	teardown, err := setup(ctx)
	if err != nil {
		logrus.WithError(err).Fatalf("failed to setup integration tests")
		os.Exit(1)
	}
	exitValue := m.Run()
	teardown(ctx)
	os.Exit(exitValue)
}

func setup(ctx context.Context) (func(context.Context), error) {
	// Parse override flags.
	var host string
	var port int
	flag.StringVar(&host, "targetHost", "", "target API host")
	flag.IntVar(&port, "targetPort", 0, "target API port")
	flag.Parse()
	if host == "" || port == 0 {
		return func(ctx context.Context) {}, errors.New("host can't be empty and port can't be 0")
	}
	basePath = fmt.Sprintf("http://%v:%v", host, port)
	apiClient = NewAPIClient(basePath)
	return func(ctx context.Context) {}, nil
}
