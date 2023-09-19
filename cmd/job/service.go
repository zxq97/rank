package main

import (
	"context"

	"github.com/zxq97/gokit/pkg/mq"
	"github.com/zxq97/rank/idl/activity"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type activityService struct {
	activity.UnimplementedActivityServer
}

func (s *activityService) DeltaScore(ctx context.Context, req *activity.DeltaScoreRequest) (*emptypb.Empty, error) {
	// todo 入榜的条件判断
	return &emptypb.Empty{}, jobDAL.DeltaScore(ctx, req.RankId, req.Score, req.Uid, req.TxId)
}

func (s *activityService) DelUser(ctx context.Context, req *activity.DelUserRequest) (*emptypb.Empty, error) {
	// todo 权限校验
	return &emptypb.Empty{}, jobDAL.DelUser(ctx, req.RankId, req.Uid)
}

func mqDeltaScore(ctx context.Context, msg *mq.MqMessage) error {
	// todo 入榜的条件判断
	req := &activity.DeltaScoreRequest{}
	if err := proto.Unmarshal(msg.Message, req); err != nil {
		return err
	}
	return jobDAL.DeltaScore(ctx, req.RankId, req.Score, req.Uid, req.TxId)
}

func mqDelUser(ctx context.Context, msg *mq.MqMessage) error {
	// todo 权限校验
	req := &activity.DelUserRequest{}
	if err := proto.Unmarshal(msg.Message, req); err != nil {
		return err
	}
	return jobDAL.DelUser(ctx, req.RankId, req.Uid)
}
