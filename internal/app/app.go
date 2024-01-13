package app

import (
	"encoding/json"
	"log"
	"webapp/internal/models"
	"webapp/internal/storages"
)

type App struct {
	ch *storages.Cache
	pg *storages.PG
}

func InitApp(ch *storages.Cache, pg *storages.PG) *App {
	return &App{ch, pg}
}

func (app *App) RestoreCache() {
	orders := app.pg.SelectAllOrders()
	app.ch.AddOrders(orders)
}

func (app *App) ProcNatsMessages(data []byte) {
	log.Println("Processing gets new NATS message")
	order := models.Order{}
	if err := json.Unmarshal(data, &order); err != nil {
		log.Println("(App): error unmarshalling the data", err)
		return
	}
	if len(order.ID) < 1 {
		log.Println("(App): error - order_uid tag was not found in data")
		return
	}
	order.Data = data
	if err := app.pg.InsertOrder(order); err != nil {
		app.ch.Add(order.ID, order.Data)
	}
}
