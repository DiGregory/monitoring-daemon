package main

import (
	"github.com/DiGregory/daemon/parser"
	"github.com/sirupsen/logrus"
	"time"
	pb "github.com/DiGregory/daemon/proto"
)

func (s server) GetStatistic(in *pb.Request, srv pb.StreamService_GetStatisticServer) error {
	//var wg sync.WaitGroup
	m := int(in.M)
	n := int(in.N)
	var pd parser.ParsedData
	for {
		if len(statistics) >= m {
			pd = parser.GetAverageInfo(statistics, m)
		} else {
			continue
		}
		response := pb.Response{
			LoadAverage: &pb.LoadAverage{
				Load1:  pd.LoadAverage.Load1,
				Load5:  pd.LoadAverage.Load5,
				Load15: pd.LoadAverage.Load15,
			},
			CpuLoad: &pb.CpuLoad{
				UserLoad:   pd.CpuLoad.UserLoad,
				SystemLoad: pd.CpuLoad.SystemLoad,
				IDLELoad:   pd.CpuLoad.IDLELoad,
			},
			DiskLoad: &pb.DiskLoad{
				Names:     pd.DiskLoad.Names,
				TPS:       pd.DiskLoad.TPS,
				WriteLoad: pd.DiskLoad.WriteLoad,
				ReadLoad:  pd.DiskLoad.ReadLoad,
				CPU: &pb.CpuLoad{
					UserLoad:   pd.DiskLoad.CPU.UserLoad,
					SystemLoad: pd.DiskLoad.CPU.SystemLoad,
					IDLELoad:   pd.DiskLoad.CPU.IDLELoad,
				},
			},
			DiskFree: &pb.DiskFree{
				Names:     pd.DiskFree.Names,
				MBFree:    pd.DiskFree.MBFree,
				MBUses:    pd.DiskFree.MBUses,
				INodeFree: pd.DiskFree.INodeFree,
				INodeUses: pd.DiskFree.INodeUses,
			},
			TopTalkers: &pb.TopTalkers{
				ProtocolBytes: pd.TopTalkers.ProtocolBytes,
				Source:        pd.TopTalkers.Source,
				Destination:   pd.TopTalkers.Destination,
				Protocol:      pd.TopTalkers.Protocol,
				BPS:           pd.TopTalkers.BPS,
			},
			NetStat: &pb.NetStat{
				Protocols: pd.NetStat.Protocols,
				Address:   pd.NetStat.Address,
				PID:       pd.NetStat.PID,
				States:    pd.NetStat.States,
			},
			NeededStatistic: s.statistics,
		}
		if err := srv.Send(&response); err != nil {
			logrus.WithError(err).Fatal("send err ")
			return err
		}
		time.Sleep(time.Duration(n) * time.Second)

	}
	return nil
}
