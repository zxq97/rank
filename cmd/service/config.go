package main

import (
	"github.com/zxq97/gokit/pkg/cache/xredis"
	"github.com/zxq97/gokit/pkg/database/mysql"
	"github.com/zxq97/gokit/pkg/etcd"
	"github.com/zxq97/gokit/pkg/rpc"
)

type conf struct {
	Server *rpc.Config    `yaml:"server"`
	Etcd   *etcd.Config   `yaml:"etcd"`
	Redis  *xredis.Config `yaml:"redis"`
	Mysql  *mysql.Config  `yaml:"mysql"`
}
