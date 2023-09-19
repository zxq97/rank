package dal

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/zxq97/gokit/pkg/cast"
)

const (
	redisKeyZRank     = "act_%d" // randid 排行榜
	redisKeySBlackSet = "blc_%d" // rankid 被踢榜的黑名单
	endtime           = 0        // 活动结束时间 可以从其他地方获取到
)

// 踢榜 加入黑名单
func (dal *DAL) cacheDelUser(ctx context.Context, rankID int32, uid int64) error {
	key := fmt.Sprintf(redisKeyZRank, rankID)
	bkey := fmt.Sprintf(redisKeySBlackSet, rankID)
	pipe := dal.redis.Pipeline()
	pipe.ZRem(ctx, key, uid)
	pipe.SAdd(ctx, bkey, uid)
	_, err := pipe.Exec(ctx)
	return err
}

// 因为同分数需要按时间排序 zset整数部分存储真是分数 小数存储活动结束时间到当前时间的差值
// 结合zset天然有序 可以解决
func (dal *DAL) cacheDeltaScore(ctx context.Context, rankID, score int32, uid int64) error {
	lua := redis.NewScript(`
        local score = redis.call("ZIncrby", KEYS[1], ARGV[2], ARGV[1])
        if score then
            score = math.floor(score)
        else
            score = 0
        end
        score = score + ARGV[3]
        return redis.call("ZAdd", KEYS[1], score, ARGV[1])
    `)
	key := fmt.Sprintf(redisKeyZRank, rankID)
	now := time.Now().Unix()
	return lua.Run(ctx, dal.redis, []string{key}, uid, score, float64(endtime-now)/1e10).Err()
}

func (dal *DAL) checkBlack(ctx context.Context, rankID int32, uid int64) (bool, error) {
	key := fmt.Sprintf(redisKeySBlackSet, rankID)
	return dal.redis.SIsMember(ctx, key, uid).Result()
}

func (dal *DAL) cacheGetRank(ctx context.Context, rankID int32, uid int64) (int64, int32, error) {
	key := fmt.Sprintf(redisKeyZRank, rankID)
	rank, err := dal.redis.ZRank(ctx, key, cast.FormatInt(uid)).Result()
	if err != nil && err != redis.Nil {
		return 0, 0, err
	} else if err == redis.Nil {
		return 0, 0, nil
	}
	score, err := dal.redis.ZScore(ctx, key, cast.FormatInt(uid)).Result()
	return rank, int32(score), err
}

// 通过lua保证原子性 先获取该用户在榜上的名次 之后再通过名次的偏移获得前/后用户的数据
// 如果不在榜用户 直接展示当前头部用户
// 如果不用lua取到的名次再之后通过名次获取上下数据时可能已经发生了榜单的变化 造成数据不准确 下同
func (dal *DAL) lowRankList(ctx context.Context, rankID int32, uid int64) ([]redis.Z, error) {
	key := fmt.Sprintf(redisKeyZRank, rankID)
	lua := redis.NewScript(`
		local rank = redis.call("ZRevRank", KEYS[1], ARGV[1])
		if not rank then
			rank = 0
		else
			rank = rank + 1
		end
		return redis.call("ZRevRange", KEYS[1], rank, rank + ARGV[2], 'withscores')
	`)
	res, err := lua.Run(ctx, dal.redis, []string{key}, uid, 10).Result()
	if err != nil {
		return nil, err
	}
	val, ok := res.([]interface{})
	if !ok || len(val) == 0 || len(val)&1 != 0 {
		return nil, redis.Nil
	}
	zs := make([]redis.Z, 0, len(val)>>1)
	for i := 0; i < len(val); i += 2 {
		id := val[i+1].(string)
		zs = append(zs, redis.Z{Member: val[i], Score: float64(cast.ParseInt(id, 0))})
	}
	return zs, nil
}

func (dal *DAL) highRankList(ctx context.Context, rankID int32, uid int64) ([]redis.Z, error) {
	key := fmt.Sprintf(redisKeyZRank, rankID)
	lua := redis.NewScript(`
		local rank = redis.call("ZRevRank", KEYS[1], ARGV[1])
		if not rank then
			rank = 0
		else
			rank = rank + 1
		end
		return redis.call("ZRevRange", KEYS[1], rank - ARGV[2], rank, 'withscores')
	`)
	res, err := lua.Run(ctx, dal.redis, []string{key}, uid, 10).Result()
	if err != nil {
		return nil, err
	}
	val, ok := res.([]interface{})
	if !ok || len(val) == 0 || len(val)&1 != 0 {
		return nil, redis.Nil
	}
	zs := make([]redis.Z, 0, len(val)>>1)
	for i := 0; i < len(val); i += 2 {
		id := val[i+1].(string)
		zs = append(zs, redis.Z{Member: val[i], Score: float64(cast.ParseInt(id, 0))})
	}
	return zs, nil
}
