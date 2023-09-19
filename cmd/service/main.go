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

	"github.com/zxq97/gokit/pkg/cache/xredis"
	"github.com/zxq97/gokit/pkg/config"
	"github.com/zxq97/gokit/pkg/database/mysql"
	"github.com/zxq97/gokit/pkg/etcd"
	"github.com/zxq97/gokit/pkg/rpc"
	"github.com/zxq97/rank/dal"
	"github.com/zxq97/rank/idl/rank"
)

var (
	apiDAL   *dal.DAL
	flagConf string
	appConf  conf
)

func init() {
	flag.StringVar(&flagConf, "conf", "api.yaml", "config path, eg: -conf config.yaml")
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
	apiDAL = dal.NewDAL(dbCli, redisCli)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	svc, err := rpc.NewGrpcServer(ctx, appConf.Server, etcdCli)
	if err != nil {
		panic(err)
	}

	rank.RegisterRankServer(svc, &rankService{})
	lis, err := net.Listen("tcp", appConf.Server.Bind)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/rank/get", HandleGetUserRank)
	http.HandleFunc("/rank/get_low", HandleGetLowRankList)
	http.HandleFunc("/rank/get_high", HandleGetHighRankList)

	errCh := make(chan error, 1)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	go func() {
		// grpc服务启动
		errCh <- svc.Serve(lis)
	}()
	go func() {
		// http服务启动
		// todo 超时控制 链路追踪信息
		errCh <- http.ListenAndServe(appConf.Server.HttpBind, nil)
	}()

	select {
	case err = <-errCh:
		svc.Stop()
		cancel()
		log.Println("service stop err", err)
	case sig := <-sigCh:
		svc.Stop()
		cancel()
		log.Println("service stop sign", sig)
	}
}
