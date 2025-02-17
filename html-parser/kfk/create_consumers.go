package kfk

import (
	"context"
	"errors"
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
    log.Info().Msg("Начинаем запускать Kafka-consumer-ы")

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
    config.Net.DialTimeout = 1 * time.Minute

    log.Info().Msgf(`Создаём consumer group с host=%s`, host)

    brokers := []string{host}

    var master sarama.ConsumerGroup
    var err error

    // Цикл для повторного создания consumer group
    for attempt := 0; attempt < 10; attempt++ {
        master, err = sarama.NewConsumerGroup(brokers, "html_parser-0", config)
        if err == nil {
            log.Info().Msg("Master-консюмер успешно создался")
            return master, nil // Успешное создание, возвращаем master
        }

        // Проверяем, что это ошибка, связанная с недоступностью брокеров
        if errors.Is(err, sarama.ErrOutOfBrokers) || err.Error() == "kafka: client has run out of available brokers to talk to" {
            log.Error().Err(err).Msg("Не смогли создать consumer-group. Попытка подключения...")
        } else {
            log.Error().Err(err).Msg("Не смогли создать consumer group по другой причине.")
            return nil, err // Если это другая ошибка, выходим
        }

        // Ждем перед следующей попыткой
        time.Sleep(10 * time.Second)
    }

    return nil, fmt.Errorf("не удалось создать consumer group после %d попыток: %w", 10, err)
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


