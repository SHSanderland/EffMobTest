package main

import (
	"context"
	"flag"

	"github.com/SHSanderland/EffMobTest/pkg/config"
	"github.com/SHSanderland/EffMobTest/pkg/logger"
	"github.com/SHSanderland/EffMobTest/pkg/server"
	"github.com/SHSanderland/EffMobTest/pkg/storage/psql"
)

var configPath = flag.String("config", "", "path to config file")

func main() {
	flag.Parse()

	cfg := config.InitConfig(*configPath)
	log := logger.InitLogger(cfg.Env)
	db, err := psql.InitDB(context.TODO(), log, cfg)

	if err != nil {
		panic(err)
	}

	server.InitServer(log, cfg, db)
}
