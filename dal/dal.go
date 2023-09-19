package dal

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/zxq97/gokit/pkg/cache/xredis"
	"github.com/zxq97/rank/idl/rank"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

type DAL struct {
	db    sqlbuilder.Database
	redis *xredis.XRedis
}

func NewDAL(db sqlbuilder.Database, redis *xredis.XRedis) *DAL {
	return &DAL{db: db, redis: redis}
}

func (dal *DAL) DeltaScore(ctx context.Context, rankID, score int32, uid int64, txID string) error {
	// 先判断是否在黑名单 已踢榜用户不可在入榜 需根据业务而定
	ok, err := dal.checkBlack(ctx, rankID, uid)
	if err != nil || ok {
		return err
	}
	if err = dal.deltaScore(ctx, rankID, score, uid, txID); err != nil {
		return err
	}
	return dal.cacheDeltaScore(ctx, rankID, score, uid)
}

func (dal *DAL) DelUser(ctx context.Context, rankID int32, uid int64) error {
	if err := dal.delUser(ctx, rankID, uid); err != nil {
		return err
	}
	return dal.cacheDelUser(ctx, rankID, uid)
}

func (dal *DAL) GetRank(ctx context.Context, rankID int32, uid int64) (int64, int32, error) {
	rk, score, err := dal.cacheGetRank(ctx, rankID, uid)
	if err != nil {
		// 判断下活动是否结束 如果结束可以做下降级到db拿名次
		// 因为没结束的话 db中rank字段是不会有值的
		now := time.Now().Unix()
		if now > endtime {
			rk, score, err = dal.getRank(ctx, rankID, uid)
			// 如果结束之后依然没有再缓存中 将一个零值加入
			if err == db.ErrNoMoreRows {
				_ = dal.cacheDeltaScore(ctx, rankID, 0, uid)
				return -1, 0, nil
			}
		}
	}
	return rk, score, err
}

func (dal *DAL) GetLowRankList(ctx context.Context, rankID int32, uid int64) ([]*rank.RankItem, error) {
	zs, err := dal.lowRankList(ctx, rankID, uid)
	return zs2item(zs), err
}

func (dal *DAL) GetHighRankList(ctx context.Context, rankID int32, uid int64) ([]*rank.RankItem, error) {
	zs, err := dal.highRankList(ctx, rankID, uid)
	return zs2item(zs), err
}

func zs2item(zs []redis.Z) []*rank.RankItem {
	l := make([]*rank.RankItem, len(zs))
	for i := range zs {
		l[i] = &rank.RankItem{
			Uid:   zs[i].Member.(int64),
			Score: int32(zs[i].Score),
		}
	}
	return l
}
