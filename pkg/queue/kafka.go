package queue

import (
	"context"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/akordium-id/waqfwise/pkg/config"
	"go.uber.org/zap"
)

// Kafka wraps Kafka producer and consumer
type Kafka struct {
	producer sarama.SyncProducer
	consumer sarama.ConsumerGroup
	config   *config.KafkaConfig
	logger   *zap.Logger
}

// New creates a new Kafka client
func New(cfg *config.KafkaConfig, log *zap.Logger) (*Kafka, error) {
	// Create Sarama config
	saramaConfig := sarama.NewConfig()
	saramaConfig.Version = sarama.V3_5_0_0
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Retry.Max = 5
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetNewest

	if cfg.AutoOffsetReset == "earliest" {
		saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	// Create producer
	producer, err := sarama.NewSyncProducer(cfg.Brokers, saramaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	// Create consumer group
	consumer, err := sarama.NewConsumerGroup(cfg.Brokers, cfg.GroupID, saramaConfig)
	if err != nil {
		producer.Close()
		return nil, fmt.Errorf("failed to create Kafka consumer group: %w", err)
	}

	if log != nil {
		log.Info("kafka connection established",
			zap.Strings("brokers", cfg.Brokers),
			zap.String("group_id", cfg.GroupID),
		)
	}

	return &Kafka{
		producer: producer,
		consumer: consumer,
		config:   cfg,
		logger:   log,
	}, nil
}

// Produce sends a message to a Kafka topic
func (k *Kafka) Produce(topic string, key, value []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(value),
	}

	partition, offset, err := k.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to produce message: %w", err)
	}

	if k.logger != nil {
		k.logger.Debug("message produced",
			zap.String("topic", topic),
			zap.Int32("partition", partition),
			zap.Int64("offset", offset),
		)
	}

	return nil
}

// ProduceJSON sends a JSON message to a Kafka topic
func (k *Kafka) ProduceJSON(topic string, key string, value interface{}) error {
	// Marshal value to JSON in the calling code
	// This is a placeholder for demonstration
	return fmt.Errorf("not implemented: use Produce with JSON marshaled bytes")
}

// Consume consumes messages from Kafka topics
func (k *Kafka) Consume(ctx context.Context, topics []string, handler ConsumerHandler) error {
	consumer := &consumerGroupHandler{
		handler: handler,
		logger:  k.logger,
	}

	for {
		if err := k.consumer.Consume(ctx, topics, consumer); err != nil {
			return fmt.Errorf("error from consumer: %w", err)
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

// Close closes Kafka connections
func (k *Kafka) Close() error {
	var errs []error

	if err := k.producer.Close(); err != nil {
		errs = append(errs, fmt.Errorf("failed to close producer: %w", err))
	}

	if err := k.consumer.Close(); err != nil {
		errs = append(errs, fmt.Errorf("failed to close consumer: %w", err))
	}

	if k.logger != nil {
		k.logger.Info("kafka connections closed")
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing Kafka: %v", errs)
	}

	return nil
}

// ConsumerHandler is a function that handles consumed messages
type ConsumerHandler func(ctx context.Context, message *sarama.ConsumerMessage) error

// consumerGroupHandler implements sarama.ConsumerGroupHandler
type consumerGroupHandler struct {
	handler ConsumerHandler
	logger  *zap.Logger
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (h *consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (h *consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages()
func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			if message == nil {
				return nil
			}

			start := time.Now()
			err := h.handler(session.Context(), message)
			duration := time.Since(start)

			if err != nil {
				if h.logger != nil {
					h.logger.Error("error handling message",
						zap.String("topic", message.Topic),
						zap.Int32("partition", message.Partition),
						zap.Int64("offset", message.Offset),
						zap.Error(err),
						zap.Duration("duration", duration),
					)
				}
				// Don't mark as consumed on error - will be retried
				continue
			}

			session.MarkMessage(message, "")

			if h.logger != nil {
				h.logger.Debug("message processed",
					zap.String("topic", message.Topic),
					zap.Int32("partition", message.Partition),
					zap.Int64("offset", message.Offset),
					zap.Duration("duration", duration),
				)
			}

		case <-session.Context().Done():
			return nil
		}
	}
}
