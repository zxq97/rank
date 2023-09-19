package dal

import (
	"context"
	"fmt"

	"golang.org/x/sync/singleflight"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

const (
	tableRankScore   = "rank_score"
	tableRankScoreTx = "rank_score_tx"

	sfKeyGetRank = "rank_%d_%d" // rankid uid
)

var (
	sfg singleflight.Group
)

type rankScore struct {
	RankID int32 `db:"rank_id"`
	UID    int64 `db:"uid"`
	Score  int32 `db:"score"`
	Rank   int32 `db:"rank"`
}

func (dal *DAL) delUser(ctx context.Context, rankID int32, uid int64) error {
	sql := fmt.Sprintf("UPDATE %s SET `state` = 0 WHERE `rank_id` = ? AND `uid` = ? AND `state` = 1 LIMIT 1", tableRankScore)
	_, err := dal.db.ExecContext(ctx, sql, rankID, uid)
	return err
}

func (dal *DAL) deltaScore(ctx context.Context, rankID, score int32, uid int64, txID string) error {
	return dal.db.Tx(ctx, func(sess sqlbuilder.Tx) error {
		sql := fmt.Sprintf("INSERT INTO %s (`tx_id`, `rank_id`, `uid`, `delta_score`) VALUES (?, ?, ?, ?)", tableRankScoreTx)
		_, err := dal.db.ExecContext(ctx, sql, txID, rankID, uid, score)
		if err != nil {
			return err
		}
		sql = fmt.Sprintf("INSERT INTO %s (`rank_id`, `uid`, `score`, `state`) VALUES (?, ?, ?, 1) ON DUPLICATE KEY UPDATE `score` = `score` + ?", tableRankScore)
		_, err = dal.db.ExecContext(ctx, sql, rankID, uid, score, score)
		return err
	})
}

// rand字段只有活动结束后才会填充
// 防止缓存失效当量请求涌入db 利用singleflight归并回源来解决
func (dal *DAL) getRank(ctx context.Context, rankID int32, uid int64) (int64, int32, error) {
	key := fmt.Sprintf(sfKeyGetRank, rankID, uid)
	val, err, _ := sfg.Do(key, func() (interface{}, error) {
		rs := &rankScore{}
		if err := dal.db.WithContext(ctx).Select("`rank`", "`score`").From(tableRankScore).Where(db.Cond{"rand_id": rankID, "uid": uid}).One(rs); err != nil {
			return 0, err
		}
		return rs, nil
	})
	if err != nil {
		return 0, 0, err
	}
	rs, ok := val.(*rankScore)
	if !ok {
		return 0, 0, err
	}
	return int64(rs.Rank), rs.Score, nil
}
