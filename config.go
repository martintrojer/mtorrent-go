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
		MaxConnections int `gcfg:"max-connections"`
		MaxUploads int `gcfg:"max-uploads"`
		Ratio float32
		UploadLimit float32 `gcfg:"upload-limit"`
		DownloadLimit float32 `gcfg:"download-limit"`
		ResolveCountries bool `gcfg:"resolve-countries"`
		DhtPort int `gcfg:"dht-port"`
		DhtRouters []string `gcfg:"dht-router"`
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
