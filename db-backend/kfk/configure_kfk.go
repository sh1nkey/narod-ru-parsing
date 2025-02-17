package kfk

import (
	"github.com/rs/zerolog/log"
)

type Conf struct {
	ProducerTopicName string
	Server Serve
	Service KfkService
	Host string
}

var serversMap = map[string]Serve{
	//FailureTopic:            ServeFailure,
	RequestedSaveUrlTopic:   ServeRequestedSaveUrl,
	RequestedMarkEmptyTopic: ServeRequestedMarkEmpty,
	RequestedSetDescTopic:   ServeRequestedSetDesc,
	AiSummarizedTopic:       ServeAiSummarized,
}


func ConfigKfk(service KfkService, host string) {
	log.Info().Msg("Начинаем запускать Kafka-consumer-ы")
	for topic, worker := range serversMap {
		conf := ConsumerProducerConfig{
			ConsumerTopicName: topic,
			Server: worker,
			Service: service,
			Host: host,
		}
		go CreateProducerConsumer(&conf)
	}
	log.Info().Msg("Успешно запустили Kafka-consumer-ы")
}