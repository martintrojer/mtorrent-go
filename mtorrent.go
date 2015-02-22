package main

import(
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err, cfg := GetConfig(); err == nil {
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-sigc
			StopSession()
			os.Exit(0)
		}()

		StartSession(cfg)
		StartWebServer(cfg)
	}
}
