package kfk

import (
	"encoding/json"
	"os"

	"github.com/IBM/sarama"
	"github.com/rs/zerolog/log"
)

// let it be empty now. later gonna do smth with it
// func ServeFailure(service KfkService, event *sarama.ConsumerMessage) {}

func ServeHtmlParsed(p *ServerParams) {
	log.Info().Msg("Приняли задачу")
	var data SavedUrlEventDTO

	event := p.event

	if err := json.Unmarshal(event.Value, &data); err != nil {
		log.Err(err).Msg("Ошибка при десериализации сообщения, event: ServeRequestedSaveUrl")
		return
	}
	log.Info().Msgf("Десериализовали задачу, вот текст: %s", event.Value)
	
	cleanedText, err := p.service(&data.HtmlContent)
	if err != nil {
		log.Err(err).Msg("Ошибка при очистке HTML-страницы, event: ServeHtmlParsed")
		return
	}

	log.Info().Msgf("Очистили текст. Результат: %s", cleanedText)

	newData := &HtmlParsedTopicEventDTO{
		Url:            data.Url,
		ParsedContent: cleanedText,
		baseEventDTO: baseEventDTO{
			CorrelationUuid: data.CorrelationUuid,
		},
	}
	newData.FillBaseData()

	if p.sess == nil {
		log.Error().Msgf("Почему-то сессия ConsumerGroup оказалась пустой...")
		return
	}

	(*p.sess).MarkMessage(event, "")
	log.Info().Msg("Пометили задачу как выполненную")

	sigCh := make(chan *os.Signal)
	go ProduceMessage(
		p.producer, 
		p.producerTopic,
		newData,
		&sigCh,
	)
}

type ServerParams struct {
	service cleanText
	event *sarama.ConsumerMessage
	sess *sarama.ConsumerGroupSession
	producer *sarama.AsyncProducer
	producerTopic string
}
type server = func(p *ServerParams)

type serveHtmlParsedHandler struct {
	service cleanText
	server server
	producer *sarama.AsyncProducer
	producerTopic string
}


func (h *serveHtmlParsedHandler) Setup(_ sarama.ConsumerGroupSession) error { return nil }
func (h *serveHtmlParsedHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *serveHtmlParsedHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Info().Msgf(`Message topic:%q partition:%d offset:%d value:%s`, msg.Topic, msg.Partition, msg.Offset, msg.Value)

		p := ServerParams{
			service: h.service,
			event: msg,
			producer: h.producer,
			producerTopic: h.producerTopic,
			sess: &sess,
		}
		go h.server(&p)
	}
	return nil
}
