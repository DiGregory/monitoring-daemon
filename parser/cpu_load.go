package parser

import (
	"regexp"
	"sync"
)

var reCPU = regexp.MustCompile(`%Cpu\(s\): (.*)`)

type CpuLoad struct {
	UserLoad   float32
	SystemLoad float32
	IDLELoad   float32
}

func GetCPULoad(outChan chan interface{}, errChan chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	rawDataString, err := execCommand(parseRequestParams("top -b -d1 -n1"))
	if err != nil {
		errChan <- err
		outChan <- &CpuLoad{}
		return
	}
	//поиск строки c нагрузкой
	s := reCPU.FindStringSubmatch(*rawDataString)
	if s == nil {
		errChan <- err
		outChan <- &CpuLoad{}
		return

	}
	//извлечение чисел
	loadValues, err := parseDataRow(s[0], floatNumRegExp)
	if err != nil {
		errChan <- err
		outChan <- &CpuLoad{}
		return
	}
	if len(loadValues) < 4 {
		errChan <- err
		outChan <- &CpuLoad{}
		return
	}
	outChan <- &CpuLoad{loadValues[0],
		loadValues[1],
		loadValues[3]}

}
