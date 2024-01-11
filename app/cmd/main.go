package main

import (
	"github.com/spf13/viper"
	"log"
	"time"
	"webapp/pkg/handler"
	"webapp/pkg/server"
)

var (
	defaultExpiration time.Duration = 10 * time.Second
	cleanupInterval   time.Duration = 10 * time.Second
)

func main() {
	if err := InitConfig(); err != nil {
		log.Fatalf("error init config: %s", err.Error())
	}
	h := new(handler.Handler)
	srv := new(server.Server)
	if err := srv.Run("8080", h.InitRoutes()); err != nil {
		log.Fatalf("error run http server: %s", err.Error())
	}
	//ch := cache.InitCache(defaultExpiration, cleanupInterval)

}

func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
