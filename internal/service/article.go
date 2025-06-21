package service

import (
	"context"
	"github.com/jym0818/webook/internal/domain"
	"github.com/jym0818/webook/internal/repository"
)

type ArticleService interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
}

type articleService struct {
	repo repository.ArticleRepository
}

func (svc *articleService) Save(ctx context.Context, art domain.Article) (int64, error) {
	if art.Id > 0 {
		err := svc.repo.Update(ctx, art)
		return art.Id, err
	}
	return svc.repo.Create(ctx, art)
}

func NewarticleService(repo repository.ArticleRepository) ArticleService {
	return &articleService{repo: repo}
}
