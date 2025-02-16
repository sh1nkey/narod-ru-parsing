package kfk

import (
	"context"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/rs/zerolog/log"
)




type cleanText func(htmlContent *string) (string, error)


type Serve func(p *ServerParams)

type ConsumerProducerConfig struct {
    ProducerTopicName string
    ConsumerTopicName string
    Server Serve
    Clean cleanText
    Host string
}


func CreateProducerConsumer(p *ConsumerProducerConfig) sarama.ConsumerGroup {
    group, err := createConsGroup(p.Host)
    if err != nil {
        return nil
    }
    go func() {
		for err := range group.Errors() {
            log.Err(err).Msg("Ошибка при обработке события")
		}
	}()

    producer, err := SetupProducer(p.Host)
    if err != nil {
        log.Err(err).Msg("Не удалось создать продюссера")
    }
    
    handler := serveHtmlParsedHandler{
        service: p.Clean,
        server: p.Server,
        producer: producer,
        producerTopic: p.ProducerTopicName,
    }
    
    go consumeGroup(&group, p.ConsumerTopicName, &handler)

    log.Info().Msg("Kafka-consumer для HTML-парсинга успешно запущен")
    return group
}


func createConsGroup(host string) (sarama.ConsumerGroup, error) {
    config := sarama.NewConfig()
    config.Consumer.Return.Errors = true
    config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()

    log.Info().Msg(fmt.Sprintf(`Создаём consumer group с host=%s`, host))

    brokers := []string{host}
    master, err := sarama.NewConsumerGroup(brokers, "html_parser-1", config)
    if err != nil {
        log.Error().Err(err).Msg("Не смогли создать consumer-group")
        return nil, err 
    }

    log.Info().Msg("Master-консюмер успешно создался")

    return master, nil 
}

func consumeGroup(group *sarama.ConsumerGroup, topicName string, handler sarama.ConsumerGroupHandler) {
    log.Info().Msgf("Подключается к топику %s", topicName)

    ctx := context.Background()
    go func() {
        for {
            select {
            case <-ctx.Done():
                log.Info().Msgf("Завершение работы с топиком %s", topicName)
                return
            default:
                err := (*group).Consume(ctx, []string{topicName}, handler)
                if err != nil {
                    log.Err(err).Msgf("Ошибка при обработке топика %s", topicName)
                    time.Sleep(2 * time.Second)
                }
            }
        }
    }()
    log.Info().Msgf("Подключение к топику %s завершено", topicName)
}


