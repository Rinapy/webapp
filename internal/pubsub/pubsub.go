package pubsub

import (
	"github.com/nats-io/stan.go"
	"log"
	"webapp/internal/config"
)

type NatsCon struct {
	cfg  config.Config
	con  stan.Conn
	sub  stan.Subscription
	next processor
}

type processor interface {
	ProcNatsMessages([]byte)
}

func InitNatsCon(cfg config.Config, clientID string) *NatsCon {
	sc, err := stan.Connect(cfg.NatsCluster, clientID, stan.NatsURL(cfg.Nats))
	if err != nil {
		log.Fatal("(Pubsub): error connecting to nats ", err)
	}
	log.Println("(Pubsub): connected to nats-stream as", clientID)
	return &NatsCon{con: sc, cfg: cfg}
}

func (n *NatsCon) SubscribeSubject(next processor) {
	n.next = next
	ss, err := n.con.Subscribe(
		n.cfg.NatsSubject,
		n.recNatsMsg,
		stan.DurableName(n.cfg.NatsDurable))
	if err != nil {
		log.Fatal("(Pubsub): error subscribing to nats", err, n.cfg.NatsSubject)
	}
	n.sub = ss
}

func (n *NatsCon) Publish(data []byte) {
	log.Println("(Pubsub): publish new message to nats chanel")
	err := n.con.Publish(n.cfg.NatsSubject, data)
	if err != nil {
		log.Println("(Pubsub): error publish message", err)
	}
}

func (n *NatsCon) Close() {
	if n.sub != nil {
		log.Println("(Pubsub): closeing nats subscription", n.con.Close())
	}
}

func (n *NatsCon) recNatsMsg(m *stan.Msg) {
	log.Println("got new msg from nats", m.Size(), m.Timestamp)
	n.next.ProcNatsMessages(m.Data)
}
