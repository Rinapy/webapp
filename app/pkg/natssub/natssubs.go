package natssub

import (
	stan "github.com/nats-io/stan.go"
	"log"
)

func SubscribeChanel() error {
	// Устанавливаем параметры подключения к nats-streaming
	clusterID := "test-cluster"
	clientID := "client-1"
	natsURL := "nats://localhost:4222"

	// Подключаемся к nats-streaming
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(natsURL))
	if err != nil {
		log.Fatalf("Ошибка подключения к nats-streaming: %v", err)
	}

	// Подписываемся на канал
	channel := "test-channel"
	_, err = sc.Subscribe(channel, func(msg *stan.Msg) {
		log.Printf("Получено сообщение: %s", string(msg.Data))
	}, stan.DurableName("durable-subscriber"))
	if err != nil {
		log.Fatalf("Ошибка подписки на канал: %v", err)
	}

	// Ждем завершения работы
	log.Printf("Подписка на канал '%s' запущена", channel)
	select {}
}
