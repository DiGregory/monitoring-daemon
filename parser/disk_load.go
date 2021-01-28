package parser

import (
	"regexp"
	"sync"
)

var reAVG = regexp.MustCompile(`avg-cpu: .*\n (.*)`)

type DiskLoad struct {
	Names     []string
	TPS       []float32
	WriteLoad []float32
	ReadLoad  []float32
	CPU       CpuLoad
}

//получение нагрузки дисков
func GetDiskLoad(outChan chan interface{}, errChan chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	rawDataString, err := execCommand(parseRequestParams("iostat -c -d 1 1"))
	if err != nil {
		errChan <- err
		outChan <- &DiskLoad{}
		return
	}
	//имена дисков
	diskNames, err := parseDataFrameNames(*rawDataString, textRegExp, 6, 13, []int{0})
	if err != nil {
		errChan <- err
		outChan <- &DiskLoad{}
		return
	}
	//получаем массивы нагрузки на чтение/запись
	diskDataFrame, err := parseDataFrame(*rawDataString, floatNumRegExp, 6, 13, []int{0, 1, 2})
	if err != nil {
		errChan <- err
		outChan <- &DiskLoad{}
		return
	}
	//поиск строки c нагрузкой на CPU
	s := reAVG.FindStringSubmatch(*rawDataString)
	if s == nil {
		errChan <- err
		outChan <- &DiskLoad{}
		return
	}
	//извлечение чисел
	loadValues, err := parseDataRow(s[0], floatNumRegExp)
	if err != nil {
		errChan <- err
		outChan <- &DiskLoad{}
		return
	}

	if len(diskNames) < 1 || len(diskDataFrame) < 3 || len(loadValues) < 6 {
		errChan <- err
		outChan <- &DiskLoad{}
		return
	}

	outChan <- &DiskLoad{diskNames[0], diskDataFrame[0], diskDataFrame[1], diskDataFrame[2],
		CpuLoad{loadValues[0],
			loadValues[2],
			loadValues[5]}}

	//fmt.Println("diskload was sended")

}
