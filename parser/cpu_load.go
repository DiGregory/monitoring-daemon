package parser

import (
	"regexp"
	"sync"
)

var reCPU = regexp.MustCompile(`%Cpu\(s\): (.*)`)

func GetCPULoad(outChan chan interface{}, errChan chan error, wg *sync.WaitGroup) () {
	defer wg.Done()
	rawDataString, err := execCommand(parseRequestParams("top -b -d1 -n1"))
	if err != nil {
		errChan <- err
	}
	//поиск строки c нагрузкой
	s := reCPU.FindStringSubmatch(*rawDataString)
	if s == nil {
		errChan <- err
	}
	//извлечение чисел
	loadValues, err := parseDataRow(s[0], floatNumRegExp)
	if err != nil {
		errChan <- err
	}
	outChan <- &CpuLoad{loadValues[0],
		loadValues[1],
		loadValues[3]}

}
