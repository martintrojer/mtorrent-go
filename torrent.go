package main

import (
	"os"
	"log"
	"strings"
	"regexp"
	"bytes"
	"github.com/steeve/libtorrent-go"
	"github.com/dustin/go-humanize"
)

type Instance struct {
	session libtorrent.Session
	config Config
	handles map[string]libtorrent.Torrent_handle
}

var instance = Instance{}

// ---------------------------------------------------------------
// Torrents

func AddMagnet(uri string) {
	r, _ := regexp.Compile("[a-fA-F0-9]{40}")
	if (! r.MatchString(uri)) { return }

	infoHash := strings.ToLower(r.FindStringSubmatch(uri)[0])

	_, present := instance.handles[infoHash]

	if (present) {
		log.Println(infoHash + " already added!")
		return
	}

	params := libtorrent.NewAdd_torrent_params()
	params.SetUrl(uri)
	params.SetSave_path(instance.config.Mtorrent.SavePath)

	handle := instance.session.Add_torrent(params)
	instance.handles[infoHash] = handle

	handle.Set_max_connections(instance.config.Torrent.MaxConnections)
	handle.Set_max_uploads(instance.config.Torrent.MaxUploads)
	handle.Set_upload_limit(instance.config.Torrent.UploadLimit)
	handle.Set_download_limit(instance.config.Torrent.DownloadLimit)
	handle.Auto_managed(true)

	log.Println("Magnet added: " + infoHash)
}

func RemoveTorrent(infoHash string) {
	handle, present := instance.handles[infoHash]
	if (present) {
		instance.session.Remove_torrent(handle)
		delete(instance.handles, infoHash)
		log.Println("Torrent removed: " + infoHash)
	}
}

func PauseTorrent(infoHash string) {
	handle, present := instance.handles[infoHash]
	if (present) {
		handle.Auto_managed(false)
		handle.Pause()
		log.Println("Torrent paused: " + infoHash)
	}
}

func ResumeTorrent(infoHash string) {
	handle, present := instance.handles[infoHash]
	if (present) {
		handle.Auto_managed(true)
		handle.Resume()
		log.Println("Torrent paused: " + infoHash)
	}
}

type TorrentStatus struct {
	name string
	hash string
	size string
	state string
	progress float32
	downRate string
	upRate string
	seeds int
	seedsTotal int
	peers int
	peersTotal int
	isPaused bool
}

func toHex(hash []byte)string {
	hexChars := []string {"0","1","2","3","4","5","6","7","8","9","a","b","c","d","e","f",}
	res := ""
	for _, ch := range hash {
		res += hexChars[ch >> 4]
		res += hexChars[ch & 0xf]
	}
	return res
}

func GetTorrentStatus() []TorrentStatus{
	handlesVec := instance.session.Get_torrents()
	vecSize := int(handlesVec.Size())
	res := make([]TorrentStatus, vecSize)

	states := []string {
		"queued for checking",
		"checking files",
		"downloading metadata",
		"downloading",
		"finished",
		"seeding",
		"allocating",
		"checking resume data",
	}

	for i := 0; i < vecSize; i++ {
		st := handlesVec.Get(i).Status()

		var ts TorrentStatus
		ts.name = st.GetName()
		ts.hash = toHex([]byte(st.GetInfo_hash().To_string()))
		ts.size = humanize.Bytes(uint64(st.GetTotal_wanted()))
		ts.state = states[st.GetState()]
		ts.progress = 100.0 * st.GetProgress()
		ts.downRate = humanize.Bytes(uint64(st.GetDownload_rate()))
		ts.upRate = humanize.Bytes(uint64(st.GetUpload_rate()))
		ts.seeds = st.GetNum_seeds()
		ts.seedsTotal = st.GetList_seeds()
		ts.peers = st.GetNum_peers()
		ts.peersTotal = st.GetList_peers()
		ts.isPaused = st.GetPaused()

		res[i] = ts

	}

	return res
}

// ---------------------------------------------------------------
// Session

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

func restoreSessionState() {
	fi, err := os.Open(instance.config.Mtorrent.SessionFile)
	if (err != nil) { return }
	defer fi.Close()

	buf := make([]byte, 1024*4)
	_, err2 := fi.Read(buf);
	if (err2 != nil) { return }

	n := bytes.Index(buf, []byte{0})

	var entry = libtorrent.NewLazy_entry()
	libtorrent.Lazy_bdecode(string(buf[:n]), entry)
	instance.session.Load_state(entry)
}

func saveSessionState() {
	fo, err := os.Create(instance.config.Mtorrent.SessionFile)
	if (err != nil) { return }
	defer fo.Close()

	var entry = libtorrent.NewEntry()
	instance.session.Save_state(entry)
	fo.Write([]byte(libtorrent.Bencode(entry)))
}

func startServices() {
	log.Println("Starting DHT, LSD, UPNP, NATPMP...")
	instance.session.Start_dht()
	instance.session.Start_lsd()
	instance.session.Start_upnp()
	instance.session.Start_natpmp()
}

func stopServices() {
	log.Println("Stopping DHT, LSD, UPNP, NATPMP...")
	instance.session.Stop_dht()
	instance.session.Stop_lsd()
	instance.session.Stop_upnp()
	instance.session.Stop_natpmp()
}

func StartSession(cfg Config) {
	instance.config = cfg
	instance.handles = make(map[string]libtorrent.Torrent_handle)

	s := libtorrent.NewSession()
	instance.session = s
	s.Add_extensions()

	restoreSessionState()

	log.Println("Adding DHT Routers...")
	for _, router := range cfg.Torrent.DhtRouters {
		s.Add_dht_router(libtorrent.NewStd_pair_string_int(router, cfg.Torrent.DhtPort))
	}

	var ports = cfg.Torrent.ListenPorts
	var error = libtorrent.NewError_code()
	s.Listen_on(libtorrent.NewStd_pair_int_int(ports[0], ports[1]), error)

	configureSession()
	startServices()
}

func StopSession() {
	stopServices()
	saveSessionState()
}
