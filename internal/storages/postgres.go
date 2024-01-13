package storages

import (
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"log"
	"webapp/internal/config"
	"webapp/internal/models"
)

const tabel = `
		CREATE TABLE IF NOT EXISTS orders (
			id 		   SERIAL PRIMARY KEY,
			order_uid  VARCHAR(100) UNIQUE,
			order_data JSONB);
		CREATE INDEX IF NOT EXISTS oid ON orders(order_uid);`

type PG struct {
	Sqlx *sqlx.DB
	Cfg  config.Config
}

func InitPG(cfg config.Config) *PG {
	conn, err := sqlx.Open("postgres", cfg.PgConnStr)
	if err != nil {
		log.Fatalf("(Postgres): error at connecting to database: %s", err.Error())
	}

	pg := PG{conn, cfg}

	if _, err := pg.Sqlx.Exec(tabel); err != nil {
		log.Printf("(Postgres): error at create tabel: %s", err.Error())
		log.Fatal(err)
	}
	return &pg
}

func (pg *PG) Close() error {
	return pg.Sqlx.Close()
}

func isNotUniqueInsert(err error) bool {
	if err, ok := err.(*pq.Error); ok && err.Code == "23505" {
		return true
	}
	return false
}

func (pg *PG) InsertOrder(ord models.Order) error {
	query := "INSERT INTO orders (order_uid, order_data) VALUES(:order_uid, :order_data)"
	_, err := pg.Sqlx.NamedExec(query, ord)
	if err != nil {
		if isNotUniqueInsert(err) {
			log.Printf("(Postgres): order_uid:%s already exists", ord.ID)
			return err
		}
		log.Printf("(Postgres): unexpected error: %s %v \n", ord.ID, err)
		return err
	}
	return nil
}

func (pg *PG) SelectAllOrders() (orders []models.Order) {
	q := "SELECT order_uid, order_data FROM orders"
	err := pg.Sqlx.Select(&orders, q)
	if err != nil {
		log.Println("pg.SelectAllOrders error:", err)
	}
	return orders
}
