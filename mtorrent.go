package main

import (
	"log"
	"time"
)

func main() {
  err, cfg := GetConfig()
	if err == nil {
		log.Printf("mtorrent startin on port %d\n", cfg.Mtorrent.UiPort)
		StartSession(cfg)
		for i:=0; i < 10; i++ {
			log.Println(GetTorrentStatus())
			time.Sleep(time.Second)
		}
		StopSession()
	}
}
