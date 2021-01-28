package parser

import (
	"sync"
)

func GetNetStat(outChan chan interface{}, errChan chan error, wg *sync.WaitGroup) () {
	defer wg.Done()
	//все слушающие tcp соединения
	rawDataStringTCP, err := execCommand(parseRequestParams("sudo netstat -lntp"))
	if err != nil {
		errChan <- err
	}
	//протокол, адрес, PID для tcp
	netStatNamesTCP, err := parseDataFrameNames(*rawDataStringTCP, textRegExp, 2, -1, []int{0, 3, 6})
	if err != nil {
		errChan <- err
	}

	//все активные udp соединения
	rawDataStringUDP, err := execCommand(parseRequestParams("sudo netstat -lnup"))
	if err != nil {
		errChan <- err
	}
	//протокол, адрес, PID для udp
	netStatNamesUDP, err := parseDataFrameNames(*rawDataStringUDP, textRegExp, 2, -1, []int{0, 3, 5})
	if err != nil {
		errChan <- err
	}

	//tcp соединения
	rawDataStringTCP, err = execCommand(parseRequestParams("netstat -ta"))
	if err != nil {
		errChan <- err
	}
	//состояния соединения
	tcpStates, err := parseDataFrameNames(*rawDataStringTCP, textRegExp, 2, -1, []int{5})
	if err != nil {
		errChan <- err
	}
	states := make(map[string]int32)
	for _, v := range tcpStates[0] {
		if _, ok := states[v]; ok {
			states[v] += 1
		} else {
			states[v] = 1
		}
	}

	outChan <- &NetStat{append(netStatNamesTCP[0], netStatNamesUDP[0]...),
		append(netStatNamesTCP[1], netStatNamesUDP[1]...),
		append(netStatNamesTCP[2], netStatNamesUDP[2]...),
		states}
	//fmt.Println("netstat was sended")
}
