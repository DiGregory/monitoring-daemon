package parser

import (
	"regexp"
	"strconv"
	"sync"
)

var reStringDouble = regexp.MustCompile(`.*\n.*\n`)
var reBytes = regexp.MustCompile(`bytes_sent:(\d+) `)
var reBPS = regexp.MustCompile(`send (\d+.\d+)`)
var reProtocol = regexp.MustCompile(`udp|tcp`)
var reAddress = regexp.MustCompile(`\w+\.\w+\.\w+\.\S+`)

type TopTalkers struct {
	ProtocolBytes map[string]int32
	Source        []string
	Destination   []string
	Protocol      []string
	BPS           []float32
}

func GetTopTalkers(outChan chan interface{}, errChan chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	var tt TopTalkers
	//получение процессов
	rawDataString, err := execCommand(parseRequestParams("ss -itupn"))
	if err != nil {
		errChan <- err
		outChan <- &tt
		return
	}
	//считывание строк по две, кроме заголовка
	rawRows := reStringDouble.FindAllString(*rawDataString, -1)[1:]
	for _, r := range rawRows {
		//протокол
		protocol := reProtocol.FindString(r)
		tt.Protocol = append(tt.Protocol, protocol)
		//количество посланных  байт
		sentBytes := reBytes.FindStringSubmatch(r)
		if len(sentBytes) < 2 {
			continue
		}
		sb, err := strconv.Atoi(sentBytes[1])
		if err != nil {
			errChan <- err
			outChan <- &tt
			return
		}
		//байты по протоколам
		tt.ProtocolBytes = make(map[string]int32)
		if _, ok := tt.ProtocolBytes[protocol]; ok {
			tt.ProtocolBytes[protocol] += int32(sb)
		} else {
			tt.ProtocolBytes[protocol] = int32(sb)
		}
		//BPS
		sentBPS := reBPS.FindStringSubmatch(r)
		if len(sentBPS) < 2 {
			errChan <- nilOutputErr
			outChan <- &tt
			return
		}
		bps, err := strconv.ParseFloat(sentBPS[1], 64)
		if err != nil {
			errChan <- err
			outChan <- &tt
			return
		}
		tt.BPS = append(tt.BPS, float32(bps))
		//адреса отправки и назначения
		addresses := reAddress.FindAllString(r, -1)
		if len(addresses) < 2 {
			errChan <- nilOutputErr
			outChan <- &tt
			return
		}
		tt.Source = append(tt.Source, addresses[0])
		tt.Destination = append(tt.Destination, addresses[1])

	}

	outChan <- &tt
	//fmt.Println("tt was sended:", tt)
}
