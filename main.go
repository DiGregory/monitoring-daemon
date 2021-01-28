package main

import (
	"./daemon"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sirupsen/logrus"
)

type config struct {
	ServerCfg struct {
		Port        int32 `yaml:"port" env:"PORT" env-default:"8081"`
		LoadAverage int32 `yaml:"load_average" env:"LOAD_AVERAGE"  `
		CpuLoad     int32 `yaml:"cpu_load" env:"CPU_LOAD" `
		DiskLoad    int32 `yaml:"disk_load" env:"DISK_LOAD"  `
		DiskFree    int32 `yaml:"disk_free" env:"DISK_FREE"  `
		NetStat     int32 `yaml:"net_stat" env:"NET_STAT"  `
		TopTalkers  int32 `yaml:"top_talkers" env:"TOP_TALKERS"  `
	} `yaml:"server"`
}

func main() {
	var cfg config
	err := cleanenv.ReadConfig("config.yaml", &cfg)
	if err != nil {
		logrus.WithError(err).Fatal("Read config error")
	}
	indeedStatistic := []int32{
		cfg.ServerCfg.LoadAverage,
		cfg.ServerCfg.CpuLoad,
		cfg.ServerCfg.DiskLoad,
		cfg.ServerCfg.DiskFree,
		cfg.ServerCfg.TopTalkers,
		cfg.ServerCfg.NetStat,
	}

	daemon.CreateDaemon(cfg.ServerCfg.Port, indeedStatistic)

}
