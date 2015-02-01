package main

import "github.com/steeve/libtorrent-go"

type Instance struct {
	session libtorrent.Session
	handles []libtorrent.Torrent_handle
}

var instance = Instance{}

func Init() {
	instance.session = libtorrent.NewSession();
}
