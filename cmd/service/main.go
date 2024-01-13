package main

import (
	"webapp/internal/app"
	"webapp/internal/config"
	"webapp/internal/pubsub"
	"webapp/internal/server"
	"webapp/internal/storages"
)

func main() {
	// InitConfig
	cfg := config.New()
	cfg.ParseFlags()

	//Init Storages
	ch := storages.InitCache()
	pg := storages.InitPG(cfg)
	defer pg.Close()

	// Init App
	app := app.InitApp(ch, pg)
	app.RestoreCache()

	// Init nats subs
	sub := pubsub.InitNatsCon(cfg, "subscriber")
	sub.SubscribeSubject(app)
	defer sub.Close()

	// Init webserver
	srv := server.InitServer(cfg, ch)
	go srv.ShutdownOnSignal()
	srv.Launch()
}
