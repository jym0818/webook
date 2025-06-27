package saramax

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type Handler[T any] struct {
	fn func(msg *sarama.ConsumerMessage, t T) error
}

func (h *Handler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()
	for msg := range msgs {
		var t T
		err := json.Unmarshal(msg.Value, &t)
		if err != nil {
			zap.L().Error("反序列化消息失败",
				zap.Error(err),
				zap.String("topic", msg.Topic),
				zap.Int64("partition", int64(msg.Partition)),
				zap.Int64("offset", msg.Offset))
			continue
		}
		for i := 0; i < 3; i++ {
			err = h.fn(msg, t)
			if err == nil {
				break
			}
			zap.L().Error("处理消息失败",
				zap.Error(err),
				zap.String("topic", msg.Topic),
				zap.Int64("partition", int64(msg.Partition)),
				zap.Int64("offset", msg.Offset),
			)
		}
		if err != nil {
			zap.L().Error("处理消息失败-重试次数上限",
				zap.Error(err),
				zap.String("topic", msg.Topic),
				zap.Int64("partition", int64(msg.Partition)),
				zap.Int64("offset", msg.Offset))
		} else {
			session.MarkMessage(msg, "")
		}
	}
	return nil
}

func NewHandler[T any](fn func(msg *sarama.ConsumerMessage, t T) error) *Handler[T] {
	return &Handler[T]{fn: fn}
}
