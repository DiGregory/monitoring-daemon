package main

import (
	pb "github.com/DiGregory/daemon/proto"
	"context"
	"github.com/DiGregory/daemon/parser"
	"net"
	"log"
	"google.golang.org/grpc"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/ilyakaznacheev/cleanenv"
)


type server struct {
	port       int32
	statistics []bool
}

var statistics []parser.ParsedData



// создание сервера
func Listen(port int32, indeedStatistics []bool) {
	// создание сервера
	grpcServer := server{port: port, statistics: indeedStatistics}
	s := grpc.NewServer()
	pb.RegisterStreamServiceServer(s, grpcServer)
	// create listener
	lis, err := net.Listen("tcp", ":"+fmt.Sprintf("%d", grpcServer.port))
	if err != nil {
		logrus.WithError(err).Fatal("failed to listen")
	}
	ctx, cancel := context.WithCancel(context.Background())
	dataChan := make(chan parser.ParsedData)
	// собираем статистику
	go func() {
		parser.ParseStatistic(indeedStatistics, dataChan, ctx)
		for {
			data := <-dataChan
			statistics = append(statistics, data)
		}
	}()

	log.Println("start server")
	if err := s.Serve(lis); err != nil {
		cancel()
		logrus.WithError(err).Fatal("failed to serve")
	}

}



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

	Listen(cfg.ServerCfg.Port, indeedStatistic)

}