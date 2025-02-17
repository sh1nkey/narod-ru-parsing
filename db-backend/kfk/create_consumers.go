package kfk

import (
	"context"
	"data-sender/core"
	"data-sender/core/parsenarod"
	"time"
	"github.com/IBM/sarama"
	"github.com/rs/zerolog/log"
)

type KfkService interface {
	Create(ctx context.Context, url *parsenarod.UrlReqDTO) error
	MarkAsEmpty(ctx context.Context, url string, options ...core.UpdateOptions) error
	SetDescription(ctx context.Context, url string, description string, options ...core.UpdateOptions) error
}

type Serve func(p *ServerParams)

type ConsumerProducerConfig struct {
	ConsumerTopicName string
	Server            Serve
	Service           KfkService
	Host              string
}

func CreateProducerConsumer(p *ConsumerProducerConfig) sarama.ConsumerGroup {
	group, err := createConsGroup(p.Host, p.ConsumerTopicName)
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

	handler := consumerHandler{
		service:  p.Service,
		server:   p.Server,
		producer: producer,
	}

	go consumeGroup(&group, p.ConsumerTopicName, &handler)

	log.Info().Msg("Kafka-consumer для HTML-парсинга успешно запущен")
	return group
}

func createConsGroup(host string, topicName string) (sarama.ConsumerGroup, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()

	log.Info().Msgf(`Создаём consumer group с host=%s`, host)

	brokers := []string{host}
	master, err := sarama.NewConsumerGroup(brokers, topicName + "-0", config)
	if err != nil {
		log.Error().Err(err).Msg("Не смогли создать consumer-group")
		return nil, err
	}

	log.Info().Msg("Группа успешно создалась")

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

// type Serve func(service KfkService, producer *sarama.AsyncProducer, event *sarama.ConsumerMessage, signalsProducer *chan os.Signal)

// type myConsumer struct {
//     serve Serve
//     kfkServer KfkService
//     consumer sarama.PartitionConsumer
//     signals chan os.Signal
//     signalsProducer chan os.Signal
//     producer *sarama.AsyncProducer
//     doneCh chan struct{}
// }

// func (mc *myConsumer) Consume() {
//     for {
//         select {
//         case err := <-mc.consumer.Errors():
//             log.Err(err).Msg("Consumer error")
//         case msg := <-mc.consumer.Messages():
//             log.Info().Msg("Received messages. Key: " + string(msg.Key) + " Value: " + string(msg.Value))
//             go mc.serve(mc.kfkServer, mc.producer, msg, &mc.signalsProducer)
//         case <-mc.signals:
//             log.Info().Msg("Interrupt is detected")
//             mc.doneCh <- struct{}{}
//         }
//     }
// }

// func CreateConsumer(topicName string, server Serve, kfkServer KfkService, producer *sarama.AsyncProducer, host string) myConsumer {
//     masterConsumer, err := createMasterConsumer(host)
//     if err != nil {
//         log.Fatal().Err(err).Msg("Couldn't create Kafka Master consumer")
//     }
//     defer func() {
// 		if err := masterConsumer.Close(); err != nil {
// 			log.Error().Err(err).Msg("Error stopping consumer")
//             return
// 		}
// 	}()

//     partitionConsumer, err := consumePartition(masterConsumer, topicName)
//     if err != nil {
//         log.Fatal().Err(err).Msg("Couldn't create Kafka partition consumer")
//     }
//     defer func() {
// 		if err := partitionConsumer.Close(); err != nil {
// 			log.Error().Err(err).Msg("Error stopping partition consumer")
//             return
// 		}
// 	}()

//     signals := make(chan os.Signal, 1)
//     signal.Notify(signals, os.Interrupt)

//     signalsProducer := make(chan os.Signal, 1)
//     signalsProducer.Notify(signals, os.Interrupt)

//     doneCh := make(chan struct{})

//     return myConsumer{
//         serve: server,
//         kfkServer: kfkServer,
//         consumer: partitionConsumer,
//         signals: signals,
//         doneCh: doneCh,
//         signalsProducer: signalsProducer,
//     }
// }

// func createMasterConsumer(host string) (sarama.Consumer, error) {
//     config := sarama.NewConfig()
//     config.Consumer.Return.Errors = true

//     brokers := []string{host}
//     master, err := sarama.NewConsumer(brokers, config)
//     if err != nil {
//         log.Error().Err(err).Msg("Couldn't create Kafka consumer")
//         return nil, err
//     }

//     return master, nil
// }

// func consumePartition(master sarama.Consumer, topicName string) (sarama.PartitionConsumer, error) {
//     consumer, err := master.ConsumePartition(topicName, int32(sarama.OffsetNewest), sarama.OffsetOldest)
//     if err != nil {
//         return nil, err
//     }
//     log.Info().Msg(fmt.Sprintf(`Конмюсер для $1 успешно запущен`, topicName))
//     return consumer, nil
// }
