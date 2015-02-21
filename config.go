package main

import "fmt"
import "code.google.com/p/gcfg"

type Config struct {
	Mtorrent struct {
		UiPort int `gcfg:"ui-port"`
		SavePath string `gcfg:"save-path"`
		WatchPath string `gcfg:"watch-path"`
		SessionFile string `gcfg:"session-file"`
		LockFile string `gcfg:"lock-file"`
	}
	Torrent struct {
		ListenPorts []int `gcfg:"listen"`
		DhtPort int `gcfg:"dht-port"`
		DhtRouters []string `gcfg:"dht-router"`

		UploadLimit int `gcfg:"upload-limit"`
		DownloadLimit int `gcfg:"download-limit"`
		MaxConnections int `gcfg:"max-connections"`
		MaxUploads int `gcfg:"max-uploads"`
	}
}

func GetConfig() (error, Config) {
	var cfg Config
	err := gcfg.ReadFileInto(&cfg, "mtorrent.config")
	if err != nil {
		fmt.Println("Error parsing config: " + err.Error())
	}
	return err, cfg
}
