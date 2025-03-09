package article

import (
	"context"
	"github.com/jym/webook/internal/domain"
)

type ArticleReaderRepository interface {
	//有就更新没有就创建
	Save(ctx context.Context, art domain.Article) (int64, error)
}
