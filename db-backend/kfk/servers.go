package kfk

import (
	"context"
	"data-sender/core/parsenarod"
	"encoding/json"
	"os"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type ServerParams struct {
	service       KfkService
	event         *sarama.ConsumerMessage
	sess          *sarama.ConsumerGroupSession
	producer      *sarama.AsyncProducer
	producerTopic string
}
type server = func(p *ServerParams)

type consumerHandler struct {
	service KfkService
	server server
	producer *sarama.AsyncProducer
}


func (h *consumerHandler) Setup(_ sarama.ConsumerGroupSession) error { return nil }
func (h *consumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *consumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Info().Msgf(`Message topic:%q partition:%d offset:%d\n`, msg.Topic, msg.Partition, msg.Offset)

		p := ServerParams{
			service: h.service,
			event: msg,
			producer: h.producer,
			sess: &sess,
		}
		go h.server(&p)
	}
	return nil
}


// let it be empty now. later gonna do smth with it
// func ServeFailure(service KfkService, producer *sarama.AsyncProducer, event *sarama.ConsumerMessage, signalsProducer *chan os.Signal) {}

func ServeRequestedSaveUrl(p *ServerParams) {
	var data RequestedSaveUrlEventDTO
	if err := json.Unmarshal(p.event.Value, &data); err != nil {
		log.Err(err).Msg("Ошибка при десериализации сообщения, event: ServeRequestedSaveUrl")
		return
	}
	createData := &parsenarod.UrlReqDTO{Url: data.Url}

	err := p.service.Create(context.Background(), createData)
	if err != nil {
		log.Err(err).Msg("Ошибка при создании url, event: ServeRequestedSaveUrl")
	}
	sendData := SavedUrlEventDTO{
		BaseEventDTO: BaseEventDTO{
			EventUuid:       uuid.New(),
			CorrelationUuid: data.BaseEventDTO.CorrelationUuid,
			CreatedAt:       time.Now(),
		},
		Url:         data.Url,
		HtmlContent: data.HtmlContent,
	}
	sigCh := make(chan *os.Signal)
	go ProduceMessage(
		p.producer,
		SavedUrlEvent,
		sendData,
		&sigCh,
	)
}

func ServeRequestedMarkEmpty(p *ServerParams) {
	var data RequestedMarkEmptyEventDTO
	if err := json.Unmarshal(p.event.Value, &data); err != nil {
		log.Err(err).Msg("Ошибка при десериализации сообщения, event: ServeRequestedMarkEmpty")
		return
	}
	go func() {
		err := p.service.MarkAsEmpty(context.Background(), data.Url)
		if err != nil {
			log.Err(err).Msg("Ошибка при создании url, event: ServeRequestedSaveUrl")
		}
	}()

}

func ServeRequestedSetDesc(p *ServerParams) {
	var data RequestedSetDescEventDTO
	if err := json.Unmarshal(p.event.Value, &data); err != nil {
		log.Err(err).Msg("Ошибка при десериализации сообщения, event: ServeRequestedSetDesc")
		return
	}
	go p.service.SetDescription(context.Background(), data.Url, data.Description)
}

// func ServeHtmlParsed(service KfkService, producer *sarama.AsyncProducer, event *sarama.ConsumerMessage, signalsProducer *chan os.Signal) {}

func ServeAiSummarized(p *ServerParams) {
	var data AiSummarizedEventDTO
	if err := json.Unmarshal(p.event.Value, &data); err != nil {
		log.Err(err).Msg("Ошибка при десериализации сообщения, event: ServeAiSummarized")
		return
	}

	go p.service.SetDescription(context.Background(), data.Url, data.SummarizedContent)
}
