package main

import (
	"log"
	"fmt"
	"sort"
	"net/http"
	"html/template"
)

type PageData struct {
	Title string
	TorrentStatuses []TorrentStatus
}

var indexTemplates = template.Must(template.ParseFiles(
	"templates/index.html",
	"templates/partials/head.html",
	"templates/partials/scripts.html",
	"templates/partials/navbar.html",
	"templates/partials/torrents.html",
	"templates/partials/add.html"))

func indexHandler(w http.ResponseWriter, r *http.Request) {
	st := GetTorrentStatus()
	sort.Sort(ByName(st))
	err := indexTemplates.Execute(w, PageData{"mtorrent-go", st})
	if (err != nil) { log.Println("Template error: ", err) }
}

func home (w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusFound)
}

func StartWebServer(cfg Config) {
	log.Printf("mtorrent listening on port %d\n", cfg.Mtorrent.UiPort)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/magnet", func (w http.ResponseWriter, r *http.Request) { AddMagnet(r.FormValue("magnet")); home(w,r) })
	http.HandleFunc("/pause", func (w http.ResponseWriter, r *http.Request) { PauseTorrent(r.FormValue("id")); home(w,r) })
	http.HandleFunc("/resume", func (w http.ResponseWriter, r *http.Request) { ResumeTorrent(r.FormValue("id")); home(w,r)})
	http.HandleFunc("/remove", func (w http.ResponseWriter, r *http.Request) { RemoveTorrent(r.FormValue("id")); home(w,r)})
	http.HandleFunc("/pause-all", func (w http.ResponseWriter, r *http.Request) { PauseAllTorrents(); home(w,r) })
	http.HandleFunc("/resume-all", func (w http.ResponseWriter, r *http.Request) { ResumeAllTorrents(); home(w,r)})
	http.HandleFunc("/remove-all", func (w http.ResponseWriter, r *http.Request) { RemoveAllTorrents(); home(w,r)})

	http.HandleFunc("/", indexHandler)

	http.ListenAndServe(fmt.Sprintf(":%d", cfg.Mtorrent.UiPort), nil)
}
