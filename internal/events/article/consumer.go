package article

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/jym0818/webook/internal/events"
	"github.com/jym0818/webook/internal/repository"
	"github.com/jym0818/webook/pkg/saramax"
	"go.uber.org/zap"
	"time"
)

type ReadEventArticleConsumer struct {
	client sarama.Client
	repo   repository.InteractiveRepository
}

func (r *ReadEventArticleConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive", r.client)
	if err != nil {
		return err
	}
	go func() {
		er := cg.Consume(context.Background(), []string{"read_article"}, saramax.NewHandler[ReadEvent](r.Consume))
		if er != nil {
			zap.L().Error("退出了消费循环", zap.Error(er))
		}
	}()

	return nil
}

func (r *ReadEventArticleConsumer) Consume(msg *sarama.ConsumerMessage, evt ReadEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return r.repo.IncrReadCnt(ctx, "article", evt.Aid)
}

func NewReadEventArticleConsumer(repo repository.InteractiveRepository, client sarama.Client) events.Consumer {
	return &ReadEventArticleConsumer{repo: repo, client: client}
}
