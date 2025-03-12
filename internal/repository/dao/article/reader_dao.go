package article

import (
	"context"
)

type ReaderDAO interface {
	Upsert(ctx context.Context, art Article) error
}

// 这个代表线上表
type PublishArticle struct {
	//组合Article
	Article
}
