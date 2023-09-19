package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/zxq97/gokit/pkg/cast"
	"github.com/zxq97/rank/idl/rank"
)

type rankService struct {
	rank.UnimplementedRankServer
}

func (r *rankService) GetUserRank(ctx context.Context, req *rank.GetRankRequest) (*rank.RankItem, error) {
	rk, score, err := apiDAL.GetRank(ctx, req.RankId, req.Uid)
	if err != nil {
		return nil, err
	}
	return &rank.RankItem{Uid: req.Uid, Rank: rk, Score: score}, nil
}

func (r *rankService) GetLowRankList(ctx context.Context, req *rank.GetRankRequest) (*rank.RankList, error) {
	list, err := apiDAL.GetLowRankList(ctx, req.RankId, req.Uid)
	if err != nil {
		return nil, err
	}
	return &rank.RankList{RankList: list}, nil
}

func (r *rankService) GetHighRankList(ctx context.Context, req *rank.GetRankRequest) (*rank.RankList, error) {
	list, err := apiDAL.GetHighRankList(ctx, req.RankId, req.Uid)
	if err != nil {
		return nil, err
	}
	return &rank.RankList{RankList: list}, nil
}

// 为简单所有参数都从url里面获取
func HandleGetUserRank(w http.ResponseWriter, r *http.Request) {
	vs := r.URL.Query()
	rankID := cast.ParseInt(vs.Get("rank_id"), 0)
	uid := cast.ParseInt(vs.Get("uid"), 0)
	rk, score, err := apiDAL.GetRank(r.Context(), int32(rankID), uid)
	if err != nil {
		// log
		// 返回值省略
		_, _ = w.Write([]byte(""))
		return
	}
	res := &rank.RankItem{Uid: uid, Rank: rk, Score: score}
	bs, err := json.Marshal(res)
	if err != nil {
		// log
		_, _ = w.Write([]byte(""))
		return
	}
	_, _ = w.Write(bs)
}

func HandleGetLowRankList(w http.ResponseWriter, r *http.Request) {
	vs := r.URL.Query()
	rankID := cast.ParseInt(vs.Get("rank_id"), 0)
	uid := cast.ParseInt(vs.Get("uid"), 0)
	list, err := apiDAL.GetLowRankList(r.Context(), int32(rankID), uid)
	if err != nil {
		_, _ = w.Write([]byte(""))
		return
	}
	res := &rank.RankList{RankList: list}
	bs, err := json.Marshal(res)
	if err != nil {
		_, _ = w.Write([]byte(""))
		return
	}
	_, _ = w.Write(bs)
}

func HandleGetHighRankList(w http.ResponseWriter, r *http.Request) {
	vs := r.URL.Query()
	rankID := cast.ParseInt(vs.Get("rank_id"), 0)
	uid := cast.ParseInt(vs.Get("uid"), 0)
	list, err := apiDAL.GetHighRankList(r.Context(), int32(rankID), uid)
	if err != nil {
		_, _ = w.Write([]byte(""))
		return
	}
	res := &rank.RankList{RankList: list}
	bs, err := json.Marshal(res)
	if err != nil {
		_, _ = w.Write([]byte(""))
		return
	}
	_, _ = w.Write(bs)
}
