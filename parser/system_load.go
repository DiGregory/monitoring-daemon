package parser

import (
	"regexp"
	"sync"
)

var reSystem = regexp.MustCompile(`load average: (.*)`)

type LoadAverage struct {
	Load1  float32
	Load5  float32
	Load15 float32
}

//получение средней нагрузки системы
func GetLoadAverage(outChan chan interface{}, errChan chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	rawDataString, err := execCommand(parseRequestParams("top -b -d1 -n1"))
	if err != nil {
		errChan <- err
		outChan <- &LoadAverage{}
		return
	}
	//поиск строки c нагрузкой
	la := reSystem.FindStringSubmatch(*rawDataString)
	if la == nil {
		errChan <- err
		outChan <- &LoadAverage{}
		return
	}
	//извлечение чисел
	loadValues, err := parseDataRow(la[0], floatNumRegExp)
	if err != nil {
		errChan <- err
		outChan <- &LoadAverage{}
		return
	}
	if len(loadValues) < 3 {
		errChan <- err
		outChan <- &LoadAverage{}
		return
	}

	outChan <- &LoadAverage{loadValues[0],
		loadValues[1],
		loadValues[2]}
	//fmt.Println("load average was sended",loadValues)
}
