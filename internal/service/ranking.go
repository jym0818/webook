package service

import (
	"context"
	"github.com/ecodeclub/ekit/queue"
	"github.com/ecodeclub/ekit/slice"
	"github.com/jym0818/webook/internal/domain"
	"github.com/jym0818/webook/internal/repository"
	"math"
	"time"
)

type RankingService interface {
	//不会返回[]domain.Article  因为会直接放在redis里面
	TopN(ctx context.Context) error
}
type BatchRankingService struct {
	artSvc    ArticleService
	intrSvc   InteractiveService
	batchSize int
	n         int
	scoreFunc func(t time.Time, likeCnt int64) float64
	repo      repository.RankingRepository
}

func NewBatchRankingService(artSvc ArticleService, intrSvc InteractiveService, repo repository.RankingRepository) RankingService {
	return &BatchRankingService{
		artSvc:    artSvc,
		intrSvc:   intrSvc,
		batchSize: 100,
		n:         100,
		//scoreFunc 不能返回负数
		scoreFunc: func(t time.Time, likeCnt int64) float64 {
			sec := time.Since(t).Seconds()
			return float64(likeCnt-1) / math.Pow(sec+2, 1.5)
		},
		repo: repo,
	}
}
func (svc *BatchRankingService) TopN(ctx context.Context) error {
	//取修改时间小于当前时间的帖子
	now := time.Now()
	offset := 0
	type Score struct {
		art   domain.Article
		score float64
	}
	// 创建优先级队列
	// 这里可以用非并发安全
	// svc.n就是容量 也就是排行榜数量
	// 要传入一个函数，来比较进行排序
	topN := queue.NewConcurrentPriorityQueue[Score](svc.n,
		func(src Score, dst Score) int {
			if src.score > dst.score {
				return 1
			} else if src.score == dst.score {
				return 0
			} else {
				return -1
			}
		})

	for {
		// 这里拿了一批
		arts, err := svc.artSvc.ListPub(ctx, now, offset, svc.batchSize)
		if err != nil {
			return err
		}
		ids := slice.Map[domain.Article, int64](arts,
			func(idx int, src domain.Article) int64 {
				return src.Id
			})
		// 要去找到对应的点赞数据 是一个map
		intrs, err := svc.intrSvc.GetByIds(ctx, "article", ids)
		if err != nil {
			return err
		}
		// 合并计算 score
		// 排序
		for _, art := range arts {
			//也可以省略ok 因为如果没有点赞数据 返回的是一个零值结构体 intr.LikeCnt 为0
			intr := intrs[art.Id]
			//if !ok {
			//	// 你都没有，肯定不可能是热榜
			//	continue
			//}

			score := svc.scoreFunc(art.Utime, intr.LikeCnt)
			// 我要考虑，我这个 score 在不在前一百名  分数要注意负数问题
			// 拿到热度最低的
			err = topN.Enqueue(Score{
				art:   art,
				score: score,
			})
			// 这种写法，要求 topN 已经满了
			if err == queue.ErrOutOfCapacity {
				val, _ := topN.Dequeue()
				if val.score < score {
					err = topN.Enqueue(Score{
						art:   art,
						score: score,
					})
				} else {
					_ = topN.Enqueue(val)
				}
			}
		}

		// 一批已经处理完了，问题来了，我要不要进入下一批？我怎么知道还有没有？
		if len(arts) < svc.batchSize || now.Sub(arts[len(arts)-1].Utime).Hours() > 7*24 {
			// 我这一批都没取够，我当然可以肯定没有下一批了
			// 又或者已经取到了七天之前的数据了，说明可以中断了
			break
		}
		// 这边要更新 offset
		offset = offset + len(arts)
	}
	// 最后得出结果
	res := make([]domain.Article, svc.n)
	//要从最大的开始取
	for i := svc.n - 1; i >= 0; i-- {
		val, err := topN.Dequeue()
		if err != nil {
			// 说明取完了，不够 n
			break
		}
		res[i] = val.art
	}
	return svc.repo.ReplaceTopN(ctx, res)
}
