package parser

import (
	"regexp"
	"sync"
)

var reSystem = regexp.MustCompile(`load average: (.*)`)

//получение средней нагрузки системы
func GetLoadAverage(outChan chan interface{}, errChan chan error, wg *sync.WaitGroup) () {
	defer wg.Done()
	rawDataString, err := execCommand(parseRequestParams("top -b -d1 -n1"))
	if err != nil {
		errChan <- err
	}
	//поиск строки c нагрузкой
	la := reSystem.FindStringSubmatch(*rawDataString)
	if la == nil {
		errChan <- err
	}
	//извлечение чисел
	loadValues, err := parseDataRow(la[0], floatNumRegExp)
	if err != nil {
		errChan <- err
	}

	outChan <- &LoadAverage{loadValues[0],
		loadValues[1],
		loadValues[2]}
	//fmt.Println("load average was sended",loadValues)
}
