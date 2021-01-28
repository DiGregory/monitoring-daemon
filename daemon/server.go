package daemon

import (
	pb "../proto"
	"context"
	"../parser"
	"net"
	"log"
	"google.golang.org/grpc"
	"fmt"
)

//
type server struct {
	port       int32
	statistics []int32
}

var statisticStorage []parser.ParsedData



//создание сервера
func CreateDaemon(port int32, indeedStatistics []int32) {
	//создание сервера
	grpcServer := server{port: port, statistics: indeedStatistics}
	s := grpc.NewServer()
	pb.RegisterStreamServiceServer(s, grpcServer)
	// create listener
	lis, err := net.Listen("tcp", ":"+fmt.Sprintf("%d", grpcServer.port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	dataChan := make(chan parser.ParsedData)
	//собираем статистику
	go func() {
		parser.ParseStatistic(indeedStatistics, dataChan, ctx)
		for {
			data := <-dataChan
			statisticStorage = append(statisticStorage, data)
		}
	}()

	log.Println("start server")
	if err := s.Serve(lis); err != nil {
		cancel()
		log.Fatalf("failed to serve: %v", err)
	}

}
