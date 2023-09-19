package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zxq97/gokit/pkg/cache/xredis"
	"github.com/zxq97/gokit/pkg/config"
	"github.com/zxq97/gokit/pkg/database/mysql"
	"github.com/zxq97/gokit/pkg/etcd"
	"github.com/zxq97/gokit/pkg/mq"
	"github.com/zxq97/gokit/pkg/mq/kafka"
	"github.com/zxq97/gokit/pkg/rpc"
	"github.com/zxq97/gokit/pkg/server"
	"github.com/zxq97/gokit/pkg/server/consumer"
	"github.com/zxq97/rank/dal"
	"github.com/zxq97/rank/idl/activity"
)

var (
	jobDAL   *dal.DAL
	flagConf string
	appConf  conf
)

func init() {
	flag.StringVar(&flagConf, "conf", "job.yaml", "config path, eg: -conf config.yaml")
}

func main() {
	flag.Parse()
	err := config.LoadYaml(flagConf, &appConf)
	if err != nil {
		panic(err)
	}
	// etcd用于 服务注册于发现
	etcdCli, err := etcd.NewEtcd(appConf.Etcd)
	if err != nil {
		panic(err)
	}
	dbCli, err := mysql.NewMysqlSess(appConf.Mysql)
	if err != nil {
		panic(err)
	}
	redisCli := xredis.NewXRedis(appConf.Redis)
	jobDAL = dal.NewDAL(dbCli, redisCli)
	// 可在进城内分发的消费组 n为消费的并行度
	scoreConsumer, err := kafka.NewDispatchConsumer(appConf.Kafka, []string{"activity"}, "activity_job", mqDeltaScore, 10)
	if err != nil {
		panic(err)
	}
	delConsumer, err := kafka.NewDispatchConsumer(appConf.Kafka, []string{"activity_admin"}, "activity_job", mqDelUser, 10)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	svc, err := rpc.NewGrpcServer(ctx, appConf.Server, etcdCli)
	if err != nil {
		panic(err)
	}

	activity.RegisterActivityServer(svc, &activityService{})
	lis, err := net.Listen("tcp", appConf.Server.Bind)
	if err != nil {
		panic(err)
	}

	s, err := consumer.NewServer([]mq.Consumer{scoreConsumer, delConsumer}, server.WithStartTimeout(time.Second), server.WithStopTimeout(time.Second))
	if err != nil {
		panic(err)
	}
	// 消费启动
	if err = s.Start(context.Background()); err != nil {
		panic(err)
	}

	errCh := make(chan error, 1)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	go func() {
		// grpc服务启动
		errCh <- svc.Serve(lis)
	}()
	go func() {
		// http服务启动 只用于pprof和普罗米修斯监控
		errCh <- http.ListenAndServe(appConf.Server.HttpBind, nil)
	}()

	select {
	case err = <-errCh:
		serr := s.Stop(context.Background())
		svc.Stop()
		cancel()
		log.Println("job stop err", err, serr)
	case sig := <-sigCh:
		serr := s.Stop(context.Background())
		svc.Stop()
		cancel()
		log.Println("job stop sign", sig, serr)
	}
}
