package kfk

import (
	"github.com/IBM/sarama"
)

func ConfigKfk(clean cleanText, server Serve, host string) sarama.ConsumerGroup {
	config := ConsumerProducerConfig{
		ConsumerTopicName: SavedUrlEvent,
		ProducerTopicName: HtmlParsedEvent,
		Server:    ServeHtmlParsed,
		Clean:     clean,
		Host:      host,
	}
	group := CreateProducerConsumer(&config)
	return group
}