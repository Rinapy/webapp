package natssub

import (
	"log"

	stan "github.com/nats-io/stan.go"
)

func SendingToChanel() {
	// Устанавливаем параметры подключения к nats-streaming
	clusterID := "test-cluster"
	clientID := "publisher-1"
	natsURL := "nats://localhost:4222"

	// Подключаемся к nats-streaming
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(natsURL))
	if err != nil {
		log.Fatalf("Ошибка подключения к nats-streaming: %v", err)
	}
	defer sc.Close()

	// Рассылаем данные в канал
	channel := "test-channel"
	message := "Пример сообщения"
	err = sc.Publish(channel, []byte(message))
	if err != nil {
		log.Fatalf("Ошибка публикации сообщения: %v", err)
	}

	log.Printf("Сообщение '%s' отправлено в канал '%s'", message, channel)
}
