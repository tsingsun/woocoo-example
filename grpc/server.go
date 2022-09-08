package main

import (
	"path/filepath"
	"runtime"

	"github.com/tsingsun/woocoo/pkg/conf"
	"github.com/tsingsun/woocoo/pkg/log"
	"github.com/tsingsun/woocoo/rpc/grpcx"
	_ "github.com/tsingsun/woocoo/rpc/grpcx/registry/etcd3"
)

func main() {
	_, currentFile, _, _ := runtime.Caller(0)
	basedir := filepath.Dir(currentFile)
	cfg := conf.New(conf.WithBaseDir(basedir)).Load()
	log.NewBuiltIn()
	srv := grpcx.New(grpcx.WithConfiguration(cfg.Sub("grpc")))
	if err := srv.Run(); err != nil {
		panic(err)
	}
}
