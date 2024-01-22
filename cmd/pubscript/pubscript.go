package main

import (
	"flag"
	"log"
	"os"
	"webapp/internal/config"
	"webapp/internal/pubsub"
)

// publisher это отдельный скрипт, для публикации данных в канал.
// при указании флага -j отправляет в nats-stream один указанный json-файл.
// по умолчанию, отправляет в nats-stream сразу все три проверочных json'а,
// два из которых имеют корректный order_uid и успешно считываются, третий - бракуется.

func main() {
	cfg := config.New()
	natsAddr, jsnFileName := parseFlags()
	cfg.Nats = natsAddr

	pub := pubsub.InitNatsCon(cfg, "publisher")
	defer pub.Close()

	if jsnFileName != "" {
		jdata := readJsnDataFromFile(jsnFileName)
		pub.Publish(jdata)
		return
	}

	// если программа была запущена без флага -j, тогда отправляю
	// все три заказа из локальных переменных ниже:
	pub.Publish([]byte(jsn1good))
	pub.Publish([]byte(jsn2bad))
	pub.Publish([]byte(jsn3good))
}

func parseFlags() (natsAddr, jsnFileName string) {
	flag.StringVar(&natsAddr, "n", "nats://localhost:4222", "json file name to publish")
	flag.StringVar(&jsnFileName, "j", "", "json file name to publish")
	flag.Parse()
	return natsAddr, jsnFileName
}

func readJsnDataFromFile(fpath string) []byte {
	data, err := os.ReadFile(fpath)
	if err != nil {
		log.Fatal("couldn't read file:", fpath)
	}
	return data
}

var jsn1good = `
	{"order_uid": "sdasd3es231ad23123dastest",
	"track_number": "WBILMTESTTRACK",
	"entry": "WBIL",
	"delivery": {
	"name": "Test Testov",
	"phone": "+9720000000",
	"zip": "2639809",
	"city": "Kiryat Mozkin",
	"address": "Ploshad Mira 15",
	"region": "Kraiot",
	"email": "test@gmail.com"
	},
	"payment": {
	"transaction": "b563feb7b2b84b6test",
	"request_id": "",
	"currency": "USD",
	"provider": "wbpay",
	"amount": 1817,
	"payment_dt": 1637907727,
	"bank": "alpha",
	"delivery_cost": 1500,
	"goods_total": 317,
	"custom_fee": 0
	},
	"items": [
	{
		"chrt_id": 9934930,
		"track_number": "WBILMTESTTRACK",
		"price": 453,
		"rid": "ab4219087a764ae0btest",
		"name": "Mascaras",
		"sale": 30,
		"size": "0",
		"total_price": 317,
		"nm_id": 2389212,
		"brand": "Vivienne Sabo",
		"status": 202
	}
	],
	"locale": "en",
	"internal_signature": "",
	"customer_id": "test",
	"delivery_service": "meest",
	"shardkey": "9",
	"sm_id": 99,
	"date_created": "2021-11-26T06:22:19Z",
	"oof_shard": "1"
	}
	`

var jsn2bad = `
	{
	"track_number": "WBILMTESTTRACK",
	"entry": "WBIL",
	"delivery": {
	  "name": "Test Testov",
	  "phone": "+9720000000",
	  "zip": "2639809",
	  "city": "Kiryat Mozkin",
	  "address": "Ploshad Mira 15",
	  "region": "Kraiot",
	  "email": "test@gmail.com"
	},
	"payment": {
	  "transaction": "b563feb7b2b84b6test",
	  "request_id": "",
	  "currency": "USD",
	  "provider": "wbpay",
	  "amount": 1817,
	  "payment_dt": 1637907727,
	  "bank": "alpha",
	  "delivery_cost": 1500,
	  "goods_total": 317,
	  "custom_fee": 0
	},
	"items": [
	  {
		"chrt_id": 9934930,
		"track_number": "WBILMTESTTRACK",
		"price": 453,
		"rid": "ab4219087a764ae0btest",
		"name": "Mascaras",
		"sale": 30,
		"size": "0",
		"total_price": 317,
		"nm_id": 2389212,
		"brand": "Vivienne Sabo",
		"status": 202
	  }
	],
	"locale": "en",
	"internal_signature": "",
	"customer_id": "test",
	"delivery_service": "meest",
	"shardkey": "9",
	"sm_id": 99,
	"date_created": "2021-11-26T06:22:19Z",
	"oof_shard": "1"
  }`

var jsn3good = `
	{"order_uid": "aasdasds3123fef1231sd7f2f84f6test",
	"track_number": "WBILMTESTTRACK",
	"entry": "WBIL",
	"delivery": {
	  "name": "Test Testov",
	  "phone": "+9720000000",
	  "zip": "2639809",
	  "city": "Kiryat Mozkin",
	  "address": "Ploshad Mira 15",
	  "region": "Kraiot",
	  "email": "test@gmail.com"
	},
	"payment": {
	  "transaction": "b563feb7b2b84b6test",
	  "request_id": "",
	  "currency": "USD",
	  "provider": "wbpay",
	  "amount": 1817,
	  "payment_dt": 1637907727,
	  "bank": "alpha",
	  "delivery_cost": 1500,
	  "goods_total": 317,
	  "custom_fee": 0
	},
	"items": [
	  {
		"chrt_id": 9934930,
		"track_number": "WBILMTESTTRACK",
		"price": 453,
		"rid": "ab4219087a764ae0btest",
		"name": "Mascaras",
		"sale": 30,
		"size": "0",
		"total_price": 317,
		"nm_id": 2389212,
		"brand": "Vivienne Sabo",
		"status": 202
	  }
	],
	"locale": "en",
	"internal_signature": "",
	"customer_id": "test",
	"delivery_service": "meest",
	"shardkey": "9",
	"sm_id": 99,
	"date_created": "2021-11-26T06:22:19Z",
	"oof_shard": "1"
  }`
