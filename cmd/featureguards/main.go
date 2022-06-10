package main

import (
	"context"
	"os"

	featureguards "github.com/featureguards/featureguards-go/v2"
	log "github.com/sirupsen/logrus"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	apiKey := os.Getenv("API_KEY")
	ft := featureguards.New(ctx, featureguards.WithApiKey(apiKey), featureguards.WithDefaults(map[string]bool{"BAR": true}))
	on, err := ft.IsOn("TEST")
	if err != nil {
		log.Fatal(err)
	}
	log.Infoln(on)
}
