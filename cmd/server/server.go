package server

import (
	pb "../../proto"
	"context"
	"../../parser"
	"net"
	"log"
	"google.golang.org/grpc"
	"fmt"
	"github.com/sirupsen/logrus"
)

//
type server struct {
	port       int32
	statistics []bool
}

var statistics []parser.ParsedData



//создание сервера
func Listen(port int32, indeedStatistics []bool) {
	//создание сервера
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
	//собираем статистику
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
