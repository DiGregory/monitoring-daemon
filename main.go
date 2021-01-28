package main

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sirupsen/logrus"
	"./cmd/server"
)

type config struct {
	ServerCfg struct {
		Port        int32 `yaml:"port" env:"PORT" env-default:"8081"`
		LoadAverage bool  `yaml:"load_average" env:"LOAD_AVERAGE"  `
		CpuLoad     bool  `yaml:"cpu_load" env:"CPU_LOAD" `
		DiskLoad    bool  `yaml:"disk_load" env:"DISK_LOAD"  `
		DiskFree    bool  `yaml:"disk_free" env:"DISK_FREE"  `
		NetStat     bool  `yaml:"net_stat" env:"NET_STAT"  `
		TopTalkers  bool  `yaml:"top_talkers" env:"TOP_TALKERS"  `
	} `yaml:"server"`
}

func main() {
	var cfg config
	err := cleanenv.ReadConfig("config.yaml", &cfg)
	if err != nil {
		logrus.WithError(err).Fatal("Read config error")
	}
	indeedStatistic := []bool{
		cfg.ServerCfg.LoadAverage,
		cfg.ServerCfg.CpuLoad,
		cfg.ServerCfg.DiskLoad,
		cfg.ServerCfg.DiskFree,
		cfg.ServerCfg.NetStat,
		cfg.ServerCfg.TopTalkers,
	}

	server.Listen(cfg.ServerCfg.Port, indeedStatistic)

}
