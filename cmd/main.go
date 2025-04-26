package main

import (
	"fmt"
	"os"

	"github.com/Magic-Kot/store/internal/application"
	"github.com/Magic-Kot/store/internal/config"
)

var (
	appName    = "store-backend" //nolint:gochecknoglobals
	appVersion = "v0.0.0"        //nolint:gochecknoglobals
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err = application.New(appName, appVersion, cfg).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
