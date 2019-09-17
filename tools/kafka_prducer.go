package tools

import (
	"encoding/json"
	"log"

	"github.com/zhangdaoling/marketmatchengine/order"

	"github.com/Shopify/sarama"
)

func producer() {
	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, nil)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()
	var orders []*order.Order
	var b []byte
	for _, o := range orders {
		b, err = json.Marshal(o)
		if err != nil {
			log.Printf("json marshal error:", err)
			return
		}
		msg := &sarama.ProducerMessage{
			Topic: "usdt-btc",
			Value: sarama.ByteEncoder(b),
		}
		partition, offset, err := producer.SendMessage(msg)
		if err != nil {
			log.Printf("FAILED to send message: %s\n", err)
		} else {
			log.Printf("> message sent to partition %d at offset %d\n", partition, offset)
		}
	}
}
