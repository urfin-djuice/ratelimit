package main

import (
	"github.com/urfin-djuice/ratelimit/pkg/app"
	"log"
)

func main() {
	cmd, err := app.NewApplication()
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}
	sig := cmd.Run()
	<-sig
}
