package main

import (
	"flag"

	"github.com/SHSanderland/EffMobTest/pkg/config"
	"github.com/SHSanderland/EffMobTest/pkg/logger"
	"github.com/SHSanderland/EffMobTest/pkg/server"
)

var configPath = flag.String("config", "", "path to config file")

func main() {
	flag.Parse()

	cfg := config.InitConfig(*configPath)
	log := logger.InitLogger(cfg.Env)
	server.InitServer(log, cfg)
}
