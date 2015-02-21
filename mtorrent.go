package main

import (
	"log"
)

func main() {
  err, cfg := GetConfig()
	if err == nil {
		log.Printf("mtorrent startin on port %d\n", cfg.Mtorrent.UiPort)
		StartSession(cfg)
		StopSession()
	}
}
