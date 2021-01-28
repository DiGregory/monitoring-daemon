package parser

import (
	"sync"
)

func GetDiskFree(outChan chan interface{}, errChan chan error, wg *sync.WaitGroup) () {
	defer wg.Done()
	rawDataStringMB, err := execCommand(parseRequestParams("df -BM"))
	if err != nil {
		errChan <- err
	}
	//имена файловых систем
	diskDataNames, err := parseDataFrameNames(*rawDataStringMB, textRegExp, 1, 16, []int{0})
	if err != nil {
		errChan <- err
	}
	// диск в mb
	diskDataFrameMB, err := parseDataFrame(*rawDataStringMB, integerNumRegExp, 1, 16, []int{2, 3})
	if err != nil {
		errChan <- err
	}

	rawDataStringINode, err := execCommand(parseRequestParams("df -i"))
	if err != nil {
		errChan <- err
	}
	// диск в inode
	diskDataFrameINode, err := parseDataFrame(*rawDataStringINode, integerNumRegExp, 1, 16, []int{2, 3})
	if err != nil {
		errChan <- err
	}

	outChan <- &DiskFree{diskDataNames[0],
		diskDataFrameMB[0],
		diskDataFrameMB[1],
		diskDataFrameINode[0],
		diskDataFrameINode[1],
	}
	// fmt.Println("diskfree was sended" )

}
