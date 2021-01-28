package parser

import (
	"sync"
	"time"
	"github.com/sirupsen/logrus"
	"context"
)

var statisticFunctions = []func(chan interface{}, chan error, *sync.WaitGroup){
	GetLoadAverage,
	GetCPULoad,
	GetDiskLoad,
	GetDiskFree,
	GetNetStat,
	GetTopTalkers,

}
//раз в секунду собирает нужную статистику
func ParseStatistic(indeedStatistic []bool, pd chan ParsedData, ctx context.Context) () {
	var wg sync.WaitGroup
	var parsedData ParsedData
	ticker := time.NewTicker(time.Second)

	statisticChan := make(chan interface{}, 1)
	errChan := make(chan error, 1)
	//сколько статистик собирать
	statisticNum := 0
	for _, v := range indeedStatistic {
		if v   {
			statisticNum += 1
		}
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case _ = <-ticker.C:

				wg.Add(statisticNum)
				for i, statisticFunc := range statisticFunctions {
					if indeedStatistic[i] {
						go statisticFunc(statisticChan, errChan, &wg)
					}
				}
				wg.Wait()
				//fmt.Println(parsedData)
				pd <- parsedData

			}
		}
	}()
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case err := <-errChan:
				logrus.Error("parser error: ", err)
			case stat := <-statisticChan:
				switch v := stat.(type) {
				case *LoadAverage:
					parsedData.LoadAverage = *v
				case *CpuLoad:
					parsedData.CpuLoad = *v
				case *DiskLoad:
					parsedData.DiskLoad = *v
				case *DiskFree:
					parsedData.DiskFree = *v
				case *NetStat:
					parsedData.NetStat = *v
				case *TopTalkers:
					parsedData.TopTalkers = *v
				}
			}
		}
	}()

}
