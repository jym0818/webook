package service

import (
	"context"
	"github.com/jym0818/webook/internal/domain"
	"github.com/jym0818/webook/internal/events/article"
	"github.com/jym0818/webook/internal/repository"
	"go.uber.org/zap"
)

type ArticleService interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
	Publish(ctx context.Context, art domain.Article) (int64, error)
	Withdraw(ctx context.Context, article domain.Article) error
	List(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error)
	GetById(ctx context.Context, id int64) (domain.Article, error)
	GetPublishedById(ctx context.Context, id int64, uid int64) (domain.Article, error)
}

type articleService struct {
	repo     repository.ArticleRepository
	producer article.Producer
}

func (svc *articleService) GetById(ctx context.Context, id int64) (domain.Article, error) {
	return svc.repo.GetByID(ctx, id)
}

func (svc *articleService) GetPublishedById(ctx context.Context, id int64, uid int64) (domain.Article, error) {
	// 另一个选项，在这里组装 Author，调用 UserService
	art, err := svc.repo.GetPublishedById(ctx, id)
	if err == nil {
		go func() {
			er := svc.producer.ProduceReadEvent(ctx, article.ReadEvent{Uid: uid, Aid: art.Id})
			if er != nil {
				zap.L().Error("消息计数增加失败", zap.Error(er))
			}
		}()
	}
	return art, err
}

func (svc *articleService) Publish(ctx context.Context, art domain.Article) (int64, error) {
	art.Status = domain.ArticleStatusPublished
	return svc.repo.Sync(ctx, art)
}

func (svc *articleService) List(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error) {
	return svc.repo.List(ctx, uid, offset, limit)
}

func (svc *articleService) Withdraw(ctx context.Context, art domain.Article) error {
	// art.Status = domain.ArticleStatusPrivate 然后直接把整个 art 往下传
	return svc.repo.SyncStatus(ctx, art.Id, art.Author.Id, domain.ArticleStatusPrivate)
}
func (svc *articleService) Save(ctx context.Context, art domain.Article) (int64, error) {
	art.Status = domain.ArticleStatusUnPublished
	if art.Id > 0 {
		err := svc.repo.Update(ctx, art)
		return art.Id, err
	}
	return svc.repo.Create(ctx, art)
}

func NewarticleService(repo repository.ArticleRepository, producer article.Producer) ArticleService {
	return &articleService{repo: repo, producer: producer}
}
