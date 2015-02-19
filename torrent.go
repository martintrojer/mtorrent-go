package main

import (
	"log"
	"github.com/steeve/libtorrent-go"
)

type Instance struct {
	session libtorrent.Session
	handles []libtorrent.Torrent_handle
}

var instance = Instance{}

func configureSession() {
	settings := instance.session.Settings()

	log.Println("Setting Session settings...")

	// settings.SetUser_agent("")

	settings.SetRequest_timeout(5)
	settings.SetPeer_connect_timeout(2)
	settings.SetAnnounce_to_all_trackers(true)
	settings.SetAnnounce_to_all_tiers(true)
	settings.SetConnection_speed(100)

	settings.SetTorrent_connect_boost(100)
	settings.SetRate_limit_ip_overhead(true)

	instance.session.Set_settings(settings)

	log.Println("Setting Encryption settings...")
	encryptionSettings := libtorrent.NewPe_settings()
	encryptionSettings.SetPrefer_rc4(true)
	instance.session.Set_pe_settings(encryptionSettings)
}

func startServices() {
	log.Println("Starting DHT...")
	instance.session.Start_dht()

	log.Println("Starting LSD...")
	instance.session.Start_lsd()

	log.Println("Starting UPNP...")
	instance.session.Start_upnp()

	log.Println("Starting NATPMP...")
	instance.session.Start_natpmp()
}

func stopServices() {
	log.Println("Stopping DHT...")
	instance.session.Stop_dht()

	log.Println("Stopping LSD...")
	instance.session.Stop_lsd()

	log.Println("Stopping UPNP...")
	instance.session.Stop_upnp()

	log.Println("Stopping NATPMP...")
	instance.session.Stop_natpmp()
}

func StartSession(cfg Config) {
	s := libtorrent.NewSession()
	instance.session = s
	s.Add_extensions()

	log.Println("Adding DHT Routers...")
	for _, router := range cfg.Torrent.DhtRouters {
		s.Add_dht_router(libtorrent.NewStd_pair_string_int(router, cfg.Torrent.DhtPort))
	}

	var ports = cfg.Torrent.ListenPorts
	var error = libtorrent.NewError_code()
	s.Listen_on(libtorrent.NewStd_pair_int_int(ports[0], ports[1]), error)

	startServices()

}

func StopSession(cfg Config) {
	stopServices()
}
