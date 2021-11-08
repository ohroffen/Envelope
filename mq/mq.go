package mq

import (
	"MyEnvelope/entity"
	"context"
	"encoding/binary"
	"encoding/json"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

var w *kafka.Writer

func Mq_init() {
	w = &kafka.Writer{
		Addr:     kafka.TCP(os.Getenv("KAFKA_HOST")),
		Topic:    os.Getenv("KAFKA_TOPIC"),
		Balancer: &kafka.CRC32Balancer{Consistent: false},
		Async:    true, // wait to be tested
	}
}

func Send_message(envelope *entity.Envelope) {
	msg, err := json.Marshal(envelope)
	if err != nil {
		log.Fatal(err)
	}
	key := make([]byte, 8)
	// put UserID of int64 into 8 byte
	binary.LittleEndian.PutUint64(key, uint64(envelope.UserID))
	w.WriteMessages(context.Background(),
		kafka.Message{
			Key:   key,
			Value: msg,
		},
	)
}
