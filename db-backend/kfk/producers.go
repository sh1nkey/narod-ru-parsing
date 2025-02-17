package kfk

import (
	"encoding/json"
	"os"

	"github.com/IBM/sarama"
	"github.com/rs/zerolog/log"
)

var enqueued int = 0



func setupProducer(host string) (*sarama.AsyncProducer, error){
	config := sarama.NewConfig()
	brokers := []string{host}
	producer, err := sarama.NewAsyncProducer(brokers, config)
	return &producer, err
}


func ProduceMessage[T MarshalableEvent](producer *sarama.AsyncProducer, topic string, msg T, signals *chan *os.Signal) {
	jsonData, err := json.Marshal(msg)
	if err != nil {
		log.Error().Msg("Не смогли реализовать десериализацию сообщения для отправки")
		return
	}
	message := &sarama.ProducerMessage{Topic: topic, Value: sarama.ByteEncoder(jsonData)}

	go func() {
		select {
		case (*producer).Input() <- message:
			enqueued++
			log.Info().Msg("Новое сообщение было запушено")
		case <-*signals:
			log.Error().Msg("Завершаем работу продюссера")
			(*producer).AsyncClose()
			return
		}
	}()
}