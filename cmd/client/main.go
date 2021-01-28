package main

import (
	pb "github.com/DiGregory/daemon/proto"

	"google.golang.org/grpc"
	"io"
	"context"
	"fmt"
	"time"
	"text/tabwriter"
	"os"
	"github.com/sirupsen/logrus"
)

func statisticOut(response *pb.Response) {
	w := new(tabwriter.Writer)
	fmt.Printf("\n\t\t%v\n\n", time.Now().Format(time.Stamp))
	//system load
	if response.NeededStatistic[0] {
		fmt.Printf("System load average: %v\n\n", response.LoadAverage)
	}
	//cpu load
	if response.NeededStatistic[1] {
		fmt.Printf("CPU load average: %v\n\n", response.CpuLoad)
	}
	//disk load
	if response.NeededStatistic[2] {
		fmt.Printf("Disk load: \n")
		fmt.Printf("CPU : %v\n", response.DiskLoad.CPU)
		w.Init(os.Stdout, 5, 0, 1, ' ', tabwriter.AlignRight)
		fmt.Fprintf(w, "Device\t TPS\t WriteLoad\t ReadLoad\t \n")
		for i := range response.DiskLoad.Names {
			fmt.Fprintf(w, "%v\t %v\t %v\t %v\t\n",
				response.DiskLoad.Names[i],
				response.DiskLoad.TPS[i],
				response.DiskLoad.WriteLoad[i],
				response.DiskLoad.ReadLoad[i])
		}
		w.Flush()
	}
	//disk free
	if response.NeededStatistic[3] {
		fmt.Printf("\nDisk free: \n")
		w.Init(os.Stdout, 5, 0, 1, ' ', tabwriter.AlignRight)
		fmt.Fprintf(w, "File system\t   MBFree\t MBUses\t INodeFree \t INodeUses\t \n")
		for i := range response.DiskFree.Names {
			fmt.Fprintf(w, "%v\t %v\t %v\t %v\t %v\t\n",
				response.DiskFree.Names[i],
				response.DiskFree.MBFree[i],
				response.DiskFree.MBUses[i],
				response.DiskFree.INodeFree[i],
				response.DiskFree.INodeUses[i])
		}
		w.Flush()
	}
	//network
	if response.NeededStatistic[4] {
		fmt.Printf("\nNetwork statistic: \n")
		fmt.Printf("\nSockets states: \n")
		for i, v := range response.NetStat.States {
			fmt.Printf("sockets with state %v: %v\n", i, v)
		}
		w.Init(os.Stdout, 10, 4, 1, ' ', tabwriter.AlignRight)
		fmt.Fprintf(w, "\nPID\t    Protocol\t Address\t\n")
		for i := range response.NetStat.PID {
			fmt.Fprintf(w, "%v\t %v\t %v\t \n",
				response.NetStat.PID[i],
				response.NetStat.Protocols[i],
				response.NetStat.Address[i],
			)
		}
		w.Flush()
	}
	//top talkers
	if response.NeededStatistic[5] {
		fmt.Printf("\nTopTalkers: \n")
		fmt.Printf("\nBytes by protocols: \n")
		for i, v := range response.TopTalkers.ProtocolBytes {
			fmt.Printf("Protocol: %v receive bytes: %v", i, v)
		}
		w.Init(os.Stdout, 10, 4, 1, ' ', tabwriter.AlignRight)
		fmt.Fprintf(w, "\nSource\t   Destination\t Protocol\t BPS \n")
		for i := range response.TopTalkers.Source {
			fmt.Fprintf(w, "%v\t %v\t %v\t %v\t\n",
				response.TopTalkers.Source[i],
				response.TopTalkers.Destination[i],
				response.TopTalkers.Protocol[i],
				response.TopTalkers.BPS[i], )
		}
		w.Flush()
	}
}

func main() {
	var requestN int32 = 3  //cooldown
	var requestM int32 = 10 //время усреднения
	// dial server
	conn, err := grpc.Dial(":8080", grpc.WithInsecure())
	if err != nil {
		logrus.WithError(err).Error("can not connect with server")
	}
	ctx := context.Background()
	// create stream
	client := pb.NewStreamServiceClient(conn)
	in := &pb.Request{N: requestN, M: requestM}
	stream, err := client.GetStatistic(ctx, in)
	if err != nil {
		logrus.WithError(err).Fatal("open stream error")
	}

	done := make(chan bool)

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			done <- true //means stream is finished
			return
		}
		if err != nil {
			logrus.WithError(err).Fatal("cannot receive data")
		}
		if resp != nil {
			//логирование
			statisticOut(resp)
		}

	}

	<-done
	logrus.Info("finished")
}
